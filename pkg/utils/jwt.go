package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignToken(userId int, username, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpires := os.Getenv("JWT_EXPIRES_IN")

	claims := jwt.MapClaims{
		"uid":  userId,
		"user": username,
		"role": role,
	}

	if jwtExpires != "" {
		duration, err := time.ParseDuration(jwtExpires)
		if err != nil {
			return "", ErrorHandler(err, "jwt expire error")
		}
		claims["exp"] = time.Now().Add(duration).Unix()
	} else {
		claims["exp"] = time.Now().Add(15 * time.Minute).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", ErrorHandler(err, "error siggning token")
	}

	return signedToken, nil
}
