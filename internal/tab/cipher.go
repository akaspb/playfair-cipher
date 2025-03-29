package tab

import (
	"fmt"

	"github.com/akaspb/playfair-cipher/internal/cipher"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

func NewCipher(cipherService *cipher.Cipher, separator *rune) *Cipher {
	ti := textarea.New()
	ti.Placeholder = ""
	ti.SetHeight(8)
	ti.SetWidth(50)
	ti.CharLimit = 500
	ti.Focus()

	to := textarea.New()
	to.Placeholder = ""
	to.SetHeight(8)
	to.SetWidth(50)
	ti.CharLimit = 500

	return &Cipher{
		cipherService: cipherService,
		separator:     separator,

		ti: ti,
		to: to,
	}
}

var _ Tab = &Cipher{}

type Cipher struct {
	cipherService *cipher.Cipher
	separator     *rune

	ti  textarea.Model
	to  textarea.Model
	err error
}

func (c *Cipher) Update(msg tea.Msg) {
	c.ti, _ = c.ti.Update(msg)

	var (
		ctrlV = false
		ctrlS = false
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+v":
			ctrlV = true
		case "ctrl+s":
			ctrlS = true
		}
	}

	if ctrlV {
		buff, err := clipboard.ReadAll()
		if err == nil {
			c.ti.SetValue(buff)
		}
	}

	ciphered, err := c.cipherService.Code(c.ti.Value(), *c.separator)
	if err != nil {
		c.err = err
		return
	}
	c.err = nil

	c.to.SetValue(ciphered)

	if ctrlS {
		clipboard.WriteAll(ciphered)
	}
}

func (c *Cipher) View() string {

	if c.err != nil {
		return fmt.Sprintf(`Напишите здесь свой текст:
%s
Зашифрованный текст:
* %s
          (ctrl+v - загрузить из буфера обмена)`,
			c.ti.View(),
			c.err.Error(),
		)
	}

	return fmt.Sprintf(`Напишите здесь свой текст:
%s
Зашифрованный текст:
%s
          (ctrl+v - загрузить из буфера обмена)
           (ctrl+s - загрузить в буфер обмена)`,
		c.ti.View(),
		c.to.View(),
	)
}
