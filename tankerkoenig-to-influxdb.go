package main

import (
	"flag"
	"fmt"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type ConfigurationStation struct {
	Id     string `json:"id"`
	Brand  string `json:"brand"`
	City   string `json:"city"`
	Street string `json:"street"`
}

type Configuration struct {
	Stations              []ConfigurationStation `json:"stations"`
}

func main() {
	fmt.Println("--- Hello World ---")

	// Program parameters
	var configFileName string
	flag.StringVar(&configFileName, "c","", "configuration file")
	flag.Parse()
	sourceFiles := flag.Args()
	fmt.Printf("Tail: %+q\n", sourceFiles)

	fmt.Println(configFileName)

	// Load configuration
	config, err := loadConfigurationFile(configFileName)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(config.Stations); i++ {
		station := config.Stations[i]
		fmt.Println(station.Id)
	}
}

func loadConfigurationFile(configFileName string) (*Configuration, error) {
	_, err := os.Stat(configFileName)
	if errors.Is(err, os.ErrNotExist) {
		// Config doesn't exists
		return nil, errors.New("Configuration file not found")
	}
	if err != nil {
		panic(err)
	}

	file, err := ioutil.ReadFile(configFileName)
	if err != nil {
		panic(err)
	}
	configuration := Configuration{}
	err = json.Unmarshal([]byte(file), &configuration)
	if err != nil {
		panic(err)
	}

	return &configuration, nil
}