package test

import (
	"testing"

	"github.com/ansurfen/cushion/runtime"
)

func TestTui(t *testing.T) {
	vm := runtime.NewVirtualMachine().Default()
	vm.EvalFile("tui.lua")
}
