package runtime

import lua "github.com/yuin/gopher-lua"

// MAT is abbreviation for module allocate table,
// which is used to load module into virtual machine with lazy.
type MAT interface {
	// Mount to add MCB to MAT
	Mount(LuaFuncs) MAT
	// Unmount to remove MCB from MAT
	Unmount(string) MAT
	// Collect to converge specify mcb according to cluster.
	Collect(string, []string) MAT
	// MCB returns mcb list
	MCB(string) map[string]*luaMCB
}

// MCB is abbreviation for module control block,
// which store and control module meta.
type MCB interface {
	// Mark to set module used state
	Mark()
	// Used returns whether mcb is used
	Used() bool
}

var (
	_ MAT = &LuaMAT{}
	_ MCB = &luaMCB{}
)

// MAT is abbreviation for module allocate table,
// which is used to load module into virtual machine with lazy.
type LuaMAT struct {
	mcbs    map[string]*luaMCB
	cluster map[string][]string
}

func NewLuaMAT() MAT {
	return &LuaMAT{
		mcbs:    make(map[string]*luaMCB),
		cluster: make(map[string][]string),
	}
}

// MCB returns specify mcb according to mid
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

// Mount to add MCB to MAT
func (mat *LuaMAT) Mount(funcs LuaFuncs) MAT {
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

// Unmount to remove MCB from MAT
func (mat *LuaMAT) Unmount(mid string) MAT {
	delete(mat.mcbs, mid)
	delete(mat.cluster, mid)
	return mat
}

// Collect to converge specify mcb according to cluster.
func (mat *LuaMAT) Collect(cluster string, mids []string) MAT {
	for _, mid := range mids {
		if _, ok := mat.mcbs[mid]; ok {
			mat.cluster[cluster] = append(mat.cluster[cluster], mid)
		}
	}
	return mat
}

// MCB is abbreviation for module control block,
// which store and control module meta.
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

// Fun returns function prototype
func (mcb *luaMCB) Fun() lua.LGFunction {
	return mcb.fun
}

// Used returns whether mcb is used
func (mcb *luaMCB) Used() bool {
	return mcb.used
}

// Mark to set module used state
func (mcb *luaMCB) Mark() {
	mcb.used = true
}
