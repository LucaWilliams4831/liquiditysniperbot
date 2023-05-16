package controllers

import (
	"strconv"
	"time"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"github.com/LucaWilliams4831/liquiditysniperbot/database"
	"github.com/LucaWilliams4831/liquiditysniperbot/models"
)

type Data struct {
	Address string `json:"address"`
}

func Login(c *fiber.Ctx) error {
	data := new(Data)

	if err := c.BodyParser(data); err != nil {
		return err
	}

	var user models.Admin

	database.DB.Where("LOWER(address) = ?", strings.ToLower(data.Address)).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}
	
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})

	token, err := claims.SignedString([]byte(Secretkey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		SameSite: "strict",
		// Secure: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"token": cookie,
		"message": "success",
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}