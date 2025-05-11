package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/C9b3rD3vi1/Burnr/models"
	"github.com/C9b3rD3vi1/Burnr/database"
)
// CreateLink handles the creation of a new shortened link
// It expects a JSON body with the URL, max clicks, and expiration time in minutes
// It returns the shortened URL

func CreateLink(c *fiber.Ctx) error {
	type Request struct {
		URL       string `json:"url"`
		MaxClicks int    `json:"max_clicks"`
		ExpireMin int    `json:"expire_min"`
	}

	var body Request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	link := models.Link{
		ID:        uuid.NewString(),
		TargetURL: body.URL,
		MaxClicks: body.MaxClicks,
		ExpiresAt: time.Now().Add(time.Duration(body.ExpireMin) * time.Minute),
		CreatedAt: time.Now(),
		Clicks:    0,
		Shortened: uuid.NewString(), // or generate a shorter custom string
	}
	

	if err := database.DB.Create(&link).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save link"})
	}

	return c.JSON(fiber.Map{"short_url": c.BaseURL() + "/" + link.ID})
}

func RedirectLink(c *fiber.Ctx) error {
	id := c.Params("id")
	var link models.Link

	// Prevent reserved paths from being treated as link IDs
	if id == "admin"|| id =="register" || id == "login" || id == "dashboard" {
		return c.Redirect(id)
	}

	result := database.DB.First(&link, "id = ?", id)
	if result.Error != nil {
		return c.Status(404).SendString("Link not found.")
	}

	if link.Clicks >= link.MaxClicks || time.Now().After(link.ExpiresAt) {
		return c.SendString("Link expired.")
	}

	link.Clicks++
	database.DB.Save(&link)

	return c.Redirect(link.TargetURL)
}
