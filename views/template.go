package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/codewithtoucans/goweb/context"
	"github.com/codewithtoucans/goweb/models"
	"github.com/gorilla/csrf"
)

type Template struct {
	htmlTpl *template.Template
}

func Must(template Template, err error) Template {
	if err != nil {
		log.Println("something was err in Must function")
		panic(err)
	}
	return template
}

func ParseFS(fs fs.FS, filePath ...string) (Template, error) {
	t := template.New(filePath[0])
	t.Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "implement it", fmt.Errorf("not implement csrfField")
		},
		"currentUser": func() (*models.User, error) {
			return nil, fmt.Errorf("not implement currentUser")
		},
		"errors": func() []string {
			return nil
		},
	})
	tpl, err := t.ParseFS(fs, filePath...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{htmlTpl: tpl}, nil
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any, errs ...error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("template clone was error\n")
		http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
		return
	}
	tpl = tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"currentUser": func() *models.User {
			return context.User(r.Context())
		},
		"errors": func() []string {
			var errMessages []string
			for _, e := range errs {
				errMessages = append(errMessages, e.Error())
			}
			return errMessages
		},
	})
	err = tpl.Execute(w, data)
	if err != nil {
		log.Printf("parsing template: %v\n", err)
		http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
		return
	}
}
