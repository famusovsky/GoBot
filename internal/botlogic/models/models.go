// Пакет models содержит определения данных для хранения в базе данных.
package models

import (
	"errors"
)

// ErrNoRecord - ошибка, которая возникает, когда запись не найдена.
var ErrNoRecord = errors.New("models: record is not found")

// User - структура, представляющая информацию о пользователе.
type User struct {
	Id     int64  // ID пользователя.
	Script string // Сценарий, в котором в данный момент находится пользователь.
	Step   int    // Шаг в сценарии, на котором в данный момент находится пользователь.
}
