package model

type (
	Short struct {
		Key      string `json:"short_url" db:"_key"`
		Location string `json:"original_url" db:"_location"`
		UserID   string `json:"-" db:"user_id"`
		Deleted  bool   `json:"-" db:"deleted"`
	}

	ShortUsecase interface {
		FindByKey(key string) (*Short, error)
		FindByLocation(location string) (*Short, error)
		GetUserShorts(userID string) ([]*Short, error)
		CreateNewShort(location string, userID string) (*Short, error)
		Save(*Short) error
		Delete(key ...string) error
	}

	ShortRepository interface {
		GetOneByKey(key string) (*Short, error)
		GetOneByLocation(location string) (*Short, error)
		GetByUserID(userID string) ([]*Short, error)
		Insert(*Short) error
		Delete(key ...string) error
	}
)
