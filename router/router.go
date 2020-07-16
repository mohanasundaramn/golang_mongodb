package router

import (
	"todo/handler"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
)

func SetupRoutes(app *fiber.App) {
	todo := app.Group("/todo", middleware.Logger())

	todo.Get("/:id", handler.ReadTodo)
	todo.Post("/", handler.CreateTodo)
	todo.Put("/:id", handler.UpdateTodo)
	todo.Delete("/:id", handler.DeleteTodo)
	todo.Get("/all/todos", handler.ReadAllTodo)
	todo.Get("/all", handler.DummyAllHandler)
}
