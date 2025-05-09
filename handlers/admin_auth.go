package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/C9b3rD3vi1/Burnr/models"
	"github.com/C9b3rD3vi1/Burnr/database"
	//"github.com/google/uuid"
	//"time"
	"log"
)


// This is an admin handler for creating and managing links
// It is not intended for public use It is used for testing and debugging purposes

func AdminManageHandler(c *fiber.Ctx) error {
	// This is a placeholder for the admin management page
	var links []models.Link
	result := database.DB.Find(&links)
	if result.Error != nil {
		log.Println("Error fetching links:", result.Error)
		return nil
	}

	//log.Println("Links fetched: ", links) // Log fetched links

	// Render the admin page with the links
	return c.Render("admin", fiber.Map{
		"links": links,
	})
}
