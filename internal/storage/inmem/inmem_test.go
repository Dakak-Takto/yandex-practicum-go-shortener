package inmem

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInmem(t *testing.T) {
	t.Run("Create storage test", func(t *testing.T) {
		var s = New()
		require.IsType(t, &store{}, s)
	})

	t.Run("Write storage test", func(t *testing.T) {
		var s = New()
		var err = s.Set("test-key", "test-value")
		require.NoError(t, err)
	})

	t.Run("Read storage test", func(t *testing.T) {
		var s = New()
		s.Set("test-key", "test-value")
		v, err := s.Get("test-key")
		require.NoError(t, err)
		require.Equal(t, "test-value", v)
	})
}
