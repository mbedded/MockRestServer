package main

import (
	"./services"
	"github.com/gorilla/mux"
	"net/http"
)

var _router = mux.NewRouter()
var _dbManager = services.NewDatabaseManager("mockRestServer.db")
var _httpHandler = services.NewHttpRequestHandler(_dbManager)

func main() {
	// Routes for REST-API
	_router.HandleFunc("/api/mock/key/{key}", _httpHandler.GetMock).Methods("GET")
	_router.HandleFunc("/api/mock/key/{key}", _httpHandler.DeleteMock).Methods("DELETE")
	_router.HandleFunc("/api/mock", _httpHandler.CreateMock).Methods("POST")
	_router.HandleFunc("/api/mock", _httpHandler.UpdateMock).Methods("PUT")
	_router.HandleFunc("/api/mock/all", _httpHandler.GetAllMocks).Methods("GET")

	fileServer := http.FileServer(http.Dir("assets/"))
	_router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fileServer)).Methods("GET")
	_router.HandleFunc("/raw/{key}", _httpHandler.GetMockContent).Methods("GET")

	_router.HandleFunc("/create", _httpHandler.ShowTemplate).Methods("GET")
	_router.HandleFunc("/showall", _httpHandler.ShowTemplate).Methods("GET")
	_router.HandleFunc("/", _httpHandler.ShowTemplate).Methods("GET")

	http.ListenAndServe(":5050", _router)
	// todo: graceful shutdown? Close DB Connection

}
