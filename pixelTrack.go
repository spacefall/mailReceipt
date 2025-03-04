package main

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

// this is the image: data:image/gif;base64,R0lGODlhAQABAIAAAP///wAAACH5BAEAAAAALAAAAAABAAEAAAICRAEAOw==
var img = []byte{71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 255, 255, 255, 0, 0, 0, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 2, 68, 1, 0, 59}

func pixelTrack(c *fiber.Ctx) error {
	id := c.Params("id")
	trackingJson := TrackData{
		Ip:        c.IP(),
		UserAgent: c.Get(fiber.HeaderUserAgent),
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	// Add request to db asynchronously
	go dbAppend(id, "pixel_events", trackingJson)

	go sendMail(trackingJson, id)

	// Disable caching
	c.Set("Cache-Control", "max-age=0, no-cache, must-revalidate, proxy-revalidate")
	c.Set("Expires", "0")

	// Content -> gif
	c.Set("Content-Type", "image/gif")

	// Send the pixel
	return c.Send(img)
}
