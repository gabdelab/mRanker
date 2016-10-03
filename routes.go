package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func listAlbums(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT albums.name as album, year, artists.name AS artist FROM albums JOIN artists ON artists.artist_id=albums.artist_id ORDER BY year;")
	if err != nil {
		log.Fatal("Error while querying the DB: %v", err.Error())
	}
	fmt.Fprintf(w, "year |               artist |                album\n")
	fmt.Fprintf(w, "--------------------------------------------------\n")
	for rows.Next() {
		var album string
		var year int
		var artist string
		err = rows.Scan(&album, &year, &artist)
		if err != nil {
			log.Fatal("Error while parsing row: %v", err.Error())
		}
		fmt.Fprintf(w, "%3v | %20v | %20v\n", year, artist, album)
	}
}
func addAlbum(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	artist := r.FormValue("artist")
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		log.Fatal("Could not add album: %v", err.Error())
	}
	upsertAlbum(name, artist, year)
}

func addArtist(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	insertArtist(name)
}

func addRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", listAlbums).Methods("GET")
	router.HandleFunc("/album/", addAlbum).Methods("POST")
	router.HandleFunc("/artist/", addArtist).Methods("POST")
	return router
}
