package main

import (
	"github.com/gofiber/fiber/v2"
	//"github.com/C9b3rD3vi1/Burnr/models"
	"github.com/C9b3rD3vi1/Burnr/handlers"
	"github.com/C9b3rD3vi1/Burnr/database"
	"github.com/gofiber/template/html/v2" // html template engine
	"github.com/gofiber/fiber/v2/middleware/logger"

)



func main() {
	// Initialize the Fiber app
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())
	// Initialize the database
	database.ConnectDB()

	app.Static("/", "./views") // serve index.html
	app.Post("/create", handlers.CreateLink)

	app.Get("/admin", handlers.AdminManageHandler) // admin page

	app.Get("/register", func(c *fiber.Ctx) error {
		return c.Render("register", fiber.Map{})
	})
	

	app.Post("register", handlers.UserRegister) // user registration page


	app.Get("/:id", handlers.RedirectLink)

	

	app.Listen(":3000")
}
