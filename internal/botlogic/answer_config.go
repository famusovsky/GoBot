package botlogic

import (
	"reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AnswerConfig - структура, представляющая конфигурацию сообщения, подлежащего к отправке.
type AnswerConfig struct {
	Text        string      // Текст сообщения.
	Parser      string      // Способ парсинга сообщения: "","Markdown","HTML".
	Keyboard    interface{} // Маркап клавиатуры, приложенной к сообщению.
	ChatId      int64       // Id, по которому необходимо отправить сообщение.
	MessageId   int         // Id сообщения, на которое необходимо ответить
	EditMessage bool        // Boolean, указывающий на то, нужно ли редактировать сообщение MessageId вместо отправки нового
	NextStep    int         // Следующий шаг в сценарии, в котором находится получатель.
	NextScript  string      // Следующий сценарий, на который должен перейти получатель.
	// TODO добавить другие конфиги при необходимости.
}

func (a AnswerConfig) ToMessageConfig() tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(
		a.ChatId,
		a.Text,
	)

	if a.Parser == tgbotapi.ModeMarkdown || a.Parser == tgbotapi.ModeHTML {
		msg.ParseMode = a.Parser
	}

	msg.ReplyMarkup = a.Keyboard // TODO добавить нормальную проверку

	if a.MessageId > 0 {
		msg.ReplyToMessageID = a.MessageId
	}

	return msg
}

func (a AnswerConfig) ToEditMessageConfig() tgbotapi.EditMessageTextConfig {
	msg := tgbotapi.NewEditMessageText(a.ChatId, a.MessageId, a.Text)

	if a.Parser == tgbotapi.ModeMarkdown || a.Parser == tgbotapi.ModeHTML {
		msg.ParseMode = a.Parser
	}
	if reflect.TypeOf(a.Keyboard) == reflect.TypeOf(tgbotapi.InlineKeyboardMarkup{}) {
		keyboard, ok := a.Keyboard.(tgbotapi.InlineKeyboardMarkup)
		if ok && len((keyboard.InlineKeyboard)) > 0 {
			msg.ReplyMarkup = &keyboard // TODO добавить нормальную проверку
		}
	}

	return msg
}
