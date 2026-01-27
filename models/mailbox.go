package models

import "time"

// Mailbox represents email notification log
type Mailbox struct {
	KodeMailbox    string    `json:"kode_mailbox"`
	KodeUser       string    `json:"kode_user"`
	KodePeminjaman string    `json:"kode_peminjaman"`
	JenisPesan     string    `json:"jenis_pesan"`
	CreatedAt      time.Time `json:"created_at"`
}

// JenisPesan constants
const (
	JenisPesanApproved       = "APPROVED"
	JenisPesanRejected       = "REJECTED"
	JenisPesanCancelled      = "CANCELLED"
	JenisPesanSecurityNotify = "SECURITY_NOTIFY"
	JenisPesanNewSubmission  = "NEW_SUBMISSION"
)

// MailboxWithDetails represents mailbox with joined data for email
type MailboxWithDetails struct {
	// Mailbox fields
	KodeMailbox string    `json:"kode_mailbox"`
	JenisPesan  string    `json:"jenis_pesan"`
	CreatedAt   time.Time `json:"created_at"`

	// Penerima (from users via kode_user)
	EmailTujuan string `json:"email_tujuan"`

	// Peminjaman
	KodePeminjaman    string    `json:"kode_peminjaman"`
	Status            string    `json:"status"`
	CatatanVerifikasi string    `json:"catatan_verifikasi"`
	TanggalMulai      time.Time `json:"tanggal_mulai"`
	TanggalSelesai    time.Time `json:"tanggal_selesai"`

	// Peminjam
	NamaPeminjam  string `json:"nama_peminjam"`
	EmailPeminjam string `json:"email_peminjam"`
	NoHPPeminjam  string `json:"no_hp_peminjam"`

	// Organisasi
	NamaOrganisasi string `json:"nama_organisasi"`

	// Ruangan
	NamaRuangan   string `json:"nama_ruangan"`
	LokasiRuangan string `json:"lokasi_ruangan"`
	Kapasitas     int    `json:"kapasitas"`

	// Kegiatan
	NamaKegiatan      string `json:"nama_kegiatan"`
	DeskripsiKegiatan string `json:"deskripsi_kegiatan"`

	// Verifikator
	NamaVerifikator string `json:"nama_verifikator"`
}
