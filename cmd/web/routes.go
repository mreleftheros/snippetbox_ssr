package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/", app.sessionManager.LoadAndSave(http.HandlerFunc(app.indexGet)))

	mux.Handle("GET /users/signup", app.sessionManager.LoadAndSave(http.HandlerFunc(app.usersSignupGet)))
	mux.Handle("POST /users/signup", app.sessionManager.LoadAndSave(http.HandlerFunc(app.usersSignupPost)))
	mux.Handle("GET /users/login", app.sessionManager.LoadAndSave(http.HandlerFunc(app.usersLoginGet)))
	mux.Handle("POST /users/login", app.sessionManager.LoadAndSave(http.HandlerFunc(app.usersLoginPost)))
	mux.Handle("GET /users/logout", app.sessionManager.LoadAndSave(http.HandlerFunc(app.usersLogoutGet)))

	mux.Handle("/snippets/{id}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetsIdParamGet)))
	mux.Handle("GET /snippets/new", app.sessionManager.LoadAndSave(app.user(http.HandlerFunc(app.snippetsNewGet))))
	mux.Handle("POST /snippets/new", app.sessionManager.LoadAndSave(app.user(http.HandlerFunc(app.snippetsNewPost))))

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
