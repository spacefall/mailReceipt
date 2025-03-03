package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	pxPng, err := generatePixel()
	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{AppName: "mailReceipt"})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/pixel", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "Tracking ID Missing")
	})

	app.Get("/pixel/:id", func(c *fiber.Ctx) error {
		// Log info
		log.Println(c.Params("id"))
		log.Println(c.IP(), c.IPs())
		log.Println(c.GetReqHeaders())
		// Disable caching
		c.Set(fiber.HeaderCacheControl, "max-age=0, no-cache, must-revalidate, proxy-revalidate")
		c.Set(fiber.HeaderExpires, "0")
		// Send the pixel
		return c.Send(pxPng)
	})

	fmt.Println("Started at: http://localhost:3000")
	app.Listen(":3000")
}
