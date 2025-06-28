package main

import (
	"log"
	"net/http"
)

func main() {
	//use the http.NewServeMux() i.e router function to initialize a new servemux
	mux := http.NewServeMux()

	//create a file server which serves files out of the "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//use the mux.Handle() function to register the file server as the handler for
	//all URL paths that start with "/static/"
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

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

	//find . -name "*.go" | entr -r sh -c 'echo "== Restarting =="; go run ./cmd/web'
}
