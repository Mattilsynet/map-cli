package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Mattilsynet/map-cli/plugins/component/component"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const (
	dotChar = " • "
)

// General stuff for styling the view
var (
	keywordStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	errorStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
	ticksStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	checkboxStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	focusedStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle           = focusedStyle
	blurredStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	helpStyle             = blurredStyle
	dotStyle              = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	noStyle               = lipgloss.NewStyle()
	mainStyle             = lipgloss.NewStyle().MarginLeft(2)
	focusedNextButton     = focusedStyle.Render("[ Next ]")
	blurredNextButton     = fmt.Sprintf("[ %s ]", blurredStyle.Render("Next"))
	focusedGenerateButton = focusedStyle.Render("[ Generate ]")
	blurredGenerateButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Generate"))
	cursorModeHelpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "component",
		Short:   "Component plugin",
		Aliases: []string{"c"},
	}
	generate := &cobra.Command{
		Use:     "generate",
		Short:   "Generate a WasmCloud component",
		Aliases: []string{"gen", "g"},
		Run: func(cmd *cobra.Command, args []string) {
			m, err := tea.NewProgram(initiateConfigPrompt()).Run()
			if err != nil {
				fmt.Println("error starting program:", err)
				os.Exit(1)
			}
			model := m.(*Model)
			component.GenerateApp(model.Config)
		},
	}
	rootCmd.AddCommand(generate)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type provider struct {
	name         string
	capabilities []*capability
}
type capability struct {
	name     string
	selected bool
}

type Model struct {
	Quitting                  bool
	NameAndPathEntered        bool
	NameAndPathCursor         int
	Inputs                    []textinput.Model
	CapabilityCatalogueCursor int
	ProviderCatalogue         []provider
	LenCapabilities           int
	Finished                  bool
	Config                    *component.Config
}

const (
	notSelected bool = false
	selected    bool = true
)

// start run with initiateConfigPrompt
func initiateConfigPrompt() *Model {
	m := Model{}
	// TODO: Fetch options from componentConfig instead of duplicating places where capabilities are mentioned
	m.ProviderCatalogue = []provider{
		{
			"nats-core",
			[]*capability{{"publish", notSelected}, {"subscription", notSelected}, {"request/reply", notSelected}},
		},
		{
			"nats-jetstream",
			[]*capability{{"publish", notSelected}, {"consumer", notSelected}},
		},
		{
			"nats-kv",
			[]*capability{{"key-value", notSelected}},
		},
	}
	m.LenCapabilities = countCapabilities(m.ProviderCatalogue)
	m.Inputs = make([]textinput.Model, 3)
	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32
		switch i {
		case 0:
			t.Placeholder = "Component name"
			t.SetSuggestions([]string{"my-component", "map-managed-environment"})
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Git repository"
			t.CharLimit = 64
			t.SetSuggestions([]string{"github.com/Mattilsynet/map-managed-environment"})
		case 2:
			t.Placeholder = "Path to install (blank for cwd)"
			t.SetSuggestions([]string{"/home/my-user/git/my-component/"})
			t.CharLimit = 64
		}
		m.Inputs[i] = t
	}
	return &m
}

func countCapabilities(provider []provider) int {
	count := 0
	for _, p := range provider {
		for range p.capabilities {
			count++
		}
	}
	return count
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Main update function.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}
	if m.NameAndPathEntered {
		return updateCapabilityCatalogue(msg, &m)
	}
	return updateNameAndPath(msg, m)
}

// TODO: Make capabilitiesView and namedAndPathView same style, remove cursor display in capabilitiesView and rather highlight same as namedAndPathView
func (m Model) View() string {
	var s string
	if m.Finished {
		return ""
	}
	if m.Quitting {
		return "\n Quitting!\n\n"
	}
	var enterSelect string
	if m.NameAndPathEntered {
		enterSelect = "⏎ / _ : Select"
		s = capabilitiesView(m)
	} else {
		enterSelect = "⏎ : Select"
		s = nameAndPathView(m)
	}
	tpl := "MAP - generate a wasmcloud component\n\n"
	tpl += "%s"
	tpl += subtleStyle.Render("↑/↓ : Navigate") + dotStyle +
		subtleStyle.Render(enterSelect) + dotStyle + subtleStyle.Render("q, ctrl+c : Quit")
	if err := validationErr(m); err != "" {
		tpl += "\n\n" + errorStyle.Render(err)
	}
	return mainStyle.Render("\n" + fmt.Sprintf(tpl, s) + "\n\n")
}

func validationErr(m Model) string {
	pick := m.NameAndPathCursor
	if pick > 0 {
		if m.Inputs[0].Value() == "" {
			return "error: component name is empty"
		}
	}
	if pick > 1 {
		if m.Inputs[1].Value() == "" {
			return "error: git repository is empty, e.g., github.com/Mattilsynet/my-component"
		}
	}
	if pick > 2 {
		if m.Inputs[2].Value() == "" {
			return "error: path is empty, e.g., /home/user/git/my-component"
		} else if !filepath.IsAbs(m.Inputs[2].Value()) {
			return "error: path is not absolute, e.g., /home/user/git/my-component"
		}
	}

	return ""
}

