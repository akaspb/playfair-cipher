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

func New(chars []rune, n, m int, key string) (*Cipher, error) {
	if n < 2 {
		panic("n < 2")
	}

	if m < 2 {
		panic("m < 2")
	}

	grid, positions, err := keymatrix.Calculate(chars, n, m, key)
	if err != nil {
		return nil, err
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

	n := len(c.grid)
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		_, err := sb.WriteString(string(c.grid[i]))
		if err != nil {
			panic(err)
		}

		if i+1 < n {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}

func (c *Cipher) Code(text string) string {
	const lambda = 'x' // c.grid[0][0]

	pairs := getPairs(text, lambda)
	if len(pairs)%2 == 1 {
		panic("pairs % 2 == 1")
	}

	fmt.Println(string(pairs))

	height, width := len(c.grid), len(c.grid[0])
	cipherPairs := make([]rune, 0, len(pairs))
	for i := 0; i < len(pairs); i += 2 {
		char1, char2 := pairs[i], pairs[i+1]
		pos1, ok := c.positions[char1]
		if !ok {
			panic("")
		}

		pos2, ok := c.positions[char2]
		if !ok {
			panic("")
		}

		if pos1 == pos2 {
			panic("pos1 == pos2")
		}

		pos1To, pos2To := procPair(pos1, pos2, height, width)
		char1To := c.grid[pos1To.I()][pos1To.J()]
		char2To := c.grid[pos2To.I()][pos2To.J()]

		fmt.Printf("%c %c -> %c %c\n", char1, char2, char1To, char2To)

		cipherPairs = append(cipherPairs, char1To, char2To)
	}

	return string(cipherPairs)
}

func getPairs(text string, lambda rune) []rune {
	res := make([]rune, 0, 2*len(text))
	var prevChar rune
	for _, char := range text {
		if len(res)%2 == 0 {
			res = append(res, char)
			prevChar = char
			continue
		}

		if prevChar == char {
			res = append(res, lambda)
			// prevChar = 0  не влияет на работу алгоритма
		}

		res = append(res, char)
	}

	if len(res)%2 == 1 {
		res = append(res, lambda)
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
