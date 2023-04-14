package components

import (
	"testing"
	"time"
)

func TestSpinner(t *testing.T) {
	UseSpinner(DefaultSpinnerStyle(), &SpinnerPayLoad{
		Callback: func() {
			timer := time.NewTicker(5*time.Second)
			select {
			case <-timer.C:
				break
			// default:
			}
		},
	})
}
