package middlewears

import (
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

const SecurityKey = "secret"

type ClaimWithScope struct {
	jwt.StandardClaims
	Scope string
}

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &ClaimWithScope{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecurityKey), nil
	})

	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "you are not authorized ",
		})
	}
	payload := token.Claims.(*ClaimWithScope)
	isAmbassador := strings.Contains(c.Path(), "/api/ambassador")

	if (payload.Scope == "admin" && isAmbassador) || (payload.Scope == "ambassador" && !isAmbassador) {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "you are not authorized ",
		})
	}

	return c.Next()
}

func GetUserID(c *fiber.Ctx) (uint, error) {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &ClaimWithScope{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecurityKey), nil
	})
	if err != nil {
		return 0, err
	}
	payload := token.Claims.(*ClaimWithScope)

	id, _ := strconv.Atoi(payload.Subject)

	return uint(id), nil

}

func GenerateJtw(id uint, scope string) (string, error) {
	payload := ClaimWithScope{}
	payload.Subject = strconv.Itoa(int(id))
	payload.ExpiresAt = time.Now().Add(time.Hour * 2).Unix()
	payload.Scope = scope

	return jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(SecurityKey))
}
