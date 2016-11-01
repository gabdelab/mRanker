package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func displayAlbums(w http.ResponseWriter, r *http.Request) {
	albums := listAlbums()
	fmt.Fprintf(w, "%4v | %3v | %20v | %25v\n", "rank", "year", "artist", "album")
	fmt.Fprintf(w, "--------------------------------------------------------------\n")

	for _, i := range albums {
		fmt.Fprintf(w, "%4v | %3v | %20v | %25v\n", i.ranking, i.year, i.artist, i.name)
	}
}

func addAlbum(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	artist := r.FormValue("artist")
	ranking, err := strconv.Atoi(r.FormValue("rank"))
	if err != nil {
		log.Fatal("Could not add album: %v", err.Error())
	}
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		log.Fatal("Could not add album: %v", err.Error())
	}
	upsertAlbum(name, artist, year, ranking)
}

func addArtist(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	insertArtist(name)
}

func addRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", displayAlbums).Methods("GET")
	router.HandleFunc("/album/", addAlbum).Methods("POST")
	router.HandleFunc("/artist/", addArtist).Methods("POST")
	return router
}
