package botlogic

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	Commands = map[string]AnswerConfig{
		startCommand: {
			Text: greeting,
			Keyboard: tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(
						defaultScript,
					),
				),
			),
		},
	}
)
