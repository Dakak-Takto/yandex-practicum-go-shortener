package postgresql

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/entity"
)

type database struct {
	db *sqlx.DB
}

func New(dsn string) (entity.URLRepository, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	result, err := db.Exec(`CREATE TABLE IF NOT EXISTS shorts (short VARCHAR(255), original VARCHAR(255) UNIQUE, user_id VARCHAR(255), deleted BOOLEAN default false)`)

	if err != nil {
		return nil, err
	}

	log.Println(result.RowsAffected())

	return &database{
		db: db,
	}, err
}

func (d *database) GetByShort(ctx context.Context, short string) (url *entity.URL, err error) {
	err = d.db.Get(url, "SELECT short, original, user_id, deleted FROM shorts WHERE short=$1", short)
	return url, err
}

func (d *database) SelectByUserID(ctx context.Context, uid string) (urls []*entity.URL, err error) {
	err = d.db.Select(urls, "SELECT short, original, user_id, deleted FROM shorts WHERE user_id = $1", uid)
	return urls, err
}

func (d *database) GetByOriginal(ctx context.Context, original string) (url *entity.URL, err error) {
	err = d.db.Get(url, "SELECT short, original, user_id, deleted FROM shorts WHERE original = $1", original)
	return url, err
}
func (d *database) Save(ctx context.Context, url *entity.URL) error {
	_, err := d.db.NamedExec("INSERT INTO shorts (short, original, user_id) VALUES (:short, :original, :user_id)", url)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				err = entity.ErrDuplicate
			}
		}
	}
	return err
}

func (d *database) Ping() error {
	return d.db.Ping()
}

func (d *database) Delete(uid string, keys ...string) {
	_, err := d.db.Exec("UPDATE shorts SET deleted = true WHERE short = any($1) AND user_id = $2", keys, uid)
	if err != nil {
		log.Println("error set deleted: ", err)
	}
}

func (d *database) Lock()   {}
func (d *database) Unlock() {}
