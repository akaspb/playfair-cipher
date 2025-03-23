package main

import (
	"log"
	"strings"

	"github.com/akaspb/playfair-cipher/internal/cipher"
	"github.com/akaspb/playfair-cipher/internal/decipher"
	"github.com/akaspb/playfair-cipher/internal/tab"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	cipherName   = "Шифрование"
	decipherName = "Расшифрование"
	configName   = "Настройки"
	aboutName    = "О программе"
)

func main() {
	if err := run(); err != nil {
		log.Fatal("Error running program:", err)
	}
}

func run() error {
	configTab := tab.NewConfig()

	cipherService, err := cipher.New(configTab.CipherCfg.GridConfig)
	if err != nil {
		return err
	}

	decipherService, err := decipher.New(configTab.CipherCfg.GridConfig)
	if err != nil {
		return err
	}

	tabNames := []string{cipherName, decipherName, configName, aboutName}
	tabs := map[string]tab.Tab{
		cipherName:   tab.NewCipher(cipherService, configTab.CipherCfg.Separator),
		decipherName: tab.NewDecipher(decipherService, configTab.CipherCfg.Separator),
		configName:   configTab,
		aboutName:    tab.NewAbout(),
	}
	m := app{CipherService: cipherService, TabNames: tabNames, Tabs: tabs}

	_, err = tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}

type app struct {
	CipherService *cipher.Cipher
	TabNames      []string
	Tabs          map[string]tab.Tab
	activeTab     int
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
	highlightColor    = lipgloss.AdaptiveColor{Light: "#bfff00", Dark: "#bfff00"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(1, 0).Align(lipgloss.Left).Border(lipgloss.NormalBorder()).UnsetBorderTop()
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

	tabRendered := windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(
		a.Tabs[a.TabNames[a.activeTab]].View(),
		"\n\n(tab - следующая вкладка)(shift+tab - предыдущая вкладка)",
		"\n               (esc - выход из программы)",
	)

	doc.WriteString(tabRendered)
	doc.WriteString("\n")

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
