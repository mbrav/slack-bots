package main

import (
	"flag"
)

// Global app config
var (
	cliConf CLIConfig
	// appConf AppConfig
)

// CLI config struct
type CLIConfig struct {
	AppConfig       string
	InputFiles      string
	OutputFile      string
	OutputExtension string
	Verbose         bool
}

// // YAML app config struct
// type AppConfig struct {
// 	Database struct {
// 		Host string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
// 	} `yaml:"database"`
// 	Server struct {
// 		Address string `yaml:"address" env:"SERVER_ADDRESS" env-default:"192.168.1.1"`
// 	} `yaml:"server"`
// }

// Initializes the CLIConfig struct with command-line arguments
func getCLIArgs(c *CLIConfig) *CLIConfig {
	flag.StringVar(&c.AppConfig, "c", "app.yaml", "AppConfig location")
	flag.StringVar(&c.InputFiles, "i", "test.png", "List of files to montage")
	flag.StringVar(&c.OutputFile, "o", "output.jpg", "Base name of the file to output")
	flag.StringVar(&c.OutputExtension, "e", "png", "Extension of the output file")
	flag.BoolVar(&c.Verbose, "v", false, "Verbose mode")
	flag.Parse()
	return c
}

//
// // Parse yaml config
// func getAppConfig(configPath string) *AppConfig {
// 	cfg := &AppConfig{}
//
// 	err := cleanenv.ReadConfig(configPath, cfg)
// 	if err != nil {
// 		fmt.Printf("Error reading config file: %v\n", err)
// 		os.Exit(1)
// 	}
//
// 	return cfg
// }
