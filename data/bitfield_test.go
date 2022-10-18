package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBitfield(t *testing.T) {
	bits8 := Bitfield1[uint8]{}
	bits8 = bits8.Set(2)
	assert.True(t, bits8.Has(2))
	bits8 = bits8.Clear(2)
	assert.False(t, bits8.Has(2))
	bits8 = bits8.Toggle(2)
	assert.True(t, bits8.Has(2))

	bits64 := Bitfield1[uint64]{}
	bits64 = bits64.Set(11)
	bits64 = bits64.Set(7)
	assert.True(t, bits64.Has(11))
	assert.True(t, bits64.Has(7))
	bits64 = bits64.Clear(2)
	bits64 = bits64.Clear(11)
	bits64 = bits64.Clear(7)
	assert.False(t, bits64.Has(2))
	assert.False(t, bits64.Has(11))
	assert.False(t, bits64.Has(7))
}
