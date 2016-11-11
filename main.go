package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	fmt.Println("Starting server...")
	initDB()
	defer db.Close()
	router := addRoutes()
	fmt.Println("Server successfully started!")
	log.Fatal(http.ListenAndServe(":8080", router))
}
