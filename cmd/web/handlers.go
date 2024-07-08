package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) indexGet(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, "home.tmpl", data, 200)
}

func (app *application) snippetsIdParamGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 0 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, "view.tmpl", data, 200)
}

func (app *application) snippetsNewGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, "create.tmpl", data, 200)
}

func (app *application) snippetsNewPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, 400)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, 400)
	}

	form := app.snippets.NewSnippetForm(title, content, expires)

	if snippetErrors, ok := app.snippets.Validate(form); !ok {
		data := app.newTemplateData(r)
		data.SnippetForm = form
		data.Errors = snippetErrors

		app.render(w, "create.tmpl", data, 422)
		return
	}

	id, err := app.snippets.Insert(form)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet created successfully")

	http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
	return
}

func (app *application) usersSignupGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, "signup.tmpl", data, 200)
}

func (app *application) usersSignupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, 400)
		return
	}

	name, email, password := r.PostForm.Get("name"), r.PostForm.Get("email"), r.PostForm.Get("password")

	form := app.users.NewUserSignupForm(name, email, password)
	data := app.newTemplateData(r)

	if userErrors, ok := app.users.Validate(form); !ok {
		data.UserSignupForm = form
		data.Errors = userErrors

		app.render(w, "signup.tmpl", data, 422)
		return
	}

	_, err = app.users.Signup(form)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			data.UserSignupForm = form
			data.Errors = &map[string]string{"error": "Email already exists"}

			app.render(w, "signup.tmpl", data, 400)
		}
	}

	app.sessionManager.Put(r.Context(), "flash", "User created successfully")
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func (app *application) usersLoginGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, "login.tmpl", data, 200)
}

func (app *application) usersLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, 400)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	form := app.users.NewUserLoginForm(email, password)

	user, err := app.users.Login(form)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") || strings.Contains(err.Error(), "incorrect password") {
			data := app.newTemplateData(r)
			data.Errors = &map[string]string{"error": "Invalid credentials"}

			app.render(w, "login.tmpl", data, 400)
		}
	}

	if err := app.sessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "userId", user.Id)

	http.Redirect(w, r, "/", 303)
}

func (app *application) usersLogoutGet(w http.ResponseWriter, r *http.Request) {
	if err := app.sessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "userId")
	app.sessionManager.Put(r.Context(), "flash", "Logged out successfully")

	http.Redirect(w, r, "/users/login", 303)
}
