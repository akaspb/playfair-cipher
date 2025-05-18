package tab

import (
	"fmt"
	"log"
	"strconv"

	configfile "github.com/akaspb/playfair-cipher/internal/config"
	"github.com/akaspb/playfair-cipher/internal/keymatrix"
	"github.com/akaspb/playfair-cipher/internal/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type inputIdx int

const (
	keyIn inputIdx = iota
	sepIn
	abcIn
	heightIn
	widthIn
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04c77c"))
	cursorStyle  = focusedStyle
)

func NewConfig() *Config {
	key := textinput.New()
	{
		key.Placeholder = "write your key here"
		key.Prompt = "> "
		key.CharLimit = 100
		key.Width = 50
		key.EchoMode = textinput.EchoPassword
		key.EchoCharacter = '*'

		key.Cursor.Style = cursorStyle
		key.PromptStyle = focusedStyle
		key.TextStyle = focusedStyle

		key.Focus()
	}

	abc := textinput.New()
	{
		abc.Placeholder = "alphabet"
		abc.Prompt = "> "
		abc.CharLimit = 100
		abc.Width = 50

		abc.Cursor.Style = cursorStyle
		abc.PromptStyle = focusedStyle
		abc.TextStyle = focusedStyle
	}

	sep := textinput.New()
	{
		sep.Placeholder = "separator character"
		sep.Prompt = "> "
		sep.CharLimit = 1
		sep.Width = 20

		sep.Cursor.Style = cursorStyle
		sep.PromptStyle = focusedStyle
		sep.TextStyle = focusedStyle
	}

	width := textinput.New()
	{
		width.Placeholder = "XX"
		width.Prompt = ""
		width.CharLimit = 2
		width.Width = 2

		width.Cursor.Style = cursorStyle
		width.PromptStyle = focusedStyle
		width.TextStyle = focusedStyle
	}

	height := textinput.New()
	{
		height.Placeholder = "XX"
		height.Prompt = ""
		height.CharLimit = 2
		height.Width = 2

		height.Cursor.Style = cursorStyle
		height.PromptStyle = focusedStyle
		height.TextStyle = focusedStyle
	}

	c := &Config{
		textInputs: map[inputIdx]*textinput.Model{
			keyIn:    &key,
			sepIn:    &sep,
			abcIn:    &abc,
			widthIn:  &width,
			heightIn: &height,
		},
		inputIdx: 0,
	}

	if err := c.loadConfig(); err != nil {
		log.Fatal(err)
	}

	c.Done = make(chan struct{}, 1)

	return c
}

func (c *Config) loadConfig() error {
	cfg, err := configfile.LoadConfigFile()
	if err != nil {
		return err
	}

	c.textInputs[keyIn].SetValue(cfg.Key)
	c.textInputs[sepIn].SetValue(string([]rune{*cfg.Separator}))
	c.textInputs[abcIn].SetValue(string(cfg.Chars))
	c.textInputs[widthIn].SetValue(strconv.Itoa(cfg.Width))
	c.textInputs[heightIn].SetValue(strconv.Itoa(cfg.Height))

	return nil
}

var _ Tab = &Config{}

type Config struct {
	Done       chan struct{}
	textInputs map[inputIdx]*textinput.Model
	inputIdx   inputIdx
	saveRes    string
	Grid       *[][]rune
	Positions  *map[rune]model.Pos
	Separator  *rune
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
			c.saveRes = "* settings restored"
		case "ctrl+s":
			err := c.saveConfig()
			if err != nil {
				c.saveRes = fmt.Sprintf("* %s", err.Error())
			} else {
				c.saveRes = "* settings saved"
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

	if err := textFieldValidator(key, "Key"); err != nil {
		return err
	}

	if err := textFieldValidator(sep, "Separator character"); err != nil {
		return err
	}

	if err := textFieldValidator(abc, "Alphabet"); err != nil {
		return err
	}

	if err := numFieldValidator(c.textInputs[heightIn].Value(), "Matrix height"); err != nil {
		return err
	}

	if err := numFieldValidator(c.textInputs[widthIn].Value(), "Matrix width"); err != nil {
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

	cfg := model.Config{
		Height:    height,
		Width:     width,
		Chars:     []rune(abc),
		Key:       key,
		Separator: &[]rune(sep)[0],
	}

	if err := configfile.CreateConfigFile(cfg); err != nil {
		log.Fatal(fmt.Errorf("error during creating config file: %w", err))
	}

	grid, positions, err := keymatrix.Calculate(
		cfg.Chars,
		cfg.Height,
		cfg.Width,
		cfg.Key,
	)
	if err != nil {
		return fmt.Errorf("error during grid making: %w", err)
	}

	c.Grid, c.Positions, c.Separator = &grid, &positions, cfg.Separator
	c.Done <- struct{}{}

	return nil
}

func (c *Config) View() string {
	return fmt.Sprintf(`Key:
%s %d
%s
Separator character:
%s
%s

Alphabet:
%s %d
%s
Matrix height: %s %s
Matrix width:  %s %s

                (ctrl+s - save changes)
               (ctrl+z - restore settings)    
%s`,
		c.textInputs[keyIn].View(), c.textInputs[keyIn].Position(), errorToText(textFieldValidator(c.textInputs[keyIn].Value(), "Key")),
		c.textInputs[sepIn].View(), errorToText(textFieldValidator(c.textInputs[sepIn].Value(), "Separator character")),
		c.textInputs[abcIn].View(), c.textInputs[abcIn].Position(), errorToText(textFieldValidator(c.textInputs[abcIn].Value(), "Alphabet")),
		c.textInputs[heightIn].View(), errorToText(numFieldValidator(c.textInputs[heightIn].Value(), "Matrix height")),
		c.textInputs[widthIn].View(), errorToText(numFieldValidator(c.textInputs[widthIn].Value(), "Matrix width")),
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
		return fmt.Errorf("field '%s' must be settled", field)
	}

	return nil
}

func numFieldValidator(s, _ string) error {
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("value must be digital")
	}

	if num <= 0 {
		return fmt.Errorf("value must be positive")
	}

	return nil
}
