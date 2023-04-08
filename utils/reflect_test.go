package utils

import (
	"fmt"
	"testing"
)

type TestData struct {
	a  int
	B  string
	_C float32
}

func TestReflectEngine(t *testing.T) {
	testRefData := &TestData{a: 10}
	re := NewReflectEngine(testRefData)
	re.DumpFields()
	re.Set("a", 6666)
	fmt.Println(re.Get("a"))
	fmt.Println("----------------")
	testRawData := TestData{B: "..."}
	re2 := NewReflectEngine(&testRawData)
	re2.Set("B", "aaa")
	re2.DumpFields()
}

func TestReflectWithBytes(t *testing.T) {
	testRefData := &TestData{_C: 0}
	re := NewReflectEngine(testRefData)
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
