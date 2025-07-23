package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
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

}

func (app *application) snippetCreate(res http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)

	app.render(res, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		app.clientError(res, http.StatusBadRequest)
		return
	}

	// Use the req.PostForm.Get() method to retrieve the title and content
	// from the req.PostForm map
	title := req.PostForm.Get("title")
	content := req.PostForm.Get("content")

	//manually convert form data (expires) to an integer using strconv.Atoi()
	expires, err := strconv.Atoi(req.PostForm.Get("expires"))
	if err != nil {
		app.clientError(res, http.StatusBadRequest)
	}

	// initialize a map to hold validation errors
	fieldErrors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This title field cannot be more than 100 characters long"

	}

	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "The content section cannot be blank"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors["expires"] = "This field must equal 1 ,7 or 365"
	}

	if len(fieldErrors) > 0 {
		fmt.Fprint(res, fieldErrors)
		return
	}

	//pass the data to SnippetModel.insert() method
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(res, err)
		return
	}

	//redirect the user to the relevant page for the snippet
	http.Redirect(res, req, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	res.Header().Set("Allow", "POST")
}
