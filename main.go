package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
	"time"
)

// this is the image: data:image/gif;base64,R0lGODlhAQABAIAAAP///wAAACH5BAEAAAAALAAAAAABAAEAAAICRAEAOw==
var img = []byte{71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 255, 255, 255, 0, 0, 0, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 2, 68, 1, 0, 59}

func main() {
	// Generate a 1x1 pixel
	//pxPng, err := generatePixel()
	//if err != nil {
	//	panic(err)
	//}

	app := fiber.New(fiber.Config{AppName: "mailReceipt"})

	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		c.Append("Server-Timing", "app;dur="+strconv.FormatFloat(float64(duration)/1000000, 'f', -1, 64))
		return err
	})

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
		// Content -> gif
		c.Set(fiber.HeaderContentType, "image/gif")
		// Send the pixel
		return c.Send(img)
	})

	fmt.Println("Started at: http://localhost:3000")
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
