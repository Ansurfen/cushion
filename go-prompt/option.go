package prompt

import (
	"strconv"

	"github.com/muesli/termenv"
)

// Option is the type to replace default parameters.
// prompt.New accepts any number of options (this is functional option pattern).
type Option func(prompt *Prompt) error

// OptionParser to set a custom ConsoleParser object. An argument should implement ConsoleParser interface.
func OptionParser(x ConsoleParser) Option {
	return func(p *Prompt) error {
		p.in = x
		return nil
	}
}

// OptionWriter to set a custom ConsoleWriter object. An argument should implement ConsoleWriter interface.
func OptionWriter(x ConsoleWriter) Option {
	return func(p *Prompt) error {
		registerConsoleWriter(x)
		p.renderer.out = x
		return nil
	}
}

// OptionTitle to set title displayed at the header bar of terminal.
func OptionTitle(x string) Option {
	return func(p *Prompt) error {
		p.renderer.title = x
		return nil
	}
}

// OptionPrefix to set prefix string.
func OptionPrefix(x string) Option {
	return func(p *Prompt) error {
		p.renderer.prefix = x
		return nil
	}
}

// OptionInitialBufferText to set the initial buffer text
func OptionInitialBufferText(x string) Option {
	return func(p *Prompt) error {
		p.buf.InsertText(x, false, true)
		return nil
	}
}

// OptionCompletionWordSeparator to set word separators. Enable only ' ' if empty.
func OptionCompletionWordSeparator(x string) Option {
	return func(p *Prompt) error {
		p.completion.wordSeparator = x
		return nil
	}
}

// OptionLivePrefix to change the prefix dynamically by callback function
func OptionLivePrefix(f func() (prefix string, useLivePrefix bool)) Option {
	return func(p *Prompt) error {
		p.renderer.livePrefixCallback = f
		return nil
	}
}

// OptionPrefixTextColor change a text color of prefix string
func OptionPrefixTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.prefixTextColor = color2lipglossColor(x)
		return nil
	}
}

// OptionPrefixBackgroundColor to change a background color of prefix string
func OptionPrefixBackgroundColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.prefixBGColor = color2lipglossColor(x)
		return nil
	}
}

// OptionInputTextColor to change a color of text which is input by user
func OptionInputTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.inputTextColor = color2lipglossColor(x)
		return nil
	}
}

// OptionInputBGColor to change a color of background which is input by user
func OptionInputBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.inputBGColor = color2lipglossColor(x)
		return nil
	}
}

// OptionPreviewSuggestionTextColor to change a text color which is completed
func OptionPreviewSuggestionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.previewSuggestionTextColor = color2lipglossColor(x)
		return nil
	}
}

// OptionPreviewSuggestionBGColor to change a background color which is completed
func OptionPreviewSuggestionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.previewSuggestionBGColor = color2lipglossColor(x)
		return nil
	}
}

// OptionSuggestionTextColor to change a text color in drop down suggestions.
func OptionSuggestionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.suggestionTextColor = color2lipglossColor(x)
		return nil
	}
}

// OptionSuggestionBGColor change a background color in drop down suggestions.
func OptionSuggestionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.suggestionBGColor = color2lipglossColor(x)
		return nil
	}
}

// OptionSelectedSuggestionTextColor to change a text color for completed text which is selected inside suggestions drop down box.
func OptionSelectedSuggestionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.selectedSuggestionTextColor = color2lipglossColor(x)
		return nil
	}
}

// OptionSelectedSuggestionBGColor to change a background color for completed text which is selected inside suggestions drop down box.
func OptionSelectedSuggestionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.selectedSuggestionBGColor = color2lipglossColor(x)
		return nil
	}
}

// OptionDescriptionTextColor to change a background color of description text in drop down suggestions.
func OptionDescriptionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.descriptionTextColor = color2lipglossColor(x)
		return nil
	}
}

// OptionDescriptionBGColor to change a background color of description text in drop down suggestions.
func OptionDescriptionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.descriptionBGColor = color2lipglossColor(x)
		return nil
	}
}

// OptionSelectedDescriptionTextColor to change a text color of description which is selected inside suggestions drop down box.
func OptionSelectedDescriptionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.selectedDescriptionTextColor = color2lipglossColor(x)
		return nil
	}
}

// OptionSelectedDescriptionBGColor to change a background color of description which is selected inside suggestions drop down box.
func OptionSelectedDescriptionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.selectedDescriptionBGColor = color2lipglossColor(x)
		return nil
	}
}

// OptionScrollbarThumbColor to change a thumb color on scrollbar.
func OptionScrollbarThumbColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.scrollbarThumbColor = color2lipglossColor(x)
		return nil
	}
}

// OptionScrollbarBGColor to change a background color of scrollbar.
func OptionScrollbarBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.scrollbarBGColor = color2lipglossColor(x)
		return nil
	}
}

func OptionModePrefixTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.modePrefixTextColor = color2lipglossColor(x)
		return nil
	}
}

func OptionModePrefixTtextBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.modePrefixTtextBGColor = color2lipglossColor(x)
		return nil
	}
}

func OptionModeSuffixTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.modeSuffixTextColor = color2lipglossColor(x)
		return nil
	}
}

func OptionModeSuffixBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.modeSuffixTtextBGColor = color2lipglossColor(x)
		return nil
	}
}

func OptionCommentSuggestionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.commentSuggestionTextColor = color2lipglossColor(x)
		return nil
	}
}

func OptionCommentSuggestionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.commentSuggestionBGColor = color2lipglossColor(x)
		return nil
	}
}

func OptionCommentDescriptionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.commentDescriptionTextColor = color2lipglossColor(x)
		return nil
	}
}

func OptionCommentDescriptionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.commentDescriptionBGColor = color2lipglossColor(x)
		return nil
	}
}

