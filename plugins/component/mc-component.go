package main

import (
	"fmt"
	"os"
	"strings"

	//	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Generate ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Generate"))
)

type model struct {
	focusIndex int
	inputs     []textinput.Model
	choices    []string
	selected   map[int]struct{}
	cursorMode cursor.Mode
}

func initialModel() model {
	m := model{
		inputs:   make([]textinput.Model, 2),
		choices:  []string{"Nats Jetstream", "Nats-Core", "Nats KeyValue"},
		selected: make(map[int]struct{}),
	}

	var t textinput.Model
	for i := range m.inputs {
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
			t.Placeholder = "Path to install component, leave blank for current directory"
			t.CharLimit = 64
		}
		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

			// Set focus to next input

		case "up", "shift+tab", "enter":
			if m.focusIndex < len(m.inputs) {
				return stuff(m, msg.String())
			} else {
			}
			// Handle character input and blinking
		}
	}
	cmd := m.updateInputs(msg)
	if m.focusIndex < len(m.inputs) {
		cmd := m.updateInputs(msg)
		return m, cmd
	}
	return m, cmd
}

func stuff(m model, s string) (tea.Model, tea.Cmd) {
	// Did the user press enter while the submit button was focused?
	// If so, exit.
	if s == "enter" && m.focusIndex == len(m.inputs)+len(m.choices) {
		return m, tea.Quit
	}

	// Cycle indexes
	if s == "up" || s == "shift+tab" {
		m.focusIndex--
	} else {
		m.focusIndex++
	}

	if m.focusIndex > len(m.inputs) {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = len(m.inputs)
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for i := 0; i <= len(m.inputs)-1; i++ {
		if i == m.focusIndex {
			// Set focused state
			cmds[i] = m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
			continue
		}
		// Remove focused state
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = noStyle
		m.inputs[i].TextStyle = noStyle
	}

	return m, tea.Batch(cmds...)
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *model) updateChoices(msg tea.Msg) tea.Cmd {
	return nil
}

func (m model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	for i, choice := range m.choices {
		y := i + len(m.inputs)
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.focusIndex == y {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s := fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		b.WriteString(s)
	}
	button := &blurredButton
	if m.focusIndex == len(m.inputs)+len(m.choices) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

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
			_, err := tea.NewProgram(initialModel()).Run()
			if err != nil {
				fmt.Println("error starting program:", err)
				os.Exit(1)
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
