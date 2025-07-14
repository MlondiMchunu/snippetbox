package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(res http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(res http.ResponseWriter, status int) {
	http.Error(res, http.StatusText(status), status)
}

// this helper is a convenience wrapper aroud clientError which sends a 404 Not Found response
func (app *application) notFound(res http.ResponseWriter) {
	app.clientError(res, http.StatusNotFound)
}

func (app *application) render(res http.ResponseWriter, status int, page string, data *templateData) {
	//retrieve the appropriate template set from cache based on page name
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(res, err)
		return
	}

	//write the provided HTTP status code
	res.WriteHeader(status)
}
