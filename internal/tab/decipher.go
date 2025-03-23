package tab

import (
	"fmt"

	"github.com/akaspb/playfair-cipher/internal/decipher"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

func NewDecipher(decipherService *decipher.Decipher, separator *rune) *Decipher {
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

	return &Decipher{
		decipherService: decipherService,
		separator:       separator,

		ti: ti,
		to: to,
	}
}

var _ Tab = &Decipher{}

type Decipher struct {
	decipherService *decipher.Decipher
	separator       *rune

	ti  textarea.Model
	to  textarea.Model
	err error
}

func (d *Decipher) Update(msg tea.Msg) {
	d.ti, _ = d.ti.Update(msg)

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
			d.ti.SetValue(buff)
		}
	}

	deciphered, err := d.decipherService.Decode(d.ti.Value(), *d.separator)
	if err != nil {
		d.err = err
		return
	}
	d.err = nil

	d.to.SetValue(deciphered)

	if ctrlS {
		clipboard.WriteAll(deciphered)
	}
}

func (d *Decipher) View() string {

	if d.err != nil {
		return fmt.Sprintf(`Напишите здесь свой зашифрованный текст:
%s
Расшифрованный текст:
* %s
          (ctrl+v - загрузить из буфера обмена)`,
			d.ti.View(),
			d.err.Error(),
		)
	}

	return fmt.Sprintf(`Напишите здесь свой зашифрованный текст:
%s
Расшифрованный текст:
%s
          (ctrl+v - загрузить из буфера обмена)
          (ctrl+s - загрузить в  буфера обмена)`,
		d.ti.View(),
		d.to.View(),
	)
}
