package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.mlodev.net/internal/models"
)

// define home handler functions i.e controller which writes a byte slice containing
func (app *application) home(res http.ResponseWriter, req *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(res, err)
		return
	}

	//call newTemplateData() helper to get a templateData struct
	data := app.newTemplateData(req)
	data.Snippets = snippets

	//use the render helper
	app.render(res, http.StatusOK, "home.tmpl", data)

}

func (app *application) snippetView(res http.ResponseWriter, req *http.Request) {
	//values of any named parameters will be stored in request context,
	//when httprouter parses a request
	params := httprouter.ParamsFromContext(req.Context())

	id, err := strconv.Atoi(params.ByName("id"))
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

	data := app.newTemplateData(req)
	data.Snippet = snippet

	//use the render helper
	app.render(res, http.StatusOK, "view.tmpl", data)

	/*
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

		//create an instance of a templateData struct holding the snippet data
		//add instance of a templateData struct holding the slice of snippets
		data := &templateData{
			Snippet: snippet,
		}

		err = ts.ExecuteTemplate(res, "base", data)
		if err != nil {
			app.serverError(res, err)
		}
	*/

	//write the snippet data as plain text HTTP response body
	//fmt.Fprintf(res, "%+v", snippet)

	//fmt.Fprintf(res, "Display a specific snippet with ID %d...", id)

	//res.Write([]byte("Display a specific snippet...."))
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
