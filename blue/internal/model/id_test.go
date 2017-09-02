package model

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ng-vu/go-grpc-sample/base/idgen"
	"github.com/ng-vu/go-grpc-sample/base/l"

	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	t.Run("NewID", func(t *testing.T) {
		ID1 := NewID(Infix(idgen.CalcInfix("AB")))
		ll.Info(string(ID1))

		assert.Equal(t, len(ID1), 26)
		assert.Equal(t, string(ID1[10:12]), "AB")
		assert.Equal(t, ID1.InfixString(), "AB")

		assert.True(t, regexp.MustCompile(`^[A-Z0-9]{26}$`).Match([]byte(ID1)))
	})

	t.Run("No duplicated", func(t *testing.T) {
		m := make(map[string]bool)

		for i := 0; i < 10; i++ {
			id := string(NewID(Infix(idgen.CalcInfix("AB"))))
			ll.Info(string(id), l.String("hex", fmt.Sprintf("%x", id)))

			assert.False(t, m[id], "Should not duplicated")
			m[id] = true
		}
	})
}
