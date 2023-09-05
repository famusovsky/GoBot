// Пакет botlogic реализует обработку обновлений Telegram бота по адаптации новых сотрудников к работе в компании.
package botlogic

import (
	"database/sql"
	"errors"

	"github.com/famusovsky/GoBot/internal/botlogic/models"
	"github.com/famusovsky/GoBot/internal/botlogic/models/postgres"
	"github.com/famusovsky/GoBot/pkg/tgbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Processor - структура, реализующая интерфейс обработки обновлений Telegram бота.
type Processor struct {
	bot   *tgbotapi.BotAPI       // Telegram бот.
	users *postgres.UsersDBModel // модель базы данных пользователей.
}

// NewProcessor - функция-конструктор для создания нового обработчика обновлений Telegram бота.
//
// Принимает указатель на Telegram бота.
//
// Возвращает обработчик обновлений и ошибку.
func NewProcessor(tgbot *tgbotapi.BotAPI, db *sql.DB) (tgbot.UpdateProcessor, error) {
	if tgbot == nil {
		return Processor{}, errors.New("shall input valid *tgbotapi.BotAPI")
	}

	users, err := postgres.NewUsersDBModel(db)

	// XXX
	// users.Tidy()

	if err != nil {
		return Processor{}, err
	}

	return Processor{
		bot:   tgbot,
		users: users,
	}, nil
}

// Process - метод, реализующий обработку обновлений Telegram бота.
//
// Принимает обновление Telegram бота.
//
// Возвращает ошибку.
func (p Processor) Process(upd tgbotapi.Update) error {
	id := upd.FromChat().ID
	user, err := p.users.Get(id)
	if err != nil || user.Step < 0 {
		user = models.User{Id: id, Script: defaultScript}
		p.users.AddOrUpdate(user.Id, user.Script, user.Step)
	}

	switch {
	case upd.Message != nil:
		return p.handleMessage(upd.Message, user)
	case upd.CallbackQuery != nil:
		return p.handleCallback(upd.CallbackQuery, user)
	}

	return nil
}

// handleMessage - метод, реализующий обработку сообщений Telegram бота.
//
// Принимает сообщение Telegram бота и id чата.
//
// Возвращает ошибку.
func (p Processor) handleMessage(msg *tgbotapi.Message, user models.User) error {
	var (
		ansConfig AnswerConfig
		ok        bool
	)

	if command := msg.Command(); command != "" {
		ansConfig, ok = Commands[command]
		ansConfig.ChatId = user.Id // TODO сделать нормально, через GetAnswer
	} else {
		ansConfig, ok = Scripts.GetAnswer(msg.Text, user)
	}
	if !ok {
		return errors.New("wrong request")
	}

	ans := ansConfig.ToMessageConfig()

	_, err := p.bot.Send(ans)
	if err != nil {
		return err
	}

	return p.updateUser(user, ansConfig)
}

// handleCallback - метод, реализующий обработку callback запросов Telegram бота.
//
// Принимает callback запрос Telegram бота и id чата.
//
// Возвращает ошибку.
func (p Processor) handleCallback(cq *tgbotapi.CallbackQuery, user models.User) error {
	ansConfig, ok := Scripts.GetAnswer(cq.Data, user)
	ansConfig.MessageId = cq.Message.MessageID
	if !ok {
		return errors.New("wrong request")
	}

	tmpCallback := tgbotapi.NewCallback(cq.ID, "")
	tmpEdit := tgbotapi.NewEditMessageText(user.Id, cq.Message.MessageID, cq.Message.Text)
	p.bot.Send(tmpCallback)
	p.bot.Send(tmpEdit)

	var ans tgbotapi.Chattable
	if ansConfig.EditMessage {
		ans = ansConfig.ToEditMessageConfig()
	} else {
		ans = ansConfig.ToMessageConfig()
	}

	_, err := p.bot.Send(ans)
	if err != nil {
		return err
	}

	return p.updateUser(user, ansConfig)
}

func (p Processor) updateUser(user models.User, ansConfig AnswerConfig) error {
	if ansConfig.NextScript != "" {
		user.Script = ansConfig.NextScript
		user.Step = 1
	}
	if ansConfig.NextStep > 0 {
		user.Step = ansConfig.NextStep
	}

	return p.users.AddOrUpdate(user.Id, user.Script, user.Step)
}
