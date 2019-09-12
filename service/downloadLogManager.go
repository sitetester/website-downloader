package service

import (
	"bufio"
	"log"
	"os"
)

type DownloadLogManager struct{}

const LogFilename = "parsedLinks.log"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (downloadLogManager DownloadLogManager) readUrlsFromFile() []string {
	var existingLinks []string

	file, _ := os.Open(LogFilename)
	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(file)
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		existingLinks = append(existingLinks, scanner.Text())
	}

	return existingLinks
}

func (downloadLogManager DownloadLogManager) appendUrlToFile(url string) {
	file, err := os.OpenFile(LogFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	defer file.Close()

	// fmt.Println("Writing URL ... " + url)
	if _, err := file.WriteString(url + "\n"); err != nil {
		log.Println(err)
	}
}
