package firstsheet

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	noStyle           = lipgloss.NewStyle()
	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle       = focusedStyle
	blurredStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	blurredNextButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Next"))
	focusedNextButton = focusedStyle.Render("[ Next ]")
)

type Form struct {
	cursor int

	Done   bool
	Inputs []textinput.Model
}

func (fw *Form) Init() tea.Cmd {
	return nil
}

func New() *Form {
	firstView := &Form{}
	firstView.Inputs = make([]textinput.Model, 3)
	firstView.cursor = 0
	var t textinput.Model
	for i := range firstView.Inputs {
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
		firstView.Inputs[i] = t
	}
	firstView.Inputs[0].SetValue("abc")
	firstView.Inputs[1].SetValue("github.com/Mattilsynet/abc")
	firstView.Inputs[2].SetValue("/home/solve/git/abc")
	return firstView
}

func (fw *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down", "enter", "up":
			s := msg.String()
			switch s {
			case "enter":
				if fw.cursor == len(fw.Inputs) {
					fw.Done = true
					return fw, nil
				} else {
					fw.cursor++
				}
			case "up":
				fw.cursor--
			case "down":
				fw.cursor++
			}
			if fw.cursor > len(fw.Inputs) {
				fw.cursor = 0
			} else if fw.cursor < 0 {
				fw.cursor = len(fw.Inputs)
			}
			if fw.cursor <= len(fw.Inputs) {
				cmds := make([]tea.Cmd, len(fw.Inputs))
				for i := 0; i <= len(fw.Inputs)-1; i++ {
					if i == fw.cursor {
						// Set focused state
						cmds[i] = fw.Inputs[i].Focus()
						fw.Inputs[i].PromptStyle = focusedStyle
						fw.Inputs[i].TextStyle = focusedStyle
						continue
					}
					// Remove focused state
					fw.Inputs[i].Blur()
					fw.Inputs[i].PromptStyle = noStyle
					fw.Inputs[i].TextStyle = noStyle
				}
				return fw, tea.Batch(cmds...)
			}
		}
	}
	cmd := fw.updateTextInputs(msg)

	return fw, cmd
}

func (fw *Form) updateTextInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(fw.Inputs))
	for i := range fw.Inputs {
		fw.Inputs[i], cmds[i] = fw.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (fw *Form) View() string {
	var b strings.Builder

	for i := range fw.Inputs {
		b.WriteString(fw.Inputs[i].View())
		if i < len(fw.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredNextButton
	if fw.cursor == len(fw.Inputs) {
		button = &focusedNextButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	return b.String()
}

func (fw *Form) Validate() string {
	pick := fw.cursor
	if pick > 0 {
		if fw.Inputs[0].Value() == "" {
			return "error: component name is empty"
		}
	}
	if pick > 1 {
		if fw.Inputs[1].Value() == "" {
			return "error: git repository is empty, e.g., github.com/Mattilsynet/my-component"
		}
	}
	if pick > 2 {
		if fw.Inputs[2].Value() == "" {
			return "error: path is empty, e.g., /home/user/git/my-component"
		} else if !filepath.IsAbs(fw.Inputs[2].Value()) {
			return "error: path is not absolute, e.g., /home/user/git/my-component"
		}
	}

	return ""
}
