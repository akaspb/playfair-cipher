package main

import (
	"fmt"
	"log"

	"github.com/akaspb/playfair-cipher/internal/config"
	"github.com/akaspb/playfair-cipher/internal/model"
)

func main() {
	cfg := model.Config{
		GridConfig: model.GridConfig{
			Chars: []rune{
				'a', 'b', 'c', 'd', 'e',
				'f', 'g', 'h', 'i', 'k',
				'l', 'm', 'n', 'o', 'p',
				'q', 'r', 's', 't', 'u',
				'v', 'w', 'x', 'y', 'z',
				' ', ',', '.', '!', '#',
			},
			Height: 6,
			Width:  5,
			Key:    "playfairexample",
		},
	}

	if err := config.CreateConfigFile(cfg); err != nil {
		log.Fatal(fmt.Errorf("error during creating config file: %w", err))
	}
}
