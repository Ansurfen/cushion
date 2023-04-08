package runtime

var luaExport ArkExport

func init() {
	luaExport = &LuaExport{
		raws:    make([]string, 0),
		formats: make(map[string]string),
	}
}

func GetExport() ArkExport {
	return luaExport
}

type ArkExport interface {
	Raw() []string
	Format() []string
	Set(raw string, format string)
}

type LuaExport struct {
	raws    []string
	formats map[string]string
}

func (luaExport *LuaExport) Raw() []string {
	return luaExport.raws
}

func (luaExport *LuaExport) Format() []string {
	var exports []string
	for _, key := range luaExport.raws {
		if value, ok := luaExport.formats[key]; ok {
			exports = append(exports, value)
		}
	}
	return exports
}

func (luaExport *LuaExport) Set(k, v string) {
	luaExport.raws = append(luaExport.raws, k)
	luaExport.formats[k] = v
}
