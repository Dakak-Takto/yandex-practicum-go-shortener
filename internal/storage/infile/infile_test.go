//Used for file storage. Making file and write url per line

package infile

import (
	"bufio"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInFile(t *testing.T) {

	file, err := os.OpenFile("testFilename", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	s := &store{
		file:   file,
		reader: bufio.NewReader(file),
		writer: bufio.NewWriter(file),
	}

	t.Run("Write storage test", func(t *testing.T) {
		s.Lock()
		defer s.Unlock()
		err = s.Set("test-key", "test-value")
		require.NoError(t, err)
	})

	t.Run("Read storage test", func(t *testing.T) {
		s.Lock()
		defer s.Unlock()
		s.Set("test-key", "test-value")
		v, err := s.Get("test-key")

		require.NoError(t, err)
		require.Equal(t, "test-value", v)
	})

	if err := s.file.Close(); err != nil {
		log.Println(err)
	}
	if err := os.Remove(file.Name()); err != nil {
		log.Println(err)
	}
}
