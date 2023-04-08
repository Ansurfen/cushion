package components

import (
	"bufio"
	"fmt"
	"os"
)

func UseRepl(questions []string) []string {
	var (
		res []string
		ask string
	)
	stdin := bufio.NewReader(os.Stdin)
	for _, question := range questions {
		ask = ""
		fmt.Print(question)
		_, err := fmt.Fscan(stdin, &ask)
		stdin.ReadString('\n')
		if err != nil {
			continue
		}
		res = append(res, ask)
	}
	return res
}
