package main

import (
	"strings"

	"github.com/ansurfen/cushion/go-prompt"
	"github.com/ansurfen/cushion/utils"
)

var tldr_cmd []string

func init() {
	out, err := utils.Exec("tldr", "-l")
	if err != nil {
		panic(err)
	}
	tldr_cmd = strings.Split(string(out), ", ")
}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	cmd := in.TextBeforeCursor()
	cmd = strings.TrimSpace(cmd)
	if len(cmd) == 0 {
		return s
	}
	flag := false
	for i := 0; i < len(tldr_cmd); i++ {
		if tldr_cmd[i] == cmd {
			flag = true
		}
	}
	if flag {
		out, err := utils.Exec("tldr", cmd)
		if err != nil {
			s = append(s, prompt.Suggest{Text: "err"})
		} else {
			s = append(s, prompt.Suggest{Text: string(out)})
		}
	}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func main() {
	for {
		prompt.Input(">> ", completer, prompt.OptionAsyncCompletionManager())
	}
}
