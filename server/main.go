package mRanker

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
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))
	fmt.Println("Server successfully started!")
	log.Fatal(http.ListenAndServe(":8080", router))
}
