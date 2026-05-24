package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

const (
	flagDel = "__deleted__"
)

var rdb *redis.Client

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при чтении .env:", err)
	}

	token := os.Getenv("TG_TOKEN")
	if token == "" {
		log.Fatal("TG_TOKEN environment variable is not set")
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Ошибка подключения к Redis:", err)
	}

	log.Println("Redis подключен!")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Ошибка подключения к Telegram:", err)
	}

	log.Println(bot.Self.UserName + " запущен!")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID

		var responseText string

		switch {
		case strings.HasPrefix(update.Message.Text, "/add"):
			note := strings.TrimPrefix(update.Message.Text, "/add")
			if note == "" || note == "/add" {
				responseText = "Напиши текст заметки. Например: `/add купить хлеб`"
			} else {
				key := fmt.Sprintf("notes:%d", chatID)
				rdb.RPush(ctx, key, note)
				responseText = "Заметка сохранена"
			}
		case update.Message.Text == "/list":
			key := fmt.Sprintf("notes:%d", chatID)

			notes, err := rdb.LRange(ctx, key, 0, -1).Result()
			if err != nil || len(notes) == 0 {
				responseText = "Заметок пока нет! Добавь: `/add текст`"
			} else {
				responseText = "*Твои заметки:*\n"
				for i, note := range notes {
					responseText += fmt.Sprintf("%d. %s\n", i+1, note)
				}
			}
		case strings.HasPrefix(update.Message.Text, "/delete"):
			numStr := strings.TrimPrefix(update.Message.Text, "/delete")
			numStr = strings.TrimSpace(numStr)

			var num int
			_, err := fmt.Sscan(numStr, &num)
			if err != nil || num < 1 {
				responseText = "Укажи номер заметки. Например: `/delete 2`"
			} else {
				key := fmt.Sprintf("notes:%d", chatID)
				notes, _ := rdb.LRange(ctx, key, 0, -1).Result()
				if num > len(notes) {
					responseText = fmt.Sprintf("Заметки с номером %d нет", num)
				} else {
					rdb.LSet(ctx, key, int64(num-1), flagDel)
					rdb.LRem(ctx, key, 0, flagDel)
					responseText = fmt.Sprintf("Заметка %d удалена!", num)
				}
			}

		case update.Message.Text == "/start":
			responseText = "Привет! Я бот для заметок 📝\n\n*Команды:*\n`/add Текст` — добавить заметку\n`/list` — показать все заметки\n`/delete` — удалить заметку\n`/help` — помощь"

		case update.Message.Text == "/help":
			responseText = "*Команды:*\n`/add Текст` — добавить заметку\n`/list` — показать все заметки\n`/delete` — удалить заметку\n"

		default:
			responseText = "Не понимаю эту команду. Напиши /help"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
		msg.ParseMode = "Markdown"
		bot.Send(msg)
	}
}
