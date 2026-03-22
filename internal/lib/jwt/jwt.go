package jwt_lib

import (
	"os"
	"time"

	"github.com/esquirelol/auth-grpc/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, app models.App, ttl time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(ttl).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
