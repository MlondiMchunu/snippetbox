package main

import (
	"log"
	"net/http"
)

// define home handler functions i.e controller which writes a byte slice containing
func home(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello from snippetbox"))
}
func snippetView(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Display a specific snippet...."))
}
func snippetCreate(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Create a new snippet..."))
}

func main() {
	//use the http.NewServeMux() i.e router function to initialize a new servemux, then,
	//register the home function as the handler for the "/" URL pattern
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	port := 4000

	log.Println("Starting server on : ", port)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
