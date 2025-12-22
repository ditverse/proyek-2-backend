package repositories

import (
	"backend-sarpras/models"
	"database/sql"
)

type NotifikasiRepository struct {
	DB *sql.DB
}

func NewNotifikasiRepository(db *sql.DB) *NotifikasiRepository {
	return &NotifikasiRepository{DB: db}
}

func (r *NotifikasiRepository) Create(notifikasi *models.Notifikasi) error {
	query := `
		INSERT INTO notifikasi (kode_notifikasi, kode_user, kode_peminjaman, jenis_notifikasi, pesan, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`
	var kodePeminjaman interface{}
	if notifikasi.KodePeminjaman != nil {
		kodePeminjaman = *notifikasi.KodePeminjaman
	}
	err := r.DB.QueryRow(
		query,
		notifikasi.KodeNotifikasi,
		notifikasi.KodeUser,
		kodePeminjaman,
		notifikasi.JenisNotifikasi,
		notifikasi.Pesan,
		notifikasi.Status,
	).Scan(&notifikasi.CreatedAt, &notifikasi.UpdatedAt)
	return err
}

func (r *NotifikasiRepository) GetByPenerimaID(kodeUser string) ([]models.Notifikasi, error) {
	query := `
		SELECT kode_notifikasi, kode_user, kode_peminjaman, jenis_notifikasi, pesan, status, created_at, updated_at
		FROM notifikasi
		WHERE kode_user = $1
		ORDER BY created_at DESC
		LIMIT 50
	`
	rows, err := r.DB.Query(query, kodeUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifikasi []models.Notifikasi
	for rows.Next() {
		var n models.Notifikasi
		var kodePeminjaman sql.NullString
		err := rows.Scan(
			&n.KodeNotifikasi,
			&n.KodeUser,
			&kodePeminjaman,
			&n.JenisNotifikasi,
			&n.Pesan,
			&n.Status,
			&n.CreatedAt,
			&n.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if kodePeminjaman.Valid {
			n.KodePeminjaman = &kodePeminjaman.String
		}
		notifikasi = append(notifikasi, n)
	}
	return notifikasi, nil
}

func (r *NotifikasiRepository) GetUnreadCount(kodeUser string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifikasi WHERE kode_user = $1 AND status = 'TERKIRIM'`
	err := r.DB.QueryRow(query, kodeUser).Scan(&count)
	return count, err
}

func (r *NotifikasiRepository) MarkAsRead(kodeNotifikasi string, kodeUser string) error {
	query := `
		UPDATE notifikasi
		SET status = 'DIBACA',
		    updated_at = NOW()
		WHERE kode_notifikasi = $1 AND kode_user = $2
	`
	_, err := r.DB.Exec(query, kodeNotifikasi, kodeUser)
	return err
}

func (r *NotifikasiRepository) MarkAllAsRead(kodeUser string) error {
	query := `
		UPDATE notifikasi
		SET status = 'DIBACA',
		    updated_at = NOW()
		WHERE kode_user = $1 AND status = 'TERKIRIM'
	`
	_, err := r.DB.Exec(query, kodeUser)
	return err
}
