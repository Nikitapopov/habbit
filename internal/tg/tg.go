package tg

import (
	"github.com/Nikitapopov/Habbit/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type client struct {
	bot    *tgbotapi.BotAPI
	logger *logging.Logger
}

func NewClient(bot *tgbotapi.BotAPI, logger *logging.Logger) *client {
	return &client{
		bot:    bot,
		logger: logger,
	}
}

func (c *client) Start() {
	updates := c.initUpdatesChannel()
	go c.handleUpdates(updates)
}

func (c *client) initUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return c.bot.GetUpdatesChan(u)
}

func (c *client) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			c.handleCommand(update.Message)
			continue
		}

		c.handleMessage(update.Message)
	}
}
