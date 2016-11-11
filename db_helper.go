package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Album struct {
	year    int
	name    string
	artist  string
	ranking int
}

func initDB() {
	pgdb, err := sql.Open("postgres", "user=gabrieldelaboulaye host=localhost dbname=mrankerdb sslmode=disable")
	if err != nil {
		log.Fatal("Could not connect to DB: %v", err.Error())
	}
	db = pgdb

	if err = db.Ping(); err != nil {
		log.Fatal("Can't ping DB, is the PG server up ?")
	}
}

func closeDB() {
	db.Close()
}

func listAlbums() []Album {
	var albums []Album
	rows, err := db.Query("SELECT albums.name as album, year, artists.name AS artist, ranking FROM albums JOIN artists ON artists.artist_id=albums.artist_id ORDER BY ranking ASC;")
	if err != nil {
		fmt.Println("Error while querying the DB: %v", err.Error())
		return nil
	}
	for rows.Next() {
		var album string
		var year int
		var artist string
		var ranking int
		err = rows.Scan(&album, &year, &artist, &ranking)
		if err != nil {
			fmt.Println("Error while parsing row: %v", err.Error())
			return nil
		}
		albums = append(albums, Album{year, album, artist, ranking})
	}
	return albums
}

func upsertAlbum(name string, artist string, year int, newRanking int) {

	// Check whether this is an insertion or a ranking update
	rows, err := db.Query("SELECT albums.album_id as id, albums.ranking as rank FROM albums JOIN artists ON artists.artist_id=albums.artist_id WHERE artists.name=$1 AND albums.name=$2 and year=$3;", artist, name, year)
	if err != nil {
		fmt.Println("Error while querying the DB: %v", err.Error())
		return
	}
	for rows.Next() {
		var rank int
		var id int
		err = rows.Scan(&id, &rank)
		if err != nil {
			fmt.Println("Error while parsing row: %v", err.Error())
			return
		}

		// Ranking update
		if newRanking != rank {
			_, err := db.Query("SELECT * FROM update_ranking($1, $2, $3)", id, rank, newRanking)
			if err != nil {
				fmt.Println("Could not update ranking: %v", err.Error())
				return
			}
		}
		return
	}
	// New insertion
	_, err = db.Query("SELECT * FROM insert_album($1, $2, $3, $4);", year, name, artist, newRanking)
	if err != nil {
		fmt.Println("Could not insert album: ", err.Error())
		return
	}
}

func insertArtist(name string) {
	var lastInsertedId int
	err := db.QueryRow("INSERT INTO artists (name) VALUES ($1) RETURNING artist_id;", name).Scan(&lastInsertedId)
	if err != nil {
		fmt.Println("Could not insert artist: ", err.Error())
		return
	}
}
