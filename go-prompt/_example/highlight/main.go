package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ansurfen/cushion/go-prompt"
	"github.com/charmbracelet/lipgloss"
)

const (
	SYNTAX = iota
	ROUTER
)

func completer(in prompt.Document) []prompt.Suggest {
	var s []prompt.Suggest
	switch in.GetMode() {
	case SYNTAX:
		s = []prompt.Suggest{
			{Text: "syntax err", Description: `line 1 column 22 near " COLUMN DateOfBirth, DROP COLUMN ID;" `, Comment: true},
			{Text: "select", Description: "选择语句"},
			{Text: "as", Description: "为表取别名"},
			{Text: "from", Description: "选择表"},
			{Text: "表名", Description: "以下为mysql数据库存在的表", Comment: true},
			{Text: "test", Description: ""},
			{Text: "user", Description: ""},
		}
	case ROUTER:
		s = []prompt.Suggest{
			{Text: "./main.go"},
			{Text: "./go.sum"},
			{Text: "./go.mod"},
		}
	}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func main() {
	in := prompt.Input(">>> ", completer,
		prompt.OptionTitle("sql-prompt"),
		prompt.OptionRegisterMode([]prompt.CompletionMode{
			{Name: "语法模式", Attr: prompt.NONE},
			{Name: "路由模式 按住Ctrl+Y切换模式", Attr: prompt.NODSCRIPTION},
		}),
		prompt.OptionHistory([]string{"SELECT * FROM users;"}),
		prompt.OptionSuggestionTextColor(prompt.White),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionBGColor(prompt.Blue),
		prompt.OptionSelectedSuggestionTextColor(prompt.Cyan),

		prompt.OptionDescriptionTextColor(prompt.White),
		prompt.OptionDescriptionBGColor(prompt.Black),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),

		prompt.OptionScrollbarBGColor(prompt.LightGray),
		prompt.OptionScrollbarThumbColor(prompt.Cyan),
		prompt.OptionCommentSuggestionTextColor(prompt.White),
		prompt.OptionCommentSuggestionBGColor(prompt.Green),
		prompt.OptionCommentDescriptionBGColor(prompt.Green),
		prompt.OptionColor("ModePrefixTextBGColor", "#ff7b52"),
		prompt.OptionColor("ModeSuffixTextBGColor", "#ff7b52"),
		// prompt.OptionCommentSuggestionBGLipglossColor("#ff7b52"),
		// prompt.OptionCommentDescriptionBGLiglossColor("#8866e9"),
		prompt.OptionHighlight([]prompt.HighlightRule{
			{Rule: "select", Color: "99"},
			{Rule: "alter", Color: "160"},
			{Rule: "from", Color: "35"},
			{Rule: "as", Color: "202"},
			{Rule: "where", Color: "214"},
			{Rule: "(", Color: "100"},
			{Rule: ")", Color: "100"},
		}, func(s string) string {
			return strings.ToLower(s)
		}))
	fmt.Println("Your input: " + in)
	for i := 0; i < 255; i++ {
		fmt.Print(lipgloss.NewStyle().Background(lipgloss.Color(strconv.Itoa(i))).Render(" "))
	}
}
