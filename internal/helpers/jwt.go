package helpers

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET = []byte("supersecret")

func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SECRET)
}

func ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SECRET, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*jwt.MapClaims), nil
}
