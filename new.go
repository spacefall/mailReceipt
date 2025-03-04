package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

func newTracker(c *fiber.Ctx) error {
	var uuid string
	err := dbpool.QueryRow(context.Background(), "INSERT INTO mail_receipts DEFAULT VALUES RETURNING id").Scan(&uuid)
	if err != nil {
		return err
	}
	return c.Status(201).SendString("The new UUID is: " + uuid)
}
