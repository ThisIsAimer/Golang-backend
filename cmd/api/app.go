package main

import (
	"fmt"
	"net/http"
)

func homeRoute(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "string")
	w.Write([]byte("This is home"))
	fmt.Println("someone accessed: home")
}

func teachersRoute(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "string")
	fmt.Fprintf(w,"This is teachers route")
	fmt.Println("someone accessed: Teachers route")
}

func studentsRoute(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "string")
	fmt.Fprintf(w,"This is students route")
	fmt.Println("someone accessed: Students route")
}

func execsRoute(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "string")
	w.Write([]byte("This is executives route"))
	fmt.Println("someone accessed: Execs route")
}

func main(){

	http.HandleFunc("/", homeRoute)
	http.HandleFunc("/teachers", teachersRoute)
	http.HandleFunc("/students", studentsRoute)
	http.HandleFunc("/execs", execsRoute)


	port := 3000

	server := &http.Server{
		Addr: fmt.Sprintf(":%d",port),
		Handler: nil,
	}

	fmt.Println("server is running on port:", port)


	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("error is:", err)
		return
	}

}