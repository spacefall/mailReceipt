package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"time"
)

// this is the image: data:image/gif;base64,R0lGODlhAQABAIAAAP///wAAACH5BAEAAAAALAAAAAABAAEAAAICRAEAOw==
var img = []byte{71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 255, 255, 255, 0, 0, 0, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 2, 68, 1, 0, 59}

func pixelTrack(c *fiber.Ctx) error {
	trackingJson := TrackData{
		Ip:        c.IP(),
		UserAgent: c.Get(fiber.HeaderUserAgent),
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	// Add request to db asynchronously
	go func(id string, data TrackData) {
		// Log info
		// log.Printf("New Request\nID: %s\nIP: %s\n", id, ip)
		// Check for valid UUID
		if err := uuid.Validate(id); err != nil {
			//log.Println("Invalid UUID: " + id)
			log.Errorf("Couldn't update tracking table, as an invalid id (%s) was passed", id)
			return
		}

		// Insert into DB
		_, err := dbpool.Exec(context.Background(), `UPDATE mail_receipts SET pixel_events = array_append(pixel_events, $2) WHERE id = $1`, id, data)
		if err != nil {
			log.Errorf("Couldn't update tracking for id %s: %v", id, err)
			return
		}
	}(c.Params("id"), trackingJson)

	// Disable caching
	c.Set("Cache-Control", "max-age=0, no-cache, must-revalidate, proxy-revalidate")
	c.Set("Expires", "0")

	// Content -> gif
	c.Set("Content-Type", "image/gif")

	// Send the pixel
	return c.Send(img)
}
