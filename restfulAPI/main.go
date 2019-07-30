package main

import (
	"log"
	"net/http"
	"restfulAPI/code"
)

func main() {

	router := code.NewRouter()

	log.Fatal(http.ListenAndServe(":3030", router))
}

