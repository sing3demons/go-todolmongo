package main

import (
	"todo-mongo/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	app := fiber.New()
	app.Use(cors.New())

	app.Get("", func(c *fiber.Ctx) error {
		return c.JSON("hello")
	})
	
	routes.Serve(app)

	app.Listen(":8080")
}
