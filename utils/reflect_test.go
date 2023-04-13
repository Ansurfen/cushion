package utils

import (
	"fmt"
	"reflect"
	"testing"
)

type TestData struct {
	a  int
	B  string
	_C float32
}

func TestReflectEngine(t *testing.T) {
	testRefData := &TestData{a: 10}
	re := NewReflectObject(testRefData)
	re.DumpFields()
	re.Set("a", 6666)
	fmt.Println(re.Get("a"))
	fmt.Println("----------------")
	testRawData := TestData{B: "..."}
	re2 := NewReflectObject(&testRawData)
	fmt.Println(re2.Set("Baa", "aaa"))
	re2.Set("B", "aaa")
	re2.DumpFields()
}

func TestReflectWithBytes(t *testing.T) {
	testRefData := &TestData{_C: 0}
	re := NewReflectObject(testRefData)
	SetBytesMode(BigEndian)
	bits := ByteTransfomer.Float32ToBytes(3.2)
	re.Set("_C", ByteTransfomer.AutoToType(bits, Float32Type))
	re.DumpFields()
}

func TestReflectDict(t *testing.T) {
	dict := NewReflectDict()
	dict.Load(1, 1.5, "")
	fmt.Println(dict.New("int", int(10)))
	fmt.Println(dict.New("string", "Hello World!"))
}

type TestMulData struct {
	Payload TestData
	b       int
}

func TestNestedReflect(t *testing.T) {
	mulData := TestMulData{Payload: TestData{a: 10}, b: 5}
	re := NewReflectObject(&mulData)
	re.Set("b", 10)
	if re.Get("Payload").CanInterface() {
		if re.Get("Payload").Type().Kind() == reflect.Struct {
			structField := re.Get("Payload").Type()
			for i := 0; i < structField.NumField(); i++ {
				fmt.Println(structField.Field(i).Name, structField.Field(i).Type)
			}
		}
	}
	fmt.Println(mulData)
}
