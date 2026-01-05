package repositories

import (
	"backend-sarpras/models"
	"database/sql"
)

type KegiatanRepository struct {
	DB *sql.DB
}

func NewKegiatanRepository(db *sql.DB) *KegiatanRepository {
	return &KegiatanRepository{DB: db}
}

func (r *KegiatanRepository) Create(k *models.Kegiatan) error {
	query := `INSERT INTO kegiatan (kode_kegiatan, nama_kegiatan, deskripsi, tanggal_mulai, tanggal_selesai, organisasi_kode) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING created_at`
	return r.DB.QueryRow(query, k.KodeKegiatan, k.NamaKegiatan, k.Deskripsi, k.TanggalMulai, k.TanggalSelesai, k.OrganisasiKode).Scan(&k.CreatedAt)
}

func (r *KegiatanRepository) GetByID(kode string) (*models.Kegiatan, error) {
	k := &models.Kegiatan{}
	query := `SELECT kode_kegiatan, nama_kegiatan, deskripsi, tanggal_mulai, tanggal_selesai, organisasi_kode, created_at 
			  FROM kegiatan WHERE kode_kegiatan = $1`
	err := r.DB.QueryRow(query, kode).Scan(&k.KodeKegiatan, &k.NamaKegiatan, &k.Deskripsi, &k.TanggalMulai, &k.TanggalSelesai, &k.OrganisasiKode, &k.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return k, err
}
