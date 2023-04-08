package runtime

import (
	"path"

	"github.com/ansurfen/cushion/utils"
	"github.com/vadv/gopher-lua-libs/plugin"

	lua "github.com/yuin/gopher-lua"
)

type CushionVM interface {
	Mount() CushionVM
	Call(string)
	Eval(string) error
	EvalFile(string) error
}

var _ CushionVM = &LuaVM{}

type LuaVM struct {
	state *lua.LState
	mat   *LuaMAT
}

type LuaFuncs map[string]lua.LGFunction

func NewLuaVM() *LuaVM {
	return &LuaVM{
		state: lua.NewState(),
		mat:   NewLuaMAT(),
	}
}

func (vm *LuaVM) Instance() *lua.LState {
	return vm.state
}

func (vm *LuaVM) Call(fun string) {
	if err := vm.state.CallByParam(lua.P{
		Fn:      vm.state.GetGlobal(fun),
		NRet:    0,
		Protect: true,
	}); err != nil {
		panic(err)
	}
}

func (vm *LuaVM) MAT() *LuaMAT {
	return vm.mat
}

func (vm *LuaVM) Mount() CushionVM {
	vm.mat.Mount(LuaFuncs{
		"cushion-check":   loadCheck,
		"cushion-io":      loadIO,
		"cushion-tmpl":    loadTmpl,
		"cushion-tui":     loadTui,
		"cushion-vm":      loadVM,
		"cushion-crypto":  loadCrypto,
		"cushion-time":    loadTime,
		"cushion-path":    loadPath,
		"cushion-strings": loadStrings,
	}).Collect("cushion", []string{
		"cushion-check", "cushion-io", "cushion-tmpl",
		"cushion-tui", "cushion-vm", "cushion-crypto",
		"cushion-time", "cushion-path", "cushion-strings"})
	vm.mount(LuaFuncs{
		"LoadSDK": globalLoadSDK,
		"Import":  globalImport(vm),
	})
	return vm
}

func (vm *LuaVM) mountLibs() {
	vm.mat.Mount(LuaFuncs{
		"libs-plugin": plugin.Loader,
	})
}

func Require() int {
	return 0
}

func (vm *LuaVM) LazyMount(pkgs []string) CushionVM {
	return vm
}

func (vm *LuaVM) mount(loaders LuaFuncs) {
	for name, loader := range loaders {
		vm.state.SetGlobal(name, vm.state.NewFunction(loader))
	}
}

func (vm *LuaVM) EvalFile(fullpath string) error {
	if path.Ext(fullpath) == ".lua" {
		return vm.state.DoFile(fullpath)
	}
	return nil
}

func (vm *LuaVM) Eval(script string) error {
	return vm.state.DoString(script)
}

func LuaModuleLoader(lvm *lua.LState, funcs LuaFuncs) int {
	lvm.Push(lvm.SetFuncs(lvm.NewTable(), funcs))
	return 1
}

func errHandle(lvm *lua.LState, err error) {
	if err != nil {
		lvm.Push(lua.LString(err.Error()))
	} else {
		lvm.Push(lua.LNil)
	}
}

func globalImport(vm *LuaVM) lua.LGFunction {
	return func(lvm *lua.LState) int {
		lvm.CheckTable(1).ForEach(func(idx, moudules lua.LValue) {
			mcbs := vm.mat.MCB(moudules.String())
			for mid, mcb := range mcbs {
				if !mcb.Used() {
					lvm.PreloadModule(mid, mcb.Fun())
					mcb.Mark()
				}
			}
		})
		return 0
	}
}

func globalLoadSDK(lvm *lua.LState) int {
	lvm.Push(lua.LString(path.Join(utils.GetEnv().Workdir(), "sdk", "sdk.lua")))
	return 1
}

func SafePanic() int {
	return 0
}

func LuaDoFunc(lvm *lua.LState, fun *lua.LFunction) error {
	lfunc := lvm.NewFunctionFromProto(fun.Proto)
	lvm.Push(lfunc)
	return lvm.PCall(0, lua.MultRet, nil)
}

func boolHandle(lvm *lua.LState, b bool) {
	if b {
		lvm.Push(lua.LTrue)
	} else {
		lvm.Push(lua.LFalse)
	}
}

func goStrSlice(tbl *lua.LTable) []string {
	ret := []string{}
	tbl.ForEach(func(idx, el lua.LValue) {
		ret = append(ret, el.String())
	})
	return ret
}

func strSlice2Table(str []string) *lua.LTable {
	tbl := &lua.LTable{}
	for i := 0; i < len(str); i++ {
		tbl.Insert(i+1, lua.LString(str[i]))
	}
	return tbl
}
