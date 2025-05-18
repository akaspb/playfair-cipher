package main

import (
	"log"
	"strings"
	"sync/atomic"

	"github.com/akaspb/playfair-cipher/internal/cipher"
	"github.com/akaspb/playfair-cipher/internal/decipher"
	"github.com/akaspb/playfair-cipher/internal/tab"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	cipherName   = "  Cipher  "
	decipherName = "   Decipher  "
	configName   = "  Settings "
	aboutName    = "  About  "
)

func main() {
	if err := run(); err != nil {
		log.Fatal("Error running program:", err)
	}
}

func run() error {
	configTab := tab.NewConfig()

	tabNames := []string{cipherName, decipherName, configName, aboutName}
	tabs := map[string]tab.Tab{
		cipherName:   nil,
		decipherName: nil,
		configName:   configTab,
		aboutName:    tab.NewAbout(),
	}

	a := &app{TabNames: tabNames, Tabs: tabs, ActiveTab: 2, ConfigSettled: atomic.Bool{}}

	go func() {
		for range configTab.Done {
			cipherService, err := cipher.New(configTab.Grid, configTab.Positions)
			if err != nil {
				log.Print(err.Error())
				continue
			}

			decipherService, err := decipher.New(configTab.Grid, configTab.Positions)
			if err != nil {
				log.Print(err.Error())
				continue
			}

			tabs[cipherName] = tab.NewCipher(cipherService, configTab.Separator)
			tabs[decipherName] = tab.NewDecipher(decipherService, configTab.Separator)

			a.ConfigSettled.Store(true)
		}
	}()

	_, err := tea.NewProgram(a, tea.WithAltScreen()).Run()

	close(configTab.Done)

	return err
}

type app struct {
	TabNames      []string
	Tabs          map[string]tab.Tab
	ActiveTab     int
	ConfigSettled atomic.Bool
}

func (a *app) Init() tea.Cmd { return nil }

func (a *app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	configSettled := a.ConfigSettled.Load()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			return a, tea.Quit
		case "tab":
			if !configSettled {
				break
			}
			a.ActiveTab = min(a.ActiveTab+1, len(a.Tabs)-1)
		case "shift+tab":
			if !configSettled {
				break
			}
			a.ActiveTab = max(a.ActiveTab-1, 0)
		}
	}

	a.Tabs[a.TabNames[a.ActiveTab]].Update(msg)

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

func (a *app) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range a.TabNames {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(a.Tabs)-1, i == a.ActiveTab
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
		a.Tabs[a.TabNames[a.ActiveTab]].View(),
		"\n\n       (tab - next tab)(shift+tab - previous tab)",
		"\n                   (esc - exit program)",
	)

	doc.WriteString(tabRendered)
	doc.WriteString("\n")

	return docStyle.Render(doc.String())
}
