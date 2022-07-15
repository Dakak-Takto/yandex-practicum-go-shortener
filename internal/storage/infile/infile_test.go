//Used for file storage. Making file and write url per line

package infile

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInFile(t *testing.T) {

	t.Run("Read storage test", func(t *testing.T) {
		const testFileName string = "testFile.txt"
		defer os.Remove(testFileName)

		store, err := New(testFileName)
		require.NoError(t, err)
		err = store.Save("test-key", "test-value", "0")
		require.NoError(t, err)
		v, err := store.GetByShort("test-key")
		require.NoError(t, err)
		require.Equal(t, "test-value", v.Original)
	})
}
