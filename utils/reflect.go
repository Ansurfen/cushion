package utils

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

// ReflectObject is an assembly on reflect.Value and relfect.Type.
type ReflectObject struct {
	raw    any
	rt     reflect.Type
	rv     reflect.Value
	fields map[string]reflect.Value
}

func NewReflectObject(data any) *ReflectObject {
	re := &ReflectObject{raw: data}
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

// DumpFields print fields and values for current struct
func (re *ReflectObject) DumpFields() {
	for name, value := range re.fields {
		fmt.Println(name, value)
	}
}

// Fields return all fields for current type
func (re *ReflectObject) Fields() map[string]reflect.Value {
	return re.fields
}

// Set can ignore case to set field value
func (re *ReflectObject) Set(field string, value any) error {
	v := re.fields[field]
	if !v.IsValid() {
		return errors.New("invalid field")
	}
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

// Get return value of spcify field
func (re *ReflectObject) Get(field string) reflect.Value {
	v := re.fields[field]
	if v.CanAddr() {
		reflect.ValueOf(re.raw).Elem().FieldByName(field)
		v = reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
	}
	return v
}

// Raw return raw data
func (re *ReflectObject) Raw() any {
	return re.raw
}

// ReflectDict is a type dictionary and manage all type by map.
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

// Load to register type into ReflectDict
func (dict *ReflectDict) Load(datas ...any) {
	for _, data := range datas {
		v := reflect.ValueOf(&data).Elem()
		dict.types[v.Elem().Type().Name()] = TypeMeta{addr: v.UnsafeAddr(), _type: v.Elem().Type()}
	}
}

// New create a variable according to type name to be specified
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
