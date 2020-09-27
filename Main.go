package main

import (
	"./models"
	"./services"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var _templates map[string]*template.Template
var _router = mux.NewRouter()
var _dbManager = services.NewDatabaseManager("mockRestServer.db")

func main() {
	_templates = InitializeTemplates()

	// Routes for REST-API
	_router.HandleFunc("/api/mock/key/{key}", getMock).Methods("GET")
	_router.HandleFunc("/api/mock/key/{key}", deleteMock).Methods("DELETE")
	_router.HandleFunc("/api/mock", createMock).Methods("POST")
	_router.HandleFunc("/api/mock", updateMock).Methods("PUT")
	_router.HandleFunc("/api/mock/all", getAllMocks).Methods("GET")

	fileServer := http.FileServer(http.Dir("assets/"))
	_router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fileServer))

	_router.HandleFunc("/raw/{key}", getMockContent).Methods("GET")

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
func getMockContent(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key := vars["key"]

	result, err := _dbManager.GetMock(key)

	if err != nil {
		log.Fatalf("Error receiving data. %q", err)
	}

	if result.Id <= 0 {
		http.NotFound(writer, request)
	} else {
		fmt.Fprintf(writer, "%s", result.Content)
	}
}

func createMock(writer http.ResponseWriter, request *http.Request) {
	result, err := ioutil.ReadAll(request.Body)

	var content models.JsonMockPost
	err = json.Unmarshal(result, &content)

	content.TrimFields()

	if err != nil {
		writeBadRequest("Unable to unmarshal Json", http.StatusBadRequest, writer)
		return
	}

	if len(content.Content) == 0 {
		writeBadRequest("Content must not be empty.", http.StatusBadRequest, writer)
		return
	}

	id, err := _dbManager.SaveMockToDatabase(content.Key, content.Content)

	if err != nil {
		writeBadRequest(fmt.Sprintf("Error saving data to database. %q", err),
			http.StatusInternalServerError, writer)
		return
	}

	content.Key = id
	data, err := json.Marshal(content)

	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusCreated)
	writer.Write(data)
}

func getMock(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key := vars["key"]

	result, err := _dbManager.GetMock(key)

	if err != nil {
		writeBadRequest(fmt.Sprintf("No mock with key '%s' found", key), http.StatusBadRequest, writer)
		log.Panicf("Error receiving data. %q", err)
	}

	if result.Id == 0 {
		writeBadRequest(fmt.Sprintf("No mock with key '%s' found", key), http.StatusBadRequest, writer)
		return
	}

	data, err := json.Marshal(result)

	writer.Header().Set("Content-Type", "text/plain")
	writer.Write(data)
}

func getAllMocks(writer http.ResponseWriter, request *http.Request) {
	result, err := _dbManager.GetAll()

	if err != nil {
		writeBadRequest(fmt.Sprintf("Error reading all mocks from database."), http.StatusBadRequest, writer)
		log.Panicf("Error receiving data. %q", err)
	}

	data, err := json.Marshal(result)

	writer.Header().Set("Content-Type", "text/plain")
	writer.Write(data)
}

func updateMock(writer http.ResponseWriter, request *http.Request) {
	result, err := ioutil.ReadAll(request.Body)

	var content models.JsonMockPost
	err = json.Unmarshal(result, &content)

	content.TrimFields()

	if err != nil {
		writeBadRequest("Unable to unmarshal Json", http.StatusBadRequest, writer)
		return
	}

	if len(content.Content) == 0 || len(content.Key) == 0 {
		writeBadRequest("Content and Key must not be empty.", http.StatusBadRequest, writer)
		return
	}

	isExisting, err := _dbManager.ContainsKey(content.Key)
	if isExisting == false {
		writeBadRequest(fmt.Sprintf("No mock with key '%s' existing.", content.Key), http.StatusBadRequest, writer)
		return
	}

	err = _dbManager.UpdateMock(content.Key, content.Content)

	if err != nil {
		log.Panic("Unable to update item in database")
	}

	data, err := json.Marshal(content)

	writer.WriteHeader(http.StatusAccepted)
	writer.Header().Set("Content-Type", "text/plain")
	writer.Write(data)
}

func deleteMock(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key := vars["key"]

	isExisting, err := _dbManager.ContainsKey(key)

	if err != nil {
		writeBadRequest(fmt.Sprintf("No mock with key '%s' found", key), http.StatusBadRequest, writer)
		log.Panicf("Error receiving data. %q", err)
	}

	if isExisting == false {
		writeBadRequest(fmt.Sprintf("No mock with key '%s' found", key), http.StatusBadRequest, writer)
		return
	}

	err = _dbManager.DeleteMock(key)
	if err != nil {
		writeBadRequest("Error deleting the mock from database.", http.StatusBadRequest, writer)
		log.Panicf("Error deleteing mock. %q", err)
		return
	}

	writer.WriteHeader(http.StatusAccepted)
}

func writeBadRequest(errorMessage string, statusCode int, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(statusCode)
	writer.Write([]byte(errorMessage))
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
