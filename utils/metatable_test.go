package utils

import (
	"fmt"
	"testing"
)

func TestWinMetatable(t *testing.T) {
	winTbl, _ := CreateMetaTable("./LOCAL_MACHINE/SOFTWARE/ABC")
	winTbl.SafeSetValue(MetaMap{
		"a": 10,
		"b": "10",
		"c": []string{"a", "b", "c"},
	})
	winTbl.CreateSubTable("CCC").SafeSetValue(MetaMap{})
	fmt.Println(winTbl)
}

// func TestPosixMetatable(t *testing.T) {
// 	posixTbl, _ := CreateMetaTable("./test")
// 	posixTbl.SafeSetValue(MetaMap{
// 		"a": 10,
// 		"b": "10",
// 		"c": []string{"a", "b", "c"},
// 	})
// 	posixTbl.CreateSubTable("d").SetValue(MetaArr{6, true})
// 	fmt.Println(posixTbl.fp.GetDict())
// 	posixTbl.fp.Write()
// }
