package tg_bot

import (
	"github.com/Nikitapopov/Habbit/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewBot(token string, logger *logging.Logger) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	bot.Debug = true

	logger.Infof("Authorized on tg account %s", bot.Self.UserName)

	return bot, nil
}
