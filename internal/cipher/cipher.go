package cipher

import (
	"strings"

	"github.com/Aleksandr-qefy/playfair-cipher/internal/keymatrix"
	"github.com/Aleksandr-qefy/playfair-cipher/internal/model"
)

type Cipher struct {
	grid      [][]rune
	positions map[rune]model.Pos
}

func New(chars []rune, n, m int, key string) (*Cipher, error) {
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

func (c *Cipher) Code() {

}
