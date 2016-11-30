package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Results struct {
	Albums  Albums
	Artists Artists
}

type Album struct {
	Year       int
	Name       string
	Artist     Artist
	Ranking    int
	AllRanking int
}

type Artist string

type Albums []Album

type Artists []Artist

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

func listArtists() Artists {
	var artists Artists
	rows, err := db.Query(`SELECT name FROM artists ORDER BY name;`)
	if err != nil {
		fmt.Println("Error while querying the DB: %v", err.Error())
		return nil
	}
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			fmt.Println("Error while parsing row: %v", err.Error())
			return nil
		}
		artists = append(artists, Artist(name))
	}
	return artists
}

func queryAlbums(query string) Albums {
	var albums Albums
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error while querying the DB: %v", err.Error())
		return nil
	}
	for rows.Next() {
		var album string
		var year int
		var artist Artist
		var ranking int
		var allRanking int
		err = rows.Scan(&album, &year, &artist, &ranking, &allRanking)
		if err != nil {
			fmt.Println("Error while parsing row: %v", err.Error())
			return nil
		}
		albums = append(albums, Album{year, album, artist, ranking, allRanking})
	}

	return albums
}

func listAlbums() Albums {
	return queryAlbums(`SELECT albums.name AS album,
                             year,
                             artists.name AS artist,
                             ranking,
                             ranking AS allRanking
                      FROM albums JOIN artists
                      ON artists.artist_id=albums.artist_id
                      ORDER BY ranking ASC;`)
}

func listYearAlbums(year int) Albums {
	return queryAlbums(fmt.Sprintf(`SELECT albums.name AS album,
                                         year,
                                         artists.name AS artist,
                                         rank() OVER (ORDER BY ranking ASC) AS ranking,
                                         ranking as allRanking
                                  FROM albums JOIN artists
                                  ON artists.artist_id=albums.artist_id
                                  WHERE year=%d;`, year))
}

func upsertAlbum(name string, artist string, year int, newRanking int) {

	// Check whether this is an insertion or a ranking update
	rows, err := db.Query(`SELECT albums.album_id as id, albums.ranking AS rank
                         FROM albums JOIN artists
                         ON artists.artist_id=albums.artist_id
                         WHERE artists.name=$1 AND albums.name=$2 AND year=$3;`, artist, name, year)
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
