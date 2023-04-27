package prompt

import (
	"strings"
	"time"

	"github.com/ansurfen/cushion/go-prompt/internal/debug"
	runewidth "github.com/mattn/go-runewidth"
)

const (
	shortenSuffix = "..."
	leftPrefix    = " "
	leftSuffix    = " "
	rightPrefix   = " "
	rightSuffix   = " "

	// CompletionMode attribute:

	// it's default value, which remain all field of suggest
	Attr_NONE = iota
	// Description in suggest will be not printed when completing
	Attr_NODSCRIPTION
	// Icon in suggest will be not printed when completing
	Attr_NOICON
	// only represent text in suggest
	Attr_OnlyText
)

var (
	leftMargin       = runewidth.StringWidth(leftPrefix + leftSuffix)
	rightMargin      = runewidth.StringWidth(rightPrefix + rightSuffix)
	completionMargin = leftMargin + rightMargin
)

// Suggest is printed when completing.
type Suggest struct {
	Icon        string
	Text        string
	Description string
	Comment     bool
}

// CompletionManager manages which suggestion is now selected.
type CompletionManager struct {
	selected  int // -1 means nothing one is selected.
	tmp       []Suggest
	max       uint16
	completer Completer
	modes     []CompletionMode

	verticalScroll int
	wordSeparator  string
	showAtStart    bool
}

// CompletionMode manage suggest list represent and mode detail
type CompletionMode struct {
	Name        string
	Description string
	Attr        uint8
}

// GetSelectedSuggestion returns the selected item.
func (c *CompletionManager) GetSelectedSuggestion() (s Suggest, ok bool) {
	if c.selected == -1 {
		return Suggest{}, false
	} else if c.selected < -1 {
		debug.Assert(false, "must not reach here")
		c.selected = -1
		return Suggest{}, false
	}
	return c.tmp[c.selected], true
}

// GetSuggestions returns the list of suggestion.
func (c *CompletionManager) GetSuggestions() []Suggest {
	return c.tmp
}

// Reset to select nothing.
func (c *CompletionManager) Reset() {
	c.selected = -1
	c.verticalScroll = 0
	c.Update(*NewDocument())
}

// Update to update the suggestions.
func (c *CompletionManager) Update(in Document) {
	c.tmp = c.completer(in)
}

// Previous to select the previous suggestion item.
func (c *CompletionManager) Previous() {
	if c.verticalScroll == c.selected && c.selected > 0 {
		c.verticalScroll--
	}
	c.selected--
	for i := c.selected; i < len(c.tmp); i++ {
		if i >= 0 && c.tmp[i].Comment {
			if c.verticalScroll == c.selected && c.selected > 0 {
				c.verticalScroll--
			}
			c.selected--
		} else {
			break
		}
	}
	c.update()
}

// Next to select the next suggestion item.
func (c *CompletionManager) Next() {
	if c.verticalScroll+int(c.max)-1 == c.selected {
		c.verticalScroll++
	}
	c.selected++
	for i := c.selected; i < len(c.tmp); i++ {
		if c.tmp[i].Comment {
			if c.verticalScroll+int(c.max)-1 == c.selected {
				c.verticalScroll++
			}
			c.selected++
		} else {
			break
		}
	}
	c.update()
}

// Completing returns whether the CompletionManager selects something one.
func (c *CompletionManager) Completing() bool {
	return c.selected != -1
}

func (c *CompletionManager) update() {
	max := int(c.max)
	if len(c.tmp) < max {
		max = len(c.tmp)
	}

	if c.selected >= len(c.tmp) {
		c.Reset()
	} else if c.selected < -1 {
		c.selected = len(c.tmp) - 1
		c.verticalScroll = len(c.tmp) - max
	}
}

// EventLoop is used to asynchronous load words of candidate.
// It is not implemented for the synchrous struct.
func (c *CompletionManager) EventLoop() {}

func (c *CompletionManager) getMax() uint16 {
	return c.max
}

func (c *CompletionManager) getVerticalScroll() int {
	return c.verticalScroll
}

func (c *CompletionManager) getTmp() []Suggest {
	return c.tmp
}

func (c *CompletionManager) getShowAtStart() bool {
	return c.showAtStart
}

func (c *CompletionManager) getSelected() int {
	return c.selected
}

func (c *CompletionManager) getCompleter() Completer {
	return c.completer
}

func (c *CompletionManager) getModes() []CompletionMode {
	return c.modes
}

func (c *CompletionManager) getWordSeparator() string {
	return c.wordSeparator
}

func (c *CompletionManager) setWordSeparator(x string) {
	c.wordSeparator = x
}

func (c *CompletionManager) setMax(x uint16) {
	c.max = x
}

func (c *CompletionManager) setShowAtStart() {
	c.showAtStart = true
}

func (c *CompletionManager) setModes(modes []CompletionMode) {
	c.modes = modes
}

func (c *CompletionManager) setPrompt(p *Prompt) {

}

func deleteBreakLineCharacters(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	return s
}

func formatTexts(o []string, max int, prefix, suffix string) (new []string, width int) {
	l := len(o)
	n := make([]string, l)

	lenPrefix := runewidth.StringWidth(prefix)
	lenSuffix := runewidth.StringWidth(suffix)
	lenShorten := runewidth.StringWidth(shortenSuffix)
	min := lenPrefix + lenSuffix + lenShorten
	for i := 0; i < l; i++ {
		o[i] = deleteBreakLineCharacters(o[i])

		w := runewidth.StringWidth(o[i])
		if width < w {
			width = w
		}
	}

	if width == 0 {
		return n, 0
	}
	if min >= max {
		return n, 0
	}
	if lenPrefix+width+lenSuffix > max {
		width = max - lenPrefix - lenSuffix
	}

	for i := 0; i < l; i++ {
		x := runewidth.StringWidth(o[i])
		if x <= width {
			spaces := strings.Repeat(" ", width-x)
			n[i] = prefix + o[i] + spaces + suffix
		} else if x > width {
			x := runewidth.Truncate(o[i], width, shortenSuffix)
			// When calling runewidth.Truncate("您好xxx您好xxx", 11, "...") returns "您好xxx..."
			// But the length of this result is 10. So we need fill right using runewidth.FillRight.
			n[i] = prefix + runewidth.FillRight(x, width) + suffix
		}
	}
	return n, lenPrefix + width + lenSuffix
}

