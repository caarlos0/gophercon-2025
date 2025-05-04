package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/v2/spinner"
	"github.com/charmbracelet/bubbles/v2/stopwatch"
	"github.com/charmbracelet/bubbles/v2/textinput"
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
		ssh.AllocatePty(),
		wish.WithMiddleware(
			btm.Middleware(func(ssh.Session) (tea.Model, []tea.ProgramOption) {
				return newModel(), []tea.ProgramOption{
					tea.WithFilter(func(_ tea.Model, msg tea.Msg) tea.Msg {
						switch msg.(type) {
						case tea.SuspendMsg:
							return tea.ResumeMsg{}
						default:
							return msg
						}
					}),
				}
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
	ti := textinput.New()
	ti.Placeholder = "What's your name?"
	ti.Focus()
	return model{
		sp: spinner.New(spinner.WithSpinner(spinner.Jump)),
		sw: stopwatch.New(stopwatch.WithInterval(time.Second)),
		ti: ti,
	}
}

var _ tea.CursorModel = model{}

type model struct {
	sw         stopwatch.Model
	sp         spinner.Model
	ti         textinput.Model
	quitting   bool
	suspending bool
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

func (m model) View() (string, *tea.Cursor) {
	if m.quitting {
		return byeStyle.Render(fmt.Sprintf("Bye %s!\n", m.ti.Value())), nil
	}
	if m.suspending {
		return byeStyle.Render("See you soon!\n"), nil
	}

	cur := m.ti.Cursor()

	return m.ti.View() + "\n" +
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			spinStyle.Render(m.sp.View()),
			swStyle.Render(m.sw.View()),
		), cur
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ti.SetWidth(msg.Width)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "enter":
			m.quitting = true
			return m, tea.Quit
		case "ctrl+z":
			m.suspending = true
			return m, tea.Suspend
		}
	case tea.ResumeMsg:
		m.suspending = false
	}
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.sw, cmd = m.sw.Update(msg)
	cmds = append(cmds, cmd)
	m.sp, cmd = m.sp.Update(msg)
	cmds = append(cmds, cmd)
	m.ti, cmd = m.ti.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
