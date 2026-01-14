package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetObjectTime(t *testing.T) {
	ts, err := ParseObjectIDTs("aaa")
	assert.Equal(t, int64(0), ts)
	assert.NotNil(t, err)
	ts, err = ParseObjectIDTs("653f00528a51ea00012a53b7")
	assert.Nil(t, err)
	assert.Equal(t, int64(1698627666000), ts)
}

func TestGetUUIDTime(t *testing.T) {
	id := GenUUID()
	t.Logf("uuid: %s", id)
	ts, err := ParseUUIDTs(id)
	assert.Nil(t, err)
	t.Logf("TestGetUUIDTime success, ts: %d", ts)
}
