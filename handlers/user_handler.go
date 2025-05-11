package handlers


import (
	"github.com/gofiber/fiber/v2"
	"github.com/C9b3rD3vi1/Burnr/models"
	"github.com/C9b3rD3vi1/Burnr/database"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
	"errors"
	"gorm.io/gorm"
	"strings"
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
	// log the user in and redirect to dashboard
	return c.Redirect("/dashboard")
	
}


// UserDashboard handles the user dashboard page
func UserDashboard(c *fiber.Ctx) error {
	// Fetch user ID from session or token
	userID := c.Locals("userID") // Assuming you have middleware to set this

	// Fetch user links
	var links []models.Link
	if err := database.DB.Where("user_id = ?", userID).Find(&links).Error; err != nil {
		log.Println("DB error fetching user links:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch user links",
		})
	}

	return c.Render("dashboard", fiber.Map{
		"links": links,
	})
}