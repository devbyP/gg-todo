package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
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

type tag string

// use as the second argument of the only one HandleFunc in sever.
//
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
		case http.MethodPost:
			fmt.Println(request.FormValue("message"))
			querystring := "?message=test"
			http.Redirect(writer, request, "/"+querystring, http.StatusSeeOther)
		}

	case "/tag":
		if request.Method == http.MethodPost {
			//tag := request.Body
			// write to database.
			result, err := json.Marshal(struct{ ResultOk bool }{ResultOk: true})
			if err != nil {
				http.Error(writer, "cannot parse result to json.", http.StatusInternalServerError)
			}
			setContentJsonHeader(writer)
			writer.WriteHeader(http.StatusOK)
			writer.Write(result)
		}
		writeNotFound(writer)
	default:
		writeNotFound(writer)
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
