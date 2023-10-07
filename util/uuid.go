package util

import "github.com/google/uuid"

func UuidGenerate() string {
	return uuid.NewString()
}
