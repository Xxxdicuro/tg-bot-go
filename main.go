package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при чтении .env:", err)
	}

	token := os.Getenv("TG_TOKEN")
	if token == "" {
		log.Fatal("TG_TOKEN environment variable is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Ошибка подключения к Telegram:", err)
	}

	log.Println("Bot запущен!")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		var responseText string

		switch update.Message.Text {
		case "/start":
			responseText = "Привет! Я бот для заметок 📝\n\n*Команды:*\n`/add Текст` — добавить заметку\n`/list` — показать все заметки\n`/help` — помощь"
		case "/help":
			responseText = "*Команды:*\n`/add Текст` — добавить заметку\n`/list` — показать все заметки"
		default:
			responseText = "Не понимаю эту команду. Напиши /help"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
		msg.ParseMode = "Markdown"
		bot.Send(msg)
	}
}
