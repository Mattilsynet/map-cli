package display_example

import (
	"bytes"
	"fmt"

	"github.com/Mattilsynet/map-cli/plugins/component/component-generator"
	"github.com/alecthomas/chroma/quick"
	"github.com/alecthomas/chroma/styles"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	ConfigChange component.Config
	Model        struct {
		tmpl                    string
		viewport                viewport.Model
		renderedTemplateContent string
		Config                  *component.Config
		language                string
		commentStyle            string
	}
)

func New(filepath string, language string, vpHeight int, vpWidth int) (*Model, error) {
	commentStyle := "#"
	if language == "go" {
		commentStyle = "//"
	}
	tmpl, err := component.GetTemplate(filepath)
	if err != nil {
		return nil, err
	}

	vp := viewport.New(vpWidth, vpHeight)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	// We need to adjust the width of the glamour render from our main width
	// to account for a few things:
	//
	//  * The viewport border width
	//  * The viewport padding
	//  * The viewport margins
	//  * The gutter glamour applies to the left side of the content
	//
	const glamourGutter = 2
	return &Model{
		tmpl:         tmpl,
		viewport:     vp,
		language:     language,
		commentStyle: commentStyle,
	}, nil
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (mo *Model) Init() tea.Cmd {
	return nil
}

func (mo *Model) UpdateRenderingContent(config *component.Config) {
	var err error
	mo.renderedTemplateContent, err = component.ExecuteTmplWithData(config, mo.tmpl)
	if err != nil {
		fmt.Printf("failed to render template: %v", err)
	}
	var buffer bytes.Buffer
	err = quick.Highlight(&buffer, mo.commentStyle+mo.language+"\n"+mo.renderedTemplateContent, mo.language, "terminal256", styles.DoomOne.Name)
	if err != nil {
		fmt.Printf("failed to highlight code: %v", err)
	}
	mo.renderedTemplateContent = buffer.String()
	mo.viewport.SetContent(mo.renderedTemplateContent)
}

// Update is called when a message is received. Use it to ins:pect messages
// and, in response, update the model and/or send a command.
func (mo *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	mo.viewport, cmd = mo.viewport.Update(msg)
	return mo, cmd
}

func (mo *Model) View() string {
	return mo.viewport.View()
}
