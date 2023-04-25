package utils

import "testing"

func TestEnvVar(t *testing.T) {
	env := NewEnvVar()
	env.SetPath("sys")
	env.Export("")
}
