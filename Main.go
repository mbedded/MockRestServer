package main

import (
	"./services"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var _templates map[string]*template.Template
var _router = mux.NewRouter()
var _dbManager = services.NewDatabaseManager("mockRestServer.db")
var _httpHandler = services.NewHttpRequestHandler(_dbManager)

func main() {
	_templates = InitializeTemplates()

	// Routes for REST-API
	_router.HandleFunc("/api/mock/key/{key}", _httpHandler.GetMock).Methods("GET")
	_router.HandleFunc("/api/mock/key/{key}", _httpHandler.DeleteMock).Methods("DELETE")
	_router.HandleFunc("/api/mock", _httpHandler.CreateMock).Methods("POST")
	_router.HandleFunc("/api/mock", _httpHandler.UpdateMock).Methods("PUT")
	_router.HandleFunc("/api/mock/all", _httpHandler.GetAllMocks).Methods("GET")

	fileServer := http.FileServer(http.Dir("assets/"))
	_router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fileServer))

	_router.HandleFunc("/raw/{key}", _httpHandler.GetMockContent).Methods("GET")

	_router.HandleFunc("/create", func(writer http.ResponseWriter, request *http.Request) {
		_templates["create"].Execute(writer, nil)
	}).Methods("GET")

	_router.HandleFunc("/showall", func(writer http.ResponseWriter, request *http.Request) {
		_templates["showAll"].Execute(writer, nil)
	}).Methods("GET")

	_router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_templates["index"].Execute(writer, nil)
	})

	http.ListenAndServe(":5050", _router)
	// todo: graceful shutdown? Close DB Connection
}

func InitializeTemplates() map[string]*template.Template {
	const extension = ".html"
	templates := make(map[string]*template.Template)

	content, err := filepath.Glob("templates/*" + extension)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range content {
		filename := filepath.Base(file)
		filename = strings.Replace(filename, extension, "", 1)

		templates[filename] = template.Must(template.ParseFiles(file))
	}

	return templates
}
