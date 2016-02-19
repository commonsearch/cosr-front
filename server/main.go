package main

import (
	"log"
	"net/http"
)

// SetupGlobals performs global initialization tasks at startup.
func SetupGlobals() {

	LoadConfig()
	LoadBangs()
	LoadTemplates()

	ElasticsearchConnect()

}

// main is the entry point of the server.
func main() {

	SetupGlobals()

	router := CreateRouter()

	log.Printf("Server listening at http://%s:%s\n", Config.Host, Config.Port)

	log.Fatal(http.ListenAndServe(Config.Host+":"+Config.Port, router))

}
