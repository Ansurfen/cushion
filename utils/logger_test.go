package utils

import (
	"testing"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	InitLoggerWithDefault()
	zap.S().Info("Test")
}
