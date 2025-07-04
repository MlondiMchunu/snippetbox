package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(res http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
	http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(res http.ResponseWriter, status int) {
	http.Error(res, http.StatusText(status), status)
}

// this helper is a convenience wrapper aroud clientError which sends a 404 Not Found response
func (app *application) notFound(res http.ResponseWriter) {
	app.clientError(res, http.StatusNotFound)
}
