package utils

import (
	"fmt"
	"testing"

	"github.com/ansurfen/cushion/cgo"
)

func TestPlugin(t *testing.T) {
	plugin, err := NewPlugin("hulo")
	if err != nil {
		panic(err)
	}
	if pfn, err := plugin.Func("echo"); err != nil {
		fmt.Println(pfn.Call(cgo.CastVoidPtr(cgo.CStr("Hello World"))))
	}
}
