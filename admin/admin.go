package main

import (
	"flag"
	"fmt"
	"os"
)

func parseXML(filename string) {
	xmlFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	fmt.Println("Successfully opened file")
	defer xmlFile.Close()
}

func main() {
	filename := flag.String("conf_file", "", "file path to iTunes XML")
	flag.Parse()

	parseXML(*filename)
}
