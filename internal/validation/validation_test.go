package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	t.Run("validates correctly", func(t *testing.T) {
		tests := map[string]struct {
			current   string
			available []string
			expected  bool
		}{
			"when available is nil":                   {"a", nil, false},
			"when available is empty":                 {"a", []string{}, false},
			"when available does not contain current": {"a", []string{"b", "c", "d"}, false},
			"when available does contain current":     {"a", []string{"b", "a", "c", "d"}, true},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				result := isValid(test.current, test.available)
				assert.Equal(t, test.expected, result)
			})
		}
	})
}
