package prompt

import (
	"runtime"
	"strings"

	"github.com/ansurfen/cushion/go-prompt/internal/debug"
	runewidth "github.com/mattn/go-runewidth"
)

// Render to render prompt information from state of Buffer.
type Render struct {
	out                ConsoleWriter
	prefix             string
	livePrefixCallback func() (prefix string, useLivePrefix bool)
	breakLineCallback  func(*Document)
	title              string
	row                uint16
	col                uint16
	highlightStyle     HighlightStyles
	highlightCvt       func(string) string

	previousCursor int

	// colors,
	prefixTextColor              lipglossColor
	prefixBGColor                lipglossColor
	inputTextColor               lipglossColor
	inputBGColor                 lipglossColor
	previewSuggestionTextColor   lipglossColor
	previewSuggestionBGColor     lipglossColor
	suggestionTextColor          lipglossColor
	suggestionBGColor            lipglossColor
	selectedSuggestionTextColor  lipglossColor
	selectedSuggestionBGColor    lipglossColor
	descriptionTextColor         lipglossColor
	descriptionBGColor           lipglossColor
	selectedDescriptionTextColor lipglossColor
	selectedDescriptionBGColor   lipglossColor
	scrollbarThumbColor          lipglossColor
	scrollbarBGColor             lipglossColor
	modePrefixTextColor          lipglossColor
	modePrefixTextBGColor        lipglossColor
	modeSuffixTextColor          lipglossColor
	modeSuffixTextBGColor        lipglossColor
	commentSuggestionTextColor   lipglossColor
	commentSuggestionBGColor     lipglossColor
	commentDescriptionTextColor  lipglossColor
	commentDescriptionBGColor    lipglossColor

	progress *Progress
}

const iconSize = 2

// Setup to initialize console output.
func (r *Render) Setup() {
	if r.title != "" {
		r.out.SetTitle(r.title)
		debug.AssertNoError(r.out.Flush())
	}
}

// getCurrentPrefix to get current prefix.
// If live-prefix is enabled, return live-prefix.
func (r *Render) getCurrentPrefix() string {
	if prefix, ok := r.livePrefixCallback(); ok {
		return prefix
	}
	return r.prefix
}

func (r *Render) renderPrefix() {
	// r.out.WriteColorableRawStr(r.prefixTextColor, r.prefixBGColor, false, r.getCurrentPrefix())
	r.out.WriteRawStr(r.getCurrentPrefix())
	r.out.SetColor(DefaultColor, DefaultColor, false)
}

// TearDown to clear title and erasing.
func (r *Render) TearDown() {
	r.out.ClearTitle()
	r.out.EraseDown()
	debug.AssertNoError(r.out.Flush())
}

func (r *Render) prepareArea(lines int) {
	for i := 0; i < lines; i++ {
		r.out.ScrollDown()
	}
	for i := 0; i < lines; i++ {
		r.out.ScrollUp()
	}
}

// UpdateWinSize called when window size is changed.
func (r *Render) UpdateWinSize(ws *WinSize) {
	r.row = ws.Row
	r.col = ws.Col
}

func (r *Render) renderWindowTooSmall() {
	r.out.CursorGoTo(0, 0)
	r.out.EraseScreen()
	r.out.SetColor(DarkRed, White, false)
	r.out.WriteStr("Your console window is too small...")
}

func (r *Render) renderMode(mode Suggest) {
	r.out.CursorDown(1)
	r.out.WriteColorableRawStr(r.modePrefixTextColor, r.modePrefixTextBGColor, false, mode.Text)
	r.out.WriteColorableRawStr(r.modeSuffixTextColor, r.modeSuffixTextBGColor, false, mode.Description+"   ")
	r.out.SetColor(DefaultColor, DefaultColor, false)
}

