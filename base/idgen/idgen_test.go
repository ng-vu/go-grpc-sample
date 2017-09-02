package idgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcInfix(t *testing.T) {
	x := CalcInfix("AB")
	assert.Equal(t, x, uint16(49234))

	assert.Panics(t, func() {
		CalcInfix("123")
	})

	assert.Panics(t, func() {
		CalcInfix("IO")
	})
}

type nilReader struct{}

func (r nilReader) Read(b []byte) (int, error) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
	return len(b), nil
}

func TestNewID(t *testing.T) {
	t.Run("Empty with prefix", func(t *testing.T) {
		infix := CalcInfix("CD")

		id := NewID(infix, 0, &nilReader{}).String()
		assert.Equal(t, id, "0000000000CD00000000000000")
		assert.Equal(t, id[10:12], "CD")
	})

	t.Run("Generate with prefix", func(t *testing.T) {
		infix := CalcInfix("XZ")

		for i := 0; i < 10; i++ {
			id := Generate(infix).String()
			assert.Equal(t, id[10:12], "XZ")

			ll.Info(id)
		}
	})
}
