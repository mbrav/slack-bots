package main

import (
	"fmt"
	"log"
	"path/filepath"

	"gopkg.in/gographics/imagick.v3/imagick"
)

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

	saveMontage(montage, outFile)
}
