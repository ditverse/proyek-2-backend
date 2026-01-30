package repositories

import (
	"backend-sarpras/models"
	"database/sql"
)

type LogAktivitasRepository struct {
	DB *sql.DB
}

func NewLogAktivitasRepository(db *sql.DB) *LogAktivitasRepository {
	return &LogAktivitasRepository{DB: db}
}

func (r *LogAktivitasRepository) Create(log *models.LogAktivitas) error {
	// Let database trigger generate kode_log to avoid duplicate key errors
	// The trigger (trigger_generate_kode_log) generates unique sequential codes per day
	query := `
		INSERT INTO log_aktivitas (kode_user, kode_peminjaman, aksi, keterangan)
		VALUES ($1, $2, $3, $4)
		RETURNING kode_log, waktu, updated_at
	`
	var kodeUser interface{}
	if log.KodeUser != nil {
		kodeUser = *log.KodeUser
	}
	var kodePeminjaman interface{}
	if log.KodePeminjaman != nil {
		kodePeminjaman = *log.KodePeminjaman
	}
	err := r.DB.QueryRow(
		query,
		kodeUser,
		kodePeminjaman,
		log.Aksi,
		log.Keterangan,
	).Scan(&log.KodeLog, &log.Waktu, &log.UpdatedAt)
	return err
}

func (r *LogAktivitasRepository) GetAll(filter string) ([]models.LogAktivitas, error) {
	query := `
		SELECT kode_log, kode_user, kode_peminjaman, aksi, keterangan, waktu, updated_at
		FROM log_aktivitas
	`
	if filter != "" {
		query += " WHERE aksi = $1"
	}
	query += " ORDER BY waktu DESC LIMIT 100"

	var rows *sql.Rows
	var err error
	if filter != "" {
		rows, err = r.DB.Query(query, filter)
	} else {
		rows, err = r.DB.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.LogAktivitas
	for rows.Next() {
		var (
			log            models.LogAktivitas
			kodeUser       sql.NullString
			kodePeminjaman sql.NullString
			updatedAt      sql.NullTime
		)
		err := rows.Scan(
			&log.KodeLog,
			&kodeUser,
			&kodePeminjaman,
			&log.Aksi,
			&log.Keterangan,
			&log.Waktu,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		if kodeUser.Valid {
			log.KodeUser = &kodeUser.String
		}
		if kodePeminjaman.Valid {
			log.KodePeminjaman = &kodePeminjaman.String
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			log.UpdatedAt = &t
		}
		logs = append(logs, log)
	}
	return logs, nil
}
