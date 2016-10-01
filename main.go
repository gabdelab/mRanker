package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func addRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", helloWorld).Methods("GET")
  return router
}

func main() {
	router := addRoutes()
	log.Fatal(http.ListenAndServe(":8080", router))
}