func updateCapabilityCatalogue(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			m.CapabilityCatalogueCursor++
		case "up":
			m.CapabilityCatalogueCursor--
		case "enter", " ":
			// INFO: we're at generate button and must now wrap up what's selected into a usable config
			if m.CapabilityCatalogueCursor == m.LenCapabilities {
				config := &component.Config{}
				config.ComponentName = m.Inputs[0].Value()
				config.Path = m.Inputs[1].Value()
				for _, provider := range m.ProviderCatalogue {
					for _, capability := range provider.capabilities {
						if capability.selected {
							// INFO: We append provider prefix such that the permutation is unique, think nats-core:publish vs nats-jetstream:publish
							config.Capabilities = append(config.Capabilities, provider.name+":"+capability.name)
						}
					}
				}
				/*  TODO:
				        2. generate go files according to selected capabilities, name of component and path to put them
				        3. Generate wit files
					4. Generate wadm files
					5. Generate wasmcloud.toml
				        6. add sdk files from custom capabilities to go.mod
					7. Generate working <component-name>-component.go with implemented requirements according to capabilities, e.g., handle etc
					8. Generate README.md
					9. Generate LICENSE

					Feature: Add fancy loading bar
					Feature: Add fancy display of files generated in which folder
				*/
				m.Finished = true
				m.Config = config
				return m, tea.Quit
			}
			checkOrUnheck(m)
		}
	}
	if m.CapabilityCatalogueCursor < 0 {
		m.CapabilityCatalogueCursor = m.LenCapabilities
	}
	if m.CapabilityCatalogueCursor > m.LenCapabilities {
		m.CapabilityCatalogueCursor = 0
	}
	return m, nil
}

func checkOrUnheck(m *Model) {
	count := 0
	for _, provider := range m.ProviderCatalogue {
		for _, capability := range provider.capabilities {
			if count == m.CapabilityCatalogueCursor {
				capability.selected = !capability.selected
				return
			}
			count++
		}
	}
}

// TODO: Make path give help on how to navigate, i.e., ctrl + space yields cwd and picker
func updateNameAndPath(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down", "enter", "up":
			s := msg.String()
			if s == "enter" && m.NameAndPathCursor == len(m.Inputs) {
				m.NameAndPathEntered = true
				return m, nil

			}
			if s == "up" {
				m.NameAndPathCursor--
			} else {
				m.NameAndPathCursor++
			}
			if m.NameAndPathCursor > len(m.Inputs) {
				m.NameAndPathCursor = 0
			} else if m.NameAndPathCursor < 0 {
				m.NameAndPathCursor = len(m.Inputs)
			}
			if m.NameAndPathCursor <= len(m.Inputs) {
				cmds := make([]tea.Cmd, len(m.Inputs))
				for i := 0; i <= len(m.Inputs)-1; i++ {
					if i == m.NameAndPathCursor {
						// Set focused state
						cmds[i] = m.Inputs[i].Focus()
						m.Inputs[i].PromptStyle = focusedStyle
						m.Inputs[i].TextStyle = focusedStyle
						continue
					}
					// Remove focused state
					m.Inputs[i].Blur()
					m.Inputs[i].PromptStyle = noStyle
					m.Inputs[i].TextStyle = noStyle
				}
				return m, tea.Batch(cmds...)
			}
		}
	}
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func capabilitiesView(m Model) string {
	tpl := "Select capabilities \n\n"
	tpl += "%s\n\n%s\n\n"
	var choices string
	capabilityCursor := 0
	for _, provider := range m.ProviderCatalogue {
		choices += fmt.Sprintf("%s \n", keywordStyle.Render(provider.name))
		for _, capability := range provider.capabilities {
			cursor := " "
			if capabilityCursor == m.CapabilityCatalogueCursor {
				cursor = ">"
			}
			choices += fmt.Sprintf("%s %s\n", cursor, checkbox(capability.name, capability.selected))
			capabilityCursor++
		}
	}
	button := &blurredGenerateButton
	if m.CapabilityCatalogueCursor == m.LenCapabilities {
		button = &focusedGenerateButton
	}
	return fmt.Sprintf(tpl, choices, *button)
}

// The second view, after a task has been chosen
func nameAndPathView(m Model) string {
	var b strings.Builder

	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		// TODO: put in examples when user presses 'h' to show underneath prompt text lines
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredNextButton
	if m.NameAndPathCursor == len(m.Inputs) {
		button = &focusedNextButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	// b.WriteString(helpStyle.Render("cursor is: " + strconv.Itoa(m.NameAndPathCursor)))
	return b.String()
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}
