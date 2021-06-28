package config

// Config ...
type Config struct {
	ConfigPath string

	ProjectName string

	DBConnection  string
	DBDialect     string
	ListenAddress string
	UsersTable    string
	FilesTable    string

	RedisConnection        string
	RedisPrefix            string
	RedisSessionIdPrefix   string
	RedisSessionUserPrefix string
	RedisSessionTTL        int

	SmtpHost       string
	SmtpPort       int
	SmtpUsername   string
	SmtpPassword   string
	SmtpEncryption string

	UserPasswordMinLength int
	UserPasswordMaxLength int
	UserPasswordCost      int

	MailsFromHeader         string
	MailsTechSupportSubject string
	MailsAdminEmails        []string
	MailsDeveloperEmails    []string

	TokensPrefix            string
	ConfirmationTokenTtl    int
	PasswordRequestTokenTtl int

	UserConfirmPasswordUrlPattern string
	UserConfirmUrlPattern         string

	TelegramNotifyEnabled       bool
	TelegramNotifyToken         string
	TelegramNotifyChannelId     int64
	TelegramNotifyFeedbackToken string
}
