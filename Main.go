package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"mbedded/mockserver/models"
	"mbedded/mockserver/services"
	"net/http"
	"os"
)

var version = "DEVELOPMENT"
var _args models.CommandArgs

func main() {
	InitializeFlags()
	ShowVersionAndExitIfVersionArgIsUsed()

	var router = mux.NewRouter()
	var dbManager = services.NewDatabaseManager(_args.DatabaseFile)
	var httpHandler = services.NewHttpRequestHandler(dbManager)

	// Routes for REST-API
	router.HandleFunc("/api/mock/key/{key}", httpHandler.GetMock).Methods("GET")
	router.HandleFunc("/api/mock/key/{key}", httpHandler.DeleteMock).Methods("DELETE")
	router.HandleFunc("/api/mock", httpHandler.CreateMock).Methods("POST")
	router.HandleFunc("/api/mock", httpHandler.UpdateMock).Methods("PUT")
	router.HandleFunc("/api/mock/all", httpHandler.GetAllMocks).Methods("GET")

	fileServer := http.FileServer(http.Dir("assets/"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fileServer)).Methods("GET")
	router.HandleFunc("/raw/{key}", httpHandler.GetMockContent).Methods("GET")

	router.HandleFunc("/create", httpHandler.ShowTemplate).Methods("GET")
	router.HandleFunc("/showall", httpHandler.ShowTemplate).Methods("GET")
	router.HandleFunc("/", httpHandler.ShowTemplate).Methods("GET")

	log.Printf("Webserver will be startet at http://localhost:%d", _args.HttpPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", _args.HttpPort), router)
	if err != nil {
		log.Fatal(err)
	}

}

func InitializeFlags() {
	flag.IntVar(&_args.HttpPort, "port", 8080, "HTTP-Port of the webserver")
	flag.StringVar(&_args.DatabaseFile, "database", "mockRestServer.db", "Name of Sqlite-file")
	flag.BoolVar(&_args.ShowVersion, "version", false, "Displays the version of this tool")
	flag.Parse()
}

func ShowVersionAndExitIfVersionArgIsUsed() {
	println("MockRestServer - Version: " + version)

	if _args.ShowVersion {
		os.Exit(0)
	}
}
