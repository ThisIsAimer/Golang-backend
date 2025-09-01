package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"simpleapi/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

func JwtMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println(r.Cookies())
		token, err := r.Cookie("Bearer")
		if err != nil {
			http.Error(w, "unauthorised", http.StatusBadRequest)
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")

		parsedToken, err := jwt.Parse(token.Value, func(token *jwt.Token) (any, error) {
			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(jwtSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				myErr := utils.ErrorHandler(err, "token expired")
				http.Error(w, myErr.Error(), http.StatusUnauthorized)
				return
			}
			myErr := utils.ErrorHandler(err, "unauthorised access")
			http.Error(w, myErr.Error(), http.StatusUnauthorized)
			return
		}

		if !parsedToken.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			log.Println("invalid jwt:", token.Value)
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)

		if !ok {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			log.Println("invalid login token:", token.Value)
			return
		}

		ctx := context.WithValue(r.Context(), "role", claims["role"])
		ctx = context.WithValue(ctx, "expiresAt", claims["exp"])
		ctx = context.WithValue(ctx, "username", claims["user"])
		ctx = context.WithValue(ctx, "userId", claims["uid"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
