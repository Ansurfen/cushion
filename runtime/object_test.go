package runtime

import (
	"fmt"
	"testing"

	"github.com/ansurfen/cushion/utils"
)

func TestLuaObject(t *testing.T) {
	str := NewLuaStringObject("str", "Hello World!")
	num := NewLuaNumberObject("num", "3.14")
	fmt.Println(LuaLocalVar(str))
	fmt.Println(LuaLocalVar(num))
	list := NewLuaListObject("list", LuaNumber("1.5"), LuaString("Hello"), LuaString("World"))
	fmt.Println(LuaLocalVar(list))
	tbl := NewLuaMapObject("tbl", Luamap{
		"version": LuaList(LuaString("1.7"), LuaString("1.8")),
		"images": NewLuaMapObject("", Luamap{
			"arch": LuaList(LuaString("x86"), LuaString("x64")),
		}),
	})
	tbl.Insert("describe.kind", LuaGlobalVar("ContainerKind.direct"))
	tbl.Insert("describe.kindStr", LuaString("ContainerKind.direct"))
	fmt.Println(LuaLocalVar(tbl))
	fmt.Println(LuaNil.Value(), LuaFalse.Value(), LuaTrue.Value())
	fmt.Println(LuaAssign(str), LuaAssignLR(LuaIdent("num2"), num))
}

func TestLuaTableObject(t *testing.T) {
	tbl := NewLuaTabelObject("tbl").SetMap(Luamap{
		"version": LuaList(LuaString("1.7"), LuaString("1.8")),
		"images": LuaMap(Luamap{
			"arch": LuaList(LuaString("x86"), LuaString("x64")),
		}),
	}).SetList(LuaString("Hello"), LuaString("World")).
		Insert(LuaNumber(("3.14"))).InsertKV("describe.kind", LuaGlobalVar("ContainerKind.direct"))
	fmt.Println(LuaLocalVar(tbl))
}

func TestLuaFuncObject(t *testing.T) {
	CallEvalInstallFunc := func(containers []string, order string) string {
		cons := []LuaObject{}
		for _, container := range containers {
			cons = append(cons, LuaGlobalVar(container))
		}
		return NewLuaFuncObject("EvalInstall").SetArgc(2).PCall(LuaList(cons...), LuaGlobalVar(order))
	}
	fun := NewLuaFuncObject("EvalInstall").SetArgs([]string{"container", "order"})
	fmt.Println(fun.PCall(LuaList(LuaGlobalVar("ark1"), LuaGlobalVar("ark2")), LuaTable().SetMap(Luamap{
		"debug":   LuaFalse,
		"console": LuaTrue,
		"env":     LuaMap(Luamap{}),
		"favor":   LuaMap(Luamap{}),
		"meta": LuaMap(Luamap{
			"container": LuaString("ark"),
			"alias":     LuaString("ark"),
			"dst":       LuaString("@/resource"),
		}),
	})))
	fmt.Println(LuaLocalVar(fun.Lambda()))
	fmt.Println(CallEvalInstallFunc([]string{"ark1", "ark2"}, LuaTable().SetMap(Luamap{
		"debug":   LuaFalse,
		"console": LuaTrue,
		"env":     LuaMap(Luamap{}),
		"favor":   LuaMap(Luamap{}),
		"meta": LuaMap(Luamap{
			"container": LuaString("ark"),
			"alias":     LuaString("ark"),
			"dst":       LuaString("@/resource"),
		}),
	}).Value()))
}

func TestLuaScriptGenerateByObject(t *testing.T) {
	fmt.Println(LuaProgram(
		LuaDisableFileUndefinedGlobal(),
		LuaImport("ark"),
		LuaPackagePath("../sdk.lua"),
		LuaRequire("sdk"),
		LuaGlobalVar(NewLuaFuncObject("Express").SetArgc(2).PCall(
			LuaString("Ark"),
			LuaTable().SetMap(Luamap{
				"describe": LuaMap(Luamap{
					"kind": LuaGlobalVar("ContainerKind.direct"),
					"addr": LuaString("https://github.com/ansurfen/ark"),
				}),
			}))),
	).Value())
}

type Express struct {
	Describe struct {
		Container string `yaml:"container"`
		Supplier  string `yaml:"supplier"`
	} `yaml:"describe"`
}

func TestLuaScriptGernarateByYAML(t *testing.T) {
	yaml := utils.OpenConfFromPath("ark.yaml")
	fmt.Println(yaml)
	var express Express
	if err := yaml.Unmarshal(&express); err != nil {
		panic(err)
	}
	fmt.Println(express)
}

func TestLuaConditionExpress(t *testing.T) {
	fmt.Println(LuaIf([]string{"i ~= 0", "False"}, LuaGoto("a")).Value())
	LuaIRange()
}
