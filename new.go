package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type NewReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func newTracker(c *fiber.Ctx) error {
	var reqBody NewReq

	// Get name and optionally email from request body (JSON)
	err := c.BodyParser(&reqBody)
	if err != nil {
		log.Errorf("Couldn't parse request body: %v", err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(ErrResp{
			Code: fiber.StatusUnprocessableEntity,
			Msg:  "Couldn't parse request body",
		})
	}

	// Make sure name is not empty
	if reqBody.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrResp{
			Code: fiber.StatusBadRequest,
			Msg:  "Name is required",
		})
	}

	// Add info to DB and complete creation of info object to return
	info := TrackInfo{"", reqBody.Name, reqBody.Email, "", c.IP(), nil}
	err = dbpool.QueryRow(context.Background(), "INSERT INTO mail_receipts (name, email, created_by) VALUES ($1, $2, $3) RETURNING id, created_at", info.Name, info.Email, info.CreatedBy).Scan(&info.Id, &info.CreatedAt)
	if err != nil {
		log.Errorf("Couldn't create tracking ID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrResp{
			Code: fiber.StatusInternalServerError,
			Msg:  "Couldn't create tracking ID",
		})
	}

	c.Set(fiber.HeaderLocation, "/track/"+info.Id)
	return c.Status(fiber.StatusCreated).JSON(info)
}
