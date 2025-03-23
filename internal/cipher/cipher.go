package cipher

import (
	"fmt"
	"strings"

	"github.com/akaspb/playfair-cipher/internal/keymatrix"
	"github.com/akaspb/playfair-cipher/internal/model"
)

type Cipher struct {
	grid      [][]rune
	positions map[rune]model.Pos
}

func New(cfg *model.GridConfig) (*Cipher, error) {
	grid, positions, err := keymatrix.Calculate(cfg.Chars, cfg.Height, cfg.Width, cfg.Key)
	if err != nil {
		return nil, fmt.Errorf("error during grid making: %w", err)
	}

	return &Cipher{
		grid:      grid,
		positions: positions,
	}, nil
}

func (c *Cipher) String() string {
	if c == nil {
		return "nil"
	}

	height := len(c.grid)
	sb := strings.Builder{}
	for i := 0; i < height; i++ {
		_, err := sb.WriteString(string(c.grid[i]))
		if err != nil {
			panic(err)
		}

		if i+1 < height {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}

func (c *Cipher) Code(text string, separator rune) (string, error) {
	if strings.ContainsRune(text, separator) {
		return "", fmt.Errorf("[text] must not contain [separator] '%c'", separator)
	}

	if _, ok := c.positions[separator]; !ok {
		return "", fmt.Errorf("[separator] '%c' not in grid", separator)
	}

	pairs := getPairs(text, separator)
	if len(pairs)%2 == 1 {
		panic("pairs % 2 == 1")
	}

	height, width := len(c.grid), len(c.grid[0])
	cipherPairs := make([]rune, 0, len(pairs))
	for i := 0; i < len(pairs); i += 2 {
		char1, char2 := pairs[i], pairs[i+1]
		pos1, ok := c.positions[char1]
		if !ok {
			return "", fmt.Errorf("char '%c' not found in grid", char1)
		}

		pos2, ok := c.positions[char2]
		if !ok {
			return "", fmt.Errorf("char '%c' not found in grid", char2)
		}

		if pos1 == pos2 {
			panic("pos1 == pos2")
		}

		pos1To, pos2To := procPair(pos1, pos2, height, width)
		char1To := c.grid[pos1To.I()][pos1To.J()]
		char2To := c.grid[pos2To.I()][pos2To.J()]

		cipherPairs = append(cipherPairs, char1To, char2To)
	}

	return string(cipherPairs), nil
}

func getPairs(text string, sep rune) []rune {
	chars := []rune(text)

	res := make([]rune, 0, 2*len(chars))
	var prevChar rune
	for _, char := range chars {
		if len(res)%2 == 0 {
			res = append(res, char)
			prevChar = char
			continue
		}

		if prevChar == char {
			res = append(res, sep)
			// prevChar = 0  не влияет на работу алгоритма
		}

		res = append(res, char)
	}

	if len(res)%2 == 1 {
		res = append(res, sep)
	}

	return res
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
	return model.Pos{p1.I(), (p1.J() + 1) % width}, model.Pos{p2.I(), (p2.J() + 1) % width}
}

func procVertical(p1, p2 model.Pos, height int) (_, _ model.Pos) {
	return model.Pos{(p1.I() + 1) % height, p1.J()}, model.Pos{(p2.I() + 1) % height, p2.J()}
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
