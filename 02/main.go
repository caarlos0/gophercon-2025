package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/v2/stopwatch"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func newModel() model {
	return model{
		sw: stopwatch.New(
			stopwatch.WithInterval(time.Second),
		),
	}
}

var _ tea.ViewModel = model{}

type model struct {
	sw       stopwatch.Model
	quitting bool
}

func (m model) Init() tea.Cmd { return m.sw.Start() }

func (m model) View() string {
	if m.quitting {
		return lipgloss.NewStyle().
			Foreground(lipgloss.BrightBlack).
			Render("Bye!")
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Yellow).
		Bold(true).
		Italic(true).
		Render(m.sw.View())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyPressMsg:
		m.quitting = true
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.sw, cmd = m.sw.Update(msg)
	return m, cmd
}
