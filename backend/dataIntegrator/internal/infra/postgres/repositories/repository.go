package repositories

import (
	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/usecase"
	"gorm.io/gorm"
	"time"
)

type PostgresDeps struct {
	DB *gorm.DB
}

type Postgres struct {
	db *gorm.DB
}

func NewPostgres(d *PostgresDeps) *Postgres {
	return &Postgres{
		db: d.DB,
	}
}

func (p *Postgres) InsertAction(action usecase.ActionType, deviceName usecase.Device) error {
	query := `
		INSERT INTO actions_log (device_name, action)
		VALUES ($1, $2)
	`

	if err := p.db.Exec(query, deviceName, action).Error; err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetDevices() ([]usecase.Device, error) {
	var devices []usecase.Device

	query := `
		SELECT name
		FROM devices
	`

	if err := p.db.Raw(query).Scan(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

func (p *Postgres) GetTelemetryByPeriod(deviceName string, from time.Time, to time.Time) ([]usecase.Telemetry, error) {
	var telemetry []usecase.Telemetry

	query := `
		SELECT id, device_name, value, created_at
		FROM telemetry
		WHERE device_name = $1
		AND created_at BETWEEN $2 AND $3
		ORDER BY created_at
	`

	if err := p.db.Raw(
		query,
		deviceName,
		from,
		to,
	).Scan(&telemetry).Error; err != nil {
		return nil, err
	}

	return telemetry, nil
}

func (p *Postgres) InsertTelemetry(deviceName usecase.Device, value int32) error {
	query := `
		INSERT INTO telemetry (device_name, value)
		VALUES ($1, $2)
	`

	if err := p.db.Exec(
		query,
		deviceName,
		value,
	).Error; err != nil {
		return err
	}

	return nil
}
