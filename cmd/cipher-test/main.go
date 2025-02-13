package main

import (
	"fmt"
	"log"

	"github.com/akaspb/playfair-cipher/internal/cipher"
	"github.com/akaspb/playfair-cipher/internal/decipher"
	"github.com/akaspb/playfair-cipher/internal/model"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := model.GridConfig{
		Chars: []rune{
			'a', 'b', 'c', 'd', 'e',
			'f', 'g', 'h', 'i', 'k',
			'l', 'm', 'n', 'o', 'p',
			'q', 'r', 's', 't', 'u',
			'v', 'w', 'x', 'y', 'z',
		},
		Height: 5,
		Width:  5,
		Key:    "playfairexample",
	}

	cipherService, err := cipher.New(cfg)
	if err != nil {
		return err
	}

	input := "thisistest" // "thisistest"
	cphr, err := cipherService.Code(input, 'q')
	if err != nil {
		return err
	}
	fmt.Println(cphr)

	decipherService, err := decipher.New(cfg)
	if err != nil {
		return err
	}

	dcphr, err := decipherService.Decode(cphr, 'q')
	if err != nil {
		return err
	}
	fmt.Println(dcphr)

	return nil
}
