package main

import (
	"fmt"
	"log"
	"net/http"
)

func DefaultPathHandler(http.ResponseWriter, *http.Request) {

	fmt.Println("This is the Default Path Handler")
	log.Println("Entered the default Path Handler")

}

func registerRouteHandlers() {

	http.HandleFunc("/", DefaultPathHandler)
	log.Fatal(http.ListenAndServe((":12345"), nil)) //nil means the default Router Server is used

}

func main() {

	registerRouteHandlers()
	log.Println("Registered Route Handlers")

}
