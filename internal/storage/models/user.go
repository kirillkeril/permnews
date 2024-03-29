package models

import (
	"authorization_service/internal/models/entity"
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id        uuid.UUID
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (user *User) ToEntity() *entity.User {
	return &entity.User{
		Id:    user.Id.String(),
		Email: user.Email,
	}
}
