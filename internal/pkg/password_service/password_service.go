package password_service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"log"
	"strings"
)

type PasswordService struct {
	secret string
}

func NewPasswordService(secret string) *PasswordService {
	return &PasswordService{
		secret: secret,
	}
}

func (s PasswordService) GenerateSalt(length uint32) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return bytes, nil
}

func (s PasswordService) HashPassword(password *string, salt *[]byte) (*string, error) {
	hashedPassword := argon2.IDKey([]byte(*password), *salt, 1, 64*1024, 1, 128)
	base64Hash := base64.RawStdEncoding.EncodeToString(hashedPassword)
	base64Salt := base64.RawStdEncoding.EncodeToString(*salt)

	pwd := fmt.Sprintf("$%s$%s", base64Salt, base64Hash)
	log.Println(pwd)
	return &pwd, nil
}

func (s PasswordService) ComparePasswordAndHash(hash *string, password *string) (bool, error) {
	log.Println(*hash, *password)
	salt := strings.Split(*hash, "$")[1]
	bytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return false, err
	}
	pass, err := s.HashPassword(password, &bytes)
	if err != nil {
		return false, err
	}

	return strings.EqualFold(*pass, *hash), nil
}
