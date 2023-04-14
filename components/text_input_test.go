package components

import (
	"fmt"
	"testing"
)

func TestTextInput(t *testing.T) {
	res := UseTextInput(DefaultTextStyle(), &TextInputPayLoad{
		Texts: []TextInputFormat{
			{Name: "host"},
			{Name: "password", EchoMode: true},
		},
	})
	fmt.Println(res)
}
