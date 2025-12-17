package repositories

import (
	"backend-sarpras/models"
	"database/sql"
)

type KehadiranRepository struct {
	DB *sql.DB
}

func NewKehadiranRepository(db *sql.DB) *KehadiranRepository {
	return &KehadiranRepository{DB: db}
}

func (r *KehadiranRepository) Create(kehadiran *models.KehadiranPeminjam) error {
	query := `
		INSERT INTO kehadiran_peminjam (kode_kehadiran, kode_peminjaman, status_kehadiran, keterangan, diverifikasi_oleh)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING waktu_konfirmasi, created_at, updated_at
	`
	var verifier interface{}
	if kehadiran.DiverifikasiOleh != nil {
		verifier = *kehadiran.DiverifikasiOleh
	}
	return r.DB.QueryRow(
		query,
		kehadiran.KodeKehadiran,
		kehadiran.KodePeminjaman,
		kehadiran.StatusKehadiran,
		kehadiran.Keterangan,
		verifier,
	).Scan(&kehadiran.WaktuKonfirmasi, &kehadiran.CreatedAt, &kehadiran.UpdatedAt)
}

func (r *KehadiranRepository) GetByPeminjamanID(kodePeminjaman string) (*models.KehadiranPeminjam, error) {
	kehadiran := &models.KehadiranPeminjam{}
	query := `
		SELECT kode_kehadiran, kode_peminjaman, status_kehadiran, waktu_konfirmasi, keterangan,
		       diverifikasi_oleh, created_at, updated_at
		FROM kehadiran_peminjam
		WHERE kode_peminjaman = $1
	`
	var verifier sql.NullString
	err := r.DB.QueryRow(query, kodePeminjaman).Scan(
		&kehadiran.KodeKehadiran,
		&kehadiran.KodePeminjaman,
		&kehadiran.StatusKehadiran,
		&kehadiran.WaktuKonfirmasi,
		&kehadiran.Keterangan,
		&verifier,
		&kehadiran.CreatedAt,
		&kehadiran.UpdatedAt,
	)
	if verifier.Valid {
		kehadiran.DiverifikasiOleh = &verifier.String
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return kehadiran, err
}

func (r *KehadiranRepository) GetAll(start, end string) ([]models.KehadiranPeminjam, error) {
	query := `
		SELECT kode_kehadiran, kode_peminjaman, status_kehadiran, waktu_konfirmasi, keterangan,
		       diverifikasi_oleh, created_at, updated_at
		FROM kehadiran_peminjam
		WHERE waktu_konfirmasi >= $1 AND waktu_konfirmasi <= $2
		ORDER BY waktu_konfirmasi DESC
	`
	rows, err := r.DB.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kehadiran []models.KehadiranPeminjam
	for rows.Next() {
		var k models.KehadiranPeminjam
		var verifier sql.NullString
		err := rows.Scan(
			&k.KodeKehadiran,
			&k.KodePeminjaman,
			&k.StatusKehadiran,
			&k.WaktuKonfirmasi,
			&k.Keterangan,
			&verifier,
			&k.CreatedAt,
			&k.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if verifier.Valid {
			k.DiverifikasiOleh = &verifier.String
		}
		kehadiran = append(kehadiran, k)
	}
	return kehadiran, nil
}


// GetRiwayat returns all kehadiran records with optional date filtering
func (r *KehadiranRepository) GetRiwayat(start, end *string) ([]models.KehadiranPeminjam, error) {
	query := `
		SELECT kode_kehadiran, kode_peminjaman, status_kehadiran, waktu_konfirmasi, keterangan,
		       diverifikasi_oleh, created_at, updated_at
		FROM kehadiran_peminjam
	`
	args := []interface{}{}
	
	// Add date filter if both start and end provided
	if start != nil && end != nil && *start != "" && *end != "" {
		query += ` WHERE waktu_konfirmasi BETWEEN $1 AND $2`
		args = append(args, *start, *end)
	}
	
	query += ` ORDER BY waktu_konfirmasi DESC`
	
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kehadiranList []models.KehadiranPeminjam
	for rows.Next() {
		var k models.KehadiranPeminjam
		var verifier sql.NullString
		err := rows.Scan(
			&k.KodeKehadiran,
			&k.KodePeminjaman,
			&k.StatusKehadiran,
			&k.WaktuKonfirmasi,
			&k.Keterangan,
			&verifier,
			&k.CreatedAt,
			&k.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if verifier.Valid {
			k.DiverifikasiOleh = &verifier.String
		}
		kehadiranList = append(kehadiranList, k)
	}
	
	if kehadiranList == nil {
		kehadiranList = []models.KehadiranPeminjam{}
	}
	
	return kehadiranList, nil
}

