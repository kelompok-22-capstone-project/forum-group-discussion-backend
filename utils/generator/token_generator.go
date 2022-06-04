package generator

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type TokenGenerator interface {
	GenerateToken(id, username, role string, isActive bool) (token string, err error)
	ExtractToken(c echo.Context) (id, username, role string, isActive bool)
}

type jwtTokenGenerator struct{}

func NewJWTTokenGenerator() *jwtTokenGenerator {
	return &jwtTokenGenerator{}
}

func (j *jwtTokenGenerator) GenerateToken(id, username, role string, isActive bool) (token string, err error) {
	claims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"role":     role,
		"isActive": isActive,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}

	jwtWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtWithClaims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return
}

func (j *jwtTokenGenerator) ExtractToken(c echo.Context) (id, username, role string, isActive bool) {
	user := c.Get("user").(*jwt.Token)

	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		id = claims["id"].(string)
		username = claims["username"].(string)
		role = claims["role"].(string)
		isActive = claims["isActive"].(bool)
	}

	return
}
