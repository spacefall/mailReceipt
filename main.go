package main

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

var dbpool *pgxpool.Pool

func main() {
	// Load .env file for convenience
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	dbpool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	// Set up the table
	_, err = dbpool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS mail_receipts (id UUID DEFAULT gen_random_uuid() PRIMARY KEY, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, pixel_events TEXT[][] DEFAULT '{}')")
	if err != nil {
		log.Fatalf("Unable to setup table: %v\n", err)
	}

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		AppName:     "mailReceipt",
		Prefork:     true,
		JSONDecoder: sonic.Unmarshal,
		JSONEncoder: sonic.Marshal,
	})

	// basic Server Timing API support, might put under a flag/remove
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		c.Append("Server-Timing", "app;dur="+strconv.FormatFloat(float64(duration)/1000000, 'f', -1, 64))
		return err
	})

	// Recover from panics
	app.Use(recover.New())

	// Serves a 1x1 transparent pixel for tracking
	app.Get("/pixel/:uuid?", pixelTrack)

	// Creates a new row for tracking
	app.Post("/new", newTracker)

	// Deletes a row for tracking
	app.Delete("/delete/:uuid?", removeTracking)

	// Start the server on :3000
	err = app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
