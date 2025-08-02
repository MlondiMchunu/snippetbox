package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.mlodev.net/internal/models"
	"snippetbox.mlodev.net/internal/validator"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
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

	//flash := app.sessionManager.PopString(req.Context(), "flash")

	data := app.newTemplateData(req)
	data.Snippet = snippet

	//data.Flash = flash

	//use the render helper
	app.render(res, http.StatusOK, "view.tmpl", data)

}

func (app *application) snippetCreate(res http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(res, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(res http.ResponseWriter, req *http.Request) {

	var form snippetCreateForm

	err := app.decodePostForm(req, &form)
	if err != nil {
		app.clientError(res, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "The title field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "The title field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "The content section cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1 ,7 or 365")

	if !form.Valid() {
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
	app.sessionManager.Put(req.Context(), "flash", "Snippet succesfully created!!")

	//redirect the user to the relevant page for the snippet
	http.Redirect(res, req, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	res.Header().Set("Allow", "POST")
}

func (app *application) userSignup(res http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)
	data.Form = userSignupForm{}
	app.render(res, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(res http.ResponseWriter, req *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(req, &form)
	if err != nil {
		app.clientError(res, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "The name field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "The email field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Invalid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "The password field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "The password field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(req)
		data.Form = form
		app.render(res, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	fmt.Fprintln(res, "Create a new user...")

}

func (app *application) userLogin(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Display a HTML form for logging in a user...")
}

func (app *application) userLoginPost(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Authenticate and login the user...")
}

func (app *application) userLogoutPost(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Logout the user...")
}
