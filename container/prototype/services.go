package prototype

import (
	"auth2/services/mail_service"
	"auth2/services/security_service"
	"auth2/services/session_service"
	"auth2/services/telegram_notify_service"
	"auth2/services/tokens_service"
	"auth2/services/user_service"
)

type ContainerServices interface {
	Security() *security_service.Service
	Users() *user_service.Service
	Sessions() *session_service.SessionsService
	Mail() *mail_service.MailService
	TelegramNotify() *telegram_notify_service.GenericService
	Tokens() *tokens_service.Service
}
