package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	mid "simpleapi/internal/api/middlewares"
	"simpleapi/internal/api/router"
	"simpleapi/pkg/utils"
	"github.com/joho/godotenv"
)

func main() {
	// used for loading .env variables to environment variables list
	err := godotenv.Load(`cmd\api\.env`)
	if err != nil {
		fmt.Println("error loading .env", err)
		return 
	}

	port := os.Getenv("API_PORT")

	key := `certificate\key.pem`
	cert := `certificate\certificate.pem`

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	rateLimiter := mid.NewRateLimiter(5, time.Second*5)

	whiteList := []string{
		"sortOrder", 
		"sortBy", 

		// genral
		"first_name", 
		"last_name",
		"class", 

		// teachers
		"email", 
		"subject",
	}

	hppSettings := &mid.HppOptions{
		CheckQuery:              true,
		CheckBody:               true,
		CheckBodyForContentType: "application/x-www-form-urlencoded",
		WhiteList:               whiteList,
	}

	hppMiddleware := mid.Hpp(*hppSettings)

	router := router.Router()

	// secureMux := mid.Cors(rateLimiter.Middleware(mid.ResponseTime(mid.SecurityHeaders(mid.CompMiddleware(hppMiddleware(router))))))
	// secureMux := applyMiddlewares(router,hppMiddleware,mid.CompMiddleware,mid.SecurityHeaders,mid.ResponseTime,rateLimiter.Middleware,mid.Cors)
	secureMux := utils.ApplyMiddlewares(router, hppMiddleware, rateLimiter.Middleware) // for now faster processing

	server := &http.Server{
		Addr:      port,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("server is running on port", port)

	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		fmt.Println("error is:", err)
		return
	}

}
