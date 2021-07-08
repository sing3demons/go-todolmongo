package controllers

import (
	"context"
	"time"
	"todo-mongo/helper"
	"todo-mongo/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)



type itemController struct {
	db *mongo.Database
}

func NewItemController(db *mongo.Database) itemController {
	return itemController{db: db}
}

type CreateItem struct {
	Title       string `json:"title" bson:"title" validate:"required"`
	Description string `json:"description" bson:"description" validate:"required"`
}

type UpdateItem struct {
	Title       string `json:"title" bson:"title" validate:"required"`
	Description string `json:"description" bson:"description" validate:"required"`
}

func (tx itemController) collection() *mongo.Collection {
	return tx.db.Collection("items")
}

func (tx itemController) Delete(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter, err := tx.findItemById(c)
	if err != nil {
		c.Status(fiber.StatusNotFound).JSON(err)
	}

	if err := tx.collection().FindOneAndDelete(ctx, filter).Err(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (tx itemController) Update(c *fiber.Ctx) error {
	var form UpdateItem

	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	if err := helper.ValidateStruct(&form); err != nil {
		return c.JSON(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter, err := tx.findItemById(c)
	if err != nil {
		c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	update := bson.D{
		{"$set", form},
	}

	if err := tx.collection().FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		return c.JSON(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (tx itemController) FindItems(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	items := []models.Item{}

	cursor, err := tx.collection().Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}
	defer cursor.Close(ctx)

	// for cursor.Next(ctx) {
	// 	var item models.Item
	// 	if err := cursor.Decode(&item); err != nil {
	// 		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	// 	}
	// 	items = append(items, item)
	// }

	if err := cursor.All(ctx, &items); err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"items": items})

}

func (tx itemController) FindOne(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var item models.Item

	filter, err := tx.findItemById(c)
	if err != nil {
		c.Status(fiber.StatusNotFound).JSON(err)
	}

	if err := tx.collection().FindOne(ctx, filter).Decode(&item); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"item": item})
}

func (tx itemController) Create(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var form CreateItem

	// for i := 0; i < 5000; i++ {
	// 	tx.collection().InsertOne(ctx, models.Item{
	// 		Title:       faker.Word(),
	// 		Description: faker.Paragraph(),
	// 	})
	// }

	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	if err := helper.ValidateStruct(&form); err != nil {
		return c.JSON(err)
	}

	_, err := tx.collection().InsertOne(ctx, form)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func (tx itemController) findItemById(c *fiber.Ctx) (primitive.M, error) {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	return filter, nil
}
