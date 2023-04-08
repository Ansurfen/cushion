package runtime

import (
	"fmt"
	"github.com/ansurfen/cushion/utils"
	"testing"
)

func TestExport(t *testing.T) {
	GetExport().Set(utils.RandString(8), "Ark")
	raws, formats := GetExport().Raw(), GetExport().Format()
	for i := 0; i < len(raws); i++ {
		fmt.Printf("%s -> %s\n", raws[i], formats[i])
	}
}
