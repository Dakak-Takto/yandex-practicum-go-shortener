package random

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	t.Run("Test random", func(t *testing.T) {
		for i := 0; i < 50; i++ {
			str := String(i)
			require.Equal(t, len(str), i)
		}
	})
}
