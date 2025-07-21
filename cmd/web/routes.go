package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	//mux := http.NewServeMux()

	//initialize the router
	router := httprouter.New()

	//fileServer := http.FileServer(http.Dir("./ui/static/"))
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	//mux.HandleFunc("/", app.home)
	//mux.HandleFunc("/snippet/view", app.snippetView)
	//mux.HandleFunc("/snippet/create", app.snippetCreate)

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	//return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	//return standard.Then(mux)

	return standard.Then(router)
}
