package keymatrix

import (
	"errors"

	"github.com/akaspb/playfair-cipher/internal/model"
)

func Calculate(chars []rune, n, m int, key string) (grid [][]rune, positions map[rune]model.Pos, err error) {
	count := len(chars)
	if count != n*m {
		return nil, nil, errors.New("chars count != n * m")
	}

	positions = make(map[rune]model.Pos, count)
	grid = make([][]rune, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]rune, m)
	}

	q := 0
	for _, char := range key {
		_, ok := positions[char]
		if ok {
			continue
		}

		i := q / m
		j := q % m
		grid[i][j] = char
		positions[char] = model.Pos{i, j}

		q++
	}

	for _, char := range chars {
		_, ok := positions[char]
		if ok {
			continue
		}

		i := q / m
		j := q % m
		grid[i][j] = char
		positions[char] = model.Pos{i, j}

		q++
	}

	return
}
