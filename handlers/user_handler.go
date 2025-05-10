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
