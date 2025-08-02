package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	mid "simpleapi/internal/api/middlewares"
	"simpleapi/internal/api/router"
	"simpleapi/internal/repositories/sqlconnect"
	"simpleapi/pkg/utils"
)

func main() {

	_, err := sqlconnect.ConnectDB("dbever_test")
	if err != nil {
		fmt.Println("error connecting to db", err)
		return
	}

	port := 3000

	key := `certificate\key.pem`
	cert := `certificate\certificate.pem`

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	rateLimiter := mid.NewRateLimiter(5, time.Second*5)

	hppSettings := &mid.HppOptions{
		CheckQuery:              true,
		CheckBody:               true,
		CheckBodyForContentType: "application/x-www-form-urlencoded",
		WhiteList:               []string{"allowedParam", "sortOrder", "sortBy", "name", "age", "class", "first_name", "last_name"},
	}

	hppMiddleware := mid.Hpp(*hppSettings)

	router := router.Router()

	// secureMux := mid.Cors(rateLimiter.Middleware(mid.ResponseTime(mid.SecurityHeaders(mid.CompMiddleware(hppMiddleware(router))))))
	// secureMux := applyMiddlewares(router,hppMiddleware,mid.CompMiddleware,mid.SecurityHeaders,mid.ResponseTime,rateLimiter.Middleware,mid.Cors)
	secureMux := utils.ApplyMiddlewares(router, hppMiddleware, rateLimiter.Middleware) // for now faster processing

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("server is running on port:", port)

	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		fmt.Println("error is:", err)
		return
	}

}
