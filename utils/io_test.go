package utils

import (
	"fmt"
	"testing"
)

func TestExec(t *testing.T) {
	out, _ := ExecStr("go help build")
	fmt.Println(ConvertByte2String(out, GB18030))
	fmt.Println(PathIsExist("./abc"))
}
