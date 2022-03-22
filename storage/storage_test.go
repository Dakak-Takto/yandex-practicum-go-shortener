package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	testString := "test string content"
	t.Run("Проверка на запись и чтение", func(t *testing.T) {
		key, err := Save(testString)
		require.NoError(t, err)

		value, err := Get(key)
		require.NoError(t, err)
		require.Equal(t, value, testString)
	})

	t.Run("Проверка чтения несуществующего ключа", func(t *testing.T) {
		value, err := Get("non-exists-key")
		require.Error(t, err)
		require.Empty(t, value)
	})
}
