package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("Проверка на запись и чтение", func(t *testing.T) {

		testStorage := New()

		testStorage.Set("key", "value")

		value, err := testStorage.Get("key")
		require.NoError(t, err)
		require.Equal(t, value, "value")
	})
}
