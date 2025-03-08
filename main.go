package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/wneessen/go-mail"
	"log"
	"os"
	"strconv"
	"time"
)

var dbpool *pgxpool.Pool
var mailClient *mail.Client

func main() {
	// Load .env file for convenience
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	// Connect to the database
	dbpool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	// Create a new mail client
	mailClient, err = mail.NewClient(os.Getenv("EMAIL_HOST"), mail.WithSSLPort(false), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(os.Getenv("EMAIL_USERNAME")), mail.WithPassword(os.Getenv("EMAIL_PASSWORD")))
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}
	defer mailClient.Close()

	// Set up the table
	_, err = dbpool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS mail_receipts (id UUID DEFAULT gen_random_uuid() PRIMARY KEY, name TEXT NOT NULL , email TEXT DEFAULT NULL, created_by TEXT NOT NULL , created_at TEXT DEFAULT to_char(CURRENT_TIMESTAMP, 'YYYY-MM-DD HH24:MI:SS'), pixel_events JSONB[] DEFAULT '{}', url_events JSONB[] DEFAULT '{}')")
	if err != nil {
		log.Fatalf("unable to setup table: %v\n", err)
	}

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		AppName: "mailReceipt",
		Prefork: true,
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
	//app.Use(recover.New())

	// Get tracking info
	app.Get("/track/:id", getTrackerInfo)

	// Serves a 1x1 transparent pixel and logs the request
	app.Get("/track/:id/pixel", pixelTrack)

	// Redirects to url and logs the request
	app.Get("/track/:id/url/:url", urlTrack)

	// Creates a new row for tracking
	app.Post("/track", newTracker)

	// Deletes a row for tracking
	app.Delete("/track/:id", deleteTracker)

	// hehe
	app.Get("/teapot", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusTeapot)
	})

	// Start the server on :3000
	err = app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
