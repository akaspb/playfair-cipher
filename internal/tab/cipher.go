package tab

import (
	"fmt"
	"os"
	"strings"

	"github.com/akaspb/playfair-cipher/internal/cipher"
	"github.com/akaspb/playfair-cipher/internal/file"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func NewCipher(cipherService *cipher.Cipher, separator *rune) *Cipher {
	fi := textinput.New()
	fi.Placeholder = "Название файла с расширением"
	fi.Prompt = "> "
	fi.CharLimit = 100
	fi.Width = 50
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

	return &Cipher{
		cipherService: cipherService,
		separator:     separator,

		fi: fi,
		ti: ti,
		to: to,
	}
}

var _ Tab = &Cipher{}

type Cipher struct {
	cipherService *cipher.Cipher
	separator     *rune

	fi  textinput.Model
	ti  textarea.Model
	to  textarea.Model
	err error
}

func (c *Cipher) Update(msg tea.Msg) {
	c.err = nil

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
			c.ti.SetValue("")
		case "up":
			c.fi.Focus()
			c.ti.Blur()
		case "down":
			c.fi.Blur()
			c.ti.Focus()
		default:
			if c.fi.Focused() {
				c.fi, _ = c.fi.Update(msg)
			}
			if c.ti.Focused() {
				c.ti, _ = c.ti.Update(msg)
			}
		}
	}

	if ctrlV {
		buff, err := clipboard.ReadAll()
		if err == nil {
			c.ti.SetValue(buff)
		}
	}

	if loadText {
		text, err := loadFile(c.fi.Value())
		if err != nil {
			c.err = err
		} else {
			c.ti.SetValue(text)
		}
	}

	ciphered, err := c.cipherService.Code(c.ti.Value(), *c.separator)
	if err != nil {
		c.err = err
		return
	}

	c.to.SetValue(ciphered)

	if ctrlS {
		clipboard.WriteAll(ciphered)
	}

	if saveText {
		c.err = saveFile(c.fi.Value(), ciphered)
	}
}

func loadFile(fileName string) (string, error) {
	path, err := getWorkingDir()
	if err != nil {
		return "", err
	}

	return file.Load(fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), fileName))
}

func saveFile(fileName string, text string) error {
	path, err := getWorkingDir()
	if err != nil {
		return err
	}

	return file.Save(fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), fileName), text)
}

func getWorkingDir() (string, error) {
	return os.Getwd()
}

func (c *Cipher) View() string {
	sb := strings.Builder{}
	sb.WriteString(c.fi.View())
	sb.WriteString("\n")

	if c.err != nil {
		sb.WriteString(fmt.Sprintf(`Напишите здесь свой текст:
%s
Зашифрованный текст:
* %s
			(ctrl+v - загрузить из буфера обмена)`,
			c.ti.View(),
			c.err.Error(),
		))
	} else {
		sb.WriteString(fmt.Sprintf(`Напишите здесь свой текст:
%s
Зашифрованный текст:
%s
  (ctrl+v / ctrl+r - загрузить из буфера обмена / файла)
   (ctrl+s / ctrl+w - загрузить в буфер обмена / файл)
					(ctrl+d - удалить текст)`,
			c.ti.View(),
			c.to.View(),
		))
	}

	return sb.String()
}
