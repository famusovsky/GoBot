package botlogic

import "github.com/famusovsky/GoBot/internal/botlogic/models"

// ScriptList - список сценариев Telegram-бота.
type ScriptList struct {
	scripts map[string]*Script
}

// NewScriptList - функция-конструктор для создания нового списка сценариев Telegram-бота.
//
// Принимает словарь: название сценария - структура-сценарий.
//
// Возвращает список сценариев.
func NewScriptList(scripts map[string]*Script) ScriptList {
	return ScriptList{scripts}
}

// GetAnswer - метод, возвращающий искомый ответ на запрос пользователя.
//
// Принимает запрос и структуру-пользователя.
//
// Возвращает конфигурацию ответа и ok.
func (sl ScriptList) GetAnswer(request string, user models.User) (AnswerConfig, bool) {
	s, ok := sl.scripts[user.Script]
	if !ok {
		return AnswerConfig{}, ok
	}

	ans, ok := s.GetAnswer(request, user)
	if !ok {
		user.Step = 0
		for _, script := range sl.scripts {
			ans, ok = script.GetAnswer(request, user)
			if ok {
				break
			}
		}
	}

	return ans, ok
}

// Script - структура, представляющая сценарий Telegram-бота.
type Script struct {
	answers []map[string]AnswerConfig // слайс: шаг в сценарии - словарь: запрос - конфигурация ответа.
}

// NewScriptList - функция-конструктор для создания нового сценария Telegram-бота.
//
// Принимает слайс: шаг в сценарии - словарь: запрос - конфигурация ответа.
//
// Возвращает сценарий.
func NewScript(answers []map[string]AnswerConfig) *Script {
	return &Script{answers}
}

// GetAnswer - метод, возвращающий искомый ответ на запрос пользователя.
//
// Принимает запрос, id пользователя и шаг в сценарии, на котором находится пользователь.
//
// Возвращает конфигурацию ответа и ok.
func (s Script) GetAnswer(request string, user models.User) (AnswerConfig, bool) {
	if user.Step >= len(s.answers) {
		return AnswerConfig{}, false
	}

	ans, ok := s.answers[user.Step][request]
	if !ok {
		ans, ok = s.answers[user.Step][""]
		if !ok {
			return ans, ok
		}
	}

	ans.ChatId = user.Id
	return ans, ok
}
