package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Album struct {
	year   int
	name   string
	artist string
}

func initDB() {
	pgdb, err := sql.Open("postgres", "user=gabrieldelaboulaye host=localhost dbname=mrankerdb sslmode=disable")
	if err != nil {
		log.Fatal("Could not connect to DB: %v", err.Error())
	} else {
		db = pgdb
	}
}

func listAlbums() []Album {
	var albums []Album
	rows, err := db.Query("SELECT albums.name as album, year, artists.name AS artist FROM albums JOIN artists ON artists.artist_id=albums.artist_id ORDER BY year;")
	if err != nil {
		log.Fatal("Error while querying the DB: %v", err.Error())
	}
	for rows.Next() {
		var album string
		var year int
		var artist string
		err = rows.Scan(&album, &year, &artist)
		if err != nil {
			log.Fatal("Error while parsing row: %v", err.Error())
		}
		albums = append(albums, Album{year, album, artist})
	}
	return albums
}
func upsertAlbum(name string, artist string, year int) {
	var lastInsertedId int
	err := db.QueryRow("SELECT * FROM insert_album($1, $2, $3);", year, name, artist).Scan(&lastInsertedId)
	if err != nil {
		log.Fatal("Could not upsert album: ", err.Error())
	}
}

func insertArtist(name string) {
	var lastInsertedId int
	err := db.QueryRow("INSERT INTO artists (name) VALUES ($1) RETURNING artist_id;", name).Scan(&lastInsertedId)
	if err != nil {
		log.Fatal("Could not insert artist: ", err.Error())
	}
}
