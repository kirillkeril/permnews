package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"permnews/internal/storage/models"
)

type Storage interface {
	GetDb() *sql.DB
}

type UserRepository struct {
	storage Storage
}

func NewUserRepository(storage Storage) UserRepository {
	return UserRepository{storage: storage}
}

func (u UserRepository) GetById(id string) (*models.User, error) {
	db := u.storage.GetDb()
	var model *models.User
	row := db.QueryRow(`SELECT * FROM users WHERE id=$1`, id)
	err := row.Scan(model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (u UserRepository) GetByEmail(email string) (*models.User, error) {
	db := u.storage.GetDb()
	var model *models.User
	row := db.QueryRow(`SELECT * FROM users WHERE email=$1`, email)
	err := row.Scan(model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (u UserRepository) Create(user models.User) error {
	db := u.storage.GetDb()
	_, err := db.Exec(`INSERT INTO users (id, email, password) VALUES ($1, $2, $3)`, user.Id, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u UserRepository) Update(id string, user models.User) error {
	db := u.storage.GetDb()
	_, err := db.Exec(`UPDATE users WHERE id=$1 SET email=$2, password=$3`, user.Id, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}
