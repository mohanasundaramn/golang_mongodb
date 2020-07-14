package main

import (
	"log"
	"os"
	"todo/database"
	"todo/router"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Unable to load .env", err.Error())
	} else {
		p := os.Getenv("PORT")

		database.ConnectDB()

		app := fiber.New()

		app.Use(middleware.Logger())

		router.SetupRoutes(app)

		log.Fatal(app.Listen(p))
	}
}
