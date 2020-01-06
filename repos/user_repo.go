package repos

import (
	"database/sql"
	"github.com/mlilley/gomarks/app"
	"strings"
)

type UserRepo interface {
	FindAll() ([]app.User, error)
	FindByID(id string) (*app.User, error)
	FindByEmail(email string) (*app.User, error)
	Create(u *app.User) (*app.User, error)
	Update(u *app.User) error
	DeleteByID(id string) (bool, error)
	DeleteByEmail(email string) (bool, error)
	DeleteAll() error
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &sqliteUserRepo{db: db}
}

type sqliteUserRepo struct {
	db *sql.DB
}

func (r *sqliteUserRepo) FindAll() ([]app.User, error) {
	rows, err := r.db.Query("SELECT id, email, password_hash, active FROM user ORDER BY email")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []app.User{}
	for rows.Next() {
		var user app.User
		err = rows.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Active)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *sqliteUserRepo) FindByID(id string) (*app.User, error) {
	var user app.User

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

func (r *sqliteUserRepo) FindByEmail(email string) (*app.User, error) {
	var user app.User
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


func (r *sqliteUserRepo) Create(user *app.User) (*app.User, error) {
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

func (r *sqliteUserRepo) Update(user *app.User) error {
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

func (r *sqliteUserRepo) DeleteByID(id string) (bool, error) {
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

func (r *sqliteUserRepo) DeleteByEmail(email string) (bool, error) {
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

func (r *sqliteUserRepo) DeleteAll() error {
	_, err := r.db.Exec("DELETE FROM user")
	if err != nil {
		return err
	}

	return nil
}
