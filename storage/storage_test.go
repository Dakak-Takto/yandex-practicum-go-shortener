package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("Проверка на запись и чтение", func(t *testing.T) {

		testValue := "test string content"

		key := SetValueReturnKey(testValue)
		require.NotEmpty(t, key)

		value, err := GetValueByKey(key)
		require.NoError(t, err)
		require.Equal(t, value, testValue)
	})
}
