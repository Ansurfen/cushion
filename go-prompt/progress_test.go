package prompt

import (
	"fmt"
	"testing"
	"time"
)

func TestProgress(t *testing.T) {
	process := NewProgress([]string{"|", "/", "-", "\\"})
	for {
		fmt.Print("\r" + process.Next())
		time.Sleep(100 * time.Millisecond)
	}
}
