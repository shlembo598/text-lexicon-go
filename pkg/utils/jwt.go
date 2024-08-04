package utils

import (
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/shlembo598/text-lexicon-go/internal/config"
	"github.com/shlembo598/text-lexicon-go/internal/models"
)

// JWT Claims struct
type Claims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.StandardClaims
}

func GenerateJWTToken(user *models.User, config *config.Config) (string, error) {
	claims := &Claims{
		Email: user.Email,
		ID:    user.UserID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.Server.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
