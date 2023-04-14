package test

import (
	"testing"

	"github.com/ansurfen/cushion/runtime"
)

func TestTui(t *testing.T) {
	vm := runtime.NewLuaVM()
	vm.Mount()
	vm.EvalFile("tui.lua")
}
