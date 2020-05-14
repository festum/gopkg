package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadinessLogic(t *testing.T) {
	logger := NewLogger()
	assert.Equal(t, "*zap.SugaredLogger", fmt.Sprintf("%T", logger))
}
