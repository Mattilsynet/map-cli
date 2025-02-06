package main

import (
	"fmt"
	"os"
	"strings"

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
			model, err := tea.NewProgram(initiateModel()).Run()
			if err != nil {
				fmt.Println("error starting program:", err)
				os.Exit(1)
			}
			modelI := model.(Model)
			for index := range modelI.SelectedCapabilities {
				fmt.Println("Selected capability: ", modelI.CapabilityCatalogue[index])
			}
			for _, input := range modelI.Inputs {
				fmt.Println(input.View())
			}

			// cue validate schema
			// componentProject := config.CreateComponentProject(componentConfig)
			// config.GenerateFilesPrompt(componentProject)
		},
	}
	rootCmd.AddCommand(generate)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Model struct {
	Quitting                  bool
	NameAndPathEntered        bool
	NameAndPathCursor         int
	Inputs                    []textinput.Model
	CapabilityCatalogueCursor int
	CapabilityCatalogue       []string
	SelectedCapabilities      map[int]struct{}
	Finished                  bool
}

// start run with initiateModel
func initiateModel() Model {
	m := Model{
		Inputs:               make([]textinput.Model, 2),
		SelectedCapabilities: make(map[int]struct{}),
	}
	m.CapabilityCatalogue = []string{"Nats-core", "Nats-jetstream", "Nats-kv"}
	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32
		switch i {
		case 0:
			t.Placeholder = "Component name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Path to install (blank for cwd)"
			t.CharLimit = 64
		}
		m.Inputs[i] = t
	}
	return m
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
		return updateCapabilityCatalogue(msg, m)
	}
	return updateNameAndPath(msg, m)
}

// The main view, which just calls the appropriate sub-view
func (m Model) View() string {
	var s string
	if m.Finished {
		return "\n Done!\n\n"
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
	return mainStyle.Render("\n" + fmt.Sprintf(tpl, s) + "\n\n")
}

func updateCapabilityCatalogue(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			m.CapabilityCatalogueCursor++
		case "up":
			m.CapabilityCatalogueCursor--
		case "enter", " ":
			if m.CapabilityCatalogueCursor == len(m.CapabilityCatalogue) {
				m.Finished = true
				return m, tea.Quit
			}
			if _, ok := m.SelectedCapabilities[m.CapabilityCatalogueCursor]; ok {
				delete(m.SelectedCapabilities, m.CapabilityCatalogueCursor)
			} else {
				m.SelectedCapabilities[m.CapabilityCatalogueCursor] = struct{}{}
			}
		}
	}
	if m.CapabilityCatalogueCursor < 0 {
		m.CapabilityCatalogueCursor = len(m.CapabilityCatalogue)
	}
	if m.CapabilityCatalogueCursor > len(m.CapabilityCatalogue) {
		m.CapabilityCatalogueCursor = 0
	}
	return m, nil
}

// Update loop for the second view after a choice has been made
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
	for index, capability := range m.CapabilityCatalogue {
		cursor := " "
		if index == m.CapabilityCatalogueCursor {
			cursor = ">"
		}
		_, ok := m.SelectedCapabilities[index]
		choices += fmt.Sprintf("%s %s\n", cursor, checkbox(capability, ok))
	}
	button := &blurredGenerateButton
	if m.CapabilityCatalogueCursor == len(m.CapabilityCatalogue) {
		button = &focusedGenerateButton
	}
	return fmt.Sprintf(tpl, choices, *button)
}

// The second view, after a task has been chosen
func nameAndPathView(m Model) string {
	var b strings.Builder

	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
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
