//Used for file storage. Making file and write url per line

package infile

import (
	"context"
	"os"
	"testing"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/entity"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/pkg/random"
)

func BenchmarkInfile(b *testing.B) {

	const (
		testFileName string = `repo.txt`
	)

	repo, err := New(testFileName)
	if err != nil {
		b.Error(err)
	}
	defer os.Remove(testFileName)

	b.Run("test save", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			repo.Save(context.Background(), &entity.URL{
				Original: "http://benchmark.test",
				Short:    random.String(10),
			})
		}
	})
}
