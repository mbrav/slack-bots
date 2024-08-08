package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Global app config
var (
	cliConf CLIConfig
	appConf AppConfig
)

// CLI config struct
type CLIConfig struct {
	AppConfig       string
	OutputFile      string
	OutputExtension string
	Verbose         bool
}

// YAML app config struct
type AppConfig struct {
	Name                  string   `yaml:"name" env-default:"Untitled Dashboard"`
	SlackToken            string   `yaml:"slack_token" env:"SLACK_TOKEN"`
	SlackChannel          string   `yaml:"slack_channel" env:"SLACK_CHANNEL"`
	SlackMentions         []string `yaml:"slack_mentions"`
	GrafanaUser           string   `yaml:"grafana_user" env:"GRAFANA_USER"`
	GrafanaPassword       string   `yaml:"grafana_password" env:"GRAFANA_PASSWORD"`
	MontageTile           string   `yaml:"montage_tile" env-default:"2x1"`
	MontageBackgroudColor string   `yaml:"montage_bg_color" env-default:"#111217"`
	Images                []Image  `yaml:"images"`
}

// Image config
type Image struct {
	Name   string `yaml:"name" env-default:"Untitled Image"`
	Url    string `yaml:"url" env-required:""`
	Width  int    `yaml:"width" env-default:"800"`
	Height int    `yaml:"height" env-default:"500"`
}

// Initializes the CLIConfig struct with command-line arguments
func getCLIArgs(c *CLIConfig) *CLIConfig {
	flag.StringVar(&c.AppConfig, "c", "app.yaml", "AppConfig location")
	flag.StringVar(&c.OutputFile, "o", "output.png", "Base name of the file to output")
	flag.BoolVar(&c.Verbose, "v", false, "Verbose mode")
	flag.Parse()
	return c
}

// Parse yaml config
func getAppConfig(configPath string) *AppConfig {
	cfg := &AppConfig{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	return cfg
}
