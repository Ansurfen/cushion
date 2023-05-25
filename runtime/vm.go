package runtime

import (
	"path"

	"github.com/ansurfen/cushion/utils"

	lua "github.com/yuin/gopher-lua"
)

// VirtualMachine is an interface to abstract different interpreter
type VirtualMachine interface {
	// Default to build vm with standard
	Default() VirtualMachine
	// Call to call specify function without arguments
	Call(string) ([]any, error)
	// FastCall to call specify function without arguments and not return value
	FastCall(string) error
	// Call to call specify function with arguments
	CallByParam(string, []lua.LValue) ([]any, error)
	// FastCallByParam to call specify function with arguments and not return value
	FastCallByParam(string, []lua.LValue) error
	// Eval to execute string of script
	Eval(string) error
	// EvalFile to execute file of script
	EvalFile(string) error
	// EvalFunc to execute function
	EvalFunc(lua.LValue, []lua.LValue) ([]any, error)
	// FastEvalFunc to execute function and not return value
	FastEvalFunc(lua.LValue, []lua.LValue) error
	// SetGlobalFn to set global function
	SetGlobalFn(Handles)
	// SafeSetGlobalFn to set global function when it isn't exist
	SafeSetGlobalFn(Handles)
	// GetGlobalVar returns global variable
	GetGlobalVar(string) lua.LValue
	// SetGlobalVar to set global variable
	SetGlobalVar(string, lua.LValue)
	// SafeSetGlobalVar to set global variable when variable isn't exist
	SafeSetGlobalVar(string, lua.LValue)
	// RegisterModule to register modules
	RegisterModule(Handles)
	// UnregisterModule to unregister specify module
	UnregisterModule(string)
	// LoadModule to immediately load module to be specified
	LoadModule(string, lua.LGFunction)
	// Interp returns interpreter
	Interp() *LuaInterp
}

var _ VirtualMachine = &LuaVM{}

type LuaVM struct {
	mat   MAT
	state *LuaInterp
}

type (
	LuaFuncs map[string]lua.LGFunction
	Handles  = LuaFuncs
	// lua interpreter
	LuaInterp = lua.LState
)

func NewVirtualMachine() VirtualMachine {
	return &LuaVM{
		state: lua.NewState(),
		mat:   NewLuaMAT(),
	}
}

// Interp returns interpreter
func (vm *LuaVM) Interp() *LuaInterp {
	return vm.state
}

// Call to call specify function without arguments
func (vm *LuaVM) Call(fun string) ([]any, error) {
	ret := []any{}
	if err := vm.state.CallByParam(lua.P{
		Fn:      vm.state.GetGlobal(fun),
		NRet:    0,
		Protect: true,
	}); err != nil {
		return ret, err
	}
	for i := 1; i <= vm.state.GetTop(); i++ {
		ret = append(ret, vm.state.CheckAny(i))
	}
	return ret, nil
}

// FastCall to call specify function without arguments and not return value
func (vm *LuaVM) FastCall(fun string) error {
	return vm.state.CallByParam(lua.P{
		Fn:      vm.state.GetGlobal(fun),
		NRet:    0,
		Protect: true,
	})
}

// Call to call specify function with arguments
func (vm *LuaVM) CallByParam(fn string, args []lua.LValue) ([]any, error) {
	ret := []any{}
	if err := vm.state.CallByParam(lua.P{
		Fn:      vm.state.GetGlobal(fn),
		Protect: true,
	}, args...); err != nil {
		return ret, err
	}
	for i := 1; i <= vm.state.GetTop(); i++ {
		ret = append(ret, vm.state.CheckAny(i))
	}
	return ret, nil
}

// FastCallByParam to call specify function with arguments and not return value
func (vm *LuaVM) FastCallByParam(fn string, args []lua.LValue) error {
	return vm.state.CallByParam(lua.P{
		Fn:      vm.state.GetGlobal(fn),
		Protect: true,
	}, args...)
}

// RegisterModule to register modules
func (vm *LuaVM) RegisterModule(fns LuaFuncs) {
	vm.mat.Mount(fns)
}

// UnregisterModule to unregister specify module
func (vm *LuaVM) UnregisterModule(mid string) {
	vm.mat.Unmount(mid)
}

// LoadModule to immediately load module to be specified
func (vm *LuaVM) LoadModule(name string, loader lua.LGFunction) {
	vm.state.PreloadModule(name, loader)
}

// Default to build vm with standard
func (vm *LuaVM) Default() VirtualMachine {
	vm.mountCushion()
	vm.SetGlobalFn(LuaFuncs{
		"LoadSDK": globalLoadSDK,
		"Import":  globalImport(vm),
	})
	return vm
}

func (vm *LuaVM) mountCushion() {
	vm.mat.Mount(LuaFuncs{
		"cushion-check":   LoadCheck,
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
}

// SetGlobalFn to set global function
func (vm *LuaVM) SetGlobalFn(loaders LuaFuncs) {
	for name, loader := range loaders {
		vm.state.SetGlobal(name, vm.state.NewFunction(loader))
	}
}

// SafeSetGlobalFn to set global function when it isn't exist
func (vm *LuaVM) SafeSetGlobalFn(loaders LuaFuncs) {
	for name, loader := range loaders {
		if value := vm.state.GetGlobal(name); value.String() == "nil" {
			vm.state.SetGlobal(name, vm.state.NewFunction(loader))
		}
	}
}

// Eval to execute string of script
func (vm *LuaVM) Eval(script string) error {
	return vm.state.DoString(script)
}

// EvalFile to execute file of script
func (vm *LuaVM) EvalFile(fullpath string) error {
	if path.Ext(fullpath) == ".lua" {
		return vm.state.DoFile(fullpath)
	}
	return nil
}

// EvalFunc to execute function
func (vm *LuaVM) EvalFunc(fn lua.LValue, args []lua.LValue) ([]any, error) {
	ret := []any{}
	if err := vm.state.CallByParam(lua.P{
		Fn:      fn,
		Protect: true,
	}, args...); err != nil {
		return ret, err
	}
	for i := 1; i <= vm.state.GetTop(); i++ {
		ret = append(ret, vm.state.CheckAny(i))
	}
	return ret, nil
}

// FastEvalFunc to execute function and not return value
func (vm *LuaVM) FastEvalFunc(fn lua.LValue, args []lua.LValue) error {
	return vm.state.CallByParam(lua.P{
		Fn:      fn,
		Protect: true,
	}, args...)
}

// GetGlobalVar returns global variable
func (vm *LuaVM) GetGlobalVar(name string) lua.LValue {
	return vm.state.GetGlobal(name)
}

// SetGlobalVar to set global variable
func (vm *LuaVM) SetGlobalVar(name string, v lua.LValue) {
	vm.state.SetGlobal(name, v)
}

// SafeSetGlobalVar to set global variable when variable isn't exist
func (vm *LuaVM) SafeSetGlobalVar(name string, v lua.LValue) {
	if vm.state.GetGlobal(name).Type().String() == "nil" {
		vm.state.SetGlobal(name, v)
	}
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
