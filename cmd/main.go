package main

import (
	"log"
	"strings"

	"github.com/akaspb/playfair-cipher/internal/config"
	"github.com/akaspb/playfair-cipher/internal/tab"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	cipherTab = "Cipher"
	configTab = "Config"
	aboutTab  = "About"
)

func main() {
	if err := run(); err != nil {
		log.Fatal("Error running program:", err)
	}
}

func run() error {
	cfg, err := config.LoadConfigFile()
	if err != nil {
		return err
	}

	tabNames := []string{cipherTab, configTab, aboutTab}
	tabs := map[string]tab.Tab{
		cipherTab: tab.NewAbout(),
		configTab: tab.NewConfig(&cfg),
		aboutTab:  tab.NewAbout(),
	}
	m := app{TabNames: tabNames, Tabs: tabs}

	_, err = tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}

type app struct {
	TabNames  []string
	Tabs      map[string]tab.Tab
	activeTab int
}

func (a app) Init() tea.Cmd { return nil }

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc":
			return a, tea.Quit
		case "tab" /*, "right"*/ :
			a.activeTab = min(a.activeTab+1, len(a.Tabs)-1)
		case "shift+tab" /*, "left"*/ :
			a.activeTab = max(a.activeTab-1, 0)
		}
	}

	a.Tabs[a.TabNames[a.activeTab]].Update(msg)

	return a, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Left).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func (a app) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range a.TabNames {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(a.Tabs)-1, i == a.activeTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(a.Tabs[a.TabNames[a.activeTab]].View()))
	return docStyle.Render(doc.String())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// // A simple program demonstrating the text input component from the Bubbles
// // component library.

// import (
// 	"fmt"
// 	"log"

// 	"github.com/charmbracelet/bubbles/textinput"
// 	app "github.com/charmbracelet/bubbletea"
// )

// // // Model contains the program's state as well as its core functions.
// // type Model interface {
// // 	// Init is the first function that will be called. It returns an optional
// // 	// initial command. To not perform an initial command return nil.
// // 	Init() Cmd

// // 	// Update is called when a message is received. Use it to inspect messages
// // 	// and, in response, update the model and/or send a command.
// // 	Update(Msg) (Model, Cmd)

// // 	// View renders the program's UI, which is just a string. The view is
// // 	// rendered after every Update.
// // 	View() string
// // }

// type stateT uint8

// const (
// 	stateInputKey = stateT(iota)
// 	stateCipherProgram
// )

// func main() {
// 	p := app.NewProgram(initialModel(), app.WithAltScreen())
// 	if _, err := p.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// type model struct {
// 	textInput textinput.Model
// 	err       error
// 	state     stateT
// }

// func initialModel() model {
// 	ti := textinput.New()
// 	ti.Placeholder = "Шифротекст"
// 	ti.Focus()
// 	ti.CharLimit = 156
// 	ti.Width = 20

// 	return model{
// 		textInput: ti,
// 		err:       nil,
// 	}
// }

// func (m model) Init() app.Cmd {
// 	return textinput.Blink
// }

// func (m model) Update(msg app.Msg) (app.Model, app.Cmd) {
// 	var cmd app.Cmd

// 	switch msg := msg.(type) {
// 	case app.KeyMsg:
// 		switch msg.Type {
// 		case app.KeyCtrlC, app.KeyEsc:
// 			return m, app.Quit
// 		}

// 	case error:
// 		m.err = msg
// 		return m, nil
// 	}

// 	m.textInput, cmd = m.textInput.Update(msg)
// 	return m, cmd
// }

// func (m model) View() string {
// 	switch m.state {

// 	}
// 	return fmt.Sprintf(
// 		"Введите шифротекст\n\n%s\n\n%s",
// 		m.textInput.View(),
// 		"(esc для выхода)",
// 	) + "\n"
// }
