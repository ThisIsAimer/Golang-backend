package main

import (
	"crypto/tls"
	//"crypto/x509"
	"fmt"
	"log"
	"net/http"
	//"os"

	"golang.org/x/net/http2"
)

func homeResp(w http.ResponseWriter, r *http.Request){
	logRequestDetails(r)
	w.Header().Set("Content-Type", "string")
	fmt.Fprintf(w,"our home page")
}

func ordersResp(w http.ResponseWriter, r *http.Request){
	logRequestDetails(r)
	w.Header().Set("Content-Type", "string")
	fmt.Fprintf(w,"handling incomming orders")
}

func usersResp(w http.ResponseWriter, r *http.Request){
	logRequestDetails(r)
	w.Header().Set("Content-Type", "string")
	fmt.Fprintf(w,"handling our users")
}

//------------------------------------------------------------------
// for reading client Ca
// Now client must also have the pem certificate to access the page!

// func loadClientCas() *x509.CertPool{
// 	clientCas := x509.NewCertPool()
// 	file, err := os.ReadFile(`certification\certificate.pem`)
// 	if err != nil {
// 		log.Fatalln("couldnt load client cert:", err)
// 	}

// 	clientCas.AppendCertsFromPEM(file)

// 	return clientCas

// }
//------------------------------------------------------------------

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
		//------------------------------------------------------------------
		// client certification

		// ClientAuth: tls.RequireAndVerifyClientCert,
		// ClientCAs: loadClientCas(),
		//------------------------------------------------------------------
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


func logRequestDetails(r *http.Request){
	httpVersion := r.Proto
	fmt.Println("-------------------------------------------------------------")
	log.Println("recieved request with http version:", httpVersion)
	if r.TLS != nil{
		tlsVersion := tlsVersionDetails(r.TLS.Version)
		log.Println("Tls version:", tlsVersion)
	} else {
		log.Println("no tls used")
	}
}

func tlsVersionDetails(version uint16) string{
	switch version{
	case tls.VersionTLS10:
		return "Tls 1.0"
	case tls.VersionTLS11:
		return "Tls 1.1"
	case tls.VersionTLS12:
		return "Tls 1.2"
	case tls.VersionTLS13:
		return "Tls 1.3"
	default:
		return "unknown TLS version"

	}
}


// http: TLS handshake error from [::X]:XXXXX: EOF
// TLS handshake error happens as there is smth wrong with the certification
// EOF in networking means connection termination or protocal error