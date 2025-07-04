package handlers

import (
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"time"

	"github.com/C9b3rD3vi1/Burnr/database"
	"github.com/C9b3rD3vi1/Burnr/middleware"
	"github.com/C9b3rD3vi1/Burnr/models"
	"github.com/gofiber/fiber/v2"
	//"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"github.com/mssola/user_agent"
)

// CreateLink handles the creation of a new shortened link
// It expects a JSON body with the URL, max clicks, and expiration time in minutes
// It returns the shortened URL

func CreateLink(c *fiber.Ctx) error {
	// Check if the user is logged in
	sess, _ := middleware.Store.Get(c)
    userID := sess.Get("userID") // This will be nil if not logged in


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
	
	 // If user is logged in, associate the link with their user ID
	 if userID != nil {
        uid, ok := userID.(uint)
        if ok {
            link.UserID = uid
        }
    }
	

	if err := database.DB.Create(&link).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save link"})
	}

	return c.JSON(fiber.Map{"short_url": c.BaseURL() + "/" + link.ID})
}


// RedirectLink handles the redirection of a shortened link
// It checks if the link exists, if it has expired, and if the max clicks have been reached

func RedirectLink(c *fiber.Ctx) error {
	id := c.Params("id")
	var link models.Link

	// Prevent reserved paths from being treated as link IDs
	if id == "admin" || id == "register" || id == "login" || id =="price" || id == "dashboard" {
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

	// Extract all necessary data BEFORE goroutine
	clickData := struct {
		LinkID    string
		IP        string
		Referer   string
		UserAgent string
	}{
		LinkID:    id,
		IP:        c.IP(),
		Referer:   c.Get("Referer"),
		UserAgent: c.Get("User-Agent"),
	}

	// Pass only the extracted data to the goroutine
	go func(data struct {
		LinkID    string
		IP        string
		Referer   string
		UserAgent string
	}) {
		_ = SaveLinkClick(data.LinkID, data.IP, data.Referer, data.UserAgent)
	}(clickData)

	return c.Redirect(link.TargetURL)
}


func SaveLinkClick(linkID, ip, referer, userAgent string) error {
	if linkID == "" {
		return fmt.Errorf("Link ID is missing")
	}

	device := ParseUserAgent(userAgent)
	country := LookupCountryByIP(ip)

	// save clicks
	clicks := models.LinkClick{
		LinkID:   linkID,
		Time:     time.Now(),
		IP:       ip,
		Country:  country,
		Referrer: referer,
		Device:   device,
	}

	if err := database.DB.Create(&clicks).Error; err != nil {
		return err
	}
	return nil
}


func ParseUserAgent(ua string) string {
	uaParser := user_agent.New(ua)
	name, version := uaParser.Browser()
	return name + " " + version
}


type IPAPIResponse struct {
	Country string `json:"country"`
}

func LookupCountryByIP(ip string) string {
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		log.Println("Geo lookup failed:", err)
		return "Unknown"
	}
	defer resp.Body.Close()

	var result IPAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Error decoding geo response:", err)
		return "Unknown"
	}

	return result.Country
}


func DeleteExpiredLinksHandler(c *fiber.Ctx) error {
	now := time.Now()

	// Fetch expired links first
	var expiredLinks []models.Link
	if err := database.DB.
		Where("expires_at IS NOT NULL AND expires_at < ?", now).
		Find(&expiredLinks).Error; err != nil {
		log.Println("Error fetching expired links:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch expired links",
		})
	}

	if len(expiredLinks) == 0 {
		return c.JSON(fiber.Map{
			"message":       "No expired links found",
			"rows_affected": 0,
			"ids":           []string{},
		})
	}

	// Collect IDs for front-end removal
	var ids []string
	for _, link := range expiredLinks {
		if link.ID != "" {
			ids = append(ids, link.ID)
		}
	}

	// Check if any valid IDs were collected
	if len(ids) == 0 {
		return c.JSON(fiber.Map{
			"message":       "No valid expired links found for deletion",
			"rows_affected": 0,
			"ids":           []string{},
		})
	}

	// Delete expired links
	if err := database.DB.Delete(&expiredLinks).Error; err != nil {
		log.Println("Error deleting expired links:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete expired links",
		})
	}

	return c.JSON(fiber.Map{
		"message":       "Expired links deleted successfully",
		"rows_affected": len(expiredLinks),
		"ids":           ids,
	})
}
