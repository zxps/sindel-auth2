package tokens_service

import (
	"auth2/models"
	"auth2/storages"
	"auth2/utils"
	"fmt"
	"strconv"
	"time"
)

const (
	UserConfirmationPrefix = "confirm"
	PasswordRequestPrefix  = "password"
)

type ServiceOptions struct {
	Storage            *storages.RedisStorage
	KeyPrefix          string
	ConfirmationTtl    int
	PasswordRequestTtl int
}

type Service struct {
	storage           *storages.RedisStorage
	keyPrefix         string
	confirmationTtl   int
	passwordRequesTtl int
}

func New(opts *ServiceOptions) *Service {
	return &Service{
		storage:           opts.Storage,
		keyPrefix:         opts.KeyPrefix,
		confirmationTtl:   opts.ConfirmationTtl,
		passwordRequesTtl: opts.PasswordRequestTtl,
	}
}

// {{Password request token methods}}

func (s *Service) SetPasswordRequest(user *models.User, token string) (string, error) {
	return s.setToken(PasswordRequestPrefix, token, user.Id, s.passwordRequesTtl)
}

func (s *Service) FindPasswordRequest(token string) (bool, int) {
	status, value := s.findToken(PasswordRequestPrefix, token)
	userId, _ := strconv.Atoi(value)
	return status, userId
}

func (s *Service) DeletePasswordRequest(token string) (int64, error) {
	return s.deleteToken(PasswordRequestPrefix, token)
}

func (s *Service) IsPasswordRequestExists(token string) bool {
	return s.isTokenExists(PasswordRequestPrefix, token)
}

func (s *Service) GeneratePasswordRequest(user *models.User) string {
	return s.generateToken(PasswordRequestPrefix, user, s.passwordRequesTtl)
}

// {{ /Password request token methods }}

// {{ User confirmation token methods }}

func (s *Service) SetUserConfirmation(user *models.User, token string) (string, error) {
	return s.setToken(UserConfirmationPrefix, token, user.Id, s.confirmationTtl)
}

func (s *Service) FindUserConfirmation(token string) (bool, int) {
	status, value := s.findToken(UserConfirmationPrefix, token)
	userId, _ := strconv.Atoi(value)
	return status, userId
}

func (s *Service) DeleteUserConfirmation(token string) (int64, error) {
	return s.deleteToken(UserConfirmationPrefix, token)
}

func (s *Service) IsUserConfirmationExists(token string) bool {
	return s.isTokenExists(UserConfirmationPrefix, token)
}

func (s *Service) GenerateUserConfirmation(user *models.User) string {
	return s.generateToken(UserConfirmationPrefix, user, s.confirmationTtl)
}

// {{ /User confirmation token methods }}

func (s *Service) setToken(prefixType string, token string, value interface{}, ttl int) (string, error) {
	tokenKey := s.buildTokenKey(prefixType, token)
	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", ttl))
	return s.storage.Set(tokenKey, value, duration)
}

func (s *Service) isTokenExists(prefixType string, token string) bool {
	return s.storage.Has(s.buildTokenKey(prefixType, token))
}

func (s *Service) deleteToken(tokenType string, token string) (int64, error) {
	return s.storage.Delete(s.buildTokenKey(tokenType, token))
}

func (s *Service) findToken(typePrefix string, token string) (bool, string) {
	tokenKey := s.buildTokenKey(typePrefix, token)
	result, err := s.storage.Get(tokenKey)
	if err != nil || len(result) < 1 {
		return false, ""
	}

	return true, result
}

func (s *Service) generateToken(typePrefix string, user *models.User, ttl int) string {
	maxRetry := 100
	for {
		token := utils.SimpleRandomToken(64)
		if !s.isTokenExists(typePrefix, token) {
			s.setToken(typePrefix, token, user.Id, ttl)
			return token
		}

		maxRetry--
		if maxRetry <= 0 {
			break
		}
	}

	return ""
}

func (s *Service) buildTokenKey(typePrefix string, token string) string {
	return fmt.Sprintf("%s%s:%s", s.keyPrefix, typePrefix, token)
}
