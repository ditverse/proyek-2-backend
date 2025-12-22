package services

import (
	"backend-sarpras/models"
	"backend-sarpras/repositories"
	"errors"
	"fmt"
	"log"
)

type KehadiranService struct {
	KehadiranRepo  *repositories.KehadiranRepository
	PeminjamanRepo *repositories.PeminjamanRepository
	LogRepo        *repositories.LogAktivitasRepository
	NotifikasiRepo *repositories.NotifikasiRepository
	UserRepo       *repositories.UserRepository
}

func NewKehadiranService(
	kehadiranRepo *repositories.KehadiranRepository,
	peminjamanRepo *repositories.PeminjamanRepository,
	logRepo *repositories.LogAktivitasRepository,
	notifikasiRepo *repositories.NotifikasiRepository,
	userRepo *repositories.UserRepository,
) *KehadiranService {
	return &KehadiranService{
		KehadiranRepo:  kehadiranRepo,
		PeminjamanRepo: peminjamanRepo,
		LogRepo:        logRepo,
		NotifikasiRepo: notifikasiRepo,
		UserRepo:       userRepo,
	}
}

func (s *KehadiranService) CreateKehadiran(req *models.CreateKehadiranRequest, securityKode string) error {
	peminjaman, err := s.PeminjamanRepo.GetByID(req.KodePeminjaman)
	if err != nil {
		return err
	}
	if peminjaman == nil {
		return errors.New("peminjaman tidak ditemukan")
	}

	if peminjaman.Status != models.StatusPeminjamanApproved {
		return errors.New("peminjaman belum disetujui")
	}

	existing, _ := s.KehadiranRepo.GetByPeminjamanID(req.KodePeminjaman)
	if existing != nil {
		return errors.New("kehadiran sudah pernah diisi")
	}

	if req.StatusKehadiran != models.KehadiranHadir &&
		req.StatusKehadiran != models.KehadiranTidakHadir &&
		req.StatusKehadiran != models.KehadiranBatal {
		return errors.New("status kehadiran tidak valid")
	}

	kehadiran := &models.KehadiranPeminjam{
		KodeKehadiran:    generateCode("KHD"),
		KodePeminjaman:   req.KodePeminjaman,
		StatusKehadiran:  req.StatusKehadiran,
		Keterangan:       req.Keterangan,
		DiverifikasiOleh: &securityKode,
	}

	if err := s.KehadiranRepo.Create(kehadiran); err != nil {
		return err
	}

	s.LogRepo.Create(&models.LogAktivitas{
		KodeLog:        generateCode("LOG"),
		KodeUser:       &securityKode,
		KodePeminjaman: &req.KodePeminjaman,
		Aksi:           "UPDATE_KEHADIRAN",
		Keterangan:     "Status kehadiran: " + string(req.StatusKehadiran),
	})

	// Notify Sarpras staff about kehadiran verification
	if s.NotifikasiRepo != nil && s.UserRepo != nil {
		if sarprasUsers, err := s.UserRepo.GetByRole(models.RoleSarpras); err == nil && len(sarprasUsers) > 0 {
			statusLabel := "diverifikasi"
			if req.StatusKehadiran == models.KehadiranTidakHadir {
				statusLabel = "tidak hadir"
			} else if req.StatusKehadiran == models.KehadiranBatal {
				statusLabel = "dibatalkan"
			}

			for _, sarpras := range sarprasUsers {
				s.NotifikasiRepo.Create(&models.Notifikasi{
					KodeNotifikasi:  generateCode("NTF"),
					KodeUser:        sarpras.KodeUser,
					KodePeminjaman:  &req.KodePeminjaman,
					JenisNotifikasi: models.NotifKehadiranVerified,
					Pesan:           fmt.Sprintf("Kehadiran peminjaman %s", statusLabel),
					Status:          models.NotifikasiTerkirim,
				})
			}
			log.Printf("📱 Kehadiran notification sent to %d Sarpras users", len(sarprasUsers))
		}
	}

	return nil
}
