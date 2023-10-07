package db

import (
	"database/sql"
	"errors"
	sqlc "jungle-proj/db/sqlc"
	"jungle-proj/util"
	"time"
)

type MemoryStorage struct {
	users []*sqlc.User
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (s *MemoryStorage) CreateUser(user *sqlc.User) error {
	s.users = append(s.users, &sqlc.User{
		Uuid:       util.UuidGenerate(),
		Name:       user.Name,
		CreateTime: sql.NullTime{Time: time.Now().Add(time.Hour * 8), Valid: true},
	})

	return nil
}

func (s *MemoryStorage) GetAllUser() ([]*sqlc.User, error) {
	if len(s.users) > 0 {
		return s.users, nil
	}
	return nil, errors.New("no data")
}

func (s *MemoryStorage) GetUserByID(uuid string) (*sqlc.User, error) {
	for _, v := range s.users {
		if v.Uuid == uuid {
			return v, nil
		}
	}

	return nil, errors.New("can't find user")
}
