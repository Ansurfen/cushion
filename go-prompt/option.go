package prompt

import (
	"strconv"

	"github.com/ansurfen/cushion/utils"
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

// OptionPrefixTextLipglossColor change a text color of prefix string
func OptionPrefixTextLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.prefixTextColor = x
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

// OptionPrefixBackgroundLipglossColor to change a background color of prefix string
func OptionPrefixBackgroundLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.prefixBGColor = x
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

// OptionInputTextLipglossColor to change a color of text which is input by user
func OptionInputTextLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.inputTextColor = x
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

// OptionInputBGLipglossColor to change a color of background which is input by user
func OptionInputBGLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.inputBGColor = x
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

// OptionPreviewSuggestionTextLipglossColor to change a text color which is completed
func OptionPreviewSuggestionTextLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.previewSuggestionTextColor = x
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

// OptionPreviewSuggestionBGLipglossColor to change a background color which is completed
func OptionPreviewSuggestionBGLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.previewSuggestionBGColor = x
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

// OptionSuggestionTextColor to change a text color in drop down suggestions.
func OptionSuggestionTextLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.suggestionTextColor = x
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

// OptionSuggestionBGLipglossColor change a background color in drop down suggestions.
func OptionSuggestionBGLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.suggestionBGColor = x
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

// OptionSelectedSuggestionTextLipglossColor to change a text color for completed text which is selected inside suggestions drop down box.
func OptionSelectedSuggestionTextLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.selectedSuggestionTextColor = x
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

// OptionSelectedSuggestionBGLipglossColor to change a background color for completed text which is selected inside suggestions drop down box.
func OptionSelectedSuggestionBGLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.selectedSuggestionBGColor = x
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

// OptionDescriptionTextLipglossColor to change a background color of description text in drop down suggestions.
func OptionDescriptionTextLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.descriptionTextColor = x
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

// OptionDescriptionBGLipglossColor to change a background color of description text in drop down suggestions.
func OptionDescriptionBGLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.descriptionBGColor = x
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

// OptionSelectedDescriptionTextLipglossColor to change a text color of description which is selected inside suggestions drop down box.
func OptionSelectedDescriptionTextLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.selectedDescriptionTextColor = x
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

// OptionSelectedDescriptionBGLipglossColor to change a background color of description which is selected inside suggestions drop down box.
func OptionSelectedDescriptionBGLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.selectedDescriptionBGColor = x
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

// OptionScrollbarThumbLipglossColor to change a thumb color on scrollbar.
func OptionScrollbarThumbLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.scrollbarThumbColor = x
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

// OptionScrollbarBGLipglossColor to change a background color of scrollbar.
func OptionScrollbarBGLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.scrollbarBGColor = x
		return nil
	}
}

func OptionModePrefixTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.modePrefixTextColor = color2lipglossColor(x)
		return nil
	}
}

func OptionModePrefixTextLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.modePrefixTextColor = x
		return nil
	}
}

func OptionModePrefixTextBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.modePrefixTextBGColor = color2lipglossColor(x)
		return nil
	}
}

func OptionModePrefixTextBGLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.modePrefixTextBGColor = x
		return nil
	}
}

func OptionModeSuffixTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.modeSuffixTextColor = color2lipglossColor(x)
		return nil
	}
}

func OptionModeSuffixTextLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.modeSuffixTextColor = x
		return nil
	}
}

func OptionModeSuffixBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.modeSuffixTextBGColor = color2lipglossColor(x)
		return nil
	}
}

func OptionModeSuffixBGLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.modeSuffixTextBGColor = x
		return nil
	}
}

func OptionCommentSuggestionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.commentSuggestionTextColor = color2lipglossColor(x)
		return nil
	}
}

func OptionCommentSuggestionTextLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.commentSuggestionTextColor = x
		return nil
	}
}

func OptionCommentSuggestionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.commentSuggestionBGColor = color2lipglossColor(x)
		return nil
	}
}

func OptionCommentSuggestionBGLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.commentSuggestionBGColor = x
		return nil
	}
}

func OptionCommentDescriptionTextColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.commentDescriptionTextColor = color2lipglossColor(x)
		return nil
	}
}

func OptionCommentDescriptionTextLipglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.commentDescriptionTextColor = x
		return nil
	}
}

func OptionCommentDescriptionBGColor(x Color) Option {
	return func(p *Prompt) error {
		p.renderer.commentDescriptionBGColor = color2lipglossColor(x)
		return nil
	}
}

func OptionCommentDescriptionBGLiglossColor(x lipglossColor) Option {
	return func(p *Prompt) error {
		p.renderer.commentDescriptionBGColor = x
		return nil
	}
}

func OptionColor(field string, color string) Option {
	return func(p *Prompt) error {
		r := utils.NewReflectObject(p.renderer)
		r.Set(utils.FirstLower(field), color)
		return nil
	}
}

func OptionColors(colors map[string]string) Option {
	return func(p *Prompt) error {
		r := utils.NewReflectObject(p.renderer)
		for field, color := range colors {
			r.Set(utils.FirstLower(field), color)
		}
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

func OptionHistoryInstance(h *History) Option {
	return func(p *Prompt) error {
		p.history = h
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
			prefixTextColor:              color2lipglossColor(Blue),
			prefixBGColor:                color2lipglossColor(DefaultColor),
			inputTextColor:               color2lipglossColor(DefaultColor),
			inputBGColor:                 color2lipglossColor(DefaultColor),
			previewSuggestionTextColor:   color2lipglossColor(Green),
			previewSuggestionBGColor:     color2lipglossColor(DefaultColor),
			suggestionTextColor:          color2lipglossColor(White),
			suggestionBGColor:            color2lipglossColor(Cyan),
			selectedSuggestionTextColor:  color2lipglossColor(Black),
			selectedSuggestionBGColor:    color2lipglossColor(Turquoise),
			descriptionTextColor:         color2lipglossColor(DefaultColor),
			descriptionBGColor:           color2lipglossColor(Turquoise),
			selectedDescriptionTextColor: color2lipglossColor(White),
			selectedDescriptionBGColor:   color2lipglossColor(Cyan),
			scrollbarThumbColor:          color2lipglossColor(DarkGray),
			scrollbarBGColor:             color2lipglossColor(Cyan),
			highlightStyle:               make(HighlightStyles),
			modePrefixTextColor:          color2lipglossColor(DefaultColor),
			modePrefixTextBGColor:        color2lipglossColor(Purple),
			modeSuffixTextColor:          color2lipglossColor(DefaultColor),
			modeSuffixTextBGColor:        color2lipglossColor(Purple),
			commentSuggestionTextColor:   color2lipglossColor(DefaultColor),
			commentSuggestionBGColor:     color2lipglossColor(DefaultColor),
			commentDescriptionTextColor:  color2lipglossColor(DefaultColor),
			commentDescriptionBGColor:    color2lipglossColor(DefaultColor),
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
		return "-1"
	}
	return strconv.Itoa(int(c - 1))
}
