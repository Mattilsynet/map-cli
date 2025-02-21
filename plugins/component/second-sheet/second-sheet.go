package secondsheet

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProviderName string

var (
	blurredStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	blurredNextButton     = fmt.Sprintf("[ %s ]", blurredStyle.Render("Next"))
	checkboxStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	focusedGenerateButton = focusedStyle.Render("[ Generate ]")
	blurredGenerateButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Generate"))
	focusedStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	keywordStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
)

const (
	selected bool = true

	natsCore      ProviderName = "nats-core"
	natsJetstream ProviderName = "nats-jetstream"
	natsKv        ProviderName = "nats-kv"
)

type Form struct {
	cursor       int
	Catalogue    []Provider
	lenCatalogue int
	Done         bool
}

type Provider struct {
	Name         ProviderName
	Capabilities []*Capability
}
type Capability struct {
	Name     string
	Selected bool
}

func (form *Form) Init() tea.Cmd {
	return nil
}

func New() *Form {
	form := Form{}
	form.Catalogue = []Provider{
		{
			natsCore,
			[]*Capability{{"publish", !selected}, {"subscription", !selected}, {"request/reply", !selected}},
		},
		{
			natsJetstream,
			[]*Capability{{"publish", !selected}, {"consumer", !selected}},
		},
		{
			natsKv,
			[]*Capability{{"key-value", !selected}},
		},
	}
	form.lenCatalogue = countCapabilities(form.Catalogue)
	return &form
}

func (form *Form) View() string {
	tpl := "Select capabilities \n\n"
	tpl += "%s\n\n%s\n\n"
	var choices string
	capabilityCursor := 0
	for _, provider := range form.Catalogue {
		choices += fmt.Sprintf("%s \n", keywordStyle.Render(string(provider.Name)))
		for _, capability := range provider.Capabilities {
			cursor := " "
			if capabilityCursor == form.cursor {
				cursor = ">"
			}
			choices += fmt.Sprintf("%s %s\n", cursor, checkboxTemplate(capability.Name, capability.Selected))
			capabilityCursor++
		}
	}
	button := &blurredGenerateButton
	if form.cursor == form.lenCatalogue {
		button = &focusedGenerateButton
	}
	return fmt.Sprintf(tpl, choices, *button)
}

func countCapabilities(provider []Provider) int {
	count := 0
	for _, p := range provider {
		for range p.Capabilities {
			count++
		}
	}
	return count
}

func checkboxTemplate(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

func (form *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			form.cursor++
		case "up":
			form.cursor--
		case "enter", " ":
			// INFO: we're at generate button and must now wrap up what's selected into a usable config
			if form.cursor == form.lenCatalogue {
				form.Done = true
				return nil, nil
			}
			form.checkOrUnheckCapability()
		}
	}
	if form.cursor < 0 {
		form.cursor = form.lenCatalogue
	}
	if form.cursor > form.lenCatalogue {
		form.cursor = 0
	}
	return nil, nil
}

func (form *Form) checkOrUnheckCapability() {
	count := 0
	for _, provider := range form.Catalogue {
		for _, capability := range provider.Capabilities {
			if count == form.cursor {
				capability.Selected = !capability.Selected
				return
			}
			count++
		}
	}
}
