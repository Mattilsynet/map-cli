package teaModel

import (
	"fmt"

	"github.com/Mattilsynet/map-cli/plugins/component/component-generator"
	firstsheet "github.com/Mattilsynet/map-cli/plugins/component/first-sheet"
	secondsheet "github.com/Mattilsynet/map-cli/plugins/component/second-sheet"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	dotChar = " • "
)

var (
	subtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
	dotStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle   = lipgloss.NewStyle().MarginLeft(2)
)

type Model struct {
	quitting    bool
	Finished    bool
	firstSheet  *firstsheet.Form
	secondSheet *secondsheet.Form
}

func New() *Model {
	m := Model{}
	m.firstSheet = firstsheet.New()
	m.secondSheet = secondsheet.New()
	return &m
}

func (m Model) ResultConfig() *component.Config {
	config := &component.Config{}
	config.ComponentName = m.firstSheet.Inputs[0].Value()
	config.Repository = m.firstSheet.Inputs[1].Value()
	config.Path = m.firstSheet.Inputs[2].Value()
	for _, provider := range m.secondSheet.SecondSheetCatalogue {
		for _, capability := range provider.Capabilities {
			if capability.Selected {
				// INFO: We append provider prefix such that the permutation is unique, think nats-core:publish vs nats-jetstream:publish
				config.Capabilities = append(config.Capabilities, string(provider.Name)+":"+capability.Name)
			}
		}
	}
	return config
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "ctrl+c" {
			m.quitting = true
			return &m, tea.Quit
		}
	}
	if !m.firstSheet.Done {
		return m, m.firstSheet.Update(msg)
	} else if !m.secondSheet.Done {
		return m, m.secondSheet.Update(msg)
	} else {
		if m.firstSheet.Done && m.secondSheet.Done {
			m.Finished = true
			return &m, tea.Quit
		}
	}
	fmt.Printf("Model: %v\n, error state, this should never happen!!", m)
	return m, tea.Quit
}

func (model Model) View() string {
	var s string
	if model.Finished {
		return ""
	}
	if model.quitting {
		return "\n Quitting!\n\n"
	}
	var enterSelect string
	if model.firstSheet.Done {
		enterSelect = "⏎ / _ : Select"
		s = model.secondSheet.View()
	} else {
		enterSelect = "⏎ : Select"
		s = model.firstSheet.View()
	}
	tpl := "MAP - generate a wasmcloud component\n\n"
	tpl += "%s"
	tpl += subtleStyle.Render("↑/↓ : Navigate") + dotStyle +
		subtleStyle.Render(enterSelect) + dotStyle + subtleStyle.Render("q, ctrl+c : Quit")
	// TODO: Generalize such that any view can yield a validation error
	if err := model.firstSheet.Validate(); err != "" {
		tpl += "\n\n" + errorStyle.Render(err)
	}
	return mainStyle.Render("\n" + fmt.Sprintf(tpl, s) + "\n\n")
}
