package prompt

import "github.com/charmbracelet/lipgloss"

type Highlight map[string]string

type HighlightRule struct {
	Rule  string
	Color string
}

func OptionHighlight(rules []HighlightRule) Option {
	return func(prompt *Prompt) error {
		for _, rule := range rules {
			prompt.renderer.highlight[rule.Rule] = lipgloss.NewStyle().Foreground(lipgloss.Color(rule.Color)).Render(rule.Rule)
		}
		return nil
	}
}
