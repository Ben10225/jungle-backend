package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID     uint64
	CustomerID  uuid.UUID
	OrderStatus string
	CreateAt    *time.Time
}
