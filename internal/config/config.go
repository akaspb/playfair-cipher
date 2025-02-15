package config

import (
	"encoding/json"
	"os"

	"github.com/akaspb/playfair-cipher/internal/model"
)

const confFile = "config/config.json"

func CreateConfigFile(c model.Config) error {
	file, err := os.Create(confFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
}

func LoadConfigFile() (model.Config, error) {
	var c model.Config
	file, err := os.Open(confFile)
	if err != nil {
		return c, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	return c, err
}
