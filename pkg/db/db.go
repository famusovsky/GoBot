// Пакет для работы с БД
package db

import (
	"database/sql"
	"fmt"
	"os"
)

// OpenViaEnvVars - открытие БД через переменные окружения.
// Возвращает БД и ошибку.
func OpenViaEnvVars() (*sql.DB, error) {
	return OpenViaDsn(getDsnFromEnv())
}

// OpenViaDsn - открытие БД через строку DSN.
// Принимает строку DSN.
// Возвращает БД и ошибку.
func OpenViaDsn(dsn string) (*sql.DB, error) {
	if dsn == "" {
		dsn = getDsnFromEnv()
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// getDsnFromEnv - получение строки DSN из переменных окружения.
func getDsnFromEnv() string {
	dsn := fmt.Sprintf("port=%s user=%s password=%s dbname=%s sslmode=%s host=%s",
		os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_SSLMODE"), os.Getenv("DB_HOST"))

	return dsn
}
