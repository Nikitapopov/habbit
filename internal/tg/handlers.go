package tg

import (
	"encoding/json"
	"io/ioutil"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var phrases map[string]string

const (
	responsesJsonFile = "./responses.json"

	commandStart            = "start"
	commandAddHabit         = "add_habit"
	commandShowHabits       = "show_habits"
	commandDeleteHabit      = "delete_habit"
	commandSetNotifications = "set_notifications"
)

func (c *client) init() {
	data, err := c.readFile(responsesJsonFile)
	if err != nil {
		panic(err)
	}

	phrases = data
}

func (c *client) handleCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	switch message.Command() {
	case commandStart:
		msg.Text = "You've entered 'start' command"
		_, err := c.bot.Send(msg)
		return err
	case commandAddHabit:
		msg.Text = "You've entered 'add_habit' command"
		_, err := c.bot.Send(msg)
		return err
	case commandShowHabits:
		msg.Text = "You've entered 'show_habits' command"
		_, err := c.bot.Send(msg)
		return err
	case commandDeleteHabit:
		msg.Text = "You've entered 'delete_habit' command"
		_, err := c.bot.Send(msg)
		return err
	case commandSetNotifications:
		msg.Text = "You've entered 'set_notifications' command"
		_, err := c.bot.Send(msg)
		return err
	default:
		msg.Text = "I don't know such command"
		_, err := c.bot.Send(msg)
		return err
	}
}

func (c *client) handleMessage(message *tgbotapi.Message) {
	c.logger.Infof("[%s], %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	c.bot.Send(msg)
}

func (c *client) readFile(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		c.logger.Errorf("Error during reading file '%s': %v", filePath, err)
		return nil, err
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	var phrases map[string]string
	json.Unmarshal(byteValue, &phrases)

	return phrases, nil
}
