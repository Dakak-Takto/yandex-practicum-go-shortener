package database

import (
	"database/sql"
	"log"
	"yandex-practicum-go-shortener/internal/storage"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type database struct {
	db *sql.DB
}

func New(dsn string) (storage.Storage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Query(`CREATE TABLE IF NOT EXISTS shorts (short VARCHAR(255) PRIMARY KEY, original VARCHAR(255) NOT NULL, user_id VARCHAR(255) )`)

	if err != nil {
		return nil, err
	}

	return &database{
		db: db,
	}, err
}

func (d *database) First(key string) (storage.URLRecord, error) {
	var original, userID string

	row := d.db.QueryRow("SELECT original, user_id FROM shorts WHERE short=$1", key)
	err := row.Scan(&original, &userID)
	if err != nil {
		return storage.URLRecord{}, err
	}

	return storage.URLRecord{
		Short:    key,
		Original: original,
		UserID:   userID,
	}, nil
}
func (d *database) Get(key string) []storage.URLRecord {
	rows, err := d.db.Query("SELECT original, user_id FROM shorts WHERE short=$1", key)
	if err != nil {
		return nil
	}
	var result []storage.URLRecord
	var original, userID string
	for err := rows.Scan(&original, &userID); err == nil; {
		result = append(result, storage.URLRecord{
			Original: original,
			UserID:   userID,
			Short:    key,
		})
	}
	return []storage.URLRecord{}
}
func (d *database) Save(short, original, userID string) {
	sqlStr, err := d.db.Query("INSERT INTO shorts (short, original, user_id) VALUES ($1, $2, $3)", short, original, userID)
	if err != nil {
		log.Println(err)
	}
	_ = sqlStr
}
func (d *database) IsExist(key string) bool {
	return false
}
func (d *database) Lock() {

}
func (d *database) Unlock() {

}
