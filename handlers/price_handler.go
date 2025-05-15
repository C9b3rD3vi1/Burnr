package handlers

import (
	//"log"
	"github.com/gofiber/fiber/v2"
	"github.com/C9b3rD3vi1/Burnr/middleware"
	"github.com/C9b3rD3vi1/Burnr/database"
	"github.com/C9b3rD3vi1/Burnr/models"
)


// ShowPricingPage shows the pricing page
func ShowPricingPage(c *fiber.Ctx) error {
	sess, err := middleware.Store.Get(c)
	userID := sess.Get("userID")

	if err != nil || userID == nil {
		return c.Redirect("/login")
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(500).SendString("User not found")
	}

	return c.Render("price", fiber.Map{
		"Title": "Pricing",
		"User":  user,
	}, "base")
}



// HandlePricingCheckout handles the checkout process
func HandlePricingCheckout(c *fiber.Ctx) error {
	var body struct {
		Email string `json:"email"`
		Plan  string `json:"plan"` // "monthly" or "yearly"
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	url, err := middleware.CreateCheckoutSession(body.Email, body.Plan)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create session"})
	}

	return c.JSON(fiber.Map{"checkout_url": url})
}
