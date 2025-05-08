package main

import (
	"github.com/gofiber/fiber/v2"
	"Burnr/database"
	"Burnr/handlers"
)



func main() {
	app := fiber.New()

	database.ConnectDB()

	app.Static("/", "./views") // serve index.html
	app.Post("/create", handlers.CreateLink)
	app.Get("/:id", handlers.RedirectLink)

	app.Listen(":3000")
}
