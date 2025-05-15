package handlers

import (
	"errors"
	"log"
	"strings"
	"time"
	"sort"
	"github.com/C9b3rD3vi1/Burnr/database"
	"github.com/C9b3rD3vi1/Burnr/middleware"
	"github.com/C9b3rD3vi1/Burnr/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserHandler handles user authentication and authorization and registering new users
func UserRegister(c *fiber.Ctx) error {
	// Normalize and validate input
	username := strings.TrimSpace(c.FormValue("username"))
	email := strings.TrimSpace(strings.ToLower(c.FormValue("email")))
	password := c.FormValue("password")
	passwordConfirm := c.FormValue("password_confirm")

	// Basic validation
	if username == "" || email == "" || password == "" || passwordConfirm == "" {
		//log.Println("Received:", username, email, password, passwordConfirm)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All fields are required",
		})
	}

	if password != passwordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Passwords do not match",
		})
	}

	// Check if username exists
	var existingUser models.User
	if err := database.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username already exists",
		})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("DB error checking username:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server error",
		})
	}

	// Check if email exists
	if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email already exists",
		})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("DB error checking email:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server error",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Password hashing error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process password",
		})
	}

	// Create user
	user := models.User{
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Println("DB error creating user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not register user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

// user login function
func UserLogin(c *fiber.Ctx) error {
	username := strings.TrimSpace(c.FormValue("username"))
	password := c.FormValue("password")

	// Basic validation
	if username == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username and password are required",
		})
	}
	// Fetch user
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		}
		log.Println("DB error fetching user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server error",
		})
	}
	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
		
	}
	// Set session or token
	sess, err := middleware.Store.Get(c)
	if err != nil {
		log.Println("Error setting session cookie:", err)

	}

	sess.Set("userID", user.ID)
	sess.Set("username", user.Username)
	sess.Save()

	// log the user in and redirect to dashboard
	return c.Redirect("/dashboard")
	
}

func UserLogout(c *fiber.Ctx) error {
    sess, _ := middleware.Store.Get(c)
    sess.Destroy()
    return c.Redirect("/login")
}



func UserDashboard(c *fiber.Ctx) error {
    // Session handling
    sess, err := middleware.Store.Get(c)
    if err != nil {
        return c.Status(500).SendString("Session error")
    }
    userID := sess.Get("userID")
    if userID == nil {
        return c.Redirect("/login")
    }

    // Fetch user info
    var user models.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        return c.Status(500).SendString("User not found")
    }

    // Fetch links with click counts
    var links []models.Link
    if err := database.DB.
        Where("user_id = ?", userID).Find(&links).Error; err != nil {
		return c.Next()
        //return c.Status(500).SendString("Could not fetch links")
    }

    // Calculate statistics
    type Stats struct {
        TotalClicks      int
        ClicksToday     int
        PopularLinks    []PopularLink
        Devices         map[string]int
        Referrers       map[string]int
        Countries       map[string]int
        ClicksOverTime  map[string]int
    }

    stats := Stats{
        Devices:    make(map[string]int),
        Referrers:  make(map[string]int),
        Countries:  make(map[string]int),
        ClicksOverTime: make(map[string]int),
    }

    today := time.Now().Truncate(24 * time.Hour)
    var allClicks []models.LinkClick

    for _, link := range links {
        stats.TotalClicks += link.Clicks
        
        // Get all clicks for this link
        var clicks []models.LinkClick
        if err := database.DB.Where("link_id = ?", link.ID).Find(&clicks).Error; err == nil {
            allClicks = append(allClicks, clicks...)
            
            for _, click := range clicks {
                // Clicks today
                if click.Time.After(today) {
                    stats.ClicksToday++
                }
                
                // Device breakdown
                stats.Devices[click.Device]++
                
                // Referrer breakdown
                if click.Referrer != "" {
                    stats.Referrers[click.Referrer]++
                }
                
                // Country breakdown
                stats.Countries[click.Country]++
                
                // Time series data (by hour)
                hour := click.Time.Format("2006-01-02 15:00")
                stats.ClicksOverTime[hour]++
            }
        }
    }

    // Prepare popular links (top 5)
    sort.Slice(links, func(i, j int) bool {
        return links[i].Clicks > links[j].Clicks
    })
    popularLinks := make([]PopularLink, 0)
    for i, link := range links {
        if i >= 5 {
            break
        }
        popularLinks = append(popularLinks, PopularLink{
            ShortURL:   link.Shortened,
            TargetURL:  link.TargetURL,
            Clicks:     link.Clicks,
        })
    }
    stats.PopularLinks = popularLinks

    return c.Render("dashboard", fiber.Map{
        "User": fiber.Map{
            "Username": user.Username,
            "Email":    user.Email,
        },
        "Links":       links,
        "Stats":       stats,
        "TotalLinks": len(links),
    })
}

type PopularLink struct {
    ShortURL  string
    TargetURL string
    Clicks    int
}