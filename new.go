package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type NewResp struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
}

func newTracker(c *fiber.Ctx) error {
	var id string
	var timestamp string
	err := dbpool.QueryRow(context.Background(), "INSERT INTO mail_receipts DEFAULT VALUES RETURNING id, created_at").Scan(&id, &timestamp)
	if err != nil {
		log.Errorf("Couldn't create tracking ID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrResp{
			Code: fiber.StatusInternalServerError,
			Msg:  "Couldn't create tracking ID",
		})
	}

	c.Set(fiber.HeaderLocation, "/track/"+id)
	return c.Status(201).JSON(NewResp{
		Id:        id,
		CreatedAt: timestamp,
	})
}
