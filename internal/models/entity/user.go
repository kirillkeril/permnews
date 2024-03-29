package entity

import "github.com/google/uuid"

type User struct {
	Id    string
	Email string
	Roles []Role
}

func NewUser(email string) *User {
	return &User{
		Id:    uuid.New().String(),
		Email: email,
		Roles: nil,
	}
}
