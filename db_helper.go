package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func initDB() {
	pgdb, err := sql.Open("postgres", "user=gabrieldelaboulaye host=localhost dbname=mrankerdb sslmode=disable")
	if err != nil {
		log.Fatal("Could not connect to DB: %v", err.Error())
	} else {
		db = pgdb
	}
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
