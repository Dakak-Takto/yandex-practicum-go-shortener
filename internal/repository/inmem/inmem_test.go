package inmem

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/entity"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/pkg/random"
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

		url := entity.URL{
			Short:    "abcdefg",
			Original: "http://abc.defg.hij/klmop?qrst=uvw",
			UserID:   "xyz",
		}
		ctx := context.Background()

		err = s.Save(ctx, &url)
		require.NoError(t, err)
		v, err := s.GetByShort(ctx, url.Short)

		require.NoError(t, err)
		require.Equal(t, url.Original, v.Original)
	})
}

func BenchmarkInmem(b *testing.B) {

	repo, err := New()
	if err != nil {
		b.Error(err)
	}

	b.Run("test save", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			repo.Save(context.Background(), &entity.URL{
				Original: "http://benchmark.test",
				Short:    random.String(10),
			})
		}
	})
}
