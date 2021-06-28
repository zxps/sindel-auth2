package security_service

import (
	"auth2/services/tokens_service"
	"golang.org/x/crypto/bcrypt"
)

type Options struct {
	PasswordCost int
	TokenService *tokens_service.Service
}

type Service struct {
	passwordCost int
	tokenService *tokens_service.Service
}

func New(opts *Options) *Service {
	return &Service{
		passwordCost: opts.PasswordCost,
		tokenService: opts.TokenService,
	}
}

func (s *Service) HashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.passwordCost)
	return bytes, err
}

func (s *Service) HashPasswordAsString(password string) (string, error) {
	bytes, err := s.HashPassword(password)
	return (string)(bytes), err
}

func (s *Service) CompareStringPasswords(hashedPassword string, password string) error {
	return s.ComparePasswords([]byte(hashedPassword), []byte(password))
}

func (s *Service) ComparePasswords(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
