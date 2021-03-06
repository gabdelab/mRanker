package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
	"howett.net/plist"
)

const (
	HOST = "localhost"
	PORT = "8080"
)

type Song struct {
	Year        int    `plist:"Year"`
	Name        string `plist:"Name"`
	Album       string `plist:"Album"`
	PlayCount   int    `plist:"PlayCount"`
	Artist      string `plist:"Artist"`
	Compilation bool   `plist:"Compilation"`
}

type File struct {
	Tracks map[string]Song `plist:"Tracks"`
}

type BlacklistedAlbum string

// call mRanker server to insert artist
func insertArtist(song Song) {
	postURL := fmt.Sprintf("http://%s:%s/artist/", HOST, PORT)
	form := url.Values{}
	form.Add("name", song.Artist)
	req, err := http.NewRequest("POST", postURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("can't add artist: %v", err)
		return
	}
	defer resp.Body.Close()
}

// call mRanker server to insert album
func insertAlbum(song Song) {
	if song.Compilation {
		fmt.Println("not adding compilation", song.Album)
		// Don't insert compilations
		return
	}
	// Is the song blacklisted ?
	postURL := fmt.Sprintf("http://%s:%s/album/", HOST, PORT)
	form := url.Values{}
	form.Add("name", song.Album)
	form.Add("artist", song.Artist)
	form.Add("year", strconv.Itoa(song.Year))
	form.Add("rank", "0")
	req, err := http.NewRequest("POST", postURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("can't add artist: %v", err)
		return
	}
	defer resp.Body.Close()
}

func parseXML(filename string) error {
	xmlFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	fmt.Println("Successfully opened file")
	defer xmlFile.Close()

	b, _ := ioutil.ReadAll(xmlFile)
	r := bytes.NewReader(b)
	var file File

	decoder := plist.NewDecoder(r)
	decoder.Decode(&file)

	var blacklisted []string
	yamlFile, err := ioutil.ReadFile("blacklist.yml")
	if err != nil {
		log.Fatal("can't parse blacklisted file !")
	}
	err = yaml.Unmarshal([]byte(yamlFile), &blacklisted)
	if err != nil {
		log.Fatal("can't parse blacklisted file !")
	}

	for _, song := range file.Tracks {
		isBlacklisted := false
		for _, blacklistedAlbum := range blacklisted {
			if ok, _ := regexp.MatchString(blacklistedAlbum, song.Album); ok {
				fmt.Println("blacklisting", song.Album)
				isBlacklisted = true
				continue
			}
		}
		if !isBlacklisted {
			insertArtist(song)
			insertAlbum(song)
			fmt.Printf("%d - %s - %s - %s\n", song.Year, song.Name, song.Album, song.Artist)
		}
	}
	return nil
}

func main() {
	filename := flag.String("itunes_xml", os.Getenv("ITUNES_XML_FILE"), "file path to iTunes XML")
	flag.Parse()

	if err := parseXML(*filename); err != nil {
		log.Fatal("could not parse configuration file")
	}
}
