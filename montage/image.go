package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/gographics/imagick.v3/imagick"
)

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

// saveMontage saves the final montage to the specified output file.
func saveMontage(mw *imagick.MagickWand, outFile string) {
	fmt.Printf("Saving Montage to: %s\n", outFile)
	if err := mw.WriteImage(outFile); err != nil {
		log.Fatalf("Failed to save montage image: %v", err)
	}
	fmt.Println("Montage Done")
}
