package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/gographics/imagick.v3/imagick"
)

// Creates a new direcotry for downloading images to.
// Downloads images.
// Then creates a montage out of those images usign imagemagick.
func createGrafanaMontage(cliConfig *CLIConfig, appConfig *AppConfig) {
	dirPath := initDir(cliConfig.DownloadDir)

	downloadImages(*appConfig, dirPath)

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	loadImagesFromDir(mw, dirPath)

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

	saveMontage(montage, outFile, uint(appConfig.MontageQuality))
}

// Initialize directory at specified location.
func initDir(path string) string {
	dirPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Error when getting absolute path for %s: %v", path, err)
	}
	// Create the directory if it does not exist
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Fatalf("Error creating directory %s: %v", dirPath, err)
	}
	return dirPath
}

// Reads images from the specified directory.
func readImagesFromDir(dirPath string) ([]os.DirEntry, error) {
	images, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read images from directory: %v", err)
	}
	return images, nil
}

// Load images from directory into wand
func loadImagesFromDir(mw *imagick.MagickWand, dirPath string) {
	images, err := readImagesFromDir(dirPath)
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
}

// Saves the final montage to the specified output file.
func saveMontage(mw *imagick.MagickWand, outFile string, quality uint) {
	// Set the image quality (compression quality) from 0-100
	mw.SetImageCompressionQuality(quality)

	fmt.Printf("Saving Montage to: %s\n", outFile)
	if err := mw.WriteImage(outFile); err != nil {
		log.Fatalf("Failed to save montage image: %v", err)
	}
	fmt.Println("Montage Done")
}
