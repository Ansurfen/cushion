package runtime

import (
	"fmt"
	"strings"
)

type LuaObject interface {
	Ident() string
	Value() string
}

var (
	LuaNil   = NewLuaNil()
	LuaTrue  = LuaBoolean(true)
	LuaFalse = LuaBoolean(false)
)

type LuaNilObject struct {
	ident string
}

func NewLuaNil() *LuaNilObject {
	return &LuaNilObject{}
}

func (obj *LuaNilObject) Ident() string {
	return obj.ident
}

func (LuaNilObject) Value() string {
	return "nil"
}

type LuaBooleanObject struct {
	ident string
	value bool
}

func NewLuaBooleanObject(ident string, value bool) *LuaBooleanObject {
	return &LuaBooleanObject{
		ident: ident,
		value: value,
	}
}

func LuaBoolean(value bool) *LuaBooleanObject {
	return NewLuaBooleanObject("", value)
}

func (obj *LuaBooleanObject) Ident() string {
	return obj.ident
}

func (obj *LuaBooleanObject) Value() string {
	if obj.value {
		return "true"
	}
	return "false"
}

type LuaStringObject struct {
	ident string
	value string
}

func NewLuaStringObject(ident, value string) *LuaStringObject {
	return &LuaStringObject{
		ident: ident,
		value: value,
	}
}

func LuaString(value string) *LuaStringObject {
	return &LuaStringObject{
		value: value,
	}
}

func (obj *LuaStringObject) Ident() string {
	return obj.ident
}

func (obj *LuaStringObject) Value() string {
	return fmt.Sprintf(`"%s"`, obj.value)
}

type LuaNumberObject struct {
	ident string
	value string
}

func NewLuaNumberObject(ident, value string) *LuaNumberObject {
	return &LuaNumberObject{
		ident: ident,
		value: value,
	}
}

func LuaNumber(value string) *LuaNumberObject {
	return &LuaNumberObject{
		value: value,
	}
}

func (obj *LuaNumberObject) Ident() string {
	return obj.ident
}

func (obj *LuaNumberObject) Value() string {
	return obj.value
}

type LuaListObejct struct {
	ident  string
	values []LuaObject
}

func NewLuaListObject(ident string, values ...LuaObject) *LuaListObejct {
	return &LuaListObejct{
		ident:  ident,
		values: values,
	}
}

func LuaList(values ...LuaObject) *LuaListObejct {
	return &LuaListObejct{
		values: values,
	}
}

func (obj *LuaListObejct) Ident() string {
	return obj.ident
}

func (obj *LuaListObejct) Value() string {
	return fmt.Sprintf(`{ %s }`, obj.Unpack())
}

func (obj *LuaListObejct) Unpack() string {
	values := ""
	for idx, v := range obj.values {
		if idx != 0 && idx <= len(obj.values)-1 {
			values += ", "
		}
		values += v.Value()
	}
	return values
}

func (obj *LuaListObejct) Insert(value ...LuaObject) {
	obj.values = append(obj.values, value...)
}

type Luamap map[string]LuaObject

type LuaMapObject struct {
	ident  string
	values Luamap
}

func NewLuaMapObject(ident string, values Luamap) *LuaMapObject {
	tbl := &LuaMapObject{
		ident:  ident,
		values: make(Luamap),
	}
	if values != nil {
		tbl.values = values
	}
	return tbl
}

func LuaMap(values Luamap) *LuaMapObject {
	tbl := &LuaMapObject{
		values: make(Luamap),
	}
	if values != nil {
		tbl.values = values
	}
	return tbl
}

func LuaStringMap(values map[string]string) *LuaMapObject {
	m := make(Luamap)
	for key, value := range values {
		m[key] = LuaString(value)
	}
	return LuaMap(m)
}

