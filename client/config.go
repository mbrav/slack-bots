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
	AppConfig  string
	KubeConfig string
	Namespace  string
	Verbose    bool
}

// YAML app config struct
type AppConfig struct {
	Database struct {
		Host string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	} `yaml:"database"`
	Server struct {
		Address string `yaml:"address" env:"SERVER_ADDRESS" env-default:"192.168.1.1"`
	} `yaml:"server"`
}

// Initializes the CLIConfig struct with command-line arguments
func getCLIArgs(c *CLIConfig) *CLIConfig {
	flag.StringVar(&c.AppConfig, "c", "./app.yaml", "AppConfig location")
	flag.StringVar(&c.KubeConfig, "k", "~/.kube/config", "Kubeconfig location")
	flag.StringVar(&c.Namespace, "n", "slack-bots", "Namespace in which bot is launched")
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
