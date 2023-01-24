package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"ice-skating-telegram-bot/internal/app/commander"
	"ice-skating-telegram-bot/internal/app/service"
	"log"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	userServiceToken := os.Getenv("USER_SERVICE_TOKEN")
	productService := service.NewService(userServiceToken)

	config := telegramBotConfig(bot)
	botCommander := commander.NewCommander(bot, productService)

	updates := bot.GetUpdatesChan(config)
	for update := range updates {
		botCommander.HandleUpdates(update)
	}
}

func telegramBotConfig(bot *tgbotapi.BotAPI) tgbotapi.UpdateConfig {
	cfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "/today",
			Description: "Показать доступные временные слоты катания на сегодня",
		},
		tgbotapi.BotCommand{
			Command:     "/tomorrow",
			Description: "Показать доступные временные слоты катания на завтра",
		},
	)

	_, _ = bot.Request(cfg)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.UpdateConfig{
		Offset:  0,
		Limit:   0,
		Timeout: 60,
	}

	return updateConfig
}
