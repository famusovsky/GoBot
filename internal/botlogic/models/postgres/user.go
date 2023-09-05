// Пакет postgres реализует взаимодействие с базой данных PostgreSQL, хранящей данные из пакета models.
package postgres

import (
	"database/sql"
	"errors"
	"log"

	"github.com/famusovsky/GoBot/internal/botlogic/models"
)

// UsersDBModel - модель базы данных пользователей.
type UsersDBModel struct {
	db *sql.DB
}

// NewUsersDBModel - функция-конструктормодели базы данных пользователей.
// Принимает базу данных.
// Возвращает модель базы данных пользователей и ошибку.
func NewUsersDBModel(db *sql.DB) (*UsersDBModel, error) {
	err := checkTable(db)
	if err != nil {
		return nil, err
	}

	return &UsersDBModel{db}, nil
}

// createTable - создание таблицы пользователей в базе данных.
func createTable(db *sql.DB) error {
	q := `CREATE TABLE users (
        id BIGINT NOT NULL PRIMARY KEY, 
		script TEXT NOT NULL,
		step INTEGER NOT NULL
    );`

	_, err := db.Exec(q)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// checkTable - проверка таблицы пользователей в базе данных.
func checkTable(db *sql.DB) error {
	q :=
		`SELECT COUNT(*) = 3 AS proper
			FROM information_schema.columns
			WHERE table_schema = 'public'
			AND table_name = 'users'
			AND (
				(column_name = 'id' AND data_type = 'bigint')
				OR (column_name = 'script' AND data_type = 'text')
				OR (column_name = 'step' AND data_type = 'integer')
			);`
	var proper bool
	db.QueryRow(q).Scan(&proper)

	if !proper {
		// TODO сделать нормальную миграцию
		q = `DROP TABLE IF EXISTS users`
		_, err := db.Exec(q)
		if err != nil {
			return errors.New("cannot drop incorrect 'users' table in the database")
		}

		err = createTable(db)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get - получение пользователя по id.
// Принимает id пользователя.
// Возвращает модель пользователя и ошибку.
func (m *UsersDBModel) Get(id int64) (models.User, error) {
	q :=
		`SELECT script, step FROM users 
		WHERE id = $1;`

	user := models.User{Id: id}

	err := m.db.QueryRow(q, id).Scan(&user.Script, &user.Step)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, models.ErrNoRecord
		}
		return models.User{}, err
	}

	return user, nil
}

// AddOrUpdate - добавление пользователя или обновление, при наличие пользователя с данным id.
// Принимает id пользователя, его текущий сценарий и шаг.
// Возвращает ошибку.
func (m *UsersDBModel) AddOrUpdate(id int64, script string, step int) error {
	_, err := m.Get(id)
	var q string
	if err != nil {
		q =
			`INSERT INTO users (id, script, step) VALUES (
		$1, $2, $3 )`
	} else {
		q = `UPDATE users
		SET script = $2,
			step = $3
		WHERE id = $1;`
	}

	_, err = m.db.Exec(q, id, script, step)
	return err
}

// delete - удаление пользователя по id.
func (m *UsersDBModel) delete(id int) error {
	q :=
		`DELETE FROM users
		WHERE id = $1;`

	_, err := m.db.Exec(q, id)
	return err
}

// Tidy - удаление всех пользователей.
// Возвращает ошибку.
func (m *UsersDBModel) Tidy() error {
	q := `SELECT id FROM users;`
	rows, err := m.db.Query(q)
	if err != nil {
		return nil
	}
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return err
		}
		err = m.delete(id)
		if err != nil {
			return err
		}
	}

	return nil
}
