package repositories

import (
	"backend-sarpras/models"
	"database/sql"
)

type RuanganRepository struct {
	DB *sql.DB
}

func NewRuanganRepository(db *sql.DB) *RuanganRepository {
	return &RuanganRepository{DB: db}
}

func (r *RuanganRepository) GetAll() ([]models.Ruangan, error) {
	query := `SELECT kode_ruangan, nama_ruangan, lokasi, kapasitas, deskripsi FROM ruangan ORDER BY kode_ruangan`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ruangans []models.Ruangan
	for rows.Next() {
		var ruangan models.Ruangan
		err := rows.Scan(
			&ruangan.KodeRuangan,
			&ruangan.NamaRuangan,
			&ruangan.Lokasi,
			&ruangan.Kapasitas,
			&ruangan.Deskripsi,
		)
		if err != nil {
			return nil, err
		}
		ruangans = append(ruangans, ruangan)
	}
	return ruangans, nil
}

func (r *RuanganRepository) GetByID(kodeRuangan string) (*models.Ruangan, error) {
	ruangan := &models.Ruangan{}
	query := `SELECT kode_ruangan, nama_ruangan, lokasi, kapasitas, deskripsi FROM ruangan WHERE kode_ruangan = $1`
	err := r.DB.QueryRow(query, kodeRuangan).Scan(
		&ruangan.KodeRuangan,
		&ruangan.NamaRuangan,
		&ruangan.Lokasi,
		&ruangan.Kapasitas,
		&ruangan.Deskripsi,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return ruangan, err
}

func (r *RuanganRepository) Create(ruangan *models.Ruangan) error {
	query := `
		INSERT INTO ruangan (kode_ruangan, nama_ruangan, lokasi, kapasitas, deskripsi)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.DB.Exec(
		query,
		ruangan.KodeRuangan,
		ruangan.NamaRuangan,
		ruangan.Lokasi,
		ruangan.Kapasitas,
		ruangan.Deskripsi,
	)
	return err
}

func (r *RuanganRepository) Update(ruangan *models.Ruangan) error {
	query := `
		UPDATE ruangan
		SET nama_ruangan = $1, lokasi = $2, kapasitas = $3, deskripsi = $4
		WHERE kode_ruangan = $5
	`
	_, err := r.DB.Exec(query, ruangan.NamaRuangan, ruangan.Lokasi, ruangan.Kapasitas, ruangan.Deskripsi, ruangan.KodeRuangan)
	return err
}

func (r *RuanganRepository) Delete(kodeRuangan string) error {
	query := `DELETE FROM ruangan WHERE kode_ruangan = $1`
	_, err := r.DB.Exec(query, kodeRuangan)
	return err
}

