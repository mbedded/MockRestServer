package main

import (
	"./models"
	"./services"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	manager := services.NewMockManager("temp")

	router.HandleFunc("/raw/{id}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		id := vars["id"]

		// Get text for key

		fmt.Fprintf(writer, "{\n\t\"id\": \"%s\"\n}", id)
	})

	tmplIndex := template.Must(template.ParseFiles("templates/index.html"))
	tmplCreate := template.Must(template.ParseFiles("templates/create.html"))

	router.HandleFunc("/create", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet {
			tmplCreate.Execute(writer, nil)
			return
		}

		if request.Method == http.MethodPost {
			//todo:  save data and create mock
			// -> verify input. Show error if content is empty
			// -> Redirect PRG Pattern

			var data = models.JsonMockPost{
				Key:     request.FormValue("key"),
				Content: request.FormValue("content"),
			}

			result, err := manager.CreateMock(data)

			if err != nil {
				// todo display error message to user
				return
			}

			//fmt.Printf("Key: %s | content: %s", data.Key, data.Content)
			fmt.Printf("result: %s ", result)

			return
		}

		fmt.Fprintf(writer, "Method '%s' not allowed", request.Method)
	})

	//router.Handle("/assets/", http.StripPrefix("/assets/", fileServer))
	fileServer := http.FileServer(http.Dir("assets/"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fileServer))

	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		tmplIndex.Execute(writer, nil)
	})

	http.ListenAndServe(":5050", router)

}
