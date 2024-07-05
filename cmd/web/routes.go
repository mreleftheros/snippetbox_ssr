package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/", app.sessionManager.LoadAndSave(http.HandlerFunc(app.indexGet)))
	mux.Handle("/snippets/{id}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetsIdParamGet)))
	mux.Handle("GET /snippets/new", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetsNewGet)))
	mux.Handle("POST /snippets/new", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetsNewPost)))

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
