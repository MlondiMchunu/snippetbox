package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"snippetbox.mlodev.net/internal/models"
)

// define home handler functions i.e controller which writes a byte slice containing
func (app *application) home(res http.ResponseWriter, req *http.Request) {
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

	/*files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			"./ui/html/pages/home.tmpl",
		}func (m *SnippetModel) Latest() ([]*Snippet, error) {
		return nil, nil
	}

		//use template.ParseFiles() function to read the files and store
		//the templates in a template set
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(res, err)
			return
		}

		err = ts.ExecuteTemplate(res, "base", nil)
		if err != nil {
			app.serverError(res, err)
		}
	*/

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(res, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(res, "%+v\n", snippet)
	}
}
func (app *application) snippetView(res http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(res)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(res)
		} else {
			app.serverError(res, err)
		}
		return
	}

	//Initialize a slice containing the paths to the view.tmpl file
	//plus base the layout and navigation partial that we made earlier
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
	}

	//parse the template files
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(res, err)
		return
	}

	err = ts.ExecuteTemplate(res, "base", snippet)
	if err != nil {
		app.serverError(res, err)
	}

	//write the snippet data as plain text HTTP response body
	fmt.Fprintf(res, "%+v", snippet)

	fmt.Fprintf(res, "Display a specific snippet with ID %d...", id)

	res.Write([]byte("Display a specific snippet...."))
}
func (app *application) snippetCreate(res http.ResponseWriter, req *http.Request) {
	//use r.Method to check the request us using POST or not
	if req.Method != "POST" {
		res.Header().Set("Allow", http.MethodPost)
		res.Header().Set("Cache-control", "public,max-age=3135600")

		/*can use http.Error shortcut to combine
		res.WriteHeader() & res.Write() methods*/

		//res.WriteHeader(405)
		//res.Write([]byte("Method not allowed!!!"))
		app.clientError(res, http.StatusMethodNotAllowed)
		fmt.Println("Method not allowed!!!")
		return
	}

	//variables holding dummy data
	title := "o Snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	//pass the data to SnippetModel.insert() method
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(res, err)
		return
	}

	//redirect the user to the relevant page for the snippet
	http.Redirect(res, req, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	res.Write([]byte("Create a new snippet..."))
	res.Header().Set("Allow", "POST")
}
