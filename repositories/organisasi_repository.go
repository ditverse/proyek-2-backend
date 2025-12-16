package repositories

import (
	"backend-sarpras/models"
	"database/sql"
)

type OrganisasiRepository struct {
	DB *sql.DB
}

func NewOrganisasiRepository(db *sql.DB) *OrganisasiRepository {
	return &OrganisasiRepository{DB: db}
}

func (r *OrganisasiRepository) GetAll() ([]models.Organisasi, error) {
	query := `SELECT kode_organisasi, nama, jenis_organisasi, kontak, created_at FROM organisasi ORDER BY nama`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []models.Organisasi
	for rows.Next() {
		var o models.Organisasi
		err := rows.Scan(&o.KodeOrganisasi, &o.NamaOrganisasi, &o.JenisOrganisasi, &o.Kontak, &o.CreatedAt)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, o)
	}
	return orgs, nil
}

func (r *OrganisasiRepository) GetByID(kode string) (*models.Organisasi, error) {
	org := &models.Organisasi{}
	query := `SELECT kode_organisasi, nama, jenis_organisasi, kontak, created_at FROM organisasi WHERE kode_organisasi = $1`
	err := r.DB.QueryRow(query, kode).Scan(&org.KodeOrganisasi, &org.NamaOrganisasi, &org.JenisOrganisasi, &org.Kontak, &org.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return org, err
}

func (r *OrganisasiRepository) Create(org *models.Organisasi) error {
	query := `INSERT INTO organisasi (kode_organisasi, nama, jenis_organisasi, kontak) VALUES ($1, $2, $3, $4) RETURNING created_at`
	return r.DB.QueryRow(query, org.KodeOrganisasi, org.NamaOrganisasi, org.JenisOrganisasi, org.Kontak).Scan(&org.CreatedAt)
}

