package main

import (
	"log"
	"net/http"
)

func main() {
	//use the http.NewServeMux() i.e router function to initialize a new servemux, then,
	//register the home function as the handler for the "/" URL pattern
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	/*register routes without declaring a servemux
	*NB avoid on production apps for security reasons
	 */
	//http.HandleFunc("/", home)

	port := 4000

	log.Println("Starting server on : ", port)

	err := http.ListenAndServe(":4000", mux)

	/*part of registering routes withut declaring a servemux*/
	//err := http.ListenAndServe(":4000,nil")
	log.Fatal(err)
}
