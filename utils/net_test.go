package utils

import (
	"fmt"
	"testing"
)

func TestNet(t *testing.T) {
	urlStrings := []string{"https://www.example.com", "http://www.example.com", "www.example.com", "example.com", "example"}

	for _, urlString := range urlStrings {
		fmt.Println(IsURL(urlString))
	}
}
