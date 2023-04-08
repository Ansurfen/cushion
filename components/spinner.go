package components

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SpinnerStyle struct {
	SpinnerStyle spinner.Spinner
}

func DefaultSpinnerStyle() *SpinnerStyle {
	return &SpinnerStyle{
		SpinnerStyle: spinner.Dot,
	}
}

type SpinnerPayLoad struct {
	Callback func()
}

type resultMsg struct {
}

func (r resultMsg) String() string {
	return ""
}

type errMsg error

type spinnerComponent struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func (m spinnerComponent) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyQ, KeyESC, KeyCtrlC:
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	case resultMsg:
		m.quitting = true
		return m, tea.Quit
	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m spinnerComponent) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s Searching...press q to quit\n\n", m.spinner.View())
	if m.quitting {
		return ""
	}
	return str
}

func UseSpinner(style *SpinnerStyle, payload *SpinnerPayLoad) {
	sp := &spinnerComponent{spinner: spinner.New()}
	sp.spinner.Spinner = style.SpinnerStyle
	sp.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	p := tea.NewProgram(sp)
	go func() {
		payload.Callback()
		p.Send(resultMsg{})
	}()
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
