package cgo

import (
	"fmt"
	"testing"
)

func TestTypeCvt(t *testing.T) {
	fmt.Println(GStr(CStr("Hello World!")))
}
