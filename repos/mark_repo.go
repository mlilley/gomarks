package repos

import (
	"database/sql"
	"github.com/mlilley/gomarks/app"
	"strconv"
)

type MarkRepo interface {
	FindAll() ([]app.Mark, error)
	FindAllForUser(userID string) ([]app.Mark, error)
	FindByID(markID string) (*app.Mark, error)
	FindByIDForUser(markID string, userID string) (*app.Mark, error)
	Create(m *app.Mark) (*app.Mark, error)
	Update(m *app.Mark) (bool, error)
	UpdateWithResult(m *app.Mark) (*app.Mark, error)
	DeleteByID(markID string) (bool, error)
	DeleteByIDForUser(markId string, userId string) (bool, error)
	DeleteAll() error
	DeleteAllForUser(userID string) error
}

func NewMarkRepo(db *sql.DB) MarkRepo {
	return &sqliteMarkRepo{db: db}
}

type sqliteMarkRepo struct {
	db *sql.DB
}

func (r *sqliteMarkRepo) FindAll() ([]app.Mark, error) {
	rows, err := r.db.Query("SELECT id, title, url, user_id FROM mark")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var marks []app.Mark
	for rows.Next() {
		var mark app.Mark
		err = rows.Scan(&mark.ID, &mark.Title, &mark.URL, &mark.UserID)
		if err != nil {
			return nil, err
		}
		marks = append(marks, mark)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return marks, nil
}

func (r *sqliteMarkRepo) FindAllForUser(userID string) ([]app.Mark, error) {
	//intId, err := strconv.Atoi(id)
	//if err != nil {
	//	return nil, err
	//}
	//
	rows, err := r.db.Query("SELECT id, title, url, user_id FROM mark WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	marks := []app.Mark{}
	for rows.Next() {
		var mark app.Mark
		err = rows.Scan(&mark.ID, &mark.Title, &mark.URL, &mark.UserID)
		if err != nil {
			return nil, err
		}
		marks = append(marks, mark)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return marks, nil
}

func (r *sqliteMarkRepo) FindByID(markID string) (*app.Mark, error) {
	//intId, err := strconv.Atoi(markID)
	//if err != nil {
	//	return nil, err
	//}
	//
	var mark app.Mark
	err := r.db.QueryRow("SELECT id, title, url, user_id FROM mark WHERE id = ?", markID).Scan(&mark.ID, &mark.Title, &mark.URL, &mark.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			// TODO: do we really get ErrNoRows here?
			return nil, nil
		}
		return nil, err
	}

	return &mark, nil
}

func (r *sqliteMarkRepo) FindByIDForUser(markID string, userID string) (*app.Mark, error) {
	//intMarkID, err := strconv.Atoi(markID)
	//if err != nil {
	//	return nil, err
	//}
	//
	//intUserID, err := strconv.Atoi(userID)
	//if err != nil {
	//	return nil, err
	//}
	//
	var mark app.Mark
	err := r.db.QueryRow("SELECT id, title, url, user_id FROM mark WHERE id = ? AND user_id = ?", markID, userID).Scan(&mark.ID, &mark.Title, &mark.URL, &mark.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			// TODO: do we really get ErrNoRows here?
			return nil, nil
		}
		return nil, err
	}

	return &mark, nil
}

func (r *sqliteMarkRepo) Create(mark *app.Mark) (*app.Mark, error) {
	//intUserId, err := strconv.Atoi(mark.UserID)
	//if err != nil {
	//	return nil, err
	//}
	//
	result, err := r.db.Exec("INSERT INTO mark(title, url, user_id) VALUES (?, ?, ?)", mark.Title, mark.URL, mark.UserID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		// TODO: Should we really use a transaction here to rollback if error getting the last insert id?
		return nil, err
	}

	mark.ID = strconv.FormatInt(id, 10)
	return mark, nil
}

func (r *sqliteMarkRepo) Update(mark *app.Mark) (bool, error) {
	result, err := r.db.Exec("UPDATE mark SET title = ?, url = ?, user_id = ? WHERE id = ?", mark.Title, mark.URL, mark.UserID, mark.ID)
	if err != nil {
		return false, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return n > 0, nil
}

func (r *sqliteMarkRepo) UpdateWithResult(mark *app.Mark) (*app.Mark, error) {
	tx, err := r.db.Begin()

	result, err := tx.Exec("UPDATE mark SET title = ?, url = ?, user_id = ? WHERE id = ?", mark.Title, mark.URL, mark.UserID, mark.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if n == 0 {
		tx.Rollback()
		return nil, nil
	}

	err = tx.QueryRow("SELECT id, title, url, user_id FROM mark WHERE id = ?", mark.ID).Scan(&mark.ID, &mark.Title, &mark.URL, &mark.UserID)
	if err != nil {
		tx.Rollback();
		return nil, err // err could be 'ErrNoRows', but this should be interpreted as an internal error
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return mark, nil
}

func (r *sqliteMarkRepo) DeleteByID(markID string) (bool, error) {
	result, err := r.db.Exec("DELETE FROM mark WHERE id = ?", markID)
	if err != nil {
		return false, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return n > 0, nil
}

func (r *sqliteMarkRepo) DeleteByIDForUser(markId string, userId string) (bool, error) {
	result, err := r.db.Exec("DELETE FROM mark WHERE id = ? AND user_id = ?", markId, userId)
	if err != nil {
		return false, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return n > 0, nil
}

func (r *sqliteMarkRepo) DeleteAll() error {
	_, err := r.db.Exec("DELETE FROM mark")
	if err != nil {
		return err
	}

	return nil
}

func (r *sqliteMarkRepo) DeleteAllForUser(userID string) error {
	_, err := r.db.Exec("DELETE FROM mark WHERE user_id = ?", userID)
	if err != nil {
		return err
	}

	return nil
}