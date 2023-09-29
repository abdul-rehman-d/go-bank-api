package main

import (
	"fmt"
	"log"
)

func main() {
	var err error
	var apiServer *ApiServer

	store, err := NewPostgresStorage()
	if err != nil {
		log.Fatalln("Failed to connect to db")
	}

	err = store.Init()
	if err != nil {
		log.Fatalln("Failed to initialize db", err.Error())
	}

	apiServer, err = NewApiServer("0.0.0.0:5000", store)
	if err != nil {
		log.Fatalln("Failed to initialize API Server")
	}

	err = apiServer.Run()
	if err != nil {
		log.Fatalln("Failed to run API Server")
	}

	fmt.Println("All is good")
}
