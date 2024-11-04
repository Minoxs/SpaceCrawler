package main

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Bindings struct {
	Quit key.Binding
}

func (b *Bindings) ShortHelp() []key.Binding {
	return []key.Binding{b.Quit}
}

func (b *Bindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		b.ShortHelp(),
	}
}

type TreeView struct {
}

func NewTreeView() *TreeView {
	return &TreeView{}
}

func (t *TreeView) Update(msg tea.Msg) (*TreeView, tea.Cmd) {
	return t, nil
}

func (t *TreeView) View() string {
	return ""
}

type SpaceCrawler struct {
	keys *Bindings
	body *TreeView
	help help.Model
}

func NewSpaceCrawler() *SpaceCrawler {

	return &SpaceCrawler{
		body: NewTreeView(),
		help: help.New(),
		keys: &Bindings{
			Quit: key.NewBinding(
				key.WithKeys("ctrl+c"),
				key.WithHelp("ctrl+c", "quit"),
			),
		},
	}
}

func (s *SpaceCrawler) Init() tea.Cmd {
	return nil
}

func (s *SpaceCrawler) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	switch event := msg.(type) {
	case tea.KeyMsg:
		switch event.String() {
		case "ctrl+c":
			return s, tea.Quit
		default:
			s.body, cmd = s.body.Update(msg)
			return s, cmd
		}
	case cursor.BlinkMsg:
		s.body, cmd = s.body.Update(msg)
		return s, cmd
	default:
		return s, nil
	}
}

func (s *SpaceCrawler) View() string {
	return "Hello, World!\n" + s.body.View() + "\n" + s.help.View(s.keys)
}

func main() {
	var (
		model = NewSpaceCrawler()
		app   = tea.NewProgram(model, tea.WithAltScreen(), tea.WithFPS(30))
	)

	if _, err := app.Run(); err != nil {
		panic(err)
	}
}
