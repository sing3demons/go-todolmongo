package routes

import (
	"todo-mongo/config"
	"todo-mongo/controllers"

	"github.com/gofiber/fiber/v2"
)

func Serve(app *fiber.App) {
	db := config.ConnectDB()
	itemController := controllers.NewItemController(db)
	v1 := app.Group("api/v1/items")
	{
		v1.Get("", itemController.FindItems)
		v1.Post("", itemController.Create)
		v1.Delete("/:id", itemController.Delete)
		v1.Put("/:id", itemController.Update)
		v1.Get("/:id", itemController.FindOne)
	}
}
