package main

import (
	"encoding/json"
	"os"
)

// Configuration read from config.json
type Configuration struct {
	// Hash length
	Length int
	// Listening address:port
	Addr string
	// Storage path
	Storage string
	// Maximum uploaded filesize (0 = unlimited)
	MaxFilesize int
}

var configuration Configuration

const minHashLen int = 1
const maxHashLen int = 30

func loadConfiguration() {
	file, err := os.Open("mangonel-config.json")
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}

	if configuration.Length > maxHashLen {
		configuration.Length = maxHashLen
	}

	if configuration.Length < minHashLen {
		configuration.Length = minHashLen
	}

	f, err := os.Stat(configuration.Storage)
	if err != nil {
		panic(err)
	}

	if !f.IsDir() {
		panic("configuration.storage path is not a valid directory")
	}
}
