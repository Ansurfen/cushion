package components

import "testing"

func TestBatchSpinner(t *testing.T) {
	UseBatchSpinner(DefaultBatchSpinnerStyle(), &BatchSpinnerPayLoad{
		Task: []BatchTask{
			{Name: "fetch", Callback: func() bool { return true }},
			{Name: "unzip", Callback: func() bool { return true }},
			{Name: "init env", Callback: func() bool { return true }},
		},
	})
}
