package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeyDoubleUint64(t *testing.T) {
	dkey := NewKeyDoubleUint64(1, 1)
	assert.True(t, dkey.Equal(NewKeyDoubleUint64(1, 1)))
	var dkey2 KeyDoubleUint64
	assert.True(t, dkey2.FromString("1-1"))
	assert.True(t, dkey.Equal(dkey2))
	assert.Equal(t, "1-1", dkey.ToString())
}
