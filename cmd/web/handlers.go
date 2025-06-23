package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// define home handler functions i.e controller which writes a byte slice containing
func home(res http.ResponseWriter, req *http.Request) {
	// restrict root url pattern
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}
	//res.Write([]byte("Hello from snippetbox"))
	ts, err := template.ParseFiles("./ui/html/pages/home.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(res, "Internal Server Error", 500)
		return
	}
}
func snippetView(res http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(res, req)
		return
	}
	fmt.Fprintf(res, "Display a specific snippet with ID %d...", id)

	res.Write([]byte("Display a specific snippet...."))
}
func snippetCreate(res http.ResponseWriter, req *http.Request) {
	//use r.Method to check the request us using POST or not
	if req.Method != "POST" {
		res.Header().Set("Allow", "POST")
		res.Header().Set("Cache-control", "public,max-age=3135600")

		/*can use http.Error shortcut to combine
		res.WriteHeader() & res.Write() methods*/

		//res.WriteHeader(405)
		//res.Write([]byte("Method not allowed!!!"))
		http.Error(res, "Method Not Allowed!!!", 405)

		fmt.Println("Method not allowed!!!")
		return
	}

	res.Write([]byte("Create a new snippet..."))
	res.Header().Set("Allow", "POST")
}
