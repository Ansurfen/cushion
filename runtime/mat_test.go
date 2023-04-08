package runtime

import (
	"fmt"
	"testing"
)

func TestMAT(t *testing.T) {
	mat := NewLuaMAT().Mount(LuaFuncs{
		"cushion-check":  loadCheck,
		"cushion-io":     loadIO,
		"cushion-tmpl":   loadTmpl,
		"cushion-tui":    loadTui,
		"cushion-vm":     loadVM,
		"cushion-crypto": loadCrypto,
		"cushion-time":   loadTime,
	}).Collect("cushion", []string{
		"cushion-check", "cushion-io", "cushion-tmpl",
		"cushion-tui", "cushion-vm", "cushion-crypto", "cushion-time"})
	testDataset := map[string][]string{
		"cushion-io": {"cushion-io"},
		"cushion": {"cushion-check", "cushion-io",
			"cushion-tmpl", "cushion-tui", "cushion-vm", "cushion-crypto", "cushion-time"},
		"cushion-none": {},
	}
	flag := true
	for in, wants := range testDataset {
		mcb := mat.MCB(in)
		if len(mcb) != len(wants) {
			fmt.Println(len(mcb), len(wants), mcb, wants)
			flag = false
			break
		}
		for got := range mcb {
			have := false
			for _, want := range wants {
				if want == got {
					have = true
				}
			}
			if !have {
				flag = false
				fmt.Printf("want: %v, got: %s\n", wants, got)
				break
			}
		}
	}
	if flag {
		fmt.Println("TestMAT PASS")
	} else {
		fmt.Println("TestMAT UNPASS")
	}
}
