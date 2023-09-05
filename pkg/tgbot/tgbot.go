// Пакет tgbot реализует базовую логику работы с Telegram ботом.
package tgbot

import (
	"database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// UpdateProcessor - интерфейс обработчика обновлений.
type UpdateProcessor interface {
	Process(upd tgbotapi.Update) error
}

// GetUpdateProcessor - функция, возвращающая обработчик обновлений.
//
// Принимает указатель на бота, от лица которого необходимо обрабатывать обновления.
//
// Возвращает обработчик обновлений и ошибку.
type GetUpdateProcessor func(bot *tgbotapi.BotAPI, db *sql.DB) (UpdateProcessor, error)

// Run - функция, запускающая работу бота.
//
// Принимает:
// token - токен API Telegram бота;
// get - функция, принимающая указатель на бота и возвращающая обработчик обновлений и ошибку;
// errorLog - логгер ошибок.
func Run(token string, get GetUpdateProcessor, errorLog *log.Logger, db *sql.DB) {
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		errorLog.Fatalln(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	processor, err := get(bot, db)
	if err != nil {
		errorLog.Fatalln(err)
	}

	for upd := range updates {
		err = processor.Process(upd)
		if err != nil {
			errorLog.Println(err)
		}
	}
}
