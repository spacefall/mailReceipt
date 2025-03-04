package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func removeTracking(c *fiber.Ctx) error {
	id := c.Params("uuid")

	// Check if uuid is valid
	if err := uuid.Validate(id); err != nil {
		return c.Status(400).SendString("Invalid tracking id")
	}

	// Delete the row corresponding to the uuid
	res, err := dbpool.Exec(context.Background(), "DELETE FROM mail_receipts WHERE id = $1", id)
	if err != nil {
		return err
	}

	// If no rows were affected, the uuid doesn't exist
	if res.RowsAffected() == 0 {
		return c.Status(404).SendString("Tracking id not found")
	}
	return c.SendStatus(200)
}
