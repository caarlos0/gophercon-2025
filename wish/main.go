package main

import (
	"errors"
	"time"

	"github.com/charmbracelet/bubbles/v2/stopwatch"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/log/v2"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/v2"
	btm "github.com/charmbracelet/wish/v2/bubbletea"
	"github.com/charmbracelet/wish/v2/logging"
)

func main() {
	srv, err := wish.NewServer(
		wish.WithAddress("localhost:23234"),
		wish.WithHostKeyPath("./.ssh/id_ed25519"),
		wish.WithMiddleware(
			btm.Middleware(func(ssh.Session) (tea.Model, []tea.ProgramOption) {
				return newModel(), nil
			}),
			logging.StructuredMiddleware(),
		),
	)
	if err != nil {
		log.Fatal("Could not create wish server", "err", err)
	}
	log.Info("Starting", "addr", ":23234")
	if err = srv.ListenAndServe(); err != nil &&
		!errors.Is(err, ssh.ErrServerClosed) {
		log.Fatal("Could not start server", "err", err)
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
