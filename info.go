package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func getTrackerInfo(c *fiber.Ctx) error {
	info := TrackInfo{c.Params("id"), "", "", "", "", nil}
	err := dbpool.QueryRow(context.Background(), "SELECT name, email, created_at, created_by, pixel_events FROM mail_receipts WHERE id = $1", info.Id).Scan(&info.Name, &info.Email, &info.CreatedAt, &info.CreatedBy, &info.Events)
	if err != nil {
		log.Errorf("Couldn't get info (ID: %s): %v", info.Id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrResp{
			Code: fiber.StatusInternalServerError,
			Msg:  "Couldn't get tracking info",
		})
	}
	return c.JSON(info)
}
