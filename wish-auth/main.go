package main

import (
	"errors"
	"time"

	"github.com/charmbracelet/bubbles/v2/spinner"
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
	carlos, _, _, _, _ := ssh.ParseAuthorizedKey([]byte(
		"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILxWe2rXKoiO6W14LYPVfJKzRfJ1f3Jhzxrgjc/D4tU7",
	))

	srv, err := wish.NewServer(
		wish.WithAddress("0.0.0.0:23234"),
		wish.WithHostKeyPath("./.ssh/id_ed25519"),
		wish.WithPublicKeyAuth(func(_ ssh.Context, key ssh.PublicKey) bool {
			log.Info("public key")
			return ssh.KeysEqual(key, carlos)
		}),

		wish.WithPasswordAuth(func(_ ssh.Context, password string) bool {
			log.Info("password")
			return password == "how you turn this on"
		}),

		wish.WithKeyboardInteractiveAuth(func(_ ssh.Context, challenger gossh.KeyboardInteractiveChallenge) bool {
			log.Info("keyboard-interactive")
			answers, err := challenger(
				"Welcome to my server!", "Please answer these questions:",
				[]string{
					"♦ How much is 2+3: ",
					"♦ Which editor is best, vim or emacs? ",
					"♦ Tell me your best secret: ",
				},
				[]bool{true, true, false},
			)
			if err != nil {
				return false
			}
			return len(answers) == 3 &&
				answers[0] == "5" &&
				answers[1] == "vim" &&
				answers[2] != ""
		}),

		wish.WithMiddleware(
			btm.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				log.Info("ENV", "env", s.Environ())
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
		sp: spinner.New(spinner.WithSpinner(spinner.Jump)),
		sw: stopwatch.New(
			stopwatch.WithInterval(time.Second),
		),
	}
}

var _ tea.ViewModel = model{}

type model struct {
	sw       stopwatch.Model
	sp       spinner.Model
	quitting bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.sw.Start(),
		m.sp.Tick,
	)
}

var (
	byeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.BrightBlack)
	swStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Yellow).
		Bold(true).
		Italic(true)
)

var spinStyle = lipgloss.NewStyle().
	Foreground(lipgloss.BrightMagenta).
	PaddingLeft(1).
	PaddingRight(1)

func (m model) View() string {
	if m.quitting {
		return byeStyle.Render("Bye!\n")
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		spinStyle.Render(m.sp.View()),
		swStyle.Render(m.sw.View()),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyPressMsg:
		m.quitting = true
		return m, tea.Quit
	}
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.sw, cmd = m.sw.Update(msg)
	cmds = append(cmds, cmd)
	m.sp, cmd = m.sp.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
