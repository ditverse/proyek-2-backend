package services

import (
	"errors"
	"fmt"
	"time"

	"backend-sarpras/models"
	"backend-sarpras/repositories"
)

type PeminjamanService struct {
	PeminjamanRepo *repositories.PeminjamanRepository
	BarangRepo     *repositories.BarangRepository
	NotifikasiRepo *repositories.NotifikasiRepository
	LogRepo        *repositories.LogAktivitasRepository
	UserRepo       *repositories.UserRepository
}

func NewPeminjamanService(
	peminjamanRepo *repositories.PeminjamanRepository,
	barangRepo *repositories.BarangRepository,
	notifikasiRepo *repositories.NotifikasiRepository,
	logRepo *repositories.LogAktivitasRepository,
	userRepo *repositories.UserRepository,
) *PeminjamanService {
	return &PeminjamanService{
		PeminjamanRepo: peminjamanRepo,
		BarangRepo:     barangRepo,
		NotifikasiRepo: notifikasiRepo,
		LogRepo:        logRepo,
		UserRepo:       userRepo,
	}
}

func (s *PeminjamanService) CreatePeminjaman(req *models.CreatePeminjamanRequest, kodeUser string) (*models.Peminjaman, error) {
	if req.PathSuratDigital == "" {
		return nil, errors.New("path surat digital wajib diisi")
	}

	tanggalMulai, err := time.Parse(time.RFC3339, req.TanggalMulai)
	if err != nil {
		return nil, errors.New("format tanggal_mulai tidak valid")
	}

	tanggalSelesai, err := time.Parse(time.RFC3339, req.TanggalSelesai)
	if err != nil {
		return nil, errors.New("format tanggal_selesai tidak valid")
	}

	if tanggalSelesai.Before(tanggalMulai) {
		return nil, errors.New("tanggal_selesai harus setelah tanggal_mulai")
	}

	for _, item := range req.Barang {
		barang, err := s.BarangRepo.GetByID(item.KodeBarang)
		if err != nil {
			return nil, err
		}
		if barang == nil {
			return nil, fmt.Errorf("barang %s tidak ditemukan", item.KodeBarang)
		}
	}

	peminjaman := &models.Peminjaman{
		KodeUser:         kodeUser,
		KodeRuangan:      req.KodeRuangan,
		KodeKegiatan:     req.KodeKegiatan,
		TanggalMulai:     tanggalMulai,
		TanggalSelesai:   tanggalSelesai,
		Keperluan:        req.Keperluan,
		Status:           models.StatusPeminjamanPending,
		PathSuratDigital: req.PathSuratDigital,
	}

	if err := s.PeminjamanRepo.Create(peminjaman); err != nil {
		return nil, err
	}

	for _, item := range req.Barang {
		if err := s.PeminjamanRepo.CreatePeminjamanBarang(
			generateCode("PMB"),
			peminjaman.KodePeminjaman,
			item.KodeBarang,
			item.Jumlah,
		); err != nil {
			return nil, err
		}
	}

	s.LogRepo.Create(&models.LogAktivitas{
		KodeLog:        generateCode("LOG"),
		KodeUser:       &kodeUser,
		KodePeminjaman: &peminjaman.KodePeminjaman,
		Aksi:           "CREATE_PEMINJAMAN",
		Keterangan:     "Pengajuan peminjaman baru dibuat",
	})

	kodePeminjaman := peminjaman.KodePeminjaman
	if petugas, err := s.UserRepo.GetByRole(models.RoleSarpras); err == nil && len(petugas) > 0 {
		for _, u := range petugas {
			s.NotifikasiRepo.Create(&models.Notifikasi{
				KodeNotifikasi:  generateCode("NTF"),
				KodeUser:        u.KodeUser,
				KodePeminjaman:  &kodePeminjaman,
				JenisNotifikasi: models.NotifPengajuanDibuat,
				Pesan:           "Pengajuan peminjaman baru menunggu verifikasi",
				Status:          models.NotifikasiTerkirim,
			})
		}
	}

	return peminjaman, nil
}

func (s *PeminjamanService) VerifikasiPeminjaman(kodePeminjaman string, verifierKode string, req *models.VerifikasiPeminjamanRequest) error {
	peminjaman, err := s.PeminjamanRepo.GetByID(kodePeminjaman)
	if err != nil {
		return err
	}
	if peminjaman == nil {
		return errors.New("peminjaman tidak ditemukan")
	}

	if peminjaman.Status != models.StatusPeminjamanPending {
		return errors.New("peminjaman sudah diverifikasi")
	}

	if req.Status != models.StatusPeminjamanApproved && req.Status != models.StatusPeminjamanRejected {
		return errors.New("status verifikasi tidak valid")
	}

	if err := s.PeminjamanRepo.UpdateStatus(kodePeminjaman, req.Status, &verifierKode, req.CatatanVerifikasi); err != nil {
		return err
	}

	s.LogRepo.Create(&models.LogAktivitas{
		KodeLog:        generateCode("LOG"),
		KodeUser:       &verifierKode,
		KodePeminjaman: &kodePeminjaman,
		Aksi:           "UPDATE_STATUS",
		Keterangan:     fmt.Sprintf("Status peminjaman diubah menjadi %s", req.Status),
	})

	pesan := "Pengajuan peminjaman Anda telah " + string(req.Status)
	if req.Status == models.StatusPeminjamanRejected && req.CatatanVerifikasi != "" {
		pesan += ". Catatan: " + req.CatatanVerifikasi
	}

	peminjamKode := peminjaman.KodeUser
	var jenis models.NotifikasiJenisEnum
	if req.Status == models.StatusPeminjamanApproved {
		jenis = models.NotifStatusApproved
	} else {
		jenis = models.NotifStatusRejected
	}
	s.NotifikasiRepo.Create(&models.Notifikasi{
		KodeNotifikasi:  generateCode("NTF"),
		KodeUser:        peminjamKode,
		KodePeminjaman:  &kodePeminjaman,
		JenisNotifikasi: jenis,
		Pesan:           pesan,
		Status:          models.NotifikasiTerkirim,
	})

	return nil
}
