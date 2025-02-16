package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	MeditationsCSVPath string
	StopwordsPath string
	VocabPath string
	VectorsPath string
	NgramsPath string
	SystemPrompt string
	Preface string
}

func GetConfig() Config {
	file, err := os.Open("../config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	return config
}
