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

	type request struct {
		Email string `json:"email"`
		Plan  string `json:"plan"` // "monthly" or "yearly"
	}

	var body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	url, err := middleware.CreateCheckoutSession(body.Email, body.Plan)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create session"})
	}

	return c.JSON(fiber.Map{"checkout_url": url})
}