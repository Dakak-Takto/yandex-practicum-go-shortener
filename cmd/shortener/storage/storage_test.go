package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	testString := "test string content"
	t.Run("Write and Read storage", func(t *testing.T) {
		key, err := Save(testString)
		require.NoError(t, err)

		value, err := Get(key)
		require.NoError(t, err)
		require.Equal(t, value, testString)
	})

	t.Run("Read non-exists", func(t *testing.T) {
		value, err := Get("non-exists-key")
		require.Error(t, err)
		require.Empty(t, value)
	})

	t.Run("Concurrent read and write", func(t *testing.T) {
		ok := true
		go func() {
			var i int
			for i = 0; ok; i++ {
				Get(time.Now().String())
			}
		}()
		go func() {
			var i int
			for i = 0; ok; i++ {
				Save(time.Now().String())
			}
		}()
		time.Sleep(time.Second)
		ok = false
	})
}
