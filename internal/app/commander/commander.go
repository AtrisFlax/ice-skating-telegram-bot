package commander

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ice-skating-telegram-bot/internal/app/service"
	"log"
)

type Commander struct {
	bot     *tgbotapi.BotAPI
	service *service.Service
}

func NewCommander(bot *tgbotapi.BotAPI, productService *service.Service) *Commander {
	return &Commander{
		bot:     bot,
		service: productService,
	}
}

func (c *Commander) SendMessage(messageChatId int64, messageText string) {
	msg := tgbotapi.NewMessage(messageChatId, messageText)
	_, _ = c.bot.Send(msg)
}

func (c *Commander) SendMessageWithKeyboard(messageChatId int64, messageText string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(messageChatId, messageText)
	msg.ReplyMarkup = keyboard
	_, _ = c.bot.Send(msg)
}

func (c *Commander) HandleUpdates(update tgbotapi.Update) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			log.Printf("Recovered from panic %v", panicValue)
		}
	}()

	if update.Message != nil {
		switch update.Message.Command() {
		case "today":
			c.Today(update.Message)
		case "tomorrow":
			c.Tomorrow(update.Message)
		default:
			c.Default(update.Message)
		}
	}
}