func (r *Render) renderCompletion(buf *Buffer, completions Completion) {
	suggestions := completions.GetSuggestions()
	if len(completions.GetSuggestions()) == 0 {
		return
	}
	prefix := r.getCurrentPrefix()
	var (
		formatted   []Suggest
		modeSuggest Suggest
		mode        CompletionMode
		width       int
	)
	if modes := completions.getModes(); modes != nil {
		modeNum := buf.Document().GetMode()
		mode = modes[modeNum]
		desc := mode.Description
		if len(desc) == 0 {
			desc = "Press Ctrl + Y to switch mode"
		}
		if mode.Attr == NODSCRIPTION {
			formatted, width = formatSuggetionsWithModeWithouDesc(
				suggestions,
				int(r.col)-runewidth.StringWidth(prefix)-1, // -1 means a width of scrollbar
				Suggest{Text: mode.Name, Description: desc},
			)
		} else {
			formatted, width = formatSuggetionsWithMode(
				suggestions,
				int(r.col)-runewidth.StringWidth(prefix)-1, // -1 means a width of scrollbar
				Suggest{Text: mode.Name, Description: desc},
			)
		}
		modeSuggest = formatted[0]
		formatted = formatted[1:]
	} else {
		formatted, width = formatSuggestions(suggestions,
			int(r.col)-runewidth.StringWidth(prefix)-1, // -1 means a width of scrollbar
			false,
		)
	}

	// +1 means a width of scrollbar.
	width++

	windowHeight := len(formatted)
	if windowHeight > int(completions.getMax()) {
		windowHeight = int(completions.getMax())
	}
	formatted = formatted[completions.getVerticalScroll() : completions.getVerticalScroll()+windowHeight]
	r.prepareArea(windowHeight)

	cursor := runewidth.StringWidth(prefix) + runewidth.StringWidth(buf.Document().TextBeforeCursor())
	x, _ := r.toPos(cursor)
	if x+width >= int(r.col) {
		cursor = r.backward(cursor, x+width-int(r.col))
	}

	contentHeight := len(completions.getTmp())

	fractionVisible := float64(windowHeight) / float64(contentHeight)
	fractionAbove := float64(completions.getVerticalScroll()) / float64(contentHeight)

	scrollbarHeight := int(clamp(float64(windowHeight), 1, float64(windowHeight)*fractionVisible))
	scrollbarTop := int(float64(windowHeight) * fractionAbove)

	isScrollThumb := func(row int) bool {
		return scrollbarTop <= row && row <= scrollbarTop+scrollbarHeight
	}

	if completions.getModes() != nil {
		r.renderMode(modeSuggest)
		r.lineWrap(cursor + width)
		r.backward(cursor+width+iconSize, width+iconSize)
	}

	selected := completions.getSelected() - completions.getVerticalScroll()
	r.out.SetColor(White, Cyan, false)
	icon := "  "
	for i := 0; i < windowHeight; i++ {
		r.out.CursorDown(1)
		if len(formatted[i].Icon) == 3 {
			icon = " " + formatted[i].Icon
		}
		if formatted[i].Comment {
			r.out.WriteColorableRawStr(r.commentSuggestionTextColor, r.commentSuggestionBGColor, false, icon+formatted[i].Text)
			r.out.WriteColorableRawStr(r.commentDescriptionTextColor, r.commentDescriptionBGColor, false, formatted[i].Description)
		} else {
			if i == selected {
				r.out.WriteColorableRawStr(r.selectedSuggestionTextColor, r.selectedSuggestionBGColor, true, icon+formatted[i].Text)
			} else {
				r.out.WriteColorableRawStr(r.suggestionTextColor, r.suggestionBGColor, false, icon+formatted[i].Text)
			}

			if i == selected {
				r.out.WriteColorableRawStr(r.selectedDescriptionTextColor, r.selectedDescriptionBGColor, false, formatted[i].Description)
			} else {
				r.out.WriteColorableRawStr(r.descriptionTextColor, r.descriptionBGColor, false, formatted[i].Description)
			}
		}

		if isScrollThumb(i) {
			r.out.WriteColorableRawStr(color2lipglossColor(DefaultColor), r.scrollbarThumbColor, false, " ")
		} else {
			r.out.WriteColorableRawStr(color2lipglossColor(DefaultColor), r.scrollbarBGColor, false, " ")
		}
		r.out.SetColor(DefaultColor, DefaultColor, false)
		icon = "  "
		r.lineWrap(cursor + width)
		r.backward(cursor+width+iconSize, width+iconSize)
	}

	if x+width >= int(r.col) {
		r.out.CursorForward(x + width - int(r.col))
	}
	if completions.getModes() != nil {
		r.out.CursorUp(windowHeight + 1)
	} else {
		r.out.CursorUp(windowHeight)
	}
	r.out.SetColor(DefaultColor, DefaultColor, false)
}

