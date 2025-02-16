package tab

import tea "github.com/charmbracelet/bubbletea"

type About struct{}

func NewAbout() About {
	return About{}
}

func (a About) Update(tea.Msg) {}

func (a About) View() string {
	return `Playfair cipher
Version: 1.0.0
Piter A. V.`
}
