package cipher

import (
	"errors"
	"math"
)

type pos [2]int

func (p pos) i() int {
	return p[0]
}

func (p pos) j() int {
	return p[1]
}

type Cipher struct {
	grid    [][]rune
	charPos map[rune]pos
}

func New(chars []rune, n, m int, key string) (*Cipher, error) {
	count := len(chars)
	if count != n*m {
		return nil, errors.New("chars count != n * m")
	}

	maxInt := math.MaxInt
	charsOrder := make(map[rune]int, count)
	for i := 0; i < count; i++ {
		charsOrder[chars[i]] = maxInt - i
	}

	k := 0
	for _, char := range chars {
		if charsOrder[char] < k {
			continue
		}

		charsOrder[char] = k
		k++
	}

	grid := make([][]rune, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]rune, m)
	}

	charPos := make(map[rune]pos, count)
	for _, char := range chars {
		k := charsOrder[char]
		i := k / m
		j := k % m
		grid[i][j] = char
		charPos[char] = pos{i, j}
	}

	return &Cipher{
		grid:    grid,
		charPos: charPos,
	}, nil
}
