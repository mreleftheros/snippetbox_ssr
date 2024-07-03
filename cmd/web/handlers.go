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
		app.serverError(w, err)
		return
	}
	f := r.PostForm
	title := f.Get("title")
	content := f.Get("content")
	expires, err := strconv.Atoi(f.Get("expires"))
	if err != nil || expires < 0 {
		app.clientError(w, 400)
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
