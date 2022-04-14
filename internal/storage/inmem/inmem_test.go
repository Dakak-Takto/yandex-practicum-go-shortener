package inmem

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInmem(t *testing.T) {
	t.Run("Create storage test", func(t *testing.T) {
		s, err := New()
		require.NoError(t, err)
		require.IsType(t, &store{}, s)
	})

	t.Run("Write storage test", func(t *testing.T) {
		s, err := New()
		require.NoError(t, err)

		err = s.Set("test-key", "test-value")
		require.NoError(t, err)
	})

	t.Run("Read storage test", func(t *testing.T) {
		s, err := New()
		require.NoError(t, err)

		s.Set("test-key", "test-value")
		v, err := s.Get("test-key")

		require.NoError(t, err)
		require.Equal(t, "test-value", v)
	})
}
