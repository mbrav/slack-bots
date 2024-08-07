package main

import (
	"fmt"
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

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// err := mw.ReadImage(cliConfig.InputFiles)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// // Get original logo size
	// width := mw.GetImageWidth()
	// height := mw.GetImageHeight()
	//
	// // Calculate half the size
	// hWidth := uint(width / 10)
	// hHeight := uint(height / 10)
	//
	// // Resize the image using the Lanczos filter
	// // The blur factor is a float, where > 1 is blurry, < 1 is sharp
	// err = mw.ResizeImage(hWidth, hHeight, imagick.FILTER_LANCZOS)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// // Set the compression quality to 95 (high quality = low compression)
	// err = mw.SetImageCompressionQuality(10)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// if err = mw.WriteImage(cliConfig.OutputFile); err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Printf("Wrote: %s\n", cliConfig.OutputFile)
	//

	for i, img := range appConfig.Images {
		wg.Add(1)
		go downloadImage(
			&wg,
			img.Url,
			fmt.Sprintf("test/img-%d.png", i),
			appConfig.GrafanaUser,
			appConfig.GrafanaPassword)
	}

	wg.Wait()
	fmt.Println("All downloads completed.")
}
