package controllers

import (
	"github.com/codewithtoucans/goweb/views"
	"html/template"
	"net/http"
)

func StaticHandler(t views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, r, nil)
	}
}

func FAQ(t views.Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "Is there a free version?",
			Answer:   "We offer a free trial for 30 days!",
		},
		{
			Question: "What are you support?",
			Answer:   "Go Java Python etc.",
		},
		{
			Question: "How can I contact support?",
			Answer:   `Email me <a href="https://www.google.com">sandycorinn@gmail.com</a>`,
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, r, questions)
	}
}
