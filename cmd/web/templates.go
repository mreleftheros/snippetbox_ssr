package main

import (
	"github.com/mreleftheros/snippetbox_ssr/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

type templateData struct {
	Snippet        *models.Snippet
	Snippets       []*models.Snippet
	CurrentYear    int
	SnippetForm    *models.SnippetForm
	Errors         *map[string]string
	UserSignupForm *models.UserSignupForm
	Flash          string
	User           *models.User
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	tD := templateData{
		CurrentYear: time.Now().Year(),
		Flash:       app.sessionManager.PopString(r.Context(), "flash"),
	}

	user, ok := r.Context().Value("user").(*models.User)
	if ok {
		tD.User = user
	}

	return &tD
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		tmpl, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		tmpl, err = tmpl.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		tmpl, err = tmpl.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = tmpl
	}

	return cache, nil
}
