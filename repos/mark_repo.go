package repos

import (
	"database/sql"
	"github.com/mlilley/gomarks/app"
	"strconv"
)

type MarkRepo interface {
	FindAll() ([]app.Mark, error)
	FindAllByUserID(id string) ([]app.Mark, error)
	FindByID(id string) (*app.Mark, error)
	Create(m *app.Mark) (*app.Mark, error)
	Update(m *app.Mark) error
	DeleteByID(id string) (bool, error)
	DeleteAll() error
}

func NewMarkRepo(db *sql.DB) MarkRepo {
	return &sqliteMarkRepo{db: db}
}

type sqliteMarkRepo struct {
	db *sql.DB
}

func (r *sqliteMarkRepo) FindAll() ([]app.Mark, error) {
	rows, err := r.db.Query("SELECT id, title, url FROM mark")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	marks := []app.Mark{}

	for rows.Next() {
		var mark app.Mark
		err = rows.Scan(&mark.ID, &mark.Title, &mark.URL)
		if err != nil {
			return nil, err
		}
		marks = append(marks, mark)
	}

	return marks, nil
}

func (r *sqliteMarkRepo) FindByID(id string) (*app.Mark, error) {
	var mark app.Mark

	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow("SELECT id, title, url FROM mark WHERE id = ?", intId).Scan(&mark.ID, &mark.Title, &mark.URL)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &mark, nil
}

func (r *sqliteMarkRepo) FindAllByUserID(id string) ([]app.Mark, error) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query("SELECT id, title, url FROM mark WHERE user_id = ?", intId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	marks := []app.Mark{}

	for rows.Next() {
		var mark app.Mark
		err = rows.Scan(&mark.ID, &mark.Title, &mark.URL)
		if err != nil {
			return nil, err
		}
		marks = append(marks, mark)
	}

	return marks, nil
}

func (r *sqliteMarkRepo) Create(mark *app.Mark) (*app.Mark, error) {
	result, err := r.db.Exec("INSERT INTO mark(title, url) VALUES (?, ?)", mark.Title, mark.URL)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	mark.ID = strconv.FormatInt(id, 10)
	return mark, nil
}

func (r *sqliteMarkRepo) Update(mark *app.Mark) error {
	result, err := r.db.Exec("UPDATE mark SET title = ?, url = ? WHERE id = ?", mark.Title, mark.URL, mark.ID)
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

func (r *sqliteMarkRepo) DeleteByID(id string) (bool, error) {
	result, err := r.db.Exec("DELETE FROM mark WHERE id = ?", id)
	if err != nil {
		return false, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return n != 0, nil
}

func (r *sqliteMarkRepo) DeleteAll() error {
	_, err := r.db.Exec("DELETE FROM mark")
	if err != nil {
		return err
	}

	return nil
}