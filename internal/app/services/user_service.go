package services

import (
	"authorization_service/internal/app/models/entity"
	"authorization_service/internal/app/storage/models"
	"errors"
	"github.com/google/uuid"
	"log"
)

type UserRepository interface {
	GetById(id string) (*entity.User, error)
	Create(user models.User) error
	Update(id string, user models.User) error
}

type PasswordService interface {
	GenerateSalt(length uint32) ([]byte, error)
	HashPassword(password *string, salt *[]byte) (*string, error)
	ComparePasswordAndHash(hash *string, password *string) (bool, error)
}

type UserService struct {
	userRepository  UserRepository
	passwordService PasswordService
}

func NewUserService(userRepository UserRepository, passwordService PasswordService) *UserService {
	return &UserService{
		userRepository:  userRepository,
		passwordService: passwordService,
	}
}

func (s *UserService) GetById(id string) (*entity.User, error) {
	user, err := s.userRepository.GetById(id)
	if err != nil {
		log.Println(err)
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) Create(user *entity.User, password string) error {
	salt, err := s.passwordService.GenerateSalt(uint32(len(user.Email)))
	if err != nil {
		return err
	}
	hashedPassword, err := s.passwordService.HashPassword(&password, &salt)
	if err != nil {
		return err
	}
	model := models.User{
		Id:       uuid.MustParse(user.Id),
		Email:    user.Email,
		Password: *hashedPassword,
	}
	return s.userRepository.Create(model)
}
