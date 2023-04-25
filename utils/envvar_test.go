package utils

import "testing"

func TestEnvVar(t *testing.T) {
	env := NewWinEnvVar()
	env.SetPath("sys")
	env.Export("")
}
