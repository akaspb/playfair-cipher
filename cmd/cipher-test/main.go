package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/akaspb/playfair-cipher/internal/cipher"
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
	// fmt.Println(cipherService.Code(strings.ToLower("IDIOCYOFTENLOOKSLIKEINTELLIGENCE")))

	fmt.Println(cipherService.Code(strings.ToLower("HIDETHEGOLDINTHETREXESTUMP")))
	// BMODZ BXDNA BEKUD MUIXM MOUVI F
	// bmodz bxdna bekud muixm mouvi f

	return nil
}
