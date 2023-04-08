package utils

import "testing"

func TestPrintf(t *testing.T) {
	Prinf(PrintfOpt{
		MaxLen: 20,
	}, []string{"Container", "Supplier", "Addr", "State", "Created"}, [][]string{
		{"go", "ark", "https://github.com/ansurfen/ark", "active", "2023-01-01"},
		{"nodejs", "ark", "https://github.com/ansurfen/ark", "unknown", "2023-01-08"},
	})
}
