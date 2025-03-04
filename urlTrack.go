package main

import (
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"time"
)

func urlTrack(c *fiber.Ctx) error {
	trackingJson := TrackData{
		Ip:        c.IP(),
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
	go dbAppend(c.Params("id"), "url_events", trackingJson)

	return c.Redirect(trackingJson.Url)
}
