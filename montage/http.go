package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Downloads the image and saves it to the specified path.
func (img *Image) download(savePath string, username string, password string) error {
	client := &http.Client{Timeout: 60 * time.Second}
	fullUrl, err := url.Parse(img.Url)
	if err != nil {
		return fmt.Errorf("error parsing URL: %v", err)
	}

	queryURL := fullUrl.Query()
	queryURL.Set("width", strconv.Itoa(img.Width))
	queryURL.Set("height", strconv.Itoa(img.Height))
	fullUrl.RawQuery = queryURL.Encode()

	req, err := http.NewRequest("GET", fullUrl.String(), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.SetBasicAuth(username, password)
	fmt.Printf("Downloading url %s to %s\n", fullUrl.String(), savePath)
	r, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing request: %v", err)
	}
	defer r.Body.Close()

	fmt.Printf("Saving url %s to %s\n", fullUrl.String(), savePath)
	f, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, r.Body); err != nil {
		return fmt.Errorf("error saving file: %v", err)
	}

	return nil
}
