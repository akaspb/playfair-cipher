package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/akaspb/playfair-cipher/internal/model"
)

func CreateConfigFile(c model.Config) error {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	confFile := filepath.Join(execPath, "..", "..", "config", "config.txt")

	file, err := os.Create(confFile)
	if err != nil {
		return err
	}
	defer file.Close()

	cfgText, err := createConfigText(c)
	if err != nil {
		return err
	}

	_, err = file.WriteString(cfgText)

	return err
}

func createConfigText(c model.Config) (string, error) {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("%d %d\n", c.Height, c.Width))

	for i := 0; i < c.Height; i++ {
		for j := 0; j < c.Width; j++ {
			sb.WriteRune(c.Chars[i*c.Width+j])
		}
		sb.WriteRune('\n')
	}

	idx := findEl(c.Chars, *c.Separator)
	if idx == -1 {
		return "", fmt.Errorf("can't find separator in chars")
	}

	sb.WriteString(fmt.Sprintf("%d\n", idx))

	return sb.String(), nil
}

func findEl[T comparable](slc []T, el T) int {
	for i, slcEl := range slc {
		if slcEl == el {
			return i
		}
	}

	return -1
}

func LoadConfigFile() (model.Config, error) {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	confFile := filepath.Join(execPath, "..", "..", "config", "config.txt")

	confData, err := os.ReadFile(confFile)
	if err != nil {
		return model.Config{}, err
	}

	return loadConfigText(string(confData))
}

func loadConfigText(cfgText string) (model.Config, error) {
	c := model.Config{}

	lines := strings.Split(cfgText, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimFunc(line, func(r rune) bool {
			return !unicode.IsGraphic(r)
		})
	}

	if len(lines) <= 2 {
		return model.Config{}, fmt.Errorf("incorrect config file")
	}

	_, err := fmt.Sscanf(lines[0], "%d %d", &c.Height, &c.Width)
	if err != nil {
		return model.Config{}, fmt.Errorf("can't read matrix height and width from config file: %w", err)
	}

	if c.Height < 1 {
		return model.Config{}, fmt.Errorf("incorrect height value")
	}

	if c.Width < 1 {
		return model.Config{}, fmt.Errorf("incorrect width valuee")
	}

	if len(lines) < c.Height+2 {
		return model.Config{}, fmt.Errorf("incorrext config file")
	}

	for _, line := range lines[1 : c.Height+1] {
		c.Chars = append(c.Chars, []rune(line)...)
	}

	if len(c.Chars) != c.Height*c.Width {
		return model.Config{}, fmt.Errorf("incorrext matrix")
	}

	var idx int
	_, err = fmt.Sscanf(lines[c.Height+1], "%d", &idx)
	if err != nil {
		return model.Config{}, fmt.Errorf("can't read separator position in matrix from config file: %w", err)
	}

	if !(0 <= idx && idx < len(c.Chars)) {
		return model.Config{}, fmt.Errorf("incorrect separator position in matrix from config file")
	}

	c.Separator = &(c.Chars[idx])

	return c, nil
}
