package middleware

import (
	"time"
	"github.com/gofiber/fiber/v2/middleware/session"
)


var Store = session.New(session.Config{
		Expiration: 24 * time.Hour,
})