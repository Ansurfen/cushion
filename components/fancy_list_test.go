package components

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/bubbles/list"
)

func TestFancyList(t *testing.T) {
	choice := UseFancyList(DefaultFancyListStyle(), &FancyListPayLoad{
		Title: "",
		Choices: []list.Item{
			FancyListItem{ChoiceTitle: "a", ChoiceDetail: "b"},
		},
	})
	fmt.Println("your choice: ", choice)
}
