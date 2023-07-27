package helper

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var store *session.Store

func init(){
	store = session.New()
}

func isAdmin(c *fiber.Ctx) error {
	tokenString := c.GetReqHeaders()["Token"]
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to Access",
			"errors":  "You are unauthorized",
			"data":    nil,
		})
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secretkey"), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to Access",
			"errors":  "You are unauthorized",
			"data":    nil,
		})
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to Access",
			"errors":  "Invalid token claims",
			"data":    nil,
		})
	}
	isAdmin, exists := claims["isAdmin"].(bool)
	if !exists || !isAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to Access",
			"errors":  "You are unauthorized",
			"data":    nil,
		})
	}
	return c.Next()
}

func isLoggedIn(c *fiber.Ctx) error {
	tokenString := c.GetReqHeaders()["Token"]
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to Access",
			"errors":  "You are unauthorized",
			"data":    nil,
		})
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secretkey"), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to Access",
			"errors":  "You are unauthorized",
			"data":    nil,
		})
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to Access",
			"errors":  "Invalid token claims",
			"data":    nil,
		})
	}
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to Access",
			"errors":  "There is an internal server error",
			"data":    nil,
		})
	}
	sess.Set("id", uint(claims["id"].(float64)))
	sess.SetExpiry(time.Hour)
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to Access",
			"errors":  "There is an internal server error",
			"data":    nil,
		})
	}
	return c.Next()
}