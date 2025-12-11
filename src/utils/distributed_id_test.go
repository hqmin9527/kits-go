package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetObjectTime(t *testing.T) {
	ts, ok := GetObjectTime("aaa")
	assert.Equal(t, int64(0), ts)
	assert.False(t, ok)
	ts, ok = GetObjectTime("653f00528a51ea00012a53b7")
	assert.True(t, ok)
	assert.Equal(t, int64(1698627666000), ts)
}
