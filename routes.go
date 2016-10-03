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
	fmt.Fprintf(w, "%3v | %20v | %20v\n", "year", "artist", "album")
	fmt.Fprintf(w, "--------------------------------------------------\n")

	for _, i := range albums {
		fmt.Fprintf(w, "%3v | %20v | %20v\n", i.year, i.artist, i.name)
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
	router.HandleFunc("/", displayAlbums).Methods("GET")
	router.HandleFunc("/album/", addAlbum).Methods("POST")
	router.HandleFunc("/artist/", addArtist).Methods("POST")
	return router
}
