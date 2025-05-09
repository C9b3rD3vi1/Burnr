package main

import (
	"github.com/gofiber/fiber/v2"
	//"github.com/C9b3rD3vi1/Burnr/models"
	"github.com/C9b3rD3vi1/Burnr/handlers"
	"github.com/C9b3rD3vi1/Burnr/database"
)



func main() {
	app := fiber.New()

	database.ConnectDB()

	app.Static("/", "./views") // serve index.html
	app.Post("/create", handlers.CreateLink)
	app.Get("/:id", handlers.RedirectLink)

	app.Listen(":3000")
}
