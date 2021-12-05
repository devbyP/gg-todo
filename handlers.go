package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"gghub.com/model"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	allTodo, err := model.GetAllFullTodo()
	if err != nil {
		log.Printf("Error cannot get todolist: %v\n", err)
		http.Error(w, "cannot get todolist.", http.StatusInternalServerError)
		return
	}
	log.Println(allTodo)
	attr := map[string]interface{}{
		"Checklist": allTodo,
	}
	getHTMLTemplate("index.html").ExecuteTemplate(w, "index.html", attr)
}

type tagsRequest struct {
	TagId   int    `json:"tagId"`
	Name    string `json:"name"`
	ColorId int    `json:"colorId"`
}

type todoRequest struct {
	Todo string        `json:"todo"`
	Tags []tagsRequest `json:"tags"`
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	todos, err := model.GetAllFullTodo()
	if err != nil {
		log.Printf("Error getting todolist from database: %v\n", err)
		http.Error(w, "can't get todolist", http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(todos)
	if err != nil {
		log.Printf("Error parsing data to json: %v\n", err)
		http.Error(w, "can't parsing json", http.StatusInternalServerError)
		return
	}
	setContentJsonHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func postTodo(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v\n", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	var todo todoRequest
	err = json.Unmarshal(body, &todo)
	if err != nil {
		log.Printf("Error cannot unmarshal body: %v\n", err)
		http.Error(w, "can't unmarshal", http.StatusInternalServerError)
		return
	}
	id, err := model.AddTodo(todo.Todo)
	if err != nil {
		log.Printf("Error cannot add todo to database: %v\n", err)
		http.Error(w, "can't insert todo", http.StatusInternalServerError)
		return
	}
	for i := 0; i < len(todo.Tags); i++ {
		tag := todo.Tags[i]
		err = model.InsertTodoTag(id, tag.TagId, tag.ColorId)
		if err != nil {
			model.DeleteTodoById(id)
			log.Printf("Error cannot add tag in todo %v\n", err)
			http.Error(w, "can't insert todotag", http.StatusInternalServerError)
			return
		}
	}
	setContentJsonHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
