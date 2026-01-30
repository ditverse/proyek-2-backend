package repositories

import (
	"backend-sarpras/models"
	"database/sql"
)

type BarangRepository struct {
	DB *sql.DB
}

func NewBarangRepository(db *sql.DB) *BarangRepository {
	return &BarangRepository{DB: db}
}

func (r *BarangRepository) GetAll() ([]models.Barang, error) {
	query := `
		SELECT kode_barang, nama_barang, deskripsi, jumlah_total, ruangan_kode
		FROM barang
		ORDER BY kode_barang
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var barangs []models.Barang
	for rows.Next() {
		var barang models.Barang
		var ruanganKode sql.NullString
		err := rows.Scan(
			&barang.KodeBarang,
			&barang.NamaBarang,
			&barang.Deskripsi,
			&barang.JumlahTotal,
			&ruanganKode,
		)
		if err != nil {
			return nil, err
		}
		if ruanganKode.Valid {
			barang.RuanganKode = &ruanganKode.String
		}
		barangs = append(barangs, barang)
	}
	return barangs, nil
}

func (r *BarangRepository) GetByID(kodeBarang string) (*models.Barang, error) {
	barang := &models.Barang{}
	query := `
		SELECT kode_barang, nama_barang, deskripsi, jumlah_total, ruangan_kode
		FROM barang
		WHERE kode_barang = $1
	`
	var ruanganKode sql.NullString
	err := r.DB.QueryRow(query, kodeBarang).Scan(
		&barang.KodeBarang,
		&barang.NamaBarang,
		&barang.Deskripsi,
		&barang.JumlahTotal,
		&ruanganKode,
	)
	if ruanganKode.Valid {
		barang.RuanganKode = &ruanganKode.String
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return barang, err
}

func (r *BarangRepository) Create(barang *models.Barang) error {
	// Let database trigger generate kode_barang to avoid duplicate key errors
	// The trigger (trigger_generate_kode_barang) generates unique sequential codes
	query := `
		INSERT INTO barang (nama_barang, deskripsi, jumlah_total, ruangan_kode)
		VALUES ($1, $2, $3, $4)
		RETURNING kode_barang
	`
	var ruangan interface{}
	if barang.RuanganKode != nil {
		ruangan = *barang.RuanganKode
	}
	return r.DB.QueryRow(
		query,
		barang.NamaBarang,
		barang.Deskripsi,
		barang.JumlahTotal,
		ruangan,
	).Scan(&barang.KodeBarang)
}

func (r *BarangRepository) Update(barang *models.Barang) error {
	query := `
		UPDATE barang
		SET nama_barang = $1, deskripsi = $2, jumlah_total = $3, ruangan_kode = $4
		WHERE kode_barang = $5
	`
	var ruangan interface{}
	if barang.RuanganKode != nil {
		ruangan = *barang.RuanganKode
	}
	_, err := r.DB.Exec(query, barang.NamaBarang, barang.Deskripsi, barang.JumlahTotal, ruangan, barang.KodeBarang)
	return err
}

func (r *BarangRepository) Delete(kodeBarang string) error {
	query := `DELETE FROM barang WHERE kode_barang = $1`
	_, err := r.DB.Exec(query, kodeBarang)
	return err
}
