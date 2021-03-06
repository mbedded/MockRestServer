package services

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"mbedded/mockserver/models"
	"net/http"
	"path/filepath"
	"strings"
)

type HttpRequestHandler struct {
	DatabaseManager *DatabaseManager
	Templates       map[string]*template.Template
}

func NewHttpRequestHandler(dbManager *DatabaseManager) *HttpRequestHandler {
	instance := &HttpRequestHandler{
		DatabaseManager: dbManager,
	}

	instance.Templates = initializeTemplates()

	return instance
}

func initializeTemplates() map[string]*template.Template {
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

func (handler *HttpRequestHandler) ShowTemplate(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	templateName := getTemplateForPath(path)

	handler.Templates[templateName].Execute(writer, nil)
}

func getTemplateForPath(path string) string {
	switch path {
	case "/create":
		return "create"

	case "/showall":
		return "showAll"

	default:
		return "index"
	}
}

func (handler *HttpRequestHandler) GetMock(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key := vars["key"]

	result, err := handler.DatabaseManager.GetMock(key)

	if err != nil {
		writeHttpRequest(fmt.Sprintf("No mock with key '%s' found", key), http.StatusBadRequest, writer)
		log.Panicf("Error receiving data. %q", err)
	}

	if result.Id == 0 {
		writeHttpRequest(fmt.Sprintf("No mock with key '%s' found", key), http.StatusBadRequest, writer)
		return
	}

	data, err := json.Marshal(result)

	writer.Header().Set("Content-Type", "text/plain")
	writer.Write(data)
}

func (handler *HttpRequestHandler) DeleteMock(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key := vars["key"]

	isExisting, err := handler.DatabaseManager.ContainsKey(key)

	if err != nil {
		writeHttpRequest(fmt.Sprintf("No mock with key '%s' found", key), http.StatusBadRequest, writer)
		log.Panicf("Error receiving data. %q", err)
	}

	if isExisting == false {
		writeHttpRequest(fmt.Sprintf("No mock with key '%s' found", key), http.StatusBadRequest, writer)
		return
	}

	err = handler.DatabaseManager.DeleteMock(key)
	if err != nil {
		writeHttpRequest("Error deleting the mock from database.", http.StatusBadRequest, writer)
		log.Panicf("Error deleteing mock. %q", err)
		return
	}

	writer.WriteHeader(http.StatusAccepted)
}

func (handler *HttpRequestHandler) CreateMock(writer http.ResponseWriter, request *http.Request) {
	result, err := ioutil.ReadAll(request.Body)

	var content models.JsonMockPost
	err = json.Unmarshal(result, &content)

	content.TrimFields()

	if err != nil {
		writeHttpRequest("Unable to unmarshal Json", http.StatusBadRequest, writer)
		return
	}

	if len(content.Content) == 0 {
		writeHttpRequest("Content must not be empty.", http.StatusBadRequest, writer)
		return
	}

	isExisting, err := handler.DatabaseManager.ContainsKey(content.Key)
	if isExisting {
		message := fmt.Sprintf("A mock with key  '%s' is already existing.", content.Key)
		writeHttpRequest(message, http.StatusConflict, writer)
		return
	}

	id, err := handler.DatabaseManager.SaveMockToDatabase(content.Key, content.Content)

	if err != nil {
		message := fmt.Sprintf("Error saving data to database. %q", err)
		writeHttpRequest(message, http.StatusInternalServerError, writer)
		return
	}

	content.Key = id
	data, err := json.Marshal(content)

	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusCreated)
	writer.Write(data)
}

func (handler *HttpRequestHandler) UpdateMock(writer http.ResponseWriter, request *http.Request) {
	result, err := ioutil.ReadAll(request.Body)

	var content models.JsonMockPost
	err = json.Unmarshal(result, &content)

	content.TrimFields()

	if err != nil {
		writeHttpRequest("Unable to unmarshal Json", http.StatusBadRequest, writer)
		return
	}

	if len(content.Content) == 0 || len(content.Key) == 0 {
		writeHttpRequest("Content and Key must not be empty.", http.StatusBadRequest, writer)
		return
	}

	isExisting, err := handler.DatabaseManager.ContainsKey(content.Key)
	if isExisting == false {
		writeHttpRequest(fmt.Sprintf("No mock with key '%s' existing.", content.Key), http.StatusBadRequest, writer)
		return
	}

	err = handler.DatabaseManager.UpdateMock(content.Key, content.Content)

	if err != nil {
		log.Panic("Unable to update item in database")
	}

	data, err := json.Marshal(content)

	writer.WriteHeader(http.StatusAccepted)
	writer.Header().Set("Content-Type", "text/plain")
	writer.Write(data)
}

func (handler *HttpRequestHandler) GetAllMocks(writer http.ResponseWriter, request *http.Request) {
	result, err := handler.DatabaseManager.GetAll()

	if err != nil {
		writeHttpRequest(fmt.Sprintf("Error reading all mocks from database."), http.StatusBadRequest, writer)
		log.Panicf("Error receiving data. %q", err)
	}

	data, err := json.Marshal(result)

	writer.Header().Set("Content-Type", "text/plain")
	writer.Write(data)
}

func (handler *HttpRequestHandler) GetMockContent(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key := vars["key"]

	result, err := handler.DatabaseManager.GetMock(key)

	if err != nil {
		message := fmt.Sprintf("Error receiving data. %q", err)
		log.Print(message)
		writeHttpRequest(message, http.StatusInternalServerError, writer)
		return
	}

	if result.Id <= 0 {
		http.NotFound(writer, request)
	} else {
		fmt.Fprintf(writer, "%s", result.Content)
	}
}

func writeHttpRequest(errorMessage string, statusCode int, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(statusCode)
	writer.Write([]byte(errorMessage))
}
