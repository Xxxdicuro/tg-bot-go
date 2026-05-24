# tg-notes-bot

Telegram бот для заметок на Go с хранением данных в Redis.

## Технологии
- Go
- Redis
- Docker
- Telegram Bot API

## Команды
- `/add Текст` — добавить заметку
- `/list` — показать все заметки
- `/delete N` — удалить заметку по номеру

## Запуск

1. Создай `.env` файл:
TG_TOKEN=твой_токен

2. Запусти Redis:
docker run -d --name redis-notes -p 6379:6379 -v redis-data:/data redis redis-server --appendonly yes

3. Запусти бота:
go run main.go
