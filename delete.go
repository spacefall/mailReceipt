package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func deleteTracker(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if uuid is valid
	if err := uuid.Validate(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrResp{
			Code: fiber.StatusBadRequest,
			Msg:  "ID is not valid",
		})
	}

	// Delete the row corresponding to the uuid
	res, err := dbpool.Exec(context.Background(), "DELETE FROM mail_receipts WHERE id = $1", id)
	if err != nil {
		log.Errorf("Couldn't delete tracking ID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrResp{
			Code: fiber.StatusInternalServerError,
			Msg:  "Couldn't delete tracking ID",
		})
	}

	// If no rows were affected, the uuid doesn't exist
	if res.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(ErrResp{
			Code: fiber.StatusNotFound,
			Msg:  "Tracking ID not found",
		})
	}
	return c.SendStatus(204)
}
