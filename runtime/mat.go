package runtime

import lua "github.com/yuin/gopher-lua"

type CushionMAT interface {
	Collect(string, []string) CushionMAT
}

type CushionMCB interface {
	Mark()
	Used() bool
}

// ? modules alloc table
type LuaMAT struct {
	mcbs    map[string]*luaMCB
	cluster map[string][]string
}

// ? module control block
type luaMCB struct {
	fun  lua.LGFunction
	used bool
}

func NewLuaMCB(fun lua.LGFunction) *luaMCB {
	return &luaMCB{
		fun:  fun,
		used: false,
	}
}

func (mcb *luaMCB) Fun() lua.LGFunction {
	return mcb.fun
}

func (mcb *luaMCB) Used() bool {
	return mcb.used
}

func (mcb *luaMCB) Mark() {
	mcb.used = true
}

func NewLuaMAT() *LuaMAT {
	return &LuaMAT{
		mcbs:    make(map[string]*luaMCB),
		cluster: make(map[string][]string),
	}
}

func (mat *LuaMAT) MCB(mid string) map[string]*luaMCB {
	ret := make(map[string]*luaMCB)
	if mids, ok := mat.cluster[mid]; ok {
		for _, m := range mids {
			if mcb, ok := mat.mcbs[m]; ok {
				ret[m] = mcb
			}
		}
		return ret
	}
	if mcb, ok := mat.mcbs[mid]; ok {
		ret[mid] = mcb
		return ret
	}
	return ret
}

func (mat *LuaMAT) Mount(funcs LuaFuncs) *LuaMAT {
	for mid, mcbFunc := range funcs {
		_, ok := mat.mcbs[mid]
		if ok {
			continue
		}
		_, ok = mat.cluster[mid]
		if ok {
			continue
		}
		mat.mcbs[mid] = NewLuaMCB(mcbFunc)
	}
	return mat
}

func (mat *LuaMAT) Unmount(mid string) {
	delete(mat.mcbs, mid)
	delete(mat.cluster, mid)
}

func (mat *LuaMAT) Collect(cluster string, mids []string) *LuaMAT {
	for _, mid := range mids {
		if _, ok := mat.mcbs[mid]; ok {
			mat.cluster[cluster] = append(mat.cluster[cluster], mid)
		}
	}
	return mat
}
