package model

import "time"

type User struct {
	ID       int
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Time     time.Time
}
