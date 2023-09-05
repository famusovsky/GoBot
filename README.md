# GoBot

Шаблон Telegram бота на языке Go.

Описание скриптов бота находится в файлах [scripts_values.go](./internal/botlogic/scripts_values.go), [texts_values.go](./internal/botlogic/texts_values.go), [commands_values.go](./internal/botlogic/commands_values.go) -- предполагается изменение этих файлов для получения рабочего бота.

Запуск:
```bash
# Среда, в которой происходит запуск, должна иметь переменные окружения:
# TGBOT_TOKEN
# DB_HOST
# DB_PORT
# DB_USER
# DB_PASSWORD
# DB_NAME
go run ./cmd/bot/main.go
```

> История коммитов была утеряна ૮(˶╥︿╥)ა