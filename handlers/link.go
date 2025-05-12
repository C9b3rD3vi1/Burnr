package handlers

import (
	"time"
	"encoding/json"
	"net/http"
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mssola/user_agent"
	"github.com/C9b3rD3vi1/Burnr/models"
	"github.com/C9b3rD3vi1/Burnr/database"
	"github.com/C9b3rD3vi1/Burnr/middleware"

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


func SaveLinkClick(c *fiber.Ctx)error{
	id := c.Params("id")
	ip := c.IP()
	referer := c.Get("Referer")
	userAgent := c.Get("User-Agent")

	device := ParseUserAgent(userAgent)
	country := LookupCountryByIP(ip)

	// save clicks
	clicks := models.LinkClick{
		LinkID: id,
		Time: time.Now(),
		IP: ip,
		Country: country,
		Referrer: referer,
		Device: device,

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