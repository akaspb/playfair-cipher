package keymatrix

import (
	"errors"

	"github.com/akaspb/playfair-cipher/internal/model"
)

func Calculate(chars []rune, height, width int, key string) (grid [][]rune, positions map[rune]model.Pos, err error) {
	if height < 2 {
		return nil, nil, errors.New("[height] must be > 1")
	}

	if width < 2 {
		return nil, nil, errors.New("[height] must be > 1")
	}

	count := len(chars)
	if count != height*width {
		return nil, nil, errors.New("chars count != height * width")
	}

	if len(key) < 1 {
		return nil, nil, errors.New("key must be non-empty string")
	}

	positions = make(map[rune]model.Pos, count)
	grid = make([][]rune, height)
	for i := 0; i < height; i++ {
		grid[i] = make([]rune, width)
	}

	q := 0
	for _, char := range key {
		_, ok := positions[char]
		if ok {
			continue
		}

		i := q / width
		j := q % width
		grid[i][j] = char
		positions[char] = model.Pos{i, j}

		q++
	}

	for _, char := range chars {
		_, ok := positions[char]
		if ok {
			continue
		}

		i := q / width
		j := q % width
		grid[i][j] = char
		positions[char] = model.Pos{i, j}

		q++
	}

	if len(positions) != count {
		return nil, nil, errors.New("some chars are duplicated")
	}

	return
}
