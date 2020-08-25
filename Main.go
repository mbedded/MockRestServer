package main

import (
	"./models"
	"./services"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	mockManager := services.NewMockManager("temp.db")

	router.HandleFunc("/raw/{key}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		key := vars["key"]

		result, err := mockManager.GetMock(key)

		if err != nil {
			log.Fatalf("Error receiving data. %q", err)
		}

		fmt.Fprintf(writer, "%s", result.Content)
	})

	tmplIndex := template.Must(template.ParseFiles("templates/index.html"))
	tmplCreate := template.Must(template.ParseFiles("templates/create.html"))

	// todo: Split requests > router.HandleFunc(..).Methods("POST")
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

			result, err := mockManager.SaveMockToDatabase(data)

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
	// todo: graceful shutdown? Close DB Connection

}
