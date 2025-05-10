package tab

import (
	"fmt"
	"strings"

	"github.com/akaspb/playfair-cipher/internal/decipher"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func NewDecipher(decipherService *decipher.Decipher, separator *rune) *Decipher {
	fi := textinput.New()
	fi.Placeholder = "Название файла с расшиением"
	fi.Prompt = "> "
	fi.CharLimit = 64
	fi.Width = 64
	fi.Cursor.Style = cursorStyle
	fi.PromptStyle = focusedStyle
	fi.TextStyle = focusedStyle

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

		fi: fi,
		ti: ti,
		to: to,
	}
}

var _ Tab = &Decipher{}

type Decipher struct {
	decipherService *decipher.Decipher
	separator       *rune

	fileIsSaved bool
	fi          textinput.Model
	ti          textarea.Model
	to          textarea.Model
	err         error
}

func (d *Decipher) Update(msg tea.Msg) {
	d.err = nil
	d.fileIsSaved = false

	var (
		ctrlV    = false
		ctrlS    = false
		loadText = false
		saveText = false
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+v":
			ctrlV = true
		case "ctrl+s":
			ctrlS = true
		case "ctrl+r":
			loadText = true
		case "ctrl+w":
			saveText = true
		case "ctrl+d":
			d.ti.SetValue("")
		case "up":
			d.fi.Focus()
			d.ti.Blur()
		case "down":
			d.fi.Blur()
			d.ti.Focus()
		default:
			if d.fi.Focused() {
				d.fi, _ = d.fi.Update(msg)
			}
			if d.ti.Focused() {
				d.ti, _ = d.ti.Update(msg)
			}
		}
	}

	if ctrlV {
		buff, err := clipboard.ReadAll()
		if err == nil {
			d.ti.SetValue(buff)
		}
	}

	if loadText {
		text, err := loadFile(d.fi.Value())
		if err != nil {
			d.err = err
		} else {
			d.ti.SetValue(text)
		}
	}

	deciphered, err := d.decipherService.Decode(d.ti.Value(), *d.separator)
	if err != nil {
		d.err = err
		return
	}

	d.to.SetValue(deciphered)

	if ctrlS {
		clipboard.WriteAll(deciphered)
	}

	if saveText {
		err = saveFile(d.fi.Value(), deciphered)
		if err != nil {
			d.err = err
		} else {
			d.fileIsSaved = true
		}
	}
}

func (d *Decipher) View() string {
	sb := strings.Builder{}
	sb.WriteString(d.fi.View())
	if d.fileIsSaved {
		sb.WriteString("* файл записан")
	}
	sb.WriteString("\n")

	if d.err != nil {
		sb.WriteString(fmt.Sprintf(`Напишите здесь свой зашифрованный текст:
%s
Расшифрованный текст:
* %s
			(ctrl+v - загрузить из буфера обмена)`,
			d.ti.View(),
			d.err.Error(),
		))
	} else {
		sb.WriteString(fmt.Sprintf(`Напишите здесь свой зашифрованный текст:
%s
Расшифрованный текст:
%s
  (ctrl+v / ctrl+r - загрузить из буфера обмена / файла)
   (ctrl+s / ctrl+w - загрузить в буфер обмена / файл)
					(ctrl+d - удалить текст)`,
			d.ti.View(),
			d.to.View(),
		))
	}

	return sb.String()
}
