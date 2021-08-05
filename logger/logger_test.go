package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadinessLogic(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("*zap.SugaredLogger", fmt.Sprintf("%T", New().Sugar()))
	assert.Equal("*zap.SugaredLogger", fmt.Sprintf("%T", New(Level("info")).Sugar()))
}
