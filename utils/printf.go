package utils


import (
	"fmt"
	"strings"
)

type PrintfOpt struct {
	MaxLen int
}

// Printf represent title and rows with tidy
func Prinf(opt PrintfOpt, title []string, rows [][]string) {
	if len(rows) <= 0 {
		for _, t := range title {
			fmt.Printf("%s ", t)
		}
		return
	}
	rowMaxLen := make([]int, len(title))
	for ri, row := range rows {
		for fi, field := range row {
			if fieldLen := len(field); opt.MaxLen <= fieldLen {
				rowMaxLen[fi] = opt.MaxLen
				rows[ri][fi] = fmt.Sprintf("%s...", field[:opt.MaxLen-3])
			} else if rowMaxLen[fi] < fieldLen {
				rowMaxLen[fi] = fieldLen
			}
		}
	}
	for ti, t := range title {
		if tLen := len(t); rowMaxLen[ti] <= tLen {
			fmt.Printf("%s ", t)
			rowMaxLen[ti] = tLen
		} else {
			fmt.Printf("%s%s ", t, strings.Repeat(" ", rowMaxLen[ti]-tLen))
		}
		if ti == len(title)-1 {
			fmt.Println()
		}
	}
	for _, row := range rows {
		for fi, field := range row {
			if fLen := len(field); rowMaxLen[fi] <= fLen {
				fmt.Printf("%s ", field)
				rowMaxLen[fi] = fLen
			} else {
				fmt.Printf("%s%s ", field, strings.Repeat(" ", rowMaxLen[fi]-fLen))
			}
			if fi == len(row) - 1 {
				fmt.Println()
			}
		}
	}
}
