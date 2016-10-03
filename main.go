package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func listAlbums(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT albums.name as album, year, artists.name AS artist FROM albums JOIN artists ON artists.artist_id=albums.artist_id;")
	if err != nil {
		log.Fatal("Error while querying the DB: %v", err.Error())
	}
	fmt.Fprintf(w, "year |           artist |            album\n")
	for rows.Next() {
		var album string
		var year int
		var artist string
		err = rows.Scan(&album, &year, &artist)
		if err != nil {
			log.Fatal("Error while parsing row: %v", err.Error())
		}
		fmt.Fprintf(w, "%3v | %16v | %16v\n", year, artist, album)
	}
}

func addRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", listAlbums).Methods("GET")
	return router
}

func initDB() {
	pgdb, err := sql.Open("postgres", "user=gabrieldelaboulaye host=localhost dbname=mrankerdb sslmode=disable")
	if err != nil {
		log.Fatal("Could not connect to DB: %v", err.Error())
	} else {
		db = pgdb
	}
}

func main() {
	initDB()
	router := addRoutes()
	log.Fatal(http.ListenAndServe(":8080", router))
}
