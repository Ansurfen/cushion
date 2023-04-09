package prompt

import "github.com/charmbracelet/lipgloss"

type HighlightStyls map[string]lipgloss.Style

type HighlightRule struct {
	Rule  string
	Color string // ANSI color range 0-255
}

func OptionHighlight(rules []HighlightRule, cvt func(s string) string) Option {
	return func(prompt *Prompt) error {
		for _, rule := range rules {
			prompt.renderer.highlightStyle[rule.Rule] = lipgloss.NewStyle().Foreground(lipgloss.Color(rule.Color))
		}
		prompt.renderer.highlightCvt = cvt
		return nil
	}
}
