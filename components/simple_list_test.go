package components

import (
	"fmt"
	"testing"
)

func TestSimpleList(t *testing.T) {
	choice := UseSimpleList(DefaultSimpleListStyle(), &SimpleListPayLoad{
		Title: "",
		Choices: []string{
			"a", "b", "c",
		},
	})
	fmt.Println(choice)
}
