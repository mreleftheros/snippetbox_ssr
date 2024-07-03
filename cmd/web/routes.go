package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.indexGet)
	mux.HandleFunc("/snippets/{id}", app.snippetsIdParamGet)
	mux.HandleFunc("GET /snippets/new", app.snippetsNewGet)
	mux.HandleFunc("POST /snippets/new", app.snippetsNewPost)

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
