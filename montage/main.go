package main

import (
	"fmt"
)

func main() {
	cliConfig := getCLIArgs(&cliConf)
	appConfig := getAppConfig(cliConf.AppConfig)

	if cliConf.Verbose {
		fmt.Printf("CLI Config: %+v\n", cliConfig)
		fmt.Printf("App Config: %+v\n", appConfig)
	}

	createGrafanaMontage(cliConfig, appConfig)
}
