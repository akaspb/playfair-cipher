package main

import (
	"fmt"
	"log"

	"github.com/akaspb/playfair-cipher/internal/cipher"
	"github.com/akaspb/playfair-cipher/internal/config"
	"github.com/akaspb/playfair-cipher/internal/decipher"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.LoadConfigFile()
	if err != nil {
		return fmt.Errorf("could't load config: %w", err)
	}

	cipherService, err := cipher.New(cfg.GridConfig)
	if err != nil {
		return err
	}

	input := "thisistest" // "thisistest"
	cphr, err := cipherService.Code(input, 'q')
	if err != nil {
		return err
	}
	fmt.Println(cphr)

	decipherService, err := decipher.New(cfg.GridConfig)
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
