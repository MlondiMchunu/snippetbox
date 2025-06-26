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
	/*ts, err := template.ParseFiles("./ui/html/pages/home.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(res, "Internal Server Error", 500)
		return
	}
	*/
	//initialize a slice containing the paths to the two files
	//gile containing base template should be the first
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	//use template.ParseFiles() function to read the files and store
	//the templates in a template set
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(res, "Inernal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(res, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(res, "Internal Server Error", 500)
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
