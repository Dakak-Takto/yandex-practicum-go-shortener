package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	testString := "test string content"
	t.Run("Проверка на запись и чтение", func(t *testing.T) {
		err := Set("test key", "test value")
		require.NoError(t, err)

		value, err := Get("test key")
		require.NoError(t, err)
		require.Equal(t, value, testString)
	})
}
