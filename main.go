package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"web-copier/service"
)

func getFileName(path string) string {
	url, err := url.Parse(path)
	if err != nil {
		panic(err)
	}

	/*file := filepath.Base(path)
	fmt.Println(file)*/

	fmt.Println(filepath.Dir(path))


	return url.Path
}


func main() {

	/*getFileName("/test/index.html")
	os.Exit(1)*/

	parsedLinks := make([]string, 0)

	const URL = "https://www.php.net"
	var downloader service.Downloader
	downloader.Start(URL, "", parsedLinks)
}
