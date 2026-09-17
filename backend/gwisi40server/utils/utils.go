package utils

import (
	"os"
	"path/filepath"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyAuthentication(c *fiber.Ctx, cookie string) (*jwt.StandardClaims, error) {

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims := token.Claims.(*jwt.StandardClaims)

	return claims, nil
}

func Running() bool {

	processes := []string{}

	myself, _ := os.Executable()
	myself = filepath.Base(myself)

	v, _ := process.Processes()

	for _, p := range v {
		proc, err := p.Name()
		if err == nil && proc == myself {
			processes = append(processes, proc)
		}
	}

	return len(processes) > 1
}
