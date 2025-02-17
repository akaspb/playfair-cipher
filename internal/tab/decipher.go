package tab

import (
	"fmt"

	"github.com/akaspb/playfair-cipher/internal/decipher"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

func NewDecipher(decipherService *decipher.Decipher) *Decipher {
	ti := textarea.New()
	ti.Placeholder = "Write your ciphered text"
	ti.SetHeight(8)
	ti.SetWidth(50)
	ti.CharLimit = 500
	ti.Focus()

	to := textarea.New()
	to.Placeholder = "Deciphered text"
	to.SetHeight(8)
	to.SetWidth(50)
	ti.CharLimit = 500

	return &Decipher{
		decipherService: decipherService,
		ti:              ti,
		to:              to,
	}
}

var _ Tab = &Decipher{}

type Decipher struct {
	decipherService *decipher.Decipher
	ti              textarea.Model
	to              textarea.Model
	err             error
}

func (d *Decipher) Update(msg tea.Msg) {
	d.ti, _ = d.ti.Update(msg)

	deciphered, err := d.decipherService.Decode(d.ti.Value(), '#')
	if err != nil {
		d.err = err
		return
	}
	d.err = nil

	d.to.SetValue(deciphered)
}

func (d *Decipher) View() string {

	if d.err != nil {
		return fmt.Sprintf(`Input:
%s
Result:
* %s
`,
			d.ti.View(),
			d.err.Error(),
		)
	}

	return fmt.Sprintf(`Input:
%s
Result:
%s
`,
		d.ti.View(),
		d.to.View(),
	)
}
