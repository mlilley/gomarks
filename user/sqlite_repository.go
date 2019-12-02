package user

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) (Repository, error) {
	return &repo{db: db}, nil
}

func (r *repo) FindAll() ([]*User, error) {
	rows, err := r.db.Query("SELECT id, email, password_hash, active FROM user ORDER BY email")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Active)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *repo) FindByID(id string) (*User, error) {
	var user User

	//sid, err := strconv.Atoi(id)
	//if err != nil {
	//	return nil, err
	//}

	err := r.db.
		QueryRow("SELECT id, email, password_hash, active FROM user WHERE id = ?", id).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Active)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repo) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.
		QueryRow("SELECT id, email, password_hash, active FROM user WHERE email = ?", strings.ToLower(email)).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Active)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}


func (r *repo) Create(user *User) (*User, error) {
	email := strings.ToLower(user.Email)
	result, err := r.db.Exec(
		"INSERT INTO user(email, password_hash, active) VALUES (?, ?, ?)",
		email, user.PasswordHash, user.Active)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.ID = string(id)
	return user, nil
}

func (r *repo) Update(user *User) error {
	email := strings.ToLower(user.Email)
	result, err := r.db.Exec(
		"UPDATE user SET email = ?, password_hash = ?, active = ? WHERE id = ?",
		email, user.PasswordHash, user.Active, user.ID)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *repo) DeleteByID(id string) (bool, error) {
	result, err := r.db.Exec("DELETE FROM user WHERE id = ?", id)
	if err != nil {
		return false, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return n != 0, nil
}

func (r *repo) DeleteByEmail(email string) (bool, error) {
	email = strings.ToLower(email)
	result, err := r.db.Exec("DELETE FROM user WHERE email = ?", email)
	if err != nil {
		return false, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return n != 0, nil
}

func (r *repo) DeleteAll() error {
	_, err := r.db.Exec("DELETE FROM user")
	if err != nil {
		return err
	}

	return nil
}