package main

import (
	"github.com/spf13/viper"
	"os"
)

// Configuration read from config.json
type Configuration struct {
	// Hash length
	HashLength int
	// Listening address:port
	ListenAddr string
	// Maximum uploaded file size (0 = unlimited)
	MaxFileSize int
	// Storage path
	StoragePath string
}

var configuration Configuration

const minHashLen int = 1
const maxHashLen int = 30

func loadConfiguration() {
	viper.SetDefault("HashLength", 5)
	viper.SetDefault("ListenAddr", "0.0.0.0:8066")
	viper.SetDefault("MaxFileSize", 10240000)
	viper.SetDefault("StoragePath", "./storage")

	viper.SetEnvPrefix("Mangonel")
	viper.AutomaticEnv()

	configuration = Configuration{
		viper.GetInt("HashLength"),
		viper.GetString("ListenAddr"),
		viper.GetInt("MaxFileSize"),
		viper.GetString("StoragePath"),
	}

	if configuration.HashLength > maxHashLen {
		configuration.HashLength = maxHashLen
	}

	if configuration.HashLength < minHashLen {
		configuration.HashLength = minHashLen
	}

	// No sanity check for ListenAddr because http.ListenAndServe will Fatal if necessary

	// No sanity check for MaxFileSize because there's no limit to madness!

	f, err := os.Stat(configuration.StoragePath)
	if err != nil {
		panic(err)
	}

	if !f.IsDir() {
		panic("configuration.storage path is not a valid directory")
	}
}
