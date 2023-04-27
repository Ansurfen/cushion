package main

import (
	"fmt"
	"os"
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
			{Icon: "\ue704", Text: "select", Description: "选择语句"},
			{Icon: "\uf1c9", Text: "from", Description: "为表取别名"},
			{Icon: "\ue608", Text: "as", Description: "选择表"},
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

// shortenDir shortens the directory path for display in the prompt.
func shortenDir(dir string, maxDepth int) string {
	parts := splitPath(dir)
	if len(parts) <= maxDepth+1 {
		return dir
	}
	for i := range parts[:len(parts)-maxDepth-1] {
		parts[i] = "."
	}
	return joinPath(parts)
}

// splitPath splits a path into its component parts.
func splitPath(path string) []string {
	var parts []string
	for {
		d, f := splitLast(path)
		parts = append([]string{f}, parts...)
		if d == "" || d == "/" {
			break
		}
		path = d
	}
	return parts
}

// splitLast splits the last element from a path.
func splitLast(path string) (dir, file string) {
	i := len(path) - 1
	for ; i >= 0 && path[i] == '/'; i-- {
	}
	j := i
	for ; j >= 0 && path[j] != '/'; j-- {
	}
	return path[:j+1], path[j+1 : i+1]
}

// joinPath joins the component parts of a path.

func joinPath(parts []string) string {
	var path string
	if len(parts) > 0 {
		path = parts[0]
		for _, part := range parts[1:] {
			path += "/" + part
		}
	}
	return path
}

func main() {
	user := os.Getenv("USER")
	if user == "" {
		user = "unknown"
	}
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}
	dir, err := os.Getwd()
	if err != nil {
		dir = "/"
	}
	dir = shortenDir(dir, 3)
	history := prompt.NewHistory()
	cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("#66b2a5"))
	dblue := lipgloss.NewStyle().Foreground(lipgloss.Color("#2880fc"))
	white := lipgloss.NewStyle().Foreground(lipgloss.Color("#e6f2ff"))
	for {
		print(cyan.Render("┌──(") + dblue.Render(fmt.Sprintf("%s\u270E %s", user, host)) + cyan.Render(")-[") + white.Render(dir) + cyan.Render("]") + "\n" + cyan.Render("└─") + dblue.Render("$") + " ")
		in := prompt.Input("", completer,
			prompt.OptionTitle("sql-prompt"),
			prompt.OptionRegisterMode([]prompt.CompletionMode{
				{Name: "语法模式", Attr: prompt.Attr_NONE},
				{Name: "路由模式 按住Ctrl+Y切换模式", Attr: prompt.Attr_NODSCRIPTION},
			}),
			prompt.OptionHistoryInstance(history),
			// prompt.OptionHistory([]string{"SELECT * FROM users;"}),
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
		if in == "exit" {
			break
		}
	}
	for i := 0; i < 255; i++ {
		fmt.Print(lipgloss.NewStyle().Background(lipgloss.Color(strconv.Itoa(i))).Render(" "))
	}
}
