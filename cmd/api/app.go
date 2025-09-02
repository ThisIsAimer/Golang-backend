package main

import (
	"crypto/tls"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	mid "simpleapi/internal/api/middlewares"
	"simpleapi/internal/api/router"
	"simpleapi/pkg/utils"

	"github.com/joho/godotenv"
)

//go:embed .env
var envFile embed.FS
func loadEnvFromEmbeddedFile() {
	//read the embedded .env file
	content, err := envFile.ReadFile(`.env`) //cmd\api\.env

	if err != nil {
		log.Fatalf("error reading .env file: %v", err)
	}

	//create a temp file with content
	tempFile, err := os.CreateTemp("", ".env")
	if err != nil {
		log.Fatalf("error creating temp .env file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// write content in the file
	_, err = tempFile.Write(content)
	if err != nil {
		log.Fatalf("error writing to the tempFile file: %v", err)
	}

	err = tempFile.Close()
	if err != nil {
		log.Fatalf("error closing tempFile file: %v", err)
	}
	err = godotenv.Load(tempFile.Name())
	if err != nil {
		log.Fatalf("error loading env variables: %v", err)
	}
}

// main function
func main() {
	// used for loading .env variables to environment variables list
	//only for development phase

	// err := godotenv.Load(`cmd\api\.env`)
	// if err != nil {
	// 	utils.ErrorHandler(fmt.Errorf("error getting env files"), "error starting server")
	// 	return
	// }

	// for production
	loadEnvFromEmbeddedFile()

	port := os.Getenv("API_PORT")

	key := os.Getenv("KEY_FILE")
	cert := os.Getenv("CERT_FILE")

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	rateLimiter := mid.NewRateLimiter(5, time.Second*5)

	whiteList := []string{
		"sortby",

		// genral
		"first_name",
		"last_name",
		"class",

		// teachers
		"email",
		"subject",

		//pagination
		"page",
		"limit",
	}

	hppSettings := &mid.HppOptions{
		CheckQuery:              true,
		CheckBody:               true,
		CheckBodyForContentType: "application/x-www-form-urlencoded",
		WhiteList:               whiteList,
	}

	hppMiddleware := mid.Hpp(*hppSettings)

	router := router.Router()

	jwtMiddleware := mid.SkipJwtRoutes(mid.JwtMiddleware, "/execs/login", "/execs/login/forgotusername", "/execs/login/forgotpassword", "/resetpassword/reset")

	// secureMux := mid.Cors(rateLimiter.Middleware(mid.ResponseTime(mid.SecurityHeaders(mid.CompMiddleware(hppMiddleware(router))))))
	// secureMux := applyMiddlewares(router,hppMiddleware,mid.CompMiddleware,mid.SecurityHeaders,mid.ResponseTime,rateLimiter.Middleware,mid.Cors)
	// final middleware sequence
	secureMux := utils.ApplyMiddlewares(router, mid.SecurityHeaders, mid.CompMiddleware, hppMiddleware, mid.XSSMiddleware, jwtMiddleware, mid.ResponseTime, rateLimiter.Middleware, mid.Cors) // for now faster processing

	server := &http.Server{
		Addr:      port,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("server is running on port", port)

	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		utils.ErrorHandler(fmt.Errorf("tls error"), "error starting server")
		return
	}

}
