package tab

import (
	"fmt"
	"log"
	"strconv"

	"github.com/akaspb/playfair-cipher/internal/config"
	configfile "github.com/akaspb/playfair-cipher/internal/config"
	"github.com/akaspb/playfair-cipher/internal/keymatrix"
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
		key.Prompt = "> "
		key.CharLimit = 100
		key.Width = 50
		key.Focus()
	}

	abc := textinput.New()
	{
		abc.Placeholder = "ABC"
		abc.Prompt = "> "
		abc.CharLimit = 100
		abc.Width = 50
	}

	sep := textinput.New()
	{
		sep.Placeholder = "Separator"
		sep.Prompt = "> "
		sep.CharLimit = 1
		sep.Width = 10
	}

	width := textinput.New()
	{
		width.Placeholder = "XX"
		width.Prompt = ""
		width.CharLimit = 2
		width.Width = 2
	}

	height := textinput.New()
	{
		height.Placeholder = "XX"
		height.Prompt = ""
		height.CharLimit = 2
		height.Width = 2
	}

	c := &Config{
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
	c.loadConfig()

	return c
}

func (c *Config) loadConfig() {
	cfg, err := configfile.LoadConfigFile()
	if err != nil {
		panic(err)
	}

	c.textInputs[keyIn].SetValue(cfg.GridConfig.Key)
	c.textInputs[sepIn].SetValue(string([]rune{cfg.Separator}))
	c.textInputs[abcIn].SetValue(string(cfg.GridConfig.Chars))
	c.textInputs[widthIn].SetValue(strconv.Itoa(cfg.GridConfig.Height))
	c.textInputs[heightIn].SetValue(strconv.Itoa(cfg.GridConfig.Width))
}

type Config struct {
	config     *model.Config
	textInputs map[inputIdx]*textinput.Model
	inputIdx   inputIdx
	saveRes    string
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
		case "ctrl+z":
			c.loadConfig()
			c.saveRes = "* reloaded"
		case "ctrl+s":
			err := c.saveConfig()
			if err != nil {
				c.saveRes = fmt.Sprintf("* %s", err.Error())
			} else {
				c.saveRes = "* saved"
			}
		}
	}

	model, _ := c.textInputs[c.inputIdx].Update(msg)
	c.textInputs[c.inputIdx] = &model
}

func (c *Config) saveConfig() error {
	var (
		key = c.textInputs[keyIn].Value()
		sep = c.textInputs[sepIn].Value()
		abc = c.textInputs[abcIn].Value()
	)

	if err := textFieldValidator(key, "key"); err != nil {
		return err
	}

	if err := textFieldValidator(sep, "separator"); err != nil {
		return err
	}

	if err := textFieldValidator(abc, "abc"); err != nil {
		return err
	}

	if err := numFieldValidator(c.textInputs[heightIn].Value(), "height"); err != nil {
		return err
	}

	if err := numFieldValidator(c.textInputs[widthIn].Value(), "width"); err != nil {
		return err
	}

	height, _ := strconv.Atoi(c.textInputs[heightIn].Value())
	width, _ := strconv.Atoi(c.textInputs[widthIn].Value())

	if _, _, err := keymatrix.Calculate(
		[]rune(c.textInputs[abcIn].Value()),
		height,
		width,
		c.textInputs[keyIn].Value(),
	); err != nil {
		return err
	}

	if err := config.CreateConfigFile(model.Config{
		GridConfig: model.GridConfig{
			Chars:  []rune(abc),
			Height: height,
			Width:  width,
			Key:    key,
		},
		Separator: []rune(sep)[0],
	}); err != nil {
		log.Fatal(fmt.Errorf("error during creating config file: %w", err))
	}

	return nil
}

func (c *Config) View() string {
	return fmt.Sprintf(`Write key word:
%s %d
%s
Separator:
%s
%s

Write ABC:
%s %d
%s
Height: %s %s
Width:  %s %s

(ctrl+s to save configs) (ctrl+z to reload configs)
%s
`,
		c.textInputs[keyIn].View(), c.textInputs[keyIn].Position(), errorToText(textFieldValidator(c.textInputs[keyIn].Value(), "key")),
		c.textInputs[sepIn].View(), errorToText(textFieldValidator(c.textInputs[sepIn].Value(), "separator")),
		c.textInputs[abcIn].View(), c.textInputs[abcIn].Position(), errorToText(textFieldValidator(c.textInputs[abcIn].Value(), "abc")),
		c.textInputs[heightIn].View(), errorToText(numFieldValidator(c.textInputs[heightIn].Value(), "height")),
		c.textInputs[widthIn].View(), errorToText(numFieldValidator(c.textInputs[widthIn].Value(), "width")),
		c.saveRes,
	)
}

func errorToText(err error) string {
	if err == nil {
		return ""
	}

	return fmt.Sprintf("* %s", err)
}

func textFieldValidator(s, field string) error {
	if len(s) == 0 {
		return fmt.Errorf("%s must be set", field)
	}

	return nil
}

func numFieldValidator(s, field string) error {
	_, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("%s must be set", field)
	}

	return nil
}
