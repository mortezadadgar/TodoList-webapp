package cachetemplate

import (
	"html/template"
	"path/filepath"
)

type templateMap map[string]*template.Template

func New(dir string) (templateMap, error) {
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	cache := make(templateMap)
	for _, page := range pages {
		tmpl, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		tmpl, err = tmpl.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		name := filepath.Base(page)
		cache[name] = tmpl
	}

	return cache, nil
}
