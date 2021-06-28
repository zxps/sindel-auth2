package server

import (
	"auth2/config"
	"auth2/container/prototype"
	pb "auth2/proto"
	"auth2/services/mail_service"
	"auth2/services/security_service"
	"auth2/services/session_service"
	"auth2/services/telegram_notify_service"
	"auth2/services/tokens_service"
	"auth2/services/user_service"
	"google.golang.org/grpc/codes"
)

const (
	UserNotEnabledStatusCode      codes.Code = 600
	UserNotModeratedStatusCode    codes.Code = 601
	UserNotConfirmedStatusCode    codes.Code = 602
	UserConfirmationTokenNotMatch codes.Code = 603
	UserAlreadyConfirmed          codes.Code = 604
)

type AuthServer struct {
	pb.UnimplementedAuth2Server

	container *prototype.Container
}

func New(c *prototype.Container) *AuthServer {
	return &AuthServer{
		container: c,
	}
}

func (s *AuthServer) reloadContainer(c *prototype.Container) {
	s.container = c
}

func (s *AuthServer) getConfig() *config.Config {
	return s.container.Config()
}

func (s *AuthServer) getUserService() *user_service.Service {
	return (*s.container.Services()).Users()
}

func (s *AuthServer) getSessionService() *session_service.SessionsService {
	return (*s.container.Services()).Sessions()
}

func (s *AuthServer) getSecurityService() *security_service.Service {
	return (*s.container.Services()).Security()
}

func (s *AuthServer) getMailService() *mail_service.MailService {
	return (*s.container.Services()).Mail()
}

func (s *AuthServer) getTelegramNotify() *telegram_notify_service.GenericService {
	service := (*s.container.Services()).TelegramNotify()
	return service
}

func (s *AuthServer) getTokenService() *tokens_service.Service {
	return (*s.container.Services()).Tokens()
}

func (s *AuthServer) sendSupportNotify(message string) {
	(*s.getTelegramNotify()).NotifySupport(message)
}
