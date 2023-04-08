package main

import (
	"fmt"

	"github.com/ansurfen/cushion/go-prompt"
)

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
		{Text: "groups", Description: "Combine users with specific rules"},
	}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func main() {
	in := prompt.Input(">>> ", completer,
		prompt.OptionTitle("sql-prompt"),
		prompt.OptionHistory([]string{"SELECT * FROM users;"}),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionHighlight([]prompt.HighlightRule{
			{Rule: "select", Color: "99"},
			{Rule: "alter", Color: "160"},
			{Rule: "from", Color: "35"},
			{Rule: "as", Color: "202"},
			{Rule: "where", Color: "214"},
			{Rule: "(", Color: "100"},
			{Rule: ")", Color: "100"},
		}))
	fmt.Println("Your input: " + in)
}
