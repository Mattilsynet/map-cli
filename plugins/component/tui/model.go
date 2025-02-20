package prompt

import (
	"fmt"

	"github.com/Mattilsynet/map-cli/plugins/component/component-generator"
	project "github.com/Mattilsynet/map-cli/plugins/component/component-template"
	display_example "github.com/Mattilsynet/map-cli/plugins/component/display-example"
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
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	dotStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle   = lipgloss.NewStyle().MarginLeft(2)
	modelStyle  = lipgloss.NewStyle().
			Width(45).
			Height(45).
			Align(lipgloss.Left, lipgloss.Left).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(45).
				Height(45).
				Align(lipgloss.Left, lipgloss.Left).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
)

type Model struct {
	quitting       bool
	Finished       bool
	firstSheet     *firstsheet.Form
	secondSheet    *secondsheet.Form
	wadmModel      display_example.Model
	componentModel display_example.Model
	swapTab        bool
}

func New() (*Model, error) {
	wadmModel, err := display_example.New(project.LocalWadmYamlPath)
	if err != nil {
		return nil, err
	}
	componentModel, err := display_example.New(project.ComponentGoPath)
	if err != nil {
		return nil, err
	}
	m := Model{}
	m.firstSheet = firstsheet.New()
	m.secondSheet = secondsheet.New()
	m.wadmModel = *wadmModel
	m.componentModel = *componentModel
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

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "ctrl+c" {
			m.quitting = true
			return &m, tea.Quit
		}
		if k == "tab" {
			m.swapTab = !m.swapTab
		}
	}
	// if !m.firstSheet.Done {
	//	return m, m.firstSheet.Update(msg)
	// } else if !m.secondSheet.Done {
	cmd := m.secondSheet.Update(msg)
	m.componentModel.Update(m.ResultConfig())
	m.wadmModel.Update(m.ResultConfig())
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
	// if model.firstSheet.Done {
	enterSelect = "⏎ / _ : Select • tab : focus next"
	if model.swapTab {
		s += lipgloss.JoinHorizontal(lipgloss.Left, focusedModelStyle.Render("yo man bbobo tatjakdsdjaskdjas dksaj dkasj dkasj dksa jdask djask"), modelStyle.Render("hfasfjaijdasidsaj idasj idas jidasj disa jdias jdasi djasi dasji"))
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Left, modelStyle.Render(fmt.Sprintf("%4s", model.secondSheet.View())), focusedModelStyle.Render(fmt.Sprintf("//local.wadm.yaml\n\n%s", model.wadmModel.View())))
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
