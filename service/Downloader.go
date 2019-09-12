package service

import (
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Downloader struct{}

type ParsedLinks struct {
	ParsedLinks []string
}

func (downloader Downloader) Start(baseUrl string, filePath string, parsedLinks []string) {
	var resp *http.Response

	filename := getFileName(filePath)
	fmt.Println("filename ------ " + filename)

	if filename == "/index.html" {
		parsedLinks = append(parsedLinks, filename)
		resp = download(baseUrl)
	} else {
		parsedLinks = append(parsedLinks, filePath)
		resp = download(baseUrl + filePath)
	}

	stringData, _ := ioutil.ReadAll(resp.Body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(stringData)))
	if err != nil {
		log.Fatal(err)
	}

	writeFile(strings.NewReader(string(stringData)), filename)

	fmt.Println("parsedLinks ----")
	fmt.Println(parsedLinks)

	bodyLinks := parseAllBodyLinks(*doc)
	/*fmt.Println("bodyLinks ---")
	fmt.Println(bodyLinks)*/

	var newLinks []string
	newLinks = findOnlyNewLinks(parsedLinks, bodyLinks)

	for _, newLink := range newLinks {
		//if newLink == "/privacy.php" || newLink == "/sites.php" || newLink == "/contact.php" {
			// fmt.Println(newLink.Link)
			go downloader.Start(baseUrl, newLink, parsedLinks)
			time.Sleep(time.Second * 5)
		//}
	}

	fmt.Println("All done!")
}

func getFileName(siteUrl string) string {
	url, err := url.Parse(siteUrl)
	if err != nil {
		panic(err)
	}

	if url.Path == "" {
		return "/index.html"
	}

	return url.Path
}

func writeFile(reader io.Reader, filename string) {
	newPath := filepath.Join("./", "downloads/"+filepath.Dir(filename))
	os.MkdirAll(newPath, os.ModePerm)

	body, _ := ioutil.ReadAll(reader)
	ioutil.WriteFile(newPath+"/"+filepath.Base(filename), body, 0644)
}

func parseAllBodyLinks(doc goquery.Document) []string {
	var bodyLinks []string

	doc.Find("body a").Each(func(index int, a *goquery.Selection) {
		href, exists := a.Attr("href")
		if exists {
			if href[0:1] == "/" {
				bodyLinks = append(bodyLinks, href)
			}
		}
	})

	return bodyLinks
}

func getIgnoredLinks(link string) []string {
	var ignoredLinks []string

	ignoredLinks = append(ignoredLinks, "")
	ignoredLinks = append(ignoredLinks, "/")
	ignoredLinks = append(ignoredLinks, "javascript:;")

	return ignoredLinks
}

func findOnlyNewLinks(parsedLinks []string, bodyLinks []string) []string {
	var newLinks []string

	for _, bodyLink := range bodyLinks {
		if !contains(parsedLinks, bodyLink) {
			newLinks = append(newLinks, bodyLink)
		}
	}

	fmt.Println("newLinks ----------------")
	fmt.Println(newLinks)
	// os.Exit(1)

	return newLinks
}

// contains tells whether a contains x.
func contains(a []string, findMe string) bool {
	for _, val := range a {
		if findMe == val {
			return true
		}
	}

	return false
}

func download(url string) *http.Response {
	fmt.Println("Downloading...", url)

	// create New http Transport
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // disable verify
		},
	}

	// create Http Client
	timeout := time.Duration(300 * time.Second)
	client := &http.Client{
		Transport: transCfg,
		Timeout:   timeout,
	}

	// request
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}

	return resp
}
