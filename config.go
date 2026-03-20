package main

import (
	"encoding/json"
	"os"
)

type ConfigFileFeed struct {
	Name string `json:"name"`
	Url string `json:"url"`
}

type ConfigFile struct {
	Method string `json:"method"`
	Feeds []ConfigFileFeed `json:"feeds"`
}

func decodeConfig() (ConfigFile, error) {
	const path = "./config.json";

	dat, err := os.ReadFile(path);
	if err != nil {
		return ConfigFile{}, err;
	}

	var c ConfigFile;
	json.Unmarshal(dat, &c);

	return c, nil;
}

