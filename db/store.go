package db

import "jungle-proj/model"

type Store interface {
	AddUser(*model.User) error
	GetAllUser() ([]*model.User, error)
	GetUserByID(int) (*model.User, error)
}
