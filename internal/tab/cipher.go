package tab

import (
	"fmt"

	"github.com/akaspb/playfair-cipher/internal/cipher"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

func NewCipher(cipherService *cipher.Cipher) *Cipher {
	ti := textarea.New()
	ti.Placeholder = "Write your text hear"
	ti.SetHeight(8)
	ti.SetWidth(50)
	ti.CharLimit = 500
	ti.Focus()

	to := textarea.New()
	to.Placeholder = "Ciphered text"
	to.SetHeight(8)
	to.SetWidth(50)
	ti.CharLimit = 500

	return &Cipher{
		cipherService: cipherService,
		ti:            ti,
		to:            to,
	}
}

var _ Tab = &Cipher{}

type Cipher struct {
	cipherService *cipher.Cipher
	ti            textarea.Model
	to            textarea.Model
	err           error
}

func (c *Cipher) Update(msg tea.Msg) {
	c.ti, _ = c.ti.Update(msg)

	ciphered, err := c.cipherService.Code(c.ti.Value(), '#')
	if err != nil {
		c.err = err
		return
	}
	c.err = nil

	c.to.SetValue(ciphered)
}

func (c *Cipher) View() string {

	if c.err != nil {
		return fmt.Sprintf(`Input:
%s
Result:
* %s
`,
			c.ti.View(),
			c.err.Error(),
		)
	}

	return fmt.Sprintf(`Input:
%s
Result:
%s
`,
		c.ti.View(),
		c.to.View(),
	)
}
