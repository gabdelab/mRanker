package main

import (
	"database/sql"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	initDB()
	router := addRoutes()
	log.Fatal(http.ListenAndServe(":8080", router))
}
