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

	t.Run("Read storage test", func(t *testing.T) {
		s, err := New()
		require.NoError(t, err)

		s.Save("test-key", "test-value", "0")
		v, err := s.GetByShort("test-key")

		require.NoError(t, err)
		require.Equal(t, "test-value", v.Original)
	})
}
