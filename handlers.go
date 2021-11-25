package main

import (
	"net/http"
	"strconv"

	"gghub.com/model"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	allTodo, err := model.GetFullTodo(id)
	if err != nil {
		//
	}
	attr := map[string]interface{}{
		"Checklist": allTodo,
	}
	getHTMLTemplate("index.html").ExecuteTemplate(w, "index.html", attr)
}

func postTodo(w http.ResponseWriter, r *http.Request) {

}
