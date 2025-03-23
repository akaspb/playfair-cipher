package decipher

import (
	"errors"
	"fmt"
	"strings"

	"github.com/akaspb/playfair-cipher/internal/keymatrix"
	"github.com/akaspb/playfair-cipher/internal/model"
)

type Decipher struct {
	grid      [][]rune
	positions map[rune]model.Pos
}

func New(cfg *model.GridConfig) (*Decipher, error) {
	grid, positions, err := keymatrix.Calculate(cfg.Chars, cfg.Height, cfg.Width, cfg.Key)
	if err != nil {
		return nil, fmt.Errorf("error during grid making: %w", err)
	}

	return &Decipher{
		grid:      grid,
		positions: positions,
	}, nil
}

func (d *Decipher) String() string {
	if d == nil {
		return "nil"
	}

	height := len(d.grid)
	sb := strings.Builder{}
	for i := 0; i < height; i++ {
		_, err := sb.WriteString(string(d.grid[i]))
		if err != nil {
			panic(err)
		}

		if i+1 < height {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}

func (d *Decipher) Decode(cipherText string, separator rune) (string, error) {
	pairs := []rune(cipherText)

	if len(pairs)%2 == 1 {
		return "", errors.New("[cipherText] must be even-length string")
	}

	height, width := len(d.grid), len(d.grid[0])
	decipherPairs := make([]rune, 0, len(pairs))
	for i := 0; i < len(pairs); i += 2 {
		char1, char2 := pairs[i], pairs[i+1]
		pos1, ok := d.positions[char1]
		if !ok {
			return "", fmt.Errorf("char '%c' not found in grid", char1)
		}

		pos2, ok := d.positions[char2]
		if !ok {
			return "", fmt.Errorf("char '%c' not found in grid", char2)
		}

		if pos1 == pos2 {
			return "", errors.New("incorrect ciphered text")
		}

		pos1To, pos2To := procPair(pos1, pos2, height, width)
		char1To := d.grid[pos1To.I()][pos1To.J()]
		char2To := d.grid[pos2To.I()][pos2To.J()]

		decipherPairs = append(decipherPairs, char1To, char2To)
	}

	if separator == 0 {
		return string(decipherPairs), nil
	}

	return string(removeSeparator(decipherPairs, separator)), nil
}

func procPair(p1, p2 model.Pos, height, width int) (_, _ model.Pos) {
	switch {
	case p1.I() == p2.I():
		return procHorizontal(p1, p2, width)
	case p1.J() == p2.J():
		return procVertical(p1, p2, height)
	default:
	}

	return procRectangle(p1, p2)
}

func procHorizontal(p1, p2 model.Pos, width int) (_, _ model.Pos) {
	return model.Pos{p1.I(), (p1.J() - 1 + width) % width}, model.Pos{p2.I(), (p2.J() - 1 + width) % width} // (-1 + X) % X == X-1
}

func procVertical(p1, p2 model.Pos, height int) (_, _ model.Pos) {
	return model.Pos{(p1.I() - 1 + height) % height, p1.J()}, model.Pos{(p2.I() - 1 + height) % height, p2.J()}
}

func procRectangle(p1, p2 model.Pos) (_, _ model.Pos) {
	var reversed bool
	if p1.I() > p2.I() {
		p1, p2 = p2, p1
		reversed = true
	}

	p1To, p2To := model.Pos{p1.I(), p2.J()}, model.Pos{p2.I(), p1.J()}
	if reversed {
		return p2To, p1To
	}

	return p1To, p2To
}

func removeSeparator(s []rune, sep rune) []rune {
	res := make([]rune, 0, len(s))
	for _, char := range s {
		if char == sep {
			continue
		}

		res = append(res, char)
	}

	return res
}
