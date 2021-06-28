package bundle

import (
	"auth2/config"
	"auth2/container/prototype"
	"auth2/services/mail_service"
	"auth2/services/security_service"
	"auth2/services/session_service"
	"auth2/services/telegram_notify_service"
	"auth2/services/tokens_service"
	"auth2/services/user_service"
	"auth2/storages"
	"github.com/doug-martin/goqu/v9"

	// Mysql dialect
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

func CreateServices(config *config.Config) *prototype.ContainerServices {
	db, err := sqlx.Connect(config.DBDialect, config.DBConnection)
	if err != nil {
		panic("Database connection failed: " + err.Error())
	}

	var services Services

	var containerServices prototype.ContainerServices = &services

	userStorage := storages.NewUserStorage(config.UsersTable, db, goqu.Dialect(config.DBDialect))
	userService := user_service.New(userStorage)
	redisStorage := storages.NewRedisStorage(config.RedisConnection, config.RedisPrefix)

	sessionService := session_service.New(
		storages.NewSessionStorage(redisStorage),
		config.RedisSessionUserPrefix,
		config.RedisSessionIdPrefix,
		config.RedisSessionTTL)

	mailService := mail_service.New(mail_service.Options{
		Host:               config.SmtpHost,
		Port:               config.SmtpPort,
		Username:           config.SmtpUsername,
		Password:           config.SmtpPassword,
		Encryption:         config.SmtpEncryption,
		ProjectName:        config.ProjectName,
		DeveloperEmails:    config.MailsDeveloperEmails,
		FromHeader:         config.MailsFromHeader,
		TechSupportSubject: config.MailsTechSupportSubject})

	telegramNotifyService := telegram_notify_service.New(&telegram_notify_service.Options{
		IsEnabled:     config.TelegramNotifyEnabled,
		Token:         config.TelegramNotifyToken,
		ChannelId:     config.TelegramNotifyChannelId,
		FeedbackToken: config.TelegramNotifyFeedbackToken})

	tokenService := tokens_service.New(&tokens_service.ServiceOptions{
		PasswordRequestTtl: config.PasswordRequestTokenTtl,
		ConfirmationTtl:    config.ConfirmationTokenTtl,
		KeyPrefix:          config.TokensPrefix,
		Storage:            redisStorage,
	})

	securityService := security_service.New(
		&security_service.Options{
			PasswordCost: config.UserPasswordCost,
			TokenService: tokenService,
		})

	services.setSecurity(securityService)
	services.setUsers(userService)
	services.setSessions(sessionService)
	services.setMailService(mailService)
	services.setTelegramNotify(telegramNotifyService)
	services.setTokens(tokenService)

	return &containerServices
}
