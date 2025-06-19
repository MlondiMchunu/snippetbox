package main

import (
	"fmt"
	"log"
	"net/http"
)

// define home handler functions i.e controller which writes a byte slice containing
func home(res http.ResponseWriter, req *http.Request) {
	// restrict root url pattern
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}
	res.Write([]byte("Hello from snippetbox"))
}
func snippetView(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Display a specific snippet...."))
}
func snippetCreate(res http.ResponseWriter, req *http.Request) {
	//use r.Method to check the request us using POST or not
	if req.Method != "POST" {
		res.Header().Set("Allow", "POST")
		res.WriteHeader(405)
		//res.WriteHeader(406)
		res.Write([]byte("Method not allowed!!!"))
		fmt.Println("Method not allowed!!!")
		return
	}

	res.Write([]byte("Create a new snippet..."))
}

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
