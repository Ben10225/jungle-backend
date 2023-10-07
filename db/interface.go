package db

import sqlc "jungle-proj/db/sqlc"

type Store interface {
	CreateUser(*sqlc.User) error
	// GetAllUser() ([]*model.User, error)
	// GetUserByID(int) (*model.User, error)
}
