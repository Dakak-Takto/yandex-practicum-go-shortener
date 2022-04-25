package database

import (
	"log"
	"yandex-practicum-go-shortener/internal/storage"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v4/stdlib"
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

	result, err := db.Exec(`CREATE TABLE IF NOT EXISTS shorts (short VARCHAR(255) PRIMARY KEY, original VARCHAR(255) NOT NULL, user_id VARCHAR(255) )`)

	if err != nil {
		return nil, err
	}

	log.Println(result)

	return &database{
		db: db,
	}, err
}

func (d *database) GetByShort(key string) (row storage.URLRecord, err error) {
	err = d.db.Get(&row, "SELECT short, original, user_id FROM shorts WHERE short=$1", key)
	return row, err
}

func (d *database) GetByUID(uid string) (rows []storage.URLRecord, err error) {
	err = d.db.Select(&rows, "SELECT short, original, user_id FROM shorts WHERE user_id = $1", uid)
	return rows, err
}
func (d *database) Save(short, original, userID string) {
	result, err := d.db.Exec("INSERT INTO shorts (short, original, user_id) VALUES ($1, $2, $3)", short, original, userID)
	if err != nil {
		log.Println(err)
	}
	log.Println(result)
}
func (d *database) IsExist(key string) bool {
	return false
}
func (d *database) Lock() {

}
func (d *database) Unlock() {

}

func (d *database) Ping() error {
	return d.db.Ping()
}
