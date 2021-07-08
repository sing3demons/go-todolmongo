package controllers

import (
	"context"
	"time"
	"todo-mongo/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Item struct {
	DB *mongo.Database
}

type CreateItem struct {
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}

type UpdateItem struct {
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}

func (tx Item) collection() *mongo.Collection {
	return tx.DB.Collection("items")
}

func (tx Item) Delete(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter, err := tx.findItemById(c)
	if err != nil {
		c.Status(fiber.StatusNotFound).JSON(err)
	}

	result, _ := tx.collection().DeleteOne(ctx, filter)
	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON("not found")
	}

	c.JSON(result)
	return c.SendStatus(fiber.StatusCreated)
}

func (tx Item) Update(c *fiber.Ctx) error {
	var form UpdateItem
	if err:=c.BodyParser(&form);err!=nil{
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
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

	result, err := tx.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	if result.UpsertedCount == 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON("cannot update document")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (tx Item) FindItems(c *fiber.Ctx) error {
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

	if err = cursor.All(ctx, &items); err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"items": items})

}

func (tx Item) FindOne(c *fiber.Ctx) error {
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
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"item": item})
}

func (tx Item) Create(c *fiber.Ctx) error {
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

	_, err := tx.collection().InsertOne(ctx, form)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func (tx Item) findItemById(c *fiber.Ctx) (primitive.M, error) {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	return filter, nil
}
