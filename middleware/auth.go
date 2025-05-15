package middleware

import (
	"log"
	"github.com/gofiber/fiber/v2"
)


// Ensure the user is logged in before accessing certain routes
func AuthMiddleware(c *fiber.Ctx) error {
	// Get user ID from session
	sess, err := Store.Get(c)
	
	if err != nil {
		log.Println("Error getting session:", err)
		return c.Redirect("/login")
	}

	userID := sess.Get("userID")
	if userID == nil {
		log.Println("User not logged in")

		return c.Redirect("/login")
	}

	return c.Next()
}