func (obj *LuaMapObject) Insert(k string, v LuaObject) *LuaMapObject {
	keys := strings.Split(k, ".")
	if len(keys) <= 0 {
		return obj
	}
	if len(keys) == 1 {
		obj.values[keys[0]] = v
		return obj
	}
	var curMap *LuaMapObject
	first := keys[0]
	if m, ok := obj.values[first]; ok {
		switch v := m.(type) {
		case *LuaMapObject:
			curMap = v
		default:
			obj.values[first] = NewLuaMapObject(first, Luamap{})
			curMap = obj.values[first].(*LuaMapObject)
		}
	} else {
		obj.values[first] = NewLuaMapObject(first, Luamap{})
		curMap = obj.values[first].(*LuaMapObject)
	}
	for idx, key := range keys {
		if idx == 0 {
			continue
		}
		if idx == len(keys)-1 {
			curMap.values[key] = v
			break
		}
		if m, ok := curMap.values[key]; ok {
			switch v := m.(type) {
			case *LuaMapObject:
				curMap = v
			default:
				obj.values[first] = NewLuaMapObject(first, Luamap{})
				curMap = curMap.values[first].(*LuaMapObject)
			}
		} else {
			curMap.values[key] = NewLuaMapObject(key, Luamap{})
			curMap = curMap.values[key].(*LuaMapObject)
		}
	}
	return obj
}

func (obj *LuaMapObject) Ident() string {
	return obj.ident
}

func (obj *LuaMapObject) Value() string {
	return fmt.Sprintf("{ %s }", obj.Unpack())
}

func (obj *LuaMapObject) Unpack() string {
	values := ""
	idx := 0
	for k, v := range obj.values {
		if idx != 0 && idx <= len(obj.values)-1 {
			values += ", "
		}
		values += fmt.Sprintf("%s = %s", k, v.Value())
		idx++
	}
	return values
}

type LuaTableObject struct {
	ident string
	m     *LuaMapObject
	l     *LuaListObejct
}

func NewLuaTabelObject(ident string) *LuaTableObject {
	return &LuaTableObject{
		ident: ident,
		m:     LuaMap(Luamap{}),
		l:     LuaList(),
	}
}

func LuaTable() *LuaTableObject {
	return NewLuaTabelObject("")
}

func (obj *LuaTableObject) Ident() string {
	return obj.ident
}

func (obj *LuaTableObject) Value() string {
	return fmt.Sprintf("{ %s }", obj.Unpack())
}

func (obj *LuaTableObject) Unpack() string {
	ml, ll := len(obj.m.Unpack()), len(obj.l.Unpack())
	if ml == ll && ml == 0 {
		return ""
	} else if ml == 0 && ll > 0 {
		return obj.l.Unpack()
	} else if ll == 0 && ml > 0 {
		return obj.m.Unpack()
	} else {
		return fmt.Sprintf("%s, %s", obj.l.Unpack(), obj.m.Unpack())
	}
}

func (obj *LuaTableObject) SetMap(kvs Luamap) *LuaTableObject {
	if obj.m == nil {
		obj.m = LuaMap(kvs)
		return obj
	}
	obj.m.values = kvs
	return obj
}

func (obj *LuaTableObject) SetList(list ...LuaObject) *LuaTableObject {
	if obj.l == nil {
		obj.l = LuaList(list...)
		return obj
	}
	obj.l.values = list
	return obj
}

func (obj *LuaTableObject) InsertKV(k string, v LuaObject) *LuaTableObject {
	obj.m.Insert(k, v)
	return obj
}

func (obj *LuaTableObject) Insert(values ...LuaObject) *LuaTableObject {
	obj.l.Insert(values...)
	return obj
}

type LuaFuncObject struct {
	lambda bool
	ident  string
	args   []string
	argv   []LuaObject
	block  []LuaObject
}

func NewLuaFuncObject(ident string) *LuaFuncObject {
	return &LuaFuncObject{
		ident: ident,
	}
}

func LuaFunc() *LuaFuncObject {
	return &LuaFuncObject{
		lambda: true,
	}
}

func (obj *LuaFuncObject) SetArgs(args []string) *LuaFuncObject {
	obj.args = args
	return obj
}

func (obj *LuaFuncObject) SetArgc(n int) *LuaFuncObject {
	obj.args = make([]string, n)
	return obj
}

