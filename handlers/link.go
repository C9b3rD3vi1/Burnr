package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"Burnr/database"
	"Burnr/models"
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
	}

	database.DB.Create(&link)

	return c.JSON(fiber.Map{"short_url": c.BaseURL() + "/" + link.ID})
}

func RedirectLink(c *fiber.Ctx) error {
	id := c.Params("id")
	var link models.Link

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
