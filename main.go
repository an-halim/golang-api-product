package main

import (
	"time"

	"github.com/an-halim/golang-api-product/database"
	"github.com/an-halim/golang-api-product/router"
	"github.com/an-halim/golang-api-product/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
)

func main()  {
	database.Connect()

	app := fiber.New()

	app.Use(logger.New())

	app.Use(cors.New())

	router.SetupRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		tin := time.Now()
		return utils.Failed(c, 404, "Not Found", tin)
	})

	app.Listen(":3000")
}