package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/turnes/go-github/internal/core/ports"
)

type TuiHandler struct {
	service ports.RepositoryService
}

func NewTuiHandler(service ports.RepositoryService) *TuiHandler {
	return &TuiHandler{service: service}
}

type model struct {
	choices   []string
	cursor    int
	selected  map[int]struct{}
	confirmed bool
	selection string
}

func initialModel() model {
	return model{
		choices:  []string{"Create repo", "List repos", "Get repo", "Update repo", "Delete repo"},
		cursor:   0,
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.selection = m.choices[m.cursor]
			m.confirmed = true
			return m, tea.Quit
		case " ":
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "GitHub Repo Manager\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	s += "\nPress enter to select, q to quit.\n"
	return s
}

// Run launches the Bubble Tea TUI and prints the selected action when confirmed.
func (h *TuiHandler) Run() error {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err == nil {
			defer f.Close()
		}
	}
	p := tea.NewProgram(initialModel())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}
	if m, ok := finalModel.(model); ok && m.selection != "" {
		fmt.Println("Selected:", m.selection)
	}
	return nil
}
