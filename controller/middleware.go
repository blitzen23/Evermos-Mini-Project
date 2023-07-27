package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func IsAdmin(c *fiber.Ctx) error {
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to Access",
			"errors":  err.Error(),
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

func IsLoggedIn(c *fiber.Ctx) error {
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to Access",
			"errors":  err.Error(),
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
	sess, err := Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to Access",
			"errors":  err.Error(),
			"data":    nil,
		})
	}
	sess.Set("id", uint(claims["id"].(float64)))
	sess.SetExpiry(time.Hour)
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to Access",
			"errors":  err.Error(),
			"data":    nil,
		})
	}
	return c.Next()
}