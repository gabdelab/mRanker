package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func displayAlbums(w http.ResponseWriter, r *http.Request) {
	albums := listAlbums()
	fmt.Printf("Found %d albums\n", len(albums))
	t, _ := template.ParseFiles("index.html")
	if err := t.Execute(w, &albums); err != nil {
		fmt.Println("Could not display albums: %v", err.Error())
	}
}

func addAlbum(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	artist := r.FormValue("artist")
	ranking, err := strconv.Atoi(r.FormValue("rank"))
	if err != nil {
		fmt.Println("Could not add album: %v", err.Error())
		return
	}
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		fmt.Println("Could not add album: %v", err.Error())
		return
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
