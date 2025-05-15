package handlers

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/C9b3rD3vi1/Burnr/middleware"
	"github.com/C9b3rD3vi1/Burnr/database"
	"github.com/C9b3rD3vi1/Burnr/models"
)

func UserPriceHandler(c *fiber.Ctx) error {
	sess, err := middleware.Store.Get(c)
	userID := sess.Get("userID")
	
	if err != nil {
		log.Println("Error getting session:", err)
		return c.Redirect("/login")
	}

	if userID == nil {
		return c.Redirect("/login")
	}

	// Fetch user info
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(500).SendString("User not found")
	}

	// Fetch pricing
	var pricing []models.Pricing
	if err := database.DB.Find(&pricing).Error; err != nil {
		return c.Status(500).SendString("Could not fetch pricing")
	}



	return c.Render("price", fiber.Map{})
}