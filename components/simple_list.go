package components

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type simpleListItem struct {
	v   string
	idx int
}

func (i simpleListItem) FilterValue() string { return "" }

type SimpleListStyle struct {
	Width, Height     int
	QuitTextStyle     lipgloss.Style
	TitleStyle        lipgloss.Style
	HelpStyle         lipgloss.Style
	ItemStyle         lipgloss.Style
	SelectedItemStyle lipgloss.Style
	PaginationStyle   lipgloss.Style
}

func DefaultSimpleListStyle() *SimpleListStyle {
	return &SimpleListStyle{
		Width:  20,
		Height: 14,
		TitleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1),
		HelpStyle:         list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1),
		ItemStyle:         lipgloss.NewStyle().PaddingLeft(4),
		SelectedItemStyle: lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")),
		PaginationStyle:   list.DefaultStyles().PaginationStyle.PaddingLeft(4),
		QuitTextStyle:     lipgloss.NewStyle().Margin(1, 0, 2, 4),
	}
}

func DefaultSimpleListPayload() *SimpleListPayLoad {
	return &SimpleListPayLoad{
		Title:              "",
		ShowQuitText:       true,
		QuitText:           "Null choice",
		QuitTextWithChoice: "Your choice: %s",
	}
}

type SimpleListPayLoad struct {
	Title              string
	Choices            []string
	QuitTextWithChoice string
	QuitText           string
	ShowQuitText       bool
}

type simpleListItemDelegate struct {
	style *SimpleListStyle
}

func (d simpleListItemDelegate) Height() int                               { return 1 }
func (d simpleListItemDelegate) Spacing() int                              { return 0 }
func (d simpleListItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d simpleListItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(simpleListItem)
	if !ok {
		return
	}
	str := string(i.v)
	fn := d.style.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return d.style.SelectedItemStyle.Render("> " + s[0])
		}
	}
	fmt.Fprint(w, fn(str))
}

type simpleList struct {
	list     list.Model
	choice   simpleListItem
	quitting bool
	style    *SimpleListStyle
	payload  *SimpleListPayLoad
}

func (sl *simpleList) Init() tea.Cmd {
	return nil
}

func (sl *simpleList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		sl.list.SetWidth(msg.Width)
		return sl, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case KeyCtrlC:
			sl.quitting = true
			return sl, tea.Quit
		case KeyEnter:
			i, ok := sl.list.SelectedItem().(simpleListItem)
			if ok {
				sl.choice = i
			}
			return sl, tea.Quit
		}
	}
	var cmd tea.Cmd
	sl.list, cmd = sl.list.Update(msg)
	return sl, cmd
}

func (sl *simpleList) View() string {
	if sl.choice.v != "" {
		if !sl.payload.ShowQuitText {
			return ""
		}
		return sl.style.QuitTextStyle.Render(fmt.Sprintf(sl.payload.QuitTextWithChoice, sl.choice))
	}
	if sl.quitting {
		return sl.style.QuitTextStyle.Render(sl.payload.QuitText)
	}
	return "\n" + sl.list.View()
}

func UseSimpleList(style *SimpleListStyle, payload *SimpleListPayLoad) int {
	var items []list.Item
	for idx, v := range payload.Choices {
		items = append(items, simpleListItem{v: v, idx: idx})
	}
	l := list.New(items, simpleListItemDelegate{style: style}, style.Width, style.Height)
	if len(payload.Title) > 0 {
		l.Title = payload.Title
	}
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = style.TitleStyle
	l.Styles.PaginationStyle = style.PaginationStyle
	l.Styles.HelpStyle = style.HelpStyle
	m := &simpleList{list: l, style: style, payload: payload}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		panic(err)
	}
	return m.choice.idx
}
