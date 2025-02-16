package tab

import tea "github.com/charmbracelet/bubbletea"

type Tab interface {
	View() string
	Update(tea.Msg)
}
