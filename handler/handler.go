package handler

import (
	"context"
	"fmt"
	"log"
	"time"
	"todo/database"
	"todo/model"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gofiber/fiber"
)

func ReadTodo(c *fiber.Ctx) {
	_id := c.Params("id")
	if docID, err := primitive.ObjectIDFromHex(_id); err != nil {
		c.Status(400).JSON(fiber.Map{"success": false})
	} else {
		q := bson.M{"_id": docID}
		todo := model.Todo{}
		result := database.TestCollection.FindOne(context.Background(), q)
		result.Decode(&todo)

		if result.Err() != nil {
			fmt.Println(result.Err().Error())
			c.Status(200).JSON(fiber.Map{"success": true, "todo": fmt.Sprintf("No todo found for give id: %s", _id)})
		} else {
			c.Status(200).JSON(fiber.Map{"success": true, "todo": todo})
		}
	}

}

func CreateTodo(c *fiber.Ctx) {
	todo := new(model.Todo)
	if err := c.BodyParser(todo); err != nil {
		c.Status(400).JSON(fiber.Map{"message": fmt.Sprintf("Invalid post body. %s", err.Error())})
	} else {
		todo.AddTimeStamps()
		if r, err := database.TestCollection.InsertOne(context.Background(), todo); err != nil {
			c.Status(500).JSON(fiber.Map{"message": "Something went wrong, please try after sometime"})
		} else {
			c.Status(200).JSON(fiber.Map{"success": true, "id": r.InsertedID})
		}
	}
}

func UpdateTodo(c *fiber.Ctx) {
	_id := c.Params("id")
	todo := new(model.Todo)
	if err := c.BodyParser(todo); err != nil {
		c.Status(400).JSON(fiber.Map{"message": fmt.Sprintf("Invalid post body. %s", err.Error())})
	} else {
		if docID, err := primitive.ObjectIDFromHex(_id); err != nil {
			c.Status(400).JSON(fiber.Map{"success": false})
		} else {
			q := bson.M{"_id": docID}

			u := bson.D{
				{"$set", bson.D{
					{"name", todo.Name},
					{"description", todo.Description},
					{"updated_at", time.Now()},
				},
				}}

			o := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

			database.TestCollection.FindOneAndUpdate(context.Background(), q, u, o).Decode(&todo)
			c.Status(200).JSON(fiber.Map{"success": true, "todo": todo})
		}
	}
}

func DeleteTodo(c *fiber.Ctx) {
	_id := c.Params("id")

	if docID, err := primitive.ObjectIDFromHex(_id); err != nil {
		c.Status(400).JSON(fiber.Map{"success": false})
	} else {
		q := bson.M{"_id": docID}
		r := database.TestCollection.FindOneAndDelete(context.Background(), q)
		if r.Err() != nil {
			c.Status(400).JSON(fiber.Map{"success": false})
		} else {
			c.Status(204)
		}
	}
}

func ReadAllTodo(c *fiber.Ctx) {
	var todos []*model.Todo
	cursor, err := database.TestCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.Status(500).JSON(fiber.Map{"message": "Something went wrong, please try after sometime"})
	}

	for cursor.Next(context.Background()) {
		var todo model.Todo
		err := cursor.Decode(&todo)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, &todo)
	}

	if cursor.Err() != nil {
		log.Fatal(err)
	}

	cursor.Close(context.Background())

	c.Status(200).JSON(fiber.Map{"success": true, "todos": todos})
}

func DummyAllHandler(c *fiber.Ctx) {
	c.Status(200).JSON(fiber.Map{
		"message": "DummyAllHandler",
	})
}
