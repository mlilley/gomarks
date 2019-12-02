package database

import (
	"database/sql"
	"github.com/mlilley/gomarks/auth"
)

func New() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./gomarks.sqlite")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS mark (id INTEGER PRIMARY KEY, title TEXT, url TEXT)")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY, email TEXT UNIQUE, password_hash TEXT, active INTEGER)")
	if err != nil {
		return nil, err
	}

	hash, err := auth.HashPassword("password")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("INSERT INTO user (email, password_hash, active) VALUES (?, ?, ?) ON CONFLICT(email) DO UPDATE SET password_hash = ?, active = ?", "test@example.com", hash, 1, hash, 1)
	if err != nil {
		return nil, err
	}

	return db, nil
}