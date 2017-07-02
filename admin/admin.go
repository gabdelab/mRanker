package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"howett.net/plist"
)

type Song struct {
	Year      int    `plist:"Year"`
	Name      string `plist:"Name"`
	Album     string `plist:"Album"`
	PlayCount int    `plist:"PlayCount"`
	Artist    string `plist:"Artist"`
}

type File struct {
	Tracks map[string]Song `plist:"Tracks"`
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

	counter := 0
	for _, song := range file.Tracks {
		if counter < 10 {
			fmt.Printf("%s - %s - %s\n", song.Name, song.Album, song.Artist)
		}
		counter++
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
