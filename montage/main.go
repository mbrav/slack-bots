package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/gographics/imagick.v3/imagick"
)

var wg sync.WaitGroup

func main() {
	cliConfig := getCLIArgs(&cliConf)
	appConfig := getAppConfig(cliConf.AppConfig)

	if cliConf.Verbose {
		// Print the configurations for debugging
		fmt.Printf("CLI Config: %+v\n", cliConfig)
		fmt.Printf("App Config: %+v\n", appConfig)
	}

	// Get the absolute path of the directory
	dirPath, err := filepath.Abs("./test/")
	if err != nil {
		log.Fatal(err)
	}

	for i, img := range appConfig.Images {
		wg.Add(1)
		go img.download(
			&wg,
			filepath.Join(dirPath, fmt.Sprintf("img-%d.png", i)),
			appConfig.GrafanaUser,
			appConfig.GrafanaPassword)
	}

	// Wait for downlod go routines to complete
	wg.Wait()
	fmt.Println("All downloads completed.")

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	images, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range images {
		fullPath := filepath.Join(dirPath, f.Name())
		fmt.Println(fullPath)

		aw := imagick.NewMagickWand()
		defer aw.Destroy()

		err = aw.ReadImage(fullPath)
		if err != nil {
			log.Fatalf("Failed to read image %s: %v", fullPath, err)
		}

		err = mw.AddImage(aw)
		if err != nil {
			log.Fatalf("Failed to add image %s to montage: %v", fullPath, err)
		}
		mw.SetLastIterator()
	}

	// Configure montage settings
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()

	// Set the background color before creating the montage
	bgColor := imagick.NewPixelWand()
	bgColor.SetColor(appConfig.MontageBackgroudColor)
	mw.SetBackgroundColor(bgColor)
	montage := mw.MontageImage(
		dw,
		appConfig.MontageTile,
		"+0+0",
		imagick.MONTAGE_MODE_FRAME,
		"0x0+0+0")

	outFile, err := filepath.Abs(cliConfig.OutputFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Saving Montage to: %s\n", outFile)
	err = montage.WriteImage(outFile)
	if err != nil {
		log.Fatalf("Failed to save montage image: %v", err)
	}

	fmt.Println("Montage Done")
}
