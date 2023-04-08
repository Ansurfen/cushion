package components

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TextInputStyle struct {
	FocusedStyle        lipgloss.Style
	BlurredStyle        lipgloss.Style
	CursorStyle         lipgloss.Style
	NoStyle             lipgloss.Style
	HelpStyle           lipgloss.Style
	CursorModeHelpStyle lipgloss.Style
	FocusedButton       string
	BlurredButton       string
}

func DefaultTextStyle() *TextInputStyle {
	focused := FontColor(lipgloss.Color("205"))
	blurred := FontColor(lipgloss.Color("240"))
	return &TextInputStyle{
		FocusedStyle:        focused,
		BlurredStyle:        blurred,
		CursorStyle:         focused.Copy(),
		NoStyle:             lipgloss.NewStyle(),
		HelpStyle:           blurred.Copy(),
		CursorModeHelpStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
		FocusedButton:       focused.Copy().Render("[ Submit ]"),
		BlurredButton:       fmt.Sprintf("[ %s ]", blurred.Render("Submit")),
	}
}

type TextInputPayLoad struct {
	Texts []TextInputFormat
}

func DefaultTextInputPayLoad() *TextInputPayLoad {
	return &TextInputPayLoad{}
}

type textInput struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	style      *TextInputStyle
}

type TextInputFormat struct {
	Name      string
	EchoMode  bool
	CharLimit int
}

func (m *textInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m *textInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyCtrlC, KeyESC:
			return m, tea.Quit
		case KeyCtrlR:
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)
		case KeyTab, KeyShiftTab, KeyEnter, KeyUp, KeyDown:
			s := msg.String()
			if s == KeyEnter && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}
			if s == KeyUp || s == KeyShiftTab {
				m.focusIndex--
			} else {
				m.focusIndex++
			}
			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = m.style.FocusedStyle
					m.inputs[i].TextStyle = m.style.FocusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = m.style.NoStyle
				m.inputs[i].TextStyle = m.style.NoStyle
			}
			return m, tea.Batch(cmds...)
		}
	}
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *textInput) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *textInput) View() string {
	var b strings.Builder
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	button := &m.style.BlurredButton
	if m.focusIndex == len(m.inputs) {
		button = &m.style.FocusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	b.WriteString(m.style.HelpStyle.Render("cursor mode is "))
	b.WriteString(m.style.CursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(m.style.HelpStyle.Render(" (ctrl+r to change style)"))
	return b.String()
}

func UseTextInput(style *TextInputStyle, payload *TextInputPayLoad) []string {
	ti := &textInput{
		inputs: make([]textinput.Model, 0),
		style:  style,
	}
	var input textinput.Model
	for i, text := range payload.Texts {
		input = textinput.New()
		input.CursorStyle = style.CursorStyle
		input.Placeholder = text.Name
		if i == 0 {
			input.Focus()
			input.PromptStyle = style.FocusedStyle
			input.TextStyle = style.FocusedStyle
		}
		if text.EchoMode {
			input.EchoMode = textinput.EchoPassword
			input.EchoCharacter = '•'
		}
		if text.CharLimit > 0 {
			input.CharLimit = text.CharLimit
		}
		ti.inputs = append(ti.inputs, input)
	}
	if _, err := tea.NewProgram(ti).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
	ret := []string{}
	for _, in := range ti.inputs {
		ret = append(ret, in.Value())
	}
	return ret
}
