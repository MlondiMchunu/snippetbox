package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	//mux := http.NewServeMux()

	//initialize the router
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		app.notFound(res)
	})

	//fileServer := http.FileServer(http.Dir("./ui/static/"))
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	//mux.HandleFunc("/", app.home)
	//mux.HandleFunc("/snippet/view", app.snippetView)
	//mux.HandleFunc("/snippet/create", app.snippetCreate)

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	//return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	//return standard.Then(mux)

	return standard.Then(router)
}
