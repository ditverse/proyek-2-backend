package models

// RoleEnum merepresentasikan enum role pengguna di database.
type RoleEnum string

const (
	RoleMahasiswa RoleEnum = "MAHASISWA"
	RoleSarpras   RoleEnum = "SARPRAS"
	RoleSecurity  RoleEnum = "SECURITY"
	RoleAdmin     RoleEnum = "ADMIN"
)

// JenisOrganisasiEnum merepresentasikan jenis organisasi.
type JenisOrganisasiEnum string

const (
	JenisOrganisasiORMAWA JenisOrganisasiEnum = "ORMAWA"
	JenisOrganisasiUKM    JenisOrganisasiEnum = "UKM"
)

// PeminjamanStatusEnum mengikuti enum peminjaman_status_enum di database.
type PeminjamanStatusEnum string

const (
	StatusPeminjamanPending  PeminjamanStatusEnum = "PENDING"
	StatusPeminjamanApproved PeminjamanStatusEnum = "APPROVED"
	StatusPeminjamanRejected PeminjamanStatusEnum = "REJECTED"
	StatusPeminjamanOngoing  PeminjamanStatusEnum = "ONGOING"
	StatusPeminjamanFinished PeminjamanStatusEnum = "FINISHED"
	StatusPeminjamanCanceled PeminjamanStatusEnum = "CANCELLED"
)

// NotifikasiStatusEnum mengikuti enum notifikasi_status_enum di database.
type NotifikasiStatusEnum string

const (
	NotifikasiTerkirim NotifikasiStatusEnum = "TERKIRIM"
	NotifikasiDibaca   NotifikasiStatusEnum = "DIBACA"
)

// NotifikasiJenisEnum mengikuti enum notifikasi_jenis_enum di database.
type NotifikasiJenisEnum string

const (
	NotifPengajuanDibuat   NotifikasiJenisEnum = "PENGAJUAN_DIBUAT"
	NotifStatusApproved    NotifikasiJenisEnum = "STATUS_APPROVED"
	NotifStatusRejected    NotifikasiJenisEnum = "STATUS_REJECTED"
	NotifReminderKehadiran NotifikasiJenisEnum = "REMINDER_KEHADIRAN"
)

// KehadiranStatusEnum mengikuti enum kehadiran_status_enum di database.
type KehadiranStatusEnum string

const (
	KehadiranHadir      KehadiranStatusEnum = "HADIR"
	KehadiranTidakHadir KehadiranStatusEnum = "TIDAK_HADIR"
	KehadiranBatal      KehadiranStatusEnum = "BATAL"
)
