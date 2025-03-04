package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

// this is the image: data:image/gif;base64,R0lGODlhAQABAIAAAP///wAAACH5BAEAAAAALAAAAAABAAEAAAICRAEAOw==
var img = []byte{71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 255, 255, 255, 0, 0, 0, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 2, 68, 1, 0, 59}

func pixelTrack(c *fiber.Ctx) error {
	// Add request to db asynchronously
	go func(id string, ip string, ua string) {
		// Log info
		// log.Printf("New Request\nID: %s\nIP: %s\n", id, ip)
		// Check for valid UUID
		if err := uuid.Validate(id); err != nil {
			//log.Println("Invalid UUID: " + id)
			log.Error("Couldn't update tracking table, as an invalid id was passed")
			return
		}
		// Insert into DB
		_, err := dbpool.Exec(context.Background(), "UPDATE mail_receipts SET pixel_events = pixel_events || ARRAY[ARRAY[$1, $2, to_char(CURRENT_TIMESTAMP, 'YYYY-MM-DD HH24:MI:SS.US')]] WHERE id = $3", ip, ua, id)
		if err != nil {
			log.Errorf("Couldn't update tracking for id %s: %v", id, err)
			return
		}
	}(c.Params("id"), c.IP(), c.Get(fiber.HeaderUserAgent))

	// Disable caching
	c.Set("Cache-Control", "max-age=0, no-cache, must-revalidate, proxy-revalidate")
	c.Set("Expires", "0")

	// Content -> gif
	c.Set("Content-Type", "image/gif")

	// Send the pixel
	return c.Send(img)
}
