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

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

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

	//manually convert form data (expires) to an integer using strconv.Atoi()
	expires, err := strconv.Atoi(req.PostForm.Get("expires"))
	if err != nil {
		app.clientError(res, http.StatusBadRequest)
	}

	form := snippetCreateForm{
		Title:       req.PostForm.Get("title"),
		Content:     req.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "The title field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "The title field cannot be more than 100 characters long"

	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "The content section cannot be blank"
	}

	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1 ,7 or 365"
	}

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(req)
		data.Form = form
		app.render(res, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	//pass the data to SnippetModel.insert() method
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(res, err)
		return
	}

	//redirect the user to the relevant page for the snippet
	http.Redirect(res, req, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	res.Header().Set("Allow", "POST")
}
