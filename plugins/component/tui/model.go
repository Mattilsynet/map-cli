package prompt

import (
	"fmt"

	"github.com/Mattilsynet/map-cli/plugins/component/component-generator"
	project "github.com/Mattilsynet/map-cli/plugins/component/component-template"
	display_example "github.com/Mattilsynet/map-cli/plugins/component/display-example"
	firstsheet "github.com/Mattilsynet/map-cli/plugins/component/first-sheet"
	secondsheet "github.com/Mattilsynet/map-cli/plugins/component/second-sheet"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	dotChar = " • "
)

var (
	width       = 45
	height      = 45
	subtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	dotStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle   = lipgloss.NewStyle().MarginLeft(2)
	modelStyle  = lipgloss.NewStyle().
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
)

type Model struct {
	quitting       bool
	Finished       bool
	firstSheet     *firstsheet.Form
	secondSheet    *secondsheet.Form
	WadmModel      *display_example.Model
	componentModel *display_example.Model
	frameSelected  tea.Model
	tabIndex       int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func New() (*Model, error) {
	wadmModel, err := display_example.New(project.LocalWadmYamlPath, "yaml", width, height)
	if err != nil {
		return nil, err
	}
	componentModel, err := display_example.New(project.ComponentGoPath, "go", width, height)
	if err != nil {
		return nil, err
	}
	m := Model{}
	m.firstSheet = firstsheet.New()
	m.secondSheet = secondsheet.New()
	m.WadmModel = wadmModel
	m.componentModel = componentModel
	m.frameSelected = m.secondSheet
	m.tabIndex = 0
	return &m, nil
}

func (m Model) ResultConfig() *component.Config {
	config := &component.Config{}
	config.ComponentName = m.firstSheet.Inputs[0].Value()
	config.Repository = m.firstSheet.Inputs[1].Value()
	config.Path = m.firstSheet.Inputs[2].Value()
	for _, provider := range m.secondSheet.Catalogue {
		for _, capability := range provider.Capabilities {
			if capability.Selected {
				// INFO: We append provider prefix such that the permutation is unique, think nats-core:publish vs nats-jetstream:publish
				config.Capabilities = append(config.Capabilities, string(provider.Name)+":"+capability.Name)
			}
		}
	}
	return config
}

// 0 = firstSheet
// 1 = secondSheet
// 2 = WadmModel
// 3 = componentModel
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "ctrl+c" {
			m.quitting = true
			return &m, tea.Quit
		}
		if k == "tab" {
			m.tabIndex++
			if m.tabIndex == 4 {
				m.tabIndex = 0
			}

		}
		if k == "shift+tab" {
			m.tabIndex--
			if m.tabIndex == -1 {
				m.tabIndex = 3
			}
		}
		switch m.tabIndex {
		case 0:
			m.frameSelected = m.firstSheet
		case 1:
			m.frameSelected = m.secondSheet
		case 2:
			m.frameSelected = m.WadmModel
		case 3:
			m.frameSelected = m.componentModel
		}
	}
	if m.tabIndex > 0 {
		m.componentModel.UpdateRenderingContent(m.ResultConfig())
		m.WadmModel.UpdateRenderingContent(m.ResultConfig())
	}
	_, cmd := m.frameSelected.Update(msg)
	return m, cmd
	// } else {
	// 	if m.firstSheet.Done && m.secondSheet.Done {
	// 		m.Finished = true

	// 		return &m, tea.Quit
	// 	}
	// }
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
		enterSelect = "⏎ / _ : Select • tab : focus next"
		switch model.tabIndex {
		case 1:
			s += lipgloss.JoinHorizontal(lipgloss.Left, focusedModelStyle.Render(fmt.Sprintf("%4s", model.secondSheet.View())), modelStyle.Render(model.WadmModel.View()), modelStyle.Render(model.componentModel.View()))
		case 2:
			s += lipgloss.JoinHorizontal(lipgloss.Left, modelStyle.Render(fmt.Sprintf("%4s", model.secondSheet.View())), focusedModelStyle.Render(model.WadmModel.View()), modelStyle.Render(model.componentModel.View()))
		case 3:
			s += lipgloss.JoinHorizontal(lipgloss.Left, modelStyle.Render(fmt.Sprintf("%4s", model.secondSheet.View())), modelStyle.Render(model.WadmModel.View()), focusedModelStyle.Render(model.componentModel.View()))

		}
	} else {
		return model.firstSheet.View()
	}
	// if model.swapTab {
	// 	s += lipgloss.JoinVertical(lipgloss.Left,
	// 		focusedModelStyle.Render(fmt.Sprintf("%s", "wtf wtf"),
	// 			modelStyle.Render(fmt.Sprintf("%s", "generation stuff"))))
	// } else {
	// 	s += lipgloss.JoinVertical(lipgloss.Left,
	// 		modelStyle.Render(fmt.Sprintf("%25s", "ehmmm"),
	// 			focusedModelStyle.Render(fmt.Sprintf("%25s", "generation stuff"))))
	// }
	// } else {
	// 	enterSelect = "⏎ : Select"
	// 	s = model.firstSheet.View()
	// }
	tpl := "MAP - generate a wasmcloud component\n\n"
	tpl += "%s"
	tpl += subtleStyle.Render("↑/↓ : Navigate") + dotStyle +
		subtleStyle.Render(enterSelect) + dotStyle + subtleStyle.Render("q, ctrl+c : Quit")
	// TODO: Generalize such that any view can yield a validation error
	if err := model.firstSheet.Validate(); err != "" {
		tpl += "\n\n" + errorStyle.Render(err)
	}
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next %s\n", "stuff"))
	return s
	// return mainStyle.Render("\n" + fmt.Sprintf(tpl, s) + "\n\n")
}
