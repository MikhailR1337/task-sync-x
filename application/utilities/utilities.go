package utilities

import (
	"errors"

	"github.com/MikhailR1337/task-sync-x/initializers"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	errUnauthorized = errors.New("Unauthorized")
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetJwtPayload(c *fiber.Ctx) (jwt.MapClaims, error) {
	jwtToken, ok := c.Context().Value(initializers.Cfg.ContextKeyUser).(*jwt.Token)
	if !ok {
		return nil, errUnauthorized
	}
	jwtPayload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errUnauthorized
	}

	return jwtPayload, nil
}
