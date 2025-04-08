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
	gossh "golang.org/x/crypto/ssh"
)

func main() {
	srv, err := wish.NewServer(
		wish.WithAddress("localhost:23234"),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			// TODO: check if key is authorized
			return true
		}),
		wish.WithPasswordAuth(func(ctx ssh.Context, password string) bool {
			// TODO: check if password is valid
			return true
		}),
		wish.WithKeyboardInteractiveAuth(func(ctx ssh.Context, challenger gossh.KeyboardInteractiveChallenge) bool {
			answers, err := challenger(
				"", "",
				[]string{
					"♦ How much is 2+3: ",
					"♦ Which editor is best, vim or emacs? ",
				},
				[]bool{true, true},
			)
			if err != nil {
				return false
			}
			// here we check for the correct answers:
			return len(answers) == 2 && answers[0] == "5" && answers[1] == "vim"
		}),
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
		stopwatch.New(
			stopwatch.WithInterval(time.Second),
		),
	}
}

var _ tea.ViewModel = model{}

type model struct {
	sw stopwatch.Model
}

func (m model) Init() tea.Cmd { return m.sw.Start() }

func (m model) View() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Yellow).
		Bold(true).
		Italic(true).
		Render(m.sw.View())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyPressMsg:
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.sw, cmd = m.sw.Update(msg)
	return m, cmd
}
