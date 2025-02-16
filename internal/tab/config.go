package tab

import (
	"fmt"

	"github.com/akaspb/playfair-cipher/internal/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputIdx int

const (
	keyIn inputIdx = iota
	sepIn
	abcIn
	heightIn
	widthIn
)

func NewConfig(config *model.Config) *Config {
	key := textinput.New()
	{
		key.Placeholder = "Key"
		key.Prompt = ""
		key.CharLimit = 156
		key.Width = 20
		key.Focus()
	}

	abc := textinput.New()
	{
		abc.Placeholder = "ABC"
		abc.Prompt = ""
		abc.CharLimit = 156
		abc.Width = 20
	}

	sep := textinput.New()
	{
		sep.Placeholder = "Separator"
		sep.Prompt = ""
		sep.CharLimit = 2
		sep.Width = 20
		// abc.Validate = func(s string) error {
		// 	_, err := strconv.ParseInt(s, 10, 64)
		// 	return err
		// }
	}

	width := textinput.New()
	{
		width.Placeholder = "XX"
		width.Prompt = ""
		width.CharLimit = 2
		width.Width = 20
		// abc.Validate = func(s string) error {
		// 	_, err := strconv.ParseInt(s, 10, 64)
		// 	return err
		// }
	}

	height := textinput.New()
	{
		height.Placeholder = "XX"
		height.Prompt = ""
		height.CharLimit = 2
		height.Width = 20
		// abc.Validate = func(s string) error {
		// 	_, err := strconv.ParseInt(s, 10, 64)
		// 	return err
		// }
	}

	return &Config{
		config: config,
		textInputs: map[inputIdx]*textinput.Model{
			keyIn:    &key,
			sepIn:    &sep,
			abcIn:    &abc,
			widthIn:  &width,
			heightIn: &height,
		},
		inputIdx: 0,
	}
}

type Config struct {
	config     *model.Config
	textInputs map[inputIdx]*textinput.Model
	inputIdx   inputIdx
}

func (c *Config) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "up":
			prevIdx := c.inputIdx
			c.inputIdx = max(c.inputIdx-1, 0)
			if prevIdx != c.inputIdx {
				c.textInputs[prevIdx].Blur()
			}

			c.textInputs[c.inputIdx].Focus()
		case "down":
			c.textInputs[c.inputIdx].Focus()
			prevIdx := c.inputIdx
			c.inputIdx = min(c.inputIdx+1, inputIdx(len(c.textInputs)-1))
			if prevIdx != c.inputIdx {
				c.textInputs[prevIdx].Blur()
			}

			c.textInputs[c.inputIdx].Focus()
		}
	}

	model, _ := c.textInputs[c.inputIdx].Update(msg)
	c.textInputs[c.inputIdx] = &model
}

func (c *Config) View() string {
	return fmt.Sprintf(`    Write key word:
    %s %d
    Separator:
    %s
    Write ABC:
    %s %d
	Height: %s
	Width: %s

    (ctrl+c to quit)
`,
		c.textInputs[keyIn].View(),
		len(c.textInputs[keyIn].Value()),
		c.textInputs[sepIn].View(),
		c.textInputs[abcIn].View(),
		len(c.textInputs[abcIn].Value()),
		c.textInputs[heightIn].View(), c.textInputs[widthIn].View(),
	)
}