func formatSuggetionsWithMode(suggests []Suggest, max int, mode Suggest) (new []Suggest, width int) {
	suggests = append([]Suggest{mode}, suggests...)
	return formatSuggestions(suggests, max, false)
}

func formatSuggetionsWithModeWithoutDesc(suggests []Suggest, max int, mode Suggest) (new []Suggest, width int) {
	suggests = append([]Suggest{mode}, suggests...)
	return formatSuggestions(suggests, max, true)
}

func formatSuggestions(suggests []Suggest, max int, no_desc bool) (new []Suggest, width int) {
	num := len(suggests)
	new = make([]Suggest, num)

	left := make([]string, num)
	for i := 0; i < num; i++ {
		left[i] = suggests[i].Text
	}
	var (
		right      = make([]string, num)
		rightWidth int
	)

	for i := 0; i < num; i++ {
		right[i] = suggests[i].Description
	}

	left, leftWidth := formatTexts(left, max, leftPrefix, leftSuffix)
	if leftWidth == 0 {
		return []Suggest{}, 0
	}

	if no_desc {
		right, rightWidth = formatTexts(right, 0, rightPrefix, rightSuffix)
	} else {
		right, rightWidth = formatTexts(right, max-leftWidth, rightPrefix, rightSuffix)
	}

	for i := 0; i < num; i++ {
		new[i] = Suggest{Text: left[i], Description: right[i], Comment: suggests[i].Comment, Icon: suggests[i].Icon}
	}
	return new, leftWidth + rightWidth
}

// NewCompletionManager returns initialized CompletionManager object.
func NewCompletionManager(completer Completer, max uint16) *CompletionManager {
	return &CompletionManager{
		selected:       -1,
		max:            max,
		completer:      completer,
		modes:          nil,
		verticalScroll: 0,
	}
}

const (
	ASYNC_COMPLETION_UPDATE = iota
	ASYNC_COMPLETION_RESET
)

// AsyncCompleterManager asynchronous manage which suggest is now selected.
type AsyncCompletionManager struct {
	*CompletionManager
	eventCh chan byte
	p       *Prompt
	lock    bool
}

// UpgradeAsyncCompletionManager to upgrade CompletionManager getting asynchronous suggests
func UpgradeAsyncCompletionManager(completion *CompletionManager) *AsyncCompletionManager {
	return &AsyncCompletionManager{
		CompletionManager: completion,
		lock:              false,
		eventCh:           make(chan byte, 64),
	}
}

func (c *AsyncCompletionManager) setPrompt(p *Prompt) {
	c.p = p
}

// Reset to select nothing through writting signal into event chan
func (c *AsyncCompletionManager) Reset() {
	c.eventCh <- ASYNC_COMPLETION_RESET
}

// Update to update the suggestions through writting signal into event chan
func (c *AsyncCompletionManager) Update(in Document) {
	c.eventCh <- ASYNC_COMPLETION_UPDATE
}

// EventLoop is used to asynchronous load words of candidate.
func (c *AsyncCompletionManager) EventLoop() {
	go func() { // asynchronous render suggest list
		for {
			if c.lock {
				c.tmp = []Suggest{{Text: c.p.renderer.progress.Next(), Comment: true}}
				c.p.renderer.Render(c.p.buf, c.p.completion)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	for { // asynchronous handle event
		select {
		case e := <-c.eventCh:
			switch e {
			case ASYNC_COMPLETION_UPDATE:
				if c.lock { // simulate mutex
					break // It'll downward execute when last goroutine complete
				}
				go func() {
					c.lock = true
					c.Update(*c.p.buf.Document())
					c.p.renderer.Render(c.p.buf, c.p.completion)
					c.lock = false
				}()
			case ASYNC_COMPLETION_RESET:
				c.tmp = []Suggest{{Text: c.p.renderer.progress.Next(), Comment: true}}
			}
		default:

		}
		time.Sleep(10 * time.Millisecond)
	}
}

// Completion is an interface to abstract CompletionManager.
type Completion interface {
	// Completing returns whether the CompletionManager selects something one.
	Completing() bool
	// GetSelectedSuggestion returns the selected item.
	GetSelectedSuggestion() (Suggest, bool)
	// GetSuggestions returns the list of suggestion.
	GetSuggestions() []Suggest
	// Next to select the next suggestion item.
	Next()
	// Previous to select the previous suggestion item.
	Previous()
	// Reset to select nothing.
	Reset()
	// Update to update the suggestions.
	Update(Document)
	// EventLoop is used to asynchronous load words of candidate.
	// It is not implemented for synchronous struct.
	EventLoop()

	// export field from struct
	getSelected() int
	getTmp() []Suggest
	getMax() uint16
	getCompleter() Completer
	getModes() []CompletionMode
	getVerticalScroll() int
	getWordSeparator() string
	getShowAtStart() bool
	setWordSeparator(string)
	setMax(uint16)
	setShowAtStart()
	setModes([]CompletionMode)
	setPrompt(*Prompt)
}

var (
	_ Completion = &CompletionManager{}
	_ Completion = &AsyncCompletionManager{}
)
