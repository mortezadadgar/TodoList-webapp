// In doubt of having this function as pkg!
package render

import (
	"errors"
	"html/template"
	"net/http"
)

var (
	ErrInvalidName = errors.New("invalid template name")
)

type Render struct {
	cachedTemplate map[string]*template.Template
}

func NewRender(cachedTemplate map[string]*template.Template) *Render {
	return &Render{
		cachedTemplate: cachedTemplate,
	}
}

func (r Render) Render(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := r.cachedTemplate[name]
	if !ok {
		return ErrInvalidName
	}

	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		return err
	}

	return nil
}
