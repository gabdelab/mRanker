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
	Year    Year
}

type Album struct {
	Year       Year
	Name       string
	Artist     Artist
	Ranking    int
	AllRanking int
	ID         int
}

type Year int

type Artist struct {
	Name string
	ID   int
}

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
	rows, err := db.Query(`SELECT name, artist_id FROM artists ORDER BY name;`)
	defer rows.Close()
	if err != nil {
		fmt.Println("Error while querying the DB: %v", err.Error())
		return nil
	}
	for rows.Next() {
		var name string
		var id int
		if err = rows.Scan(&name, &id); err != nil {
			fmt.Println("Error while parsing row: %v", err.Error())
			return nil
		}
		artists = append(artists, Artist{Name: name, ID: id})
	}
	return artists
}

func queryAlbums(query string) Albums {
	var albums Albums
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		fmt.Println("Error while querying the DB: %v", err.Error())
		return nil
	}
	for rows.Next() {
		var album string
		var year int
		var artist string
		var artistID int
		var ranking int
		var allRanking int
		var id int
		err = rows.Scan(&album, &year, &artist, &artistID, &ranking, &allRanking, &id)
		if err != nil {
			fmt.Println("Error while parsing row: %v", err.Error())
			return nil
		}
		albums = append(albums, Album{Year(year), album, Artist{artist, artistID}, ranking, allRanking, id})
	}

	return albums
}

func listAlbums() Albums {
	return queryAlbums(`SELECT albums.name AS album,
                             year,
                             artists.name AS artist,
                             artists.artist_id AS artistID,
                             ranking,
                             ranking AS allRanking,
                             albums.album_id
                      FROM albums JOIN artists
                      ON artists.artist_id=albums.artist_id
                      ORDER BY ranking ASC;`)
}

func listYearAlbums(year int) Albums {
	return queryAlbums(fmt.Sprintf(`SELECT albums.name AS album,
                                         year,
                                         artists.name AS artist,
                                         artists.artist_id AS artistID,
                                         rank() OVER (ORDER BY year_ranking ASC) AS ranking,
                                         ranking as allRanking,
                                         albums.album_id
                                  FROM albums JOIN artists
                                  ON artists.artist_id=albums.artist_id
                                  WHERE year=%d;`, year))
}

func upsertAlbum(name string, artist string, year int, newRanking int, newYearRanking int) error {
	var rank int
	var yearRank int
	var id int

	// Check whether this is an insertion or a ranking update
	err := db.QueryRow(`SELECT
												 	 albums.album_id as id,
												 	 albums.ranking AS rank,
												 	 albums.year_ranking AS yearRank
                         FROM albums JOIN artists
                         ON artists.artist_id=albums.artist_id
                         WHERE artists.name=$1 AND albums.name=$2 AND year=$3;`,
		artist, name, year).Scan(&id, &rank, &yearRank)
	if err != nil {
		if err == sql.ErrNoRows {
			// No existing row, this is an insertion
			return insertAlbum(name, artist, year, newRanking, newYearRanking)
		}
		// System error
		fmt.Println("Error while querying the DB: %v", err.Error())
		return err
	}
	return updateAlbumRanking(id, rank, newRanking, yearRank, newYearRanking)
}

// Update album - if newRanking is 0, nothing is done
func updateAlbumRanking(id, rank, newRanking, yearRank, newYearRanking int) error {
	// Ranking update
	if newRanking != rank && newRanking > 0 {
		_, err := db.Exec("SELECT * FROM update_ranking($1, $2, $3)", id, rank, newRanking)
		if err != nil {
			fmt.Println("Could not update ranking: %v", err.Error())
			return err
		}
	}

	// Year ranking update
	if newYearRanking != yearRank && newYearRanking > 0 {
		_, err := db.Exec("SELECT * FROM update_year_ranking($1, $2, $3)", id, yearRank, newYearRanking)
		if err != nil {
			fmt.Println("Could not update year ranking: %v", err.Error())
			return err
		}
	}
	return nil
}

func insertAlbum(name string, artist string, year int, newRanking int, newYearRanking int) error {
	// New insertion - if no ranking is specified, it should be inserted last
	if newRanking == 0 {
		err := db.QueryRow("SELECT max(ranking) + 1 FROM albums;").Scan(&newRanking)
		if err != nil {
			fmt.Println("Could not get max ranking")
			return err
		}
	}
	_, err := db.Exec("SELECT * FROM insert_album($1, $2, $3, $4);", year, name, artist, newRanking)
	if err != nil {
		fmt.Println("Could not insert album: ", err.Error())
		return err
	}
	return nil
}

func insertArtist(name string) {
	var lastInsertedId int
	err := db.QueryRow("INSERT INTO artists (name) VALUES ($1) RETURNING artist_id;", name).Scan(&lastInsertedId)
	if err != nil {
		fmt.Println("Could not insert artist: ", err.Error())
		return
	}
	fmt.Println("successfully inserted artist !")
}
