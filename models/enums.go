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

// KehadiranStatusEnum mengikuti enum kehadiran_status_enum di database.
type KehadiranStatusEnum string

const (
	KehadiranHadir      KehadiranStatusEnum = "HADIR"
	KehadiranTidakHadir KehadiranStatusEnum = "TIDAK_HADIR"
	KehadiranBatal      KehadiranStatusEnum = "BATAL"
)
