package telegram_notify_service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type Options struct {
	IsEnabled     bool
	Token         string
	ChannelId     int64
	FeedbackToken string
}

type GenericService interface {
	NotifySupport(string)
}

type Service struct {
	token         string
	channelId     int64
	fallbackToken string
	bot           tgbotapi.BotAPI
}

func New(opts *Options) *GenericService {
	if !opts.IsEnabled {
		return NewDisabledService()
	}

	bot, err := tgbotapi.NewBotAPI(opts.Token)
	if err != nil {
		logrus.Errorf("unable to create telegram bot (%s)", err.Error())
	}

	logrus.Infof("Authorize telegram on on account %s", bot.Self.UserName)

	var service GenericService = &Service{
		token:         opts.Token,
		channelId:     opts.ChannelId,
		fallbackToken: opts.FeedbackToken,
		bot:           *bot,
	}

	return &service
}

func (s *Service) NotifySupport(m string) {
	message := tgbotapi.NewMessage(s.channelId, m)
	_, err := s.bot.Send(message)
	if err != nil {
		logrus.Errorf("unable to send message to telegram channel %d", s.channelId)
	}
}
