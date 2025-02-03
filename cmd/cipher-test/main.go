package main

import (
	"fmt"
	"log"

	"github.com/Aleksandr-qefy/playfair-cipher/internal/cipher"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cipherService, err := cipher.New(
		[]rune{
			'a', 'b', 'c', 'd', 'e',
			'f', 'g', 'h', 'i', 'k',
			'l', 'm', 'n', 'o', 'p',
			'q', 'r', 's', 't', 'u',
			'v', 'w', 'x', 'y', 'z',
		},
		5,
		5,
		"playfairexample",
	)
	if err != nil {
		return err
	}

	fmt.Println(cipherService)

	return nil
}
