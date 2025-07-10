package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"golang.org/x/net/http2"
)

func homeResp(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"our home page")
}

func ordersResp(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"handling incomming orders")
}

func usersResp(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"handling our users")
}

func main(){

	http.HandleFunc("/", homeResp)
	http.HandleFunc("/orders", ordersResp)
	http.HandleFunc("/users", usersResp)


	port := 3000

	// for http2

	cert := `certification\certificate.pem`
	key := `certification\key.pem`

	//configure tls
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	//create a custom server
	server := &http.Server{
		Addr: fmt.Sprintf(":%d",port),
		Handler: nil,
		TLSConfig: tlsConfig,
	}

	// enable http2
	// this configures the http server that we created to have https functionality
	// its configuring the server
	http2.ConfigureServer(server, &http2.Server{})



	fmt.Println("server is running on port: ", port)


	err := server.ListenAndServeTLS(cert,key)
	if err != nil {
		fmt.Println("error is:", err)
		return
	}


	// http 1.1

	// err := http.ListenAndServe(fmt.Sprintf(":%d",port),nil)
	// if err != nil {
	// 	fmt.Println("error is:", err)
	// 	return
	// }

}