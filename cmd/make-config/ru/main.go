package main

import (
	"fmt"
	"log"

	"github.com/akaspb/playfair-cipher/internal/config"
	"github.com/akaspb/playfair-cipher/internal/model"
)

func main() {
	separator := '#'

	cfg := model.Config{
		GridConfig: &model.GridConfig{
			Chars: []rune{
				'а', 'б', 'в', 'г', 'д', 'е',
				'ё', 'ж', 'з', 'и', 'й', 'к',
				'л', 'м', 'н', 'о', 'п', 'р',
				'с', 'т', 'у', 'ф', 'х', 'ц',
				'ч', 'ш', 'щ', 'ъ', 'ы', 'ь',
				'э', 'ю', 'я', ' ', '.', '#',
			},
			Height: 6,
			Width:  6,
			Key:    "парольная фраза",
		},
		Separator: &separator,
	}

	if err := config.CreateConfigFile(cfg); err != nil {
		log.Fatal(fmt.Errorf("error during creating config file: %w", err))
	}
}