func (obj *LuaFuncObject) Ident() string {
	return obj.ident
}

func (obj *LuaFuncObject) Value() string {
	ident := ""
	if !obj.lambda {
		ident = obj.ident
	}
	return fmt.Sprintf("function %s(%s)\n%send", ident, strings.Join(obj.args, ", "), "")
}

func (obj *LuaFuncObject) Lambda() *LuaFuncObject {
	obj.lambda = true
	return obj
}

func (obj *LuaFuncObject) Bind(argv ...LuaObject) *LuaFuncObject {
	obj.argv = argv
	return obj
}

func (obj *LuaFuncObject) PCall(argv ...LuaObject) string {
	argc := len(obj.args)
	if argc <= 0 || argc > len(argv) {
		return obj.prototype()
	}
	a := make([]any, len(obj.args))
	for i := 0; i < argc; i++ {
		a[i] = argv[i].Value()
	}
	return fmt.Sprintf(obj.prototype(), a...)
}

func (obj *LuaFuncObject) prototype() string {
	argc := len(obj.args)
	if argc <= 0 {
		return fmt.Sprintf("%s()", obj.ident)
	}
	tmp := make([]string, argc)
	for i := 0; i < argc; i++ {
		tmp[i] = "%s"
	}
	return fmt.Sprintf("%s(%s)", obj.ident, strings.Join(tmp, ", "))
}

func (obj *LuaFuncObject) SetBlock(block ...LuaObject) {
	obj.block = append(obj.block, block...)
}

type LuaProgramObject struct {
	ident string
	codes []LuaObject
}

func NewLuaProgramObject(ident string, codes ...LuaObject) *LuaProgramObject {
	return &LuaProgramObject{
		ident: ident,
		codes: codes,
	}
}

func LuaProgram(codes ...LuaObject) *LuaProgramObject {
	return NewLuaProgramObject("", codes...)
}

func (obj *LuaProgramObject) Ident() string {
	return obj.ident
}

func (obj *LuaProgramObject) Value() string {
	values := ""
	for _, code := range obj.codes {
		values += code.Value() + "\n"
	}
	return values
}

func LuaLocalVar(obj LuaObject) string {
	return fmt.Sprintf("local %s", LuaAssign(obj))
}

func LuaAssign(obj LuaObject) string {
	return LuaAssignLR(obj, obj)
}

func LuaAssignLR(left, right LuaObject) string {
	return fmt.Sprintf("%s = %s", left.Ident(), right.Value())
}

func LuaGlobalVar(v string) LuaObject {
	return LuaNumber(v)
}

func LuaIdent(ident string) LuaObject {
	return &LuaNumberObject{
		ident: ident,
	}
}

func LuaComment(comment string) LuaObject {
	return &LuaNumberObject{
		value: comment,
	}
}

func LuaDisableFileUndefinedGlobal() LuaObject {
	return LuaComment("---@diagnostic disable: undefined-global")
}

func LuaDisableNextLineUndefinedGlobal() LuaObject {
	return LuaComment("---@diagnostic disable-next-line: undefined-global")
}

func LuaSpaceLine() LuaObject {
	return LuaGlobalVar("")
}

func LuaPackagePath(paths ...string) LuaObject {
	if len(paths) == 0 {
		return LuaSpaceLine()
	}
	tmp := make([]string, len(paths))
	for i := 0; i < len(paths); i++ {
		tmp[i] = fmt.Sprintf(`"%s"`, paths[i])
	}
	return LuaGlobalVar(fmt.Sprintf(`package.path = package.path .. [[;]] .. %s`, strings.Join(tmp, " .. ")))
}

func LuaRequire(modename string) LuaObject {
	if len(modename) == 0 {
		return LuaSpaceLine()
	}
	return LuaGlobalVar(fmt.Sprintf(`require("%s")`, modename))
}

func LuaImport(mods ...string) LuaObject {
	if len(mods) == 0 {
		return LuaSpaceLine()
	}
	return LuaGlobalVar(fmt.Sprintf(`Import({ %s })`, strings.Join(mods, ", ")))
}
