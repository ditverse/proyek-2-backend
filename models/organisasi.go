package models

import "time"

type Organisasi struct {
	KodeOrganisasi  string              `json:"kode_organisasi"`
	NamaOrganisasi  string              `json:"nama_organisasi"`
	JenisOrganisasi JenisOrganisasiEnum `json:"jenis_organisasi"`
	Kontak          string              `json:"kontak"`
	CreatedAt       time.Time           `json:"created_at"`
}

