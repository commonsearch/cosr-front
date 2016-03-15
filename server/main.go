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

	log.Printf(
		"Server listening on %s:%s - You should open http://%s:%s in your browser!\n",
		Config.Host, Config.Port, GetDockerDaemonIP(), Config.Port)

	log.Fatal(http.ListenAndServe(Config.Host+":"+Config.Port, router))

}
