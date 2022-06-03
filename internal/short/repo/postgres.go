package repo

import (
	"yandex-practicum-go-shortener/internal/short/model"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type postgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepository(dsn string) (model.ShortRepository, error) {

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &postgresRepo{
		db: db,
	}, nil
}

func (p *postgresRepo) GetOneByKey(key string) (*model.Short, error) {

	var short model.Short
	if err := p.db.Get(&short, "SELECT key, location, user_id, deleted FROM urls WHERE key = $1", key); err != nil {
		return nil, err
	}

	return &short, nil
}

func (p *postgresRepo) GetOneByLocation(location string) (*model.Short, error) {

	var short model.Short
	err := p.db.Get(&short, "SELECT key, location, user_id, deleted FROM urls WHERE location = $1", location)
	if err != nil {
		return nil, err
	}

	return &short, nil
}

func (p *postgresRepo) GetByUserID(userID string) ([]*model.Short, error) {

	var shorts []*model.Short
	err := p.db.Select(&shorts, "SELECT key, location, user_id, deleted FROM urls WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	return shorts, nil
}

func (p *postgresRepo) Insert(short *model.Short) error {

	_, err := p.db.NamedExec(`INSERT INTO urls (key, location, user_id) VALUES (:key, :location, :user_id)`, short)
	if err != nil {
		return err
	}
	return nil
}

func (p *postgresRepo) Delete(keys ...string) error {

	_, err := p.db.Exec(`UPDATE urls SET deleted=true WHERE key IN ANY($1)`, keys)
	if err != nil {
		return err
	}

	return nil
}
