package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
)

func downloadImage(wg *sync.WaitGroup, url_str string, savePath string, username string, password string) {
	defer wg.Done() // Signal that this goroutine is done

	client := &http.Client{}

	fullUrl, _ := url.Parse(url_str)

	// Add image dimensions params to request
	queryURL := fullUrl.Query()
	queryURL.Set("width", "2000")
	queryURL.Set("height", "1000")
	fullUrl.RawQuery = queryURL.Encode()

	req, err := http.NewRequest("GET", fullUrl.String(), nil)
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return
	}

	basicAuth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	req.Header = http.Header{
		"Host":          {"www.host.com"},
		"Authorization": {"Basic " + basicAuth},
	}

	fmt.Printf("Downloading url %s to %s \n", fullUrl.String(), savePath)
	r, err := client.Do(req)
	if err != nil {
		log.Printf("Error performing request: %v\n", err)
		return
	}
	defer r.Body.Close()

	fmt.Printf("Saving url %s to %s \n", fullUrl.String(), savePath)
	f, err := os.Create(savePath)
	if err != nil {
		log.Printf("Error creating file: %v\n", err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, r.Body)
	if err != nil {
		log.Printf("Error saving file: %v\n", err)
	}
}
