package main

import (
	"log"
	"net/http"
)

//define home handler function which writes a byte slice containing
//"Hello from Snippetbox" as the response body.
func home(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello from Snippetbox"))
}

func main() {
	//use the http.NewServeMux() function to initialize a new servemux, then,
	//register the home function as the handler for the "/" URL pattern
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	port := 4000

	log.Println("Starting server on : ", port)
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
