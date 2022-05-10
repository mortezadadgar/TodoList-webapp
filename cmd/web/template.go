package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

type templateData struct {
	ThisYear        int
	IsAuthenticated bool
	User            string
	UserTodo        todo
}

func newCacheTemplate(dir string) (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		tmpl, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// NOTE: Parse layouts with every page
		// TODO: investigate this
		tmpl, err = tmpl.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = tmpl
	}

	return cache, nil
}

func (app *application) addDefaultData(r *http.Request, t *templateData) *templateData {
	if t == nil {
		t = &app.templateData
	}

	if app.isAuthenticated(r) {
		t.IsAuthenticated = true
		t.User = app.session.GetString(r, "username")
	}

	t.ThisYear = time.Now().Year()

	return t
}
