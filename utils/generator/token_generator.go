package generator

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type TokenPayload struct {
	ID       string
	Username string
	Role     string
	IsActive bool
}

type TokenGenerator interface {
	GenerateToken(payload TokenPayload) (token string, err error)
	ExtractToken(c echo.Context) (payload TokenPayload)
}

type jwtTokenGenerator struct{}

func NewJWTTokenGenerator() *jwtTokenGenerator {
	return &jwtTokenGenerator{}
}

func (j *jwtTokenGenerator) GenerateToken(payload TokenPayload) (token string, err error) {
	claims := jwt.MapClaims{
		"id":       payload.ID,
		"username": payload.Username,
		"role":     payload.Role,
		"isActive": payload.IsActive,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}

	jwtWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtWithClaims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return
}

func (j *jwtTokenGenerator) ExtractToken(c echo.Context) (payload TokenPayload) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return
	}

	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		payload.ID = claims["id"].(string)
		payload.Username = claims["username"].(string)
		payload.Role = claims["role"].(string)
		payload.IsActive = claims["isActive"].(bool)
	}

	return
}
