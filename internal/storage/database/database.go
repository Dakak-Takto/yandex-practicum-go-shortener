package database

import (
	"errors"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"yandex-practicum-go-shortener/internal/storage"
)

type database struct {
	db *sqlx.DB
}

func New(dsn string) (storage.Storage, error) {
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

func (d *database) GetByShort(key string) (row storage.URLRecord, err error) {
	err = d.db.Get(&row, "SELECT short, original, user_id, deleted FROM shorts WHERE short=$1", key)
	return row, err
}

func (d *database) SelectByUID(uid string) (rows []storage.URLRecord, err error) {
	err = d.db.Select(&rows, "SELECT short, original, user_id, deleted FROM shorts WHERE user_id = $1", uid)
	return rows, err
}

func (d *database) GetByOriginal(original string) (row storage.URLRecord, err error) {
	err = d.db.Get(&row, "SELECT short, original, user_id, deleted FROM shorts WHERE original = $1", original)
	return row, err
}
func (d *database) Save(short, original, userID string) error {
	_, err := d.db.Exec("INSERT INTO shorts (short, original, user_id) VALUES ($1, $2, $3)", short, original, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				err = storage.ErrDuplicate
			}
		}
	}
	return err
}

func (d *database) Ping() error {
	return d.db.Ping()
}

func (d *database) Delete(keys []string) {
	_, err := d.db.Exec("UPDATE shorts SET deleted = true WHERE short = any($1)", keys)
	if err != nil {
		log.Println("error set deleted: ", err)
	}
}

func (d *database) Lock()   {}
func (d *database) Unlock() {}
