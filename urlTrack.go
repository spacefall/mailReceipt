package main

import (
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"time"
)

func urlTrack(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get IP (supporting cf proxy too)
	ip := c.Get("cf-connecting-ip")
	if ip == "" {
		ip = c.IP()
	}

	trackingJson := TrackData{
		Ip:        ip,
		UserAgent: c.Get(fiber.HeaderUserAgent),
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	plainUrl, err := base64.RawURLEncoding.DecodeString(c.Params("url"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrResp{
			Code: fiber.StatusBadRequest,
			Msg:  "Invalid url encoding",
		})
	}

	trackingJson.Url = string(plainUrl)

	// Add request to db asynchronously
	go dbAppend(id, "url_events", trackingJson)

	go sendMail(trackingJson, id)

	return c.Redirect(trackingJson.Url)
}
