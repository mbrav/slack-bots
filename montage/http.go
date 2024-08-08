package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

// Download image
func (img *Image) download(wg *sync.WaitGroup, savePath string, username string, password string) {
	// Signal that this goroutine is done
	defer wg.Done()

	client := &http.Client{Timeout: 60 * time.Second}
	fullUrl, _ := url.Parse(img.Url)

	// Add image dimensions params to request
	queryURL := fullUrl.Query()
	queryURL.Set("width", strconv.Itoa(img.Width))
	queryURL.Set("height", strconv.Itoa(img.Height))
	fullUrl.RawQuery = queryURL.Encode()

	req, err := http.NewRequest("GET", fullUrl.String(), nil)
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return
	}

	req.SetBasicAuth(username, password)
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
