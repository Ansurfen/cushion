package components

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MultiSelectStyle struct {
}

func DefaultMultiSelectStyle() *MultiSelectStyle {
	return &MultiSelectStyle{}
}

type MultiSelectPayLoad struct {
	Title   string
	Choices []string
}

func DefaultMultiSelectPayLoad() *MultiSelectPayLoad {
	return &MultiSelectPayLoad{}
}

type multiSelect struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	payload  *MultiSelectPayLoad
}

func (m multiSelect) Init() tea.Cmd {
	return nil
}

func (m multiSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyCtrlC, KeyEnter:
			return m, tea.Quit
		case KeyUp, KeyK:
			if m.cursor > 0 {
				m.cursor--
			}
		case KeyDown, KeyJ:
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case KeySpace:
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m multiSelect) View() string {
	s := m.payload.Title + "\n\n"
	var SelectedItemStyle = FontColor(lipgloss.Color(THEME_DANGER))
	var greenStyle = FontColor(lipgloss.Color(THEME_INFO))
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		if m.cursor != i {
			s += fmt.Sprintf("%s [%s] %s\n", cursor, SelectedItemStyle.Render(checked), choice)
		} else {
			s += fmt.Sprintf("%s %s%s%s %s\n", greenStyle.Render(cursor), greenStyle.Render("["), SelectedItemStyle.Render(checked), greenStyle.Render("]"), greenStyle.Render(choice))
		}
	}
	return s
}

func UseMultiSelect(style *MultiSelectStyle, payload *MultiSelectPayLoad) []int {
	ms := &multiSelect{
		choices:  payload.Choices,
		selected: make(map[int]struct{}),
		payload:  payload,
	}
	if _, err := tea.NewProgram(ms).Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	selected := []int{}
	for idx := range ms.selected {
		selected = append(selected, idx)
	}
	return selected
}
