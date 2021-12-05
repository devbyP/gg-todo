package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"gghub.com/model"
)

const (
	addr = "localhost:8000"
)

func main() {
	// target static location
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", baseHandler)
	fmt.Println("listen...")
	log.Fatal(http.ListenAndServe(addr, nil))
}

// use as the second argument of the only one HandleFunc in sever.
func baseHandler(writer http.ResponseWriter, request *http.Request) {
	// log all the incoming request.
	logRequest(request.Method, request.Host, request.URL.Path)

	// switch path(route) ->
	//	 switch method in each path ->
	// 	   control logic(business logic)
	switch request.URL.Path {
	// root route
	case "/":
		switch request.Method {
		case http.MethodGet:
			homePage(writer, request)
		default:
			writeNotFound(writer)
		}
	case "/todo":
		todoRoute(writer, request)
	case "/tags":
		tagsRoute(writer, request)
	case "/highlight":
		highlightRoute(writer, request)
	default:
		writeNotFound(writer)
	}
}

func todoRoute(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postTodo(w, r)
	}
}

func tagsRoute(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		setContentJsonHeader(w)
		tags, err := model.GetNumOfTag(10)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			w.Write([]byte("error getting tags."))
			return
		}
		result, err := json.Marshal(tags)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error cannot parse json."))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	case http.MethodPost:
		setContentJsonHeader(w)
		w.WriteHeader(http.StatusNoContent)
	default:
		writeNotFound(w)
	}
}

func highlightRoute(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		setContentJsonHeader(w)
		colors, err := model.GetAllHighlight()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			w.Write([]byte("server error cannot get highlight color."))
			return
		}
		result, err := json.Marshal(colors)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error cannot parse json."))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	default:
		writeNotFound(w)
	}
}

func setContentJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func writeNotFound(w http.ResponseWriter) {
	http.Error(w, "error 404 not found.", http.StatusNotFound)
}

func logRequest(method, host, path string) {
	log.Printf("[%s] %s%s", method, host, path)
}

func getHTMLTemplate(filename string) *template.Template {
	path := "template"
	fileLocation := path + "/" + filename
	return template.Must(template.ParseFiles(fileLocation))
}
