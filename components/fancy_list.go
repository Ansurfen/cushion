package components

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FancyListStyle struct {
	DocStyle lipgloss.Style
}

func DefaultFancyListStyle() *FancyListStyle {
	return &FancyListStyle{
		DocStyle: lipgloss.NewStyle().Margin(1, 2),
	}
}

type FancyListPayLoad struct {
	Title   string
	Choices []list.Item
}

func DefaultFancyListPayLoad() *FancyListPayLoad {
	return &FancyListPayLoad{}
}

type FancyListItem struct {
	ChoiceTitle  string
	ChoiceDetial string
}

func (i FancyListItem) Title() string       { return i.ChoiceTitle }
func (i FancyListItem) Description() string { return i.ChoiceDetial }
func (i FancyListItem) FilterValue() string { return i.ChoiceTitle }

type fancyList struct {
	list   list.Model
	style  *FancyListStyle
	choice FancyListItem
}

func (fl *fancyList) Init() tea.Cmd {
	return nil
}

func (fl *fancyList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyCtrlC:
			return fl, tea.Quit
		case KeyEnter:
			i, ok := fl.list.SelectedItem().(FancyListItem)
			if ok {
				fl.choice = i
			}
			return fl, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := fl.style.DocStyle.GetFrameSize()
		fl.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	fl.list, cmd = fl.list.Update(msg)
	return fl, cmd
}

func (fl *fancyList) View() string {
	return fl.style.DocStyle.Render(fl.list.View())
}

func UseFancyList(style *FancyListStyle, payload *FancyListPayLoad) FancyListItem {
	fl := &fancyList{list: list.New(payload.Choices, list.NewDefaultDelegate(), 0, 0), style: style}
	if len(payload.Title) > 0 {
		fl.list.Title = payload.Title
	}
	if _, err := tea.NewProgram(fl, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return fl.choice
}
