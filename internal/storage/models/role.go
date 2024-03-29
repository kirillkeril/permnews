package models

import (
	"github.com/google/uuid"
	"time"
)

type Role struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
