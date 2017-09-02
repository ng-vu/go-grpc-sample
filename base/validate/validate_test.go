package validate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalize(t *testing.T) {
	t.Run("Phone", func(t *testing.T) {
		tests := []struct {
			name  string
			phone string
			want  string
			ok    bool
		}{{
			"Empty (Invalid)",
			"", "", false,
		}, {
			"Valid",
			"0123456789", "0123456789", true,
		}, {
			"Trim dash",
			"0123-456-789 ", "0123456789", true,
		}, {
			"Too short (Invalid)",
			"1234", "01234", false,
		}, {
			"Too short (Invalid)",
			"1234", "01234", false,
		}}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, ok := NormalizePhone(tt.phone)

				require.Equal(t, tt.ok, ok)
				if ok {
					require.Equal(t, tt.want, got)
				}
			})
		}
	})

	t.Run("Code", func(t *testing.T) {
		tests := []struct {
			name string
			code string
			ok   bool
		}{{
			"Empty (Invalid)",
			"", false,
		}, {
			"Contain space (Invalid)",
			"A code", false,
		}, {
			"Valid",
			"0123456789_@#!?-[]", true,
		}, {
			"Too short (Invalid)",
			"12", false,
		}}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ok := Code(tt.code)
				require.Equal(t, tt.ok, ok)
			})
		}
	})

	t.Run("Name", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
			want  string
			ok    bool
		}{{
			"Empty (Invalid)",
			"", "", false,
		}, {
			"Trim space",
			" Sample  Name ", "Sample Name", true,
		}, {
			"Allow vietnamese",
			" Nguyễn Ngọc Minh An ", "Nguyễn Ngọc Minh An", true,
		}, {
			"Allow sign characters",
			` @#$%\/? shop `, `@#$%\/? shop`, true,
		}, {
			"Trim invalid characters",
			"One ≠ 1", "One 1", true,
		}, {
			"Too short after trim (Invalid)",
			" ≠≠≠≠≠≠A≠≠≠≠≠≠ ", "", false,
		}}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, ok := NormalizeName(tt.input)

				require.Equal(t, tt.ok, ok)
				if ok {
					require.Equal(t, tt.want, got)
				}
			})
		}
	})
}
