package tab

import tea "github.com/charmbracelet/bubbletea"

func NewAbout() About {
	return About{}
}

var _ Tab = &About{}

type About struct{}

func (a About) Update(tea.Msg) {}

func (a About) View() string {
	return `Playfair cipher
Version: 1.1.0
Aleksandr Piter`
}
