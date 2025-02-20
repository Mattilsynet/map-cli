package display_example

import (
	"bytes"
	"fmt"

	"github.com/Mattilsynet/map-cli/plugins/component/component-generator"
	"github.com/alecthomas/chroma/quick"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	ConfigChange component.Config
	Model        struct {
		renderedTemplateContent string
		tmpl                    string
	}
)

func New(filepath string) (*Model, error) {
	tmpl, err := component.GetTemplate(filepath)
	if err != nil {
		return nil, err
	}
	return &Model{
		renderedTemplateContent: "",
		tmpl:                    tmpl,
	}, nil
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (mo *Model) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to ins:pect messages
// and, in response, update the model and/or send a command.
func (mo *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case *component.Config:
		config := msg.(*component.Config)
		var err error
		mo.renderedTemplateContent, err = component.ExecuteTmplWithData(config, mo.tmpl)
		if err != nil {
			fmt.Printf("failed to render template: %v", err)
			return mo, tea.Printf("failed to render template: %v", err)
		}
		var buffer bytes.Buffer
		err = quick.Highlight(&buffer, mo.renderedTemplateContent, "go", "yaml", "monokai")
		if err != nil {
			return mo, tea.Printf("failed to highlight code: %v", err)
		}
		mo.renderedTemplateContent = buffer.String()

	case tea.QuitMsg:
		return mo, tea.Quit

	}
	return mo, nil
}

func (mo *Model) View() string {
	return mo.renderedTemplateContent
}
