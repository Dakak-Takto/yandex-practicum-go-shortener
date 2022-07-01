package repo

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/short/model"
)

type postgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepository(dsn string) (model.ShortRepository, error) {

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = initTables(db)
	if err != nil {
		return nil, err
	}

	return &postgresRepo{
		db: db,
	}, nil
}

func initTables(db *sqlx.DB) error {

	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	return nil
}

func (p *postgresRepo) GetOneByKey(key string) (*model.Short, error) {

	var short model.Short
	if err := p.db.Get(&short, "SELECT _key, _location, user_id, deleted FROM shorts WHERE _key = $1", key); err != nil {
		return nil, err
	}

	return &short, nil
}

func (p *postgresRepo) GetOneByLocation(location string) (*model.Short, error) {

	var short model.Short
	err := p.db.Get(&short, "SELECT _key, _location, user_id, deleted FROM shorts WHERE _location = $1", location)
	if err != nil {
		return nil, err
	}

	return &short, nil
}

func (p *postgresRepo) GetByUserID(userID string) ([]*model.Short, error) {

	var shorts []*model.Short
	err := p.db.Select(&shorts, "SELECT _key, _location, user_id, deleted FROM shorts WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	return shorts, nil
}

func (p *postgresRepo) Insert(short *model.Short) error {

	_, err := p.db.NamedExec(`INSERT INTO shorts (_key, _location, user_id) VALUES (:_key, :_location, :user_id)`, short)
	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return ErrDuplicate
			}
		}
		return err
	}
	return nil
}

func (p *postgresRepo) Delete(keys ...string) error {

	_, err := p.db.Exec(`UPDATE shorts SET deleted=true WHERE key IN ANY($1)`, keys)
	if err != nil {
		return err
	}

	return nil
}

const schema string = `
CREATE TABLE IF NOT EXISTS shorts (
    _key VARCHAR(255) PRIMARY KEY,
    _location VARCHAR(255) UNIQUE,
    user_id VARCHAR(255),
    deleted BOOLEAN DEFAULT false
)`
