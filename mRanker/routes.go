package main

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
		templateFile = "index.html"
	} else {
		albums = listYearAlbums(year)
		templateFile = "year_index.html"
	}
	artists = listArtists()
	t, err := template.New(templateFile).Funcs(template.FuncMap{
		"next":     func(i Year) int { return int(i) + 1 },
		"previous": func(i Year) int { return int(i) - 1 },
	}).ParseGlob("templates/*.html")
	if err != nil {
		fmt.Println("could not parse template: %v", err.Error())
	}
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
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		fmt.Println("Could not add album: %v", err.Error())
		return
	}
	if err = upsertAlbum(name, artist, year, ranking); err != nil {
		fmt.Println("Failed to upsert album: %s", err.Error())
	}
	http.Redirect(w, r, fmt.Sprintf("http://localhost:8080/?year=%d", year), 301)
}

func deleteAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	album_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("could not delete album: %s", err.Error())
		return
	}
	if err = removeAlbum(album_id); err != nil {
		fmt.Println("failed to delete album: %s", err.Error())
	}
	http.Redirect(w, r, "http://localhost:8080", 301)
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
	router.HandleFunc("/album/{id}", deleteAlbum).Methods("DELETE")
	router.HandleFunc("/artist/", addArtist).Methods("POST")
	return router
}
