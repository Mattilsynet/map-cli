package component

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/rand"
)

type Model struct {
	pathContentList []PathContent
	RootPath        string
	index           int
	width           int
	height          int
	spinner         spinner.Model
	progress        progress.Model
	Done            bool
}

var (
	currentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle           = lipgloss.NewStyle().Margin(1, 2)
	checkMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

func NewModel(config *Config) (*Model, error) {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	pathContentList, err := GetPathContentList(config)
	if err != nil {
		return nil, err
	}
	return &Model{
		pathContentList: pathContentList,
		spinner:         s,
		progress:        p,
		index:           0,
		RootPath:        config.Path,
	}, nil
}

func (m Model) Init() tea.Cmd {
	pathContent := m.pathContentList[m.index]
	return tea.Batch(installFile(m.RootPath, pathContent.Path, pathContent.Content))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case Errored:
		return m, tea.Sequence(tea.Printf("%v", msg), tea.Quit)

	case InstalledFile:
		pathContent := m.pathContentList[m.index]
		if m.index >= len(m.pathContentList)-1 {
			// Everything's been installed. We're done!
			m.Done = true
			return m, tea.Sequence(
				tea.Printf("%s", "all done"), // print the last success message
				tea.Quit,                     // exit the program
			)
		}

		m.index++
		// Update progress bar
		progressCmd := m.progress.SetPercent(float64(m.index) / float64(len(m.pathContentList)))
		pathContentNext := m.pathContentList[m.index]
		return m, tea.Batch(
			progressCmd,
			tea.Printf("%s %s", checkMark, pathContent.Path),                       // print success message above our program
			installFile(m.RootPath, pathContentNext.Path, pathContentNext.Content), // download the next package
		)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	n := len(m.pathContentList)
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.Done {
		return doneStyle.Render(fmt.Sprintf("Done! Generated %d files.\n", n))
	}

	pkgCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, n)

	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+pkgCount))

	pkgName := currentPkgNameStyle.Render(m.pathContentList[m.index].Path)
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Generate " + pkgName)

	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+pkgCount))
	gap := strings.Repeat(" ", cellsRemaining)

	return spin + info + gap + prog + pkgCount
}

type (
	InstalledFile string
	Errored       error
)

func installFile(projectRootPath, path, content string) tea.Cmd {
	err := GenerateAndInstall(projectRootPath, path, content)
	if err != nil {
		fmt.Println("error installing file:", err)
		os.Exit(1)
	}
	d := time.Millisecond * time.Duration(rand.Intn(200))
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return InstalledFile(path)
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
