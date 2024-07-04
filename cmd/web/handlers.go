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
		data.SnippetErrors = snippetErrors

		app.render(w, "create.tmpl", data, http.StatusUnprocessableEntity)
		return
	}

	id, err := app.snippets.Insert(form)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
	return
}
