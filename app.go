package main

import (
	"fmt"
	"net/http"
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

	fmt.Println("server is running on port: ", port)

	http.ListenAndServe(fmt.Sprintf(":%d",port),nil)

}