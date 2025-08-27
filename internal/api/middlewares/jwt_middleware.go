package middlewares

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"simpleapi/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

func JwtMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println(r.Cookies())
		cookie, err := r.Cookie("Bearer")
		if err != nil {
			http.Error(w, "unauthorised", http.StatusBadRequest)
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")

		parsedToken, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (any, error) {
			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(jwtSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil && r.URL.Path != "/execs/login" {
			if errors.Is(err, jwt.ErrTokenExpired) {
				myErr := utils.ErrorHandler(err, "token expired")
				http.Error(w, myErr.Error(), http.StatusUnauthorized)
				return
			}
			myErr := utils.ErrorHandler(err, "unauthorised access")
			http.Error(w, myErr.Error(), http.StatusUnauthorized)
			return
		}

		if parsedToken.Valid {
			log.Println("valid Jwt")
		}

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			fmt.Println(claims["uid"], claims["exp"])
		} else {
			fmt.Println(err)
		}

		next.ServeHTTP(w, r)
	})

}
