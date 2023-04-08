package components

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BatchSpinnerStyle struct {
	CurrentPkgNameStyle lipgloss.Style
	DoneStyle           lipgloss.Style
	CheckMark           lipgloss.Style
	ErrorMask           lipgloss.Style
}

func DefaultBatchSpinnerStyle() *BatchSpinnerStyle {
	return &BatchSpinnerStyle{
		CurrentPkgNameStyle: FontColor(lipgloss.Color("211")),
		DoneStyle:           lipgloss.NewStyle().Margin(1, 2),
		CheckMark:           FontColor(lipgloss.Color("42")).SetString("✓"),
		ErrorMask:           FontColor(lipgloss.Color("#bb1b85")).SetString("X"),
	}
}

type BatchSpinnerPayLoad struct {
	Task []BatchTask
}

func DefaultBatchSpinnerPayLoad() *BatchSpinnerPayLoad {
	return &BatchSpinnerPayLoad{}
}

type BatchTask struct {
	Name     string
	Callback func() bool //TODO: show error msg
}

type batchSpinner struct {
	packages []string
	state    bool
	index    int
	width    int
	height   int
	spinner  spinner.Model
	progress progress.Model
	done     bool
	style    *BatchSpinnerStyle
	payload  *BatchSpinnerPayLoad
}

func (m *batchSpinner) Init() tea.Cmd {
	return tea.Batch(dosomething(m.payload.Task[m.index]), m.spinner.Tick)
}

func (m *batchSpinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case KeyCtrlC, KeyESC, KeyQ:
			return m, tea.Quit
		}
	case installedPkgMsg:
		if m.index >= len(m.packages)-1 {
			m.done = true
			if msg[0] == '0' {
				m.state = false
			} else {
				m.state = true
			}
			return m, tea.Quit
		}
		var res lipgloss.Style
		if msg[0] == '0' {
			res = m.style.ErrorMask
		} else {
			res = m.style.CheckMark
		}
		progressCmd := m.progress.SetPercent(float64(m.index) / float64(len(m.packages)-1))
		m.index++
		return m, tea.Batch(
			progressCmd,
			tea.Printf("%s %s", res, m.packages[m.index-1]), // print success message above our program
			dosomething(m.payload.Task[m.index]),            // download the next package
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

func (m *batchSpinner) View() string {
	n := len(m.packages)
	w := lipgloss.Width(fmt.Sprintf("%d", n))
	if m.done {
		var res lipgloss.Style
		if !m.state {
			res = m.style.ErrorMask
		} else {
			res = m.style.CheckMark
		}
		lastTask := fmt.Sprintf("%s %s\n", res, m.packages[m.index])
		return lastTask + m.style.DoneStyle.Render(fmt.Sprintf("Done! Batched %d tasks.\n", n))
	}
	pkgCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, n)
	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+pkgCount))
	pkgName := m.style.CurrentPkgNameStyle.Render(m.packages[m.index])
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Execing " + pkgName)
	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+pkgCount))
	gap := strings.Repeat(" ", cellsRemaining)
	return spin + info + gap + prog + pkgCount
}

type installedPkgMsg string

func dosomething(task BatchTask) tea.Cmd {
	d := time.Millisecond * time.Duration(rand.Intn(1200))
	return tea.Tick(d, func(t time.Time) tea.Msg {
		var state string
		if task.Callback() {
			state = "1" + task.Name
		} else {
			state = "0" + task.Name
		}
		return installedPkgMsg(state)
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func UseBatchSpinner(style *BatchSpinnerStyle, payload *BatchSpinnerPayLoad) {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = FontColor(lipgloss.Color("63"))
	tasks := []string{}
	for _, task := range payload.Task {
		tasks = append(tasks, task.Name)
	}
	bs := &batchSpinner{
		packages: tasks,
		spinner:  s,
		progress: p,
		style:    style,
		payload:  payload,
	}
	if _, err := tea.NewProgram(bs).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
