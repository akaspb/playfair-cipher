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
	fi.Placeholder = "file name with extension"
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

	fileIsSaved bool
	fi          textinput.Model
	ti          textarea.Model
	to          textarea.Model
	err         error
}

func (c *Cipher) Update(msg tea.Msg) {
	c.err = nil
	c.fileIsSaved = false

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
		err := saveFile(c.fi.Value(), ciphered)
		if err != nil {
			c.err = err
		} else {
			c.fileIsSaved = true
		}
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
	if c.fileIsSaved {
		sb.WriteString("* ciphertext saved to file")
	}
	sb.WriteString("\n")

	if c.err != nil {
		sb.WriteString(fmt.Sprintf(`Your text:
%s
Ciphertext:
* %s
     (ctrl+v / ctrl+r - load from clipboard / file)`,
			c.ti.View(),
			c.err.Error(),
		))
	} else {
		sb.WriteString(fmt.Sprintf(`Your text:
%s
Ciphertext:
%s
     (ctrl+v / ctrl+r - load from clipboard / file)
      (ctrl+s / ctrl+w - save to clipboard / file)
			     (ctrl+d - clear text)`,
			c.ti.View(),
			c.to.View(),
		))
	}

	return sb.String()
}
