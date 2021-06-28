package bundle

import (
	"auth2/container/prototype"
	"auth2/services/mail_service"
	"auth2/services/security_service"
	"auth2/services/session_service"
	"auth2/services/telegram_notify_service"
	"auth2/services/tokens_service"
	"auth2/services/user_service"
)

type Services struct {
	prototype.ContainerServices

	security       *security_service.Service
	users          *user_service.Service
	sessions       *session_service.SessionsService
	mail           *mail_service.MailService
	tokens         *tokens_service.Service
	telegramNotify *telegram_notify_service.GenericService
}

func (s *Services) Tokens() *tokens_service.Service {
	return s.tokens
}

func (s *Services) Security() *security_service.Service {
	return s.security
}

func (s *Services) Users() *user_service.Service {
	return s.users
}

func (s *Services) Sessions() *session_service.SessionsService {
	return s.sessions
}

func (s *Services) Mail() *mail_service.MailService {
	return s.mail
}

func (s *Services) TelegramNotify() *telegram_notify_service.GenericService {
	return s.telegramNotify
}

func (s *Services) setUsers(u *user_service.Service) {
	s.users = u
}

func (s *Services) setSecurity(security *security_service.Service) {
	s.security = security
}

func (s *Services) setTokens(tokens *tokens_service.Service) {
	s.tokens = tokens
}

func (s *Services) setSessions(sessions *session_service.SessionsService) {
	s.sessions = sessions
}

func (s *Services) setMailService(mail *mail_service.MailService) {
	s.mail = mail
}

func (s *Services) setTelegramNotify(telegramNotify *telegram_notify_service.GenericService) {
	s.telegramNotify = telegramNotify
}
