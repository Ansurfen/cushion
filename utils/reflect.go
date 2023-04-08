package utils

import (
	"fmt"
	"reflect"
	"unsafe"
)

type ReflectEngine struct {
	raw    any
	rt     reflect.Type
	rv     reflect.Value
	fields map[string]reflect.Value
}

func NewReflectEngine(data any) *ReflectEngine {
	re := &ReflectEngine{raw: data}
	if rv := reflect.ValueOf(data); rv.Kind() == reflect.Ptr {
		re.rt = reflect.TypeOf(data).Elem()
		re.rv = rv.Elem()
	} else {
		re.rv = reflect.ValueOf(&data).Elem().Elem()
		re.rt = reflect.TypeOf(data)
	}
	if re.rt.NumField() > 0 {
		re.fields = make(map[string]reflect.Value)
	}
	for i := 0; i < re.rt.NumField(); i++ {
		re.fields[re.rt.Field(i).Name] = re.rv.Field(i)
	}
	return re
}

func (re *ReflectEngine) DumpFields() {
	for name, value := range re.fields {
		fmt.Println(name, value)
	}
}

func (re *ReflectEngine) Set(field string, value any) error {
	v := re.fields[field]
	if v.CanAddr() {
		rv := reflect.ValueOf(value)
		if v.Kind() != rv.Kind() {
			return fmt.Errorf("invalid kind: expected kind %v, got kind: %v", v.Kind(), rv.Kind())
		}
		reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Set(rv)
		return nil
	}
	// set raw, using Raw() to get no reference data
	if !v.CanSet() {
		vv := reflect.ValueOf(&re.raw).Elem()
		tmp := reflect.New(vv.Elem().Type()).Elem()
		tmp.Set(vv.Elem())
		tmp.FieldByName(field).Set(reflect.ValueOf(value))
		vv.Set(tmp)
	}
	return nil
}

func (re *ReflectEngine) Get(field string) reflect.Value {
	v := re.fields[field]
	if v.CanAddr() {
		reflect.ValueOf(re.raw).Elem().FieldByName(field)
		v = reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
	}
	return v
}

func (re *ReflectEngine) Raw() any {
	return re.raw
}

type ReflectDict struct {
	types map[string]TypeMeta
}

type TypeMeta struct {
	addr  uintptr
	_type reflect.Type
}

func NewReflectDict() *ReflectDict {
	return &ReflectDict{
		types: make(map[string]TypeMeta),
	}
}

func (dict *ReflectDict) Load(datas ...any) {
	for _, data := range datas {
		v := reflect.ValueOf(&data).Elem()
		dict.types[v.Elem().Type().Name()] = TypeMeta{addr: v.UnsafeAddr(), _type: v.Elem().Type()}
	}
}

func (dict *ReflectDict) New(name string, value any) any {
	if meta, ok := dict.types[name]; ok {
		rv := reflect.ValueOf(value)
		if rv.Kind() != meta._type.Kind() {
			return nil
		}
		ret := reflect.NewAt(meta._type, unsafe.Pointer(meta.addr)).Elem()
		ret.Set(rv)
		return ret.Interface()
	}
	return nil
}
