package models

import "time"

type Peminjaman struct {
	KodePeminjaman    string                   `json:"kode_peminjaman"`
	KodeUser          string                   `json:"kode_user"`
	Peminjam          *User                    `json:"peminjam,omitempty"`
	KodeRuangan       *string                  `json:"kode_ruangan"`
	Ruangan           *Ruangan                 `json:"ruangan,omitempty"`
	KodeKegiatan      *string                  `json:"kode_kegiatan"`
	Kegiatan          *Kegiatan                `json:"kegiatan,omitempty"`
	TanggalMulai      time.Time                `json:"tanggal_mulai"`
	TanggalSelesai    time.Time                `json:"tanggal_selesai"`
	Status            PeminjamanStatusEnum     `json:"status"`
	PathSuratDigital  string                   `json:"path_surat_digital"`
	VerifiedBy        *string                  `json:"verified_by"`
	Verifier          *User                    `json:"verifier,omitempty"`
	VerifiedAt        *time.Time               `json:"verified_at"`
	CatatanVerifikasi string                   `json:"catatan_verifikasi"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         *time.Time               `json:"updated_at"`
	Barang            []PeminjamanBarangDetail `json:"barang,omitempty"`
}

type PeminjamanBarangDetail struct {
	KodePeminjamanBarang string  `json:"kode_peminjaman_barang"`
	KodeBarang           string  `json:"kode_barang"`
	Barang               *Barang `json:"barang,omitempty"`
	Jumlah               int     `json:"jumlah"`
}

type CreatePeminjamanRequest struct {
	KodeRuangan      *string                  `json:"kode_ruangan"`
	NamaKegiatan     string                   `json:"nama_kegiatan"`
	Deskripsi        string                   `json:"deskripsi"`
	TanggalMulai     string                   `json:"tanggal_mulai"` // ISO 8601
	TanggalSelesai   string                   `json:"tanggal_selesai"`
	PathSuratDigital string                   `json:"path_surat_digital"` // New schema
	SuratDigitalURL  string                   `json:"surat_digital_url"`  // Old schema (backward compatibility)
	Barang           []CreatePeminjamanBarang `json:"barang"`
}

type CreatePeminjamanBarang struct {
	KodeBarang string `json:"kode_barang"`
	Jumlah     int    `json:"jumlah"`
}

type VerifikasiPeminjamanRequest struct {
	Status            PeminjamanStatusEnum `json:"status"`
	CatatanVerifikasi string               `json:"catatan_verifikasi"`
}

type KegiatanSimple struct {
	KodeKegiatan string `json:"kode_kegiatan"`
	NamaKegiatan string `json:"nama_kegiatan"`
}

type JadwalRuanganResponse struct {
	KodePeminjaman string               `json:"kode_peminjaman"`
	KodeRuangan    string               `json:"kode_ruangan"`
	NamaRuangan    string               `json:"nama_ruangan"`
	TanggalMulai   time.Time            `json:"tanggal_mulai"`
	TanggalSelesai time.Time            `json:"tanggal_selesai"`
	Status         PeminjamanStatusEnum `json:"status"`
	Peminjam       string               `json:"peminjam"`
	Organisasi     string               `json:"organisasi"`
	Kegiatan       *KegiatanSimple      `json:"kegiatan,omitempty"`
}
