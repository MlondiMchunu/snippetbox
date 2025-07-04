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
