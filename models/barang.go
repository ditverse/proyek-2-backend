package models

type Barang struct {
	KodeBarang  string   `json:"kode_barang"`
	NamaBarang  string   `json:"nama_barang"`
	Deskripsi   string   `json:"deskripsi"`
	JumlahTotal int      `json:"jumlah_total"`
	RuanganKode *string  `json:"ruangan_kode"`
	Ruangan     *Ruangan `json:"ruangan,omitempty"`
}

type CreateBarangRequest struct {
	KodeBarang  string  `json:"kode_barang"`
	NamaBarang  string  `json:"nama_barang"`
	Deskripsi   string  `json:"deskripsi"`
	JumlahTotal int     `json:"jumlah_total"`
	RuanganKode *string `json:"ruangan_kode"`
}

type UpdateBarangRequest struct {
	NamaBarang  string  `json:"nama_barang"`
	Deskripsi   string  `json:"deskripsi"`
	JumlahTotal int     `json:"jumlah_total"`
	RuanganKode *string `json:"ruangan_kode"`
}
