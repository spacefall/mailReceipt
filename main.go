package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

// this is the image: data:image/gif;base64,R0lGODlhAQABAIAAAP///wAAACH5BAEAAAAALAAAAAABAAEAAAICRAEAOw==
var img = []byte{71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 255, 255, 255, 0, 0, 0, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 2, 68, 1, 0, 59}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}
	defer dbpool.Close()

	_, err = dbpool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS mail_receipts (id UUID DEFAULT gen_random_uuid() PRIMARY KEY, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, pixel_events TEXT[][] DEFAULT '{}')")
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Unable to setup table: %v\n", err)
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}

	app := fiber.New(fiber.Config{AppName: "mailReceipt"})

	// Server Timing API support, might put under a flag later
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		c.Append("Server-Timing", "app;dur="+strconv.FormatFloat(float64(duration)/1000000, 'f', -1, 64))
		return err
	})

	app.Use(recover.New())

	// Serves a 1x1 transparent pixel for tracking
	app.Get("/pixel/:uuid?", func(c *fiber.Ctx) error {
		go func(id string, ip string, ua string) {
			// Log info
			log.Printf("New Request: \nID: %s\nIP: %s\nUser-Agent:%s\n", id, ip, ua)
			// Check for valid UUID
			if err := uuid.Validate(id); err != nil {
				log.Println("Invalid UUID: " + id)
				return
			}
			// Insert into DB
			_, err := dbpool.Exec(context.Background(), "UPDATE mail_receipts SET pixel_events = pixel_events || ARRAY[ARRAY[$1, $2, to_char(CURRENT_TIMESTAMP, 'YYYY-MM-DD HH24:MI:SS.US')]] WHERE id = $3", ip, ua, id)
			if err != nil {
				panic(fmt.Errorf("Unable to insert into database: %v\n", err))
			}
		}(c.Params("uuid"), c.IP(), c.Get(fiber.HeaderUserAgent))

		// Disable caching
		c.Set(fiber.HeaderCacheControl, "max-age=0, no-cache, must-revalidate, proxy-revalidate")
		c.Set(fiber.HeaderExpires, "0")
		// Content -> gif
		c.Set(fiber.HeaderContentType, "image/gif")
		// Send the pixel
		return c.Send(img)
	})

	app.Post("/new", func(c *fiber.Ctx) error {
		var uuid string
		err := dbpool.QueryRow(context.Background(), "INSERT INTO mail_receipts DEFAULT VALUES RETURNING id").Scan(&uuid)
		if err != nil {
			return err
		}
		return c.SendString("The new UUID is: " + uuid)
	})

	err = app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