// OptionMaxSuggestion specify the max number of displayed suggestions.
func OptionMaxSuggestion(x uint16) Option {
	return func(p *Prompt) error {
		p.completion.max = x
		return nil
	}
}

// OptionHistory to set history expressed by string array.
func OptionHistory(x []string) Option {
	return func(p *Prompt) error {
		p.history.histories = x
		p.history.Clear()
		return nil
	}
}

// OptionSwitchKeyBindMode set a key bind mode.
func OptionSwitchKeyBindMode(m KeyBindMode) Option {
	return func(p *Prompt) error {
		p.keyBindMode = m
		return nil
	}
}

// OptionCompletionOnDown allows for Down arrow key to trigger completion.
func OptionCompletionOnDown() Option {
	return func(p *Prompt) error {
		p.completionOnDown = true
		return nil
	}
}

// SwitchKeyBindMode to set a key bind mode.
// Deprecated: Please use OptionSwitchKeyBindMode.
var SwitchKeyBindMode = OptionSwitchKeyBindMode

// OptionAddKeyBind to set a custom key bind.
func OptionAddKeyBind(b ...KeyBind) Option {
	return func(p *Prompt) error {
		p.keyBindings = append(p.keyBindings, b...)
		return nil
	}
}

// OptionAddASCIICodeBind to set a custom key bind.
func OptionAddASCIICodeBind(b ...ASCIICodeBind) Option {
	return func(p *Prompt) error {
		p.ASCIICodeBindings = append(p.ASCIICodeBindings, b...)
		return nil
	}
}

// OptionShowCompletionAtStart to set completion window is open at start.
func OptionShowCompletionAtStart() Option {
	return func(p *Prompt) error {
		p.completion.showAtStart = true
		return nil
	}
}

// OptionBreakLineCallback to run a callback at every break line
func OptionBreakLineCallback(fn func(*Document)) Option {
	return func(p *Prompt) error {
		p.renderer.breakLineCallback = fn
		return nil
	}
}

// OptionSetExitCheckerOnInput set an exit function which checks if go-prompt exits its Run loop
func OptionSetExitCheckerOnInput(fn ExitChecker) Option {
	return func(p *Prompt) error {
		p.exitChecker = fn
		return nil
	}
}

func OptionRegisterMode(modes []CompletionMode) Option {
	return func(prompt *Prompt) error {
		if len(modes) > 0 {
			prompt.completion.modes = modes
			prompt.keyBindings = append(prompt.keyBindings, KeyBind{
				Key: ControlY,
				Fn: func(b *Buffer) {
					m := (b.Document().GetMode() + 1) % len(modes)
					b.Document().SetMode(m)
					prompt.mode = m
				},
			})
		}
		return nil
	}
}

// New returns a Prompt with powerful auto-completion.
func New(executor Executor, completer Completer, opts ...Option) *Prompt {
	defaultWriter := NewStdoutWriter()
	registerConsoleWriter(defaultWriter)

	pt := &Prompt{
		in: NewStandardInputParser(),
		renderer: &Render{
			prefix:                       "> ",
			out:                          defaultWriter,
			livePrefixCallback:           func() (string, bool) { return "", false },
			prefixTextColor:              ansiHex[termenv.ANSIBlue],
			prefixBGColor:                ansiHex[termenv.ANSIBlack],
			inputTextColor:               ansiHex[termenv.ANSIBlack],
			inputBGColor:                 ansiHex[termenv.ANSIBlack],
			previewSuggestionTextColor:   ansiHex[termenv.ANSIGreen],
			previewSuggestionBGColor:     ansiHex[termenv.ANSIBlack],
			suggestionTextColor:          ansiHex[termenv.ANSIBrightWhite],
			suggestionBGColor:            ansiHex[termenv.ANSIBrightCyan],
			selectedSuggestionTextColor:  ansiHex[termenv.ANSIBlack],
			selectedSuggestionBGColor:    ansiHex[termenv.ANSIBrightCyan],
			descriptionTextColor:         ansiHex[termenv.ANSIBlack],
			descriptionBGColor:           ansiHex[termenv.ANSIBrightCyan],
			selectedDescriptionTextColor: ansiHex[termenv.ANSIBrightWhite],
			selectedDescriptionBGColor:   ansiHex[termenv.ANSIBrightCyan],
			scrollbarThumbColor:          ansiHex[termenv.ANSIBrightBlack],
			scrollbarBGColor:             ansiHex[termenv.ANSIBrightCyan],
			highlightStyle:               make(HighlightStyles),
			modePrefixTextColor:          ansiHex[termenv.ANSIBlack],
			modePrefixTtextBGColor:       ansiHex[termenv.ANSIMagenta],
			modeSuffixTextColor:          ansiHex[termenv.ANSIBlack],
			modeSuffixTtextBGColor:       ansiHex[termenv.ANSIMagenta],
			commentSuggestionTextColor:   ansiHex[termenv.ANSIBlack],
			commentSuggestionBGColor:     ansiHex[termenv.ANSIBlack],
			commentDescriptionTextColor:  ansiHex[termenv.ANSIBlack],
			commentDescriptionBGColor:    ansiHex[termenv.ANSIBlack],
		},
		buf:         NewBuffer(),
		executor:    executor,
		history:     NewHistory(),
		completion:  NewCompletionManager(completer, 6),
		keyBindMode: EmacsKeyBind, // All the above assume that bash is running in the default Emacs setting
	}

	for _, opt := range opts {
		if err := opt(pt); err != nil {
			panic(err)
		}
	}
	return pt
}

func color2lipglossColor(c Color) string {
	if c <= 0 || c > 16 {
		return "0"
	}
	return strconv.Itoa(int(c-1))
}
