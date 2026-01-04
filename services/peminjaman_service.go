package services

import (
	"errors"
	"fmt"
	"time"

	storagesvc "backend-sarpras/internal/services"
	"backend-sarpras/models"
	"backend-sarpras/repositories"
)

type PeminjamanService struct {
	PeminjamanRepo *repositories.PeminjamanRepository
	BarangRepo     *repositories.BarangRepository
	NotifikasiRepo *repositories.NotifikasiRepository
	LogRepo        *repositories.LogAktivitasRepository
	UserRepo       *repositories.UserRepository
	KegiatanRepo   *repositories.KegiatanRepository
}

func NewPeminjamanService(
	peminjamanRepo *repositories.PeminjamanRepository,
	barangRepo *repositories.BarangRepository,
	notifikasiRepo *repositories.NotifikasiRepository,
	logRepo *repositories.LogAktivitasRepository,
	userRepo *repositories.UserRepository,
	kegiatanRepo *repositories.KegiatanRepository,
) *PeminjamanService {
	return &PeminjamanService{
		PeminjamanRepo: peminjamanRepo,
		BarangRepo:     barangRepo,
		NotifikasiRepo: notifikasiRepo,
		LogRepo:        logRepo,
		UserRepo:       userRepo,
		KegiatanRepo:   kegiatanRepo,
	}
}

func (s *PeminjamanService) CreatePeminjaman(req *models.CreatePeminjamanRequest, kodeUser string) (*models.Peminjaman, error) {
	// VALIDASI 1: Jam Submit (Request Time)
	// Mahasiswa/user hanya boleh submit pengajuan pada hari kerja (Senin-Jumat) jam 07:00-17:00 WIB
	if err := ValidateSubmissionTime(); err != nil {
		return nil, err
	}

	// Backward compatibility: support both old (surat_digital_url) and new (path_surat_digital) field names
	suratPath := req.PathSuratDigital
	if suratPath == "" && req.SuratDigitalURL != "" {
		suratPath = req.SuratDigitalURL // Fallback to old field name
	}

	// Path is now optional - can be uploaded later via /upload-surat endpoint
	// Common placeholder values from frontend: "uploaded-via-form", "pending", etc.
	isPlaceholder := suratPath == "" ||
		suratPath == "uploaded-via-form" ||
		suratPath == "pending" ||
		suratPath == "temp"

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

	// VALIDASI 2: Waktu Peminjaman (Rental Period)
	// Waktu peminjaman harus dalam range Senin-Sabtu, jam 07:00-17:00 WIB
	// Tidak boleh ada Minggu di antara periode peminjaman
	if err := ValidateRentalPeriod(tanggalMulai, tanggalSelesai); err != nil {
		return nil, err
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

	// Get user untuk ambil organisasi_kode
	user, err := s.UserRepo.GetByID(kodeUser)
	if err != nil {
		return nil, errors.New("gagal mengambil data user")
	}
	if user == nil {
		return nil, errors.New("user tidak ditemukan")
	}

	// Auto-create kegiatan dari data peminjaman
	var kodeKegiatan *string
	if req.NamaKegiatan != "" {
		organisasiKode := ""
		if user.OrganisasiKode != nil {
			organisasiKode = *user.OrganisasiKode
		}

		kegiatan := &models.Kegiatan{
			KodeKegiatan:   generateCode("KGT"),
			NamaKegiatan:   req.NamaKegiatan,
			Deskripsi:      req.Deskripsi,
			TanggalMulai:   tanggalMulai,
			TanggalSelesai: tanggalSelesai,
			OrganisasiKode: organisasiKode,
		}

		if err := s.KegiatanRepo.Create(kegiatan); err != nil {
			return nil, fmt.Errorf("gagal membuat kegiatan: %v", err)
		}
		kodeKegiatan = &kegiatan.KodeKegiatan
	}

	peminjaman := &models.Peminjaman{
		KodeUser:         kodeUser,
		KodeRuangan:      req.KodeRuangan,
		KodeKegiatan:     kodeKegiatan,
		TanggalMulai:     tanggalMulai,
		TanggalSelesai:   tanggalSelesai,
		Status:           models.StatusPeminjamanPending,
		PathSuratDigital: suratPath, // Use the resolved path
	}

	if err := s.PeminjamanRepo.Create(peminjaman); err != nil {
		return nil, err
	}

	// IMPORTANT: Regenerate unique path using kode_peminjaman to prevent file overwrite
	// Each peminjaman will have unique path: peminjaman/{kode_peminjaman}/surat.pdf
	// This ensures files from different peminjaman don't overwrite each other
	uniquePath := fmt.Sprintf("peminjaman/%s/surat.pdf", peminjaman.KodePeminjaman)

	// Move file from old path (frontend-provided) to new unique path in storage
	// Only if path is not a placeholder and file actually exists in storage
	if !isPlaceholder && suratPath != uniquePath {
		if err := storagesvc.MoveFile(suratPath, uniquePath); err != nil {
			// File doesn't exist or move failed - this is OK
			// User can upload later via /upload-surat endpoint
			fmt.Printf("Info: skipping file move (file may not exist yet): %v\n", err)
			// Set path to empty so frontend knows to upload
			uniquePath = ""
		}
	} else if isPlaceholder {
		// Placeholder path - set to empty, user must upload via /upload-surat
		uniquePath = ""
	}

	// Update path in database with unique path (or empty if file not uploaded yet)
	if err := s.PeminjamanRepo.UpdateSuratDigitalURL(peminjaman.KodePeminjaman, uniquePath); err != nil {
		return nil, err
	}

	// Update local object for response
	peminjaman.PathSuratDigital = uniquePath

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
