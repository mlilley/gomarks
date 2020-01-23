package repos

import (
	"database/sql"
	"github.com/mlilley/gomarks/app"
)

type DeviceRepo interface {
	FindById(deviceId string, userId string) (*app.Device, error)
	Upsert(d *app.Device) error
}

func NewDeviceRepo(db *sql.DB) DeviceRepo {
	return &sqliteDeviceRepo{db: db}
}

type sqliteDeviceRepo struct {
	db *sql.DB
}

func (r *sqliteDeviceRepo) FindById(userId string, deviceId string) (*app.Device, error) {
	var d app.Device
	err := r.db.QueryRow("SELECT user_id, device_id, token_hash FROM device WHERE user_id = ? AND device_id = ?", userId, deviceId).Scan(&d.UserId, &d.DeviceId, &d.TokenHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

func (r *sqliteDeviceRepo) Upsert(d *app.Device) error {
	_, err := r.db.Exec("REPLACE INTO device (user_id, device_id, token_hash) VALUES (?, ?, ?)", d.UserId, d.DeviceId, d.TokenHash)
	if err != nil {
		return err
	}
	return nil
}

