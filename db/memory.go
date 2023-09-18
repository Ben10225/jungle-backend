package db

import (
	"errors"
	"jungle-proj/model"
	"time"
)

type MemoryStorage struct {
	users []*model.User
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (s *MemoryStorage) AddUser(user *model.User) error {
	s.users = append(s.users, &model.User{
		ID:   len(s.users) + 1,
		Name: user.Name,
		Time: time.Now(),
	})

	return nil
}

func (s *MemoryStorage) GetAllUser() ([]*model.User, error) {
	if len(s.users) > 0 {
		return s.users, nil
	}
	return nil, errors.New("no data")
}

func (s *MemoryStorage) GetUserByID(i int) (*model.User, error) {
	for _, v := range s.users {
		if v.ID == i {
			return v, nil
		}
	}

	return nil, errors.New("can't find user")
}
