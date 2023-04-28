package utils

import (
	"fmt"
	"testing"
)

func TestCvt(t *testing.T) {
	out, _ := ExecStr("docker -v")
	fmt.Println(ConvertByte2String(out, GB18030))
}

func TestSimilarText(t *testing.T) {
	rate := 1.0
	SimilarText("aasdage", "abadanodja", &rate)
	fmt.Println(rate)
	SimilarText("xsa", "abadanodja", &rate)
	fmt.Println(rate)
}

func TestTypeCvtBytes(t *testing.T) {
	SetBytesMode(BigEndian)
	fmt.Println(ByteTransfomer.BytesToFloat32(ByteTransfomer.Float32ToBytes(3.2)))
	fmt.Println(ByteTransfomer.BytesToFloat64(ByteTransfomer.Float64ToBytes(1.23456)))
}

func TestByteWalk(t *testing.T) {
	bw := NewByteWalk([]byte("Hello world"))
	fmt.Printf("Cursor: %d\nSize: %d\n", bw.Cursor(), bw.Size())
	out, err := bw.Next(5)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	out, err = bw.Next(10)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
