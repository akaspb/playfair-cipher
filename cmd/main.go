package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	app "github.com/charmbracelet/bubbletea"
)

// // Model contains the program's state as well as its core functions.
// type Model interface {
// 	// Init is the first function that will be called. It returns an optional
// 	// initial command. To not perform an initial command return nil.
// 	Init() Cmd

// 	// Update is called when a message is received. Use it to inspect messages
// 	// and, in response, update the model and/or send a command.
// 	Update(Msg) (Model, Cmd)

// 	// View renders the program's UI, which is just a string. The view is
// 	// rendered after every Update.
// 	View() string
// }

type stateT uint8

const (
	stateInputKey = stateT(iota)
	stateCipherProgram
)

func main() {
	p := app.NewProgram(initialModel(), app.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	textInput textinput.Model
	err       error
	state     stateT
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Шифротекст"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() app.Cmd {
	return textinput.Blink
}

func (m model) Update(msg app.Msg) (app.Model, app.Cmd) {
	var cmd app.Cmd

	switch msg := msg.(type) {
	case app.KeyMsg:
		switch msg.Type {
		case app.KeyCtrlC, app.KeyEsc:
			return m, app.Quit
		}

	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	switch m.state {

	}
	return fmt.Sprintf(
		"Введите шифротекст\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc для выхода)",
	) + "\n"
}
