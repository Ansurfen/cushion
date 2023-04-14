package components

import (
	"fmt"
	"testing"
)

func TestMultiSelect(t *testing.T) {
	choice := UseMultiSelect(DefaultMultiSelectStyle(), &MultiSelectPayLoad{
		Title: "",
		Choices: []string{
			"a", "b", "c",
		},
	})
	fmt.Println("your choice: ", choice)
}
