package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

func DefaultPathHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("This is the Default Path Handler")
	log.Println("Entered the default Path Handler")

	//Convert io.ReadCloser to String

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	requestString := buf.String()

	log.Println("Request Body: \n" + requestString)

}

func registerRouteHandlers() {

	log.Println("Registering Route Handlers")

	http.HandleFunc("/", DefaultPathHandler)

	var endpoint_port string
	// If port is set an the environment variables, use that
	if endpoint_port = os.Getenv("PORT"); endpoint_port == "" {

		// Default endpoint is 12345
		endpoint_port = "12345"
		log.Println("setting port")

	}

	log.Println("Listening on port " + endpoint_port)
	log.Fatal(http.ListenAndServe(":"+endpoint_port, nil)) //nil means the default Router Server is used

}

func main() {

	log.Println("V2 Test Bot Started")
	registerRouteHandlers()
	log.Println("Registered Route Handlers")

}
