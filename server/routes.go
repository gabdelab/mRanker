package mRanker

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func displayAlbums(w http.ResponseWriter, r *http.Request) {
	var albums Albums
	var artists Artists
	var templateFile string
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		albums = listAlbums()
		templateFile = "templates/index.html"
	} else {
		albums = listYearAlbums(year)
		templateFile = "templates/year_index.html"
	}
	artists = listArtists()
	t, _ := template.ParseFiles(templateFile)

	results := Results{Artists: artists, Albums: albums, Year: Year(year)}
	if err := t.Execute(w, &results); err != nil {
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
	yearRanking, err := strconv.Atoi(r.FormValue("year_rank"))
	if err != nil {
		fmt.Println("Could not add album: %v", err.Error())
		return
	}
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		fmt.Println("Could not add album: %v", err.Error())
		return
	}
	upsertAlbum(name, artist, year, ranking, yearRanking)
	http.Redirect(w, r, "http://localhost:8080/", 301)
}

func addArtist(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	insertArtist(name)
	http.Redirect(w, r, "http://localhost:8080/", 301)
}

func addRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", displayAlbums).Methods("GET")
	router.HandleFunc("/album/", addAlbum).Methods("POST")
	router.HandleFunc("/artist/", addArtist).Methods("POST")
	return router
}
