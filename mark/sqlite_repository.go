package mark

import (
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) (Repository, error) {
	return &repo{db: db}, nil
}

func (r *repo) FindAll() ([]*Mark, error) {
	rows, err := r.db.Query("SELECT id, title, url FROM mark")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	marks := []*Mark{}

	for rows.Next() {
		var mark Mark
		err = rows.Scan(&mark.ID, &mark.Title, &mark.URL)
		if err != nil {
			return nil, err
		}
		marks = append(marks, &mark)
	}

	return marks, nil
}

func (r *repo) FindByID(id string) (*Mark, error) {
	var mark Mark

	sid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow("SELECT id, title, url FROM mark WHERE id = ?", sid).Scan(&mark.ID, &mark.Title, &mark.URL)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &mark, nil
}

func (r *repo) Create(mark *Mark) (*Mark, error) {
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

func (r *repo) Update(mark *Mark) error {
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

func (r *repo) DeleteByID(id string) (bool, error) {
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

func (r *repo) DeleteAll() error {
	_, err := r.db.Exec("DELETE FROM mark")
	if err != nil {
		return err
	}

	return nil
}