// Render renders to the console.
func (r *Render) Render(buffer *Buffer, completion Completion) {
	// In situations where a pseudo tty is allocated (e.g. within a docker container),
	// window size via TIOCGWINSZ is not immediately available and will result in 0,0 dimensions.
	if r.col == 0 {
		return
	}
	defer func() { debug.AssertNoError(r.out.Flush()) }()
	r.move(r.previousCursor, 0)

	line := buffer.Text()
	prefix := r.getCurrentPrefix()
	cursor := runewidth.StringWidth(prefix) + runewidth.StringWidth(line)
	// prepare area
	_, y := r.toPos(cursor)

	h := y + 1 + int(completion.getMax())
	if h > int(r.row) || completionMargin > int(r.col) {
		r.renderWindowTooSmall()
		return
	}

	// Rendering
	r.out.HideCursor()
	defer r.out.ShowCursor()

	r.renderPrefix()
	r.renderInput(line)
	r.lineWrap(cursor)

	r.out.EraseDown()

	cursor = r.backward(cursor, runewidth.StringWidth(line)-buffer.DisplayCursorPosition())

	r.renderCompletion(buffer, completion)
	if suggest, ok := completion.GetSelectedSuggestion(); ok {
		cursor = r.backward(cursor, runewidth.StringWidth(buffer.Document().GetWordBeforeCursorUntilSeparator(completion.getWordSeparator())))

		r.out.WriteColorableRawStr(r.previewSuggestionTextColor, r.previewSuggestionBGColor, false, suggest.Text)
		r.out.SetColor(DefaultColor, DefaultColor, false)
		cursor += runewidth.StringWidth(suggest.Text)

		rest := buffer.Document().TextAfterCursor()
		r.out.WriteStr(rest)
		cursor += runewidth.StringWidth(rest)
		r.lineWrap(cursor)

		cursor = r.backward(cursor, runewidth.StringWidth(rest))
	}
	r.previousCursor = cursor
}

// BreakLine to break line.
func (r *Render) BreakLine(buffer *Buffer) {
	// Erasing and Render
	cursor := runewidth.StringWidth(buffer.Document().TextBeforeCursor()) + runewidth.StringWidth(r.getCurrentPrefix())
	r.clear(cursor)
	r.renderPrefix()
	r.out.WriteColorableRawStr(r.inputTextColor, r.inputBGColor, false, buffer.Document().Text+"\n")
	r.out.SetColor(DefaultColor, DefaultColor, false)
	debug.AssertNoError(r.out.Flush())
	if r.breakLineCallback != nil {
		r.breakLineCallback(buffer.Document())
	}

	r.previousCursor = 0
}

// clear erases the screen from a beginning of input
// even if there is line break which means input length exceeds a window's width.
func (r *Render) clear(cursor int) {
	r.move(cursor, 0)
	r.out.EraseDown()
}

// backward moves cursor to backward from a current cursor position
// regardless there is a line break.
func (r *Render) backward(from, n int) int {
	return r.move(from, from-n)
}

// move moves cursor to specified position from the beginning of input
// even if there is a line break.
func (r *Render) move(from, to int) int {
	fromX, fromY := r.toPos(from)
	toX, toY := r.toPos(to)

	r.out.CursorUp(fromY - toY)
	r.out.CursorBackward(fromX - toX)
	return to
}

// toPos returns the relative position from the beginning of the string.
func (r *Render) toPos(cursor int) (x, y int) {
	col := int(r.col)
	return cursor % col, cursor / col
}

func (r *Render) lineWrap(cursor int) {
	if runtime.GOOS != "windows" && cursor > 0 && cursor%int(r.col) == 0 {
		r.out.WriteRaw([]byte{'\n'})
	}
}

func (r *Render) renderInput(line string) {
	for _, l := range strings.SplitAfter(line, " ") {
		raw := l
		if r.highlightCvt != nil {
			l = r.highlightCvt(l)
		}
		if style, ok := r.highlightStyle[strings.TrimRight(l, " ")]; ok {
			r.out.WriteRawStr(style.Render(strings.TrimRight(raw, " ")) + strings.Repeat(" ", len(l)-len(strings.TrimRight(l, " "))))
		} else {
			r.out.SetColor(DefaultColor, DefaultColor, false)
			r.out.WriteRawStr(raw)
			r.out.SetColor(DefaultColor, DefaultColor, false)
		}
	}
}

func clamp(high, low, x float64) float64 {
	switch {
	case high < x:
		return high
	case x < low:
		return low
	default:
		return x
	}
}
