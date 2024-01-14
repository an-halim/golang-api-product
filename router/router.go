package router

import (
	"github.com/an-halim/golang-api-product/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/api/product", handler.CreateProduct)
	app.Get("/api/product", handler.GetProducts)
	app.Get("/api/product/:id", handler.GetProduct)
	app.Put("/api/product/:id", handler.UpdateProduct)
	app.Delete("/api/product/:id", handler.DeleteProduct)
}