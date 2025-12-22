package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	storagesvc "backend-sarpras/internal/services"
	"backend-sarpras/models"
	"backend-sarpras/repositories"
)

type PeminjamanService struct {
	PeminjamanRepo  *repositories.PeminjamanRepository
	BarangRepo      *repositories.BarangRepository
	NotifikasiRepo  *repositories.NotifikasiRepository
	LogRepo         *repositories.LogAktivitasRepository
	UserRepo        *repositories.UserRepository
	KegiatanRepo    *repositories.KegiatanRepository
	RuanganRepo     *repositories.RuanganRepository
	EmailService    *EmailService
	WhatsappService *WhatsappService
}

func NewPeminjamanService(
	peminjamanRepo *repositories.PeminjamanRepository,
	barangRepo *repositories.BarangRepository,
	notifikasiRepo *repositories.NotifikasiRepository,
	logRepo *repositories.LogAktivitasRepository,
	userRepo *repositories.UserRepository,
	kegiatanRepo *repositories.KegiatanRepository,
	ruanganRepo *repositories.RuanganRepository,
	emailService *EmailService,
	whatsappService *WhatsappService,
) *PeminjamanService {
	return &PeminjamanService{
		PeminjamanRepo:  peminjamanRepo,
		BarangRepo:      barangRepo,
		NotifikasiRepo:  notifikasiRepo,
		LogRepo:         logRepo,
		UserRepo:        userRepo,
		KegiatanRepo:    kegiatanRepo,
		RuanganRepo:     ruanganRepo,
		EmailService:    emailService,
		WhatsappService: whatsappService,
	}
}

func (s *PeminjamanService) CreatePeminjaman(req *models.CreatePeminjamanRequest, kodeUser string) (*models.Peminjaman, error) {
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

	// Get ruangan name for notification
	namaRuangan := "Ruangan"
	if req.KodeRuangan != nil && s.RuanganRepo != nil {
		if ruangan, err := s.RuanganRepo.GetByID(*req.KodeRuangan); err == nil && ruangan != nil {
			namaRuangan = ruangan.NamaRuangan
		}
	}

	// Notify Sarpras staff about new peminjaman
	if petugas, err := s.UserRepo.GetByRole(models.RoleSarpras); err == nil && len(petugas) > 0 {
		for _, u := range petugas {
			// Create in-app notification
			s.NotifikasiRepo.Create(&models.Notifikasi{
				KodeNotifikasi:  generateCode("NTF"),
				KodeUser:        u.KodeUser,
				KodePeminjaman:  &kodePeminjaman,
				JenisNotifikasi: models.NotifPengajuanDibuat,
				Pesan:           "Pengajuan peminjaman baru menunggu verifikasi",
				Status:          models.NotifikasiTerkirim,
			})

			// Send email notification to Sarpras
			if s.EmailService != nil && s.EmailService.IsConfigured() {
				emailBody := EmailTemplatePengajuanBaru(
					user.Nama,
					req.NamaKegiatan,
					namaRuangan,
					tanggalMulai,
					tanggalSelesai,
				)
				s.EmailService.SendEmailHTML(
					u.Email,
					fmt.Sprintf("📋 Pengajuan Peminjaman Baru - %s", user.Nama),
					emailBody,
				)
				log.Printf("📧 Email notification sent to Sarpras: %s", u.Email)
			}
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

	// Get peminjam data for notification
	peminjam, err := s.UserRepo.GetByID(peminjaman.KodeUser)
	if err != nil || peminjam == nil {
		log.Printf("⚠️ Could not get peminjam data: %v", err)
		peminjam = &models.User{Nama: "Peminjam", Email: ""}
	}

	// Get ruangan name
	namaRuangan := "Ruangan"
	if peminjaman.KodeRuangan != nil && s.RuanganRepo != nil {
		if ruangan, err := s.RuanganRepo.GetByID(*peminjaman.KodeRuangan); err == nil && ruangan != nil {
			namaRuangan = ruangan.NamaRuangan
		}
	}

	// Get kegiatan name
	namaKegiatan := "Kegiatan"
	if peminjaman.KodeKegiatan != nil && s.KegiatanRepo != nil {
		if kegiatan, err := s.KegiatanRepo.GetByID(*peminjaman.KodeKegiatan); err == nil && kegiatan != nil {
			namaKegiatan = kegiatan.NamaKegiatan
		}
	}

	// Create in-app notification message
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

	// Create in-app notification for peminjam
	s.NotifikasiRepo.Create(&models.Notifikasi{
		KodeNotifikasi:  generateCode("NTF"),
		KodeUser:        peminjamKode,
		KodePeminjaman:  &kodePeminjaman,
		JenisNotifikasi: jenis,
		Pesan:           pesan,
		Status:          models.NotifikasiTerkirim,
	})

	// Send external notifications based on status
	if req.Status == models.StatusPeminjamanApproved {
		// ===== APPROVED: Email + WA to Mahasiswa, WA to Security =====

		// 1. Send Email to Mahasiswa
		if s.EmailService != nil && s.EmailService.IsConfigured() && peminjam.Email != "" {
			emailBody := EmailTemplateApproved(
				peminjam.Nama,
				namaKegiatan,
				namaRuangan,
				peminjaman.TanggalMulai,
				peminjaman.TanggalSelesai,
			)
			s.EmailService.SendEmailHTML(
				peminjam.Email,
				fmt.Sprintf("✅ Peminjaman Disetujui - %s", namaKegiatan),
				emailBody,
			)
			log.Printf("📧 Approval email sent to: %s", peminjam.Email)
		}

		// 2. Send WhatsApp to Mahasiswa
		if s.WhatsappService != nil && s.WhatsappService.IsConfigured() && peminjam.NoHP != nil && *peminjam.NoHP != "" {
			waMessage := WATemplateApproved(namaKegiatan, namaRuangan)
			s.WhatsappService.SendMessage(*peminjam.NoHP, waMessage)
			log.Printf("💬 Approval WhatsApp sent to: %s", *peminjam.NoHP)
		}

		// 3. Send WhatsApp to Security
		if s.WhatsappService != nil && s.WhatsappService.IsConfigured() {
			if securityUsers, err := s.UserRepo.GetByRole(models.RoleSecurity); err == nil && len(securityUsers) > 0 {
				for _, security := range securityUsers {
					if security.NoHP != nil && *security.NoHP != "" {
						waMessage := WATemplateSecurity(
							namaKegiatan,
							namaRuangan,
							peminjaman.TanggalMulai,
							peminjaman.TanggalSelesai,
						)
						s.WhatsappService.SendMessage(*security.NoHP, waMessage)
						log.Printf("💬 Security WhatsApp sent to: %s", *security.NoHP)
					}
				}
			}
		}

		// 4. Create in-app notification for Security
		if securityUsers, err := s.UserRepo.GetByRole(models.RoleSecurity); err == nil && len(securityUsers) > 0 {
			for _, security := range securityUsers {
				s.NotifikasiRepo.Create(&models.Notifikasi{
					KodeNotifikasi:  generateCode("NTF"),
					KodeUser:        security.KodeUser,
					KodePeminjaman:  &kodePeminjaman,
					JenisNotifikasi: models.NotifInfoKegiatan,
					Pesan:           fmt.Sprintf("Kegiatan baru disetujui: %s di %s", namaKegiatan, namaRuangan),
					Status:          models.NotifikasiTerkirim,
				})
			}
			log.Printf("📱 In-app notification sent to %d security users", len(securityUsers))
		}

	} else if req.Status == models.StatusPeminjamanRejected {
		// ===== REJECTED: Email to Mahasiswa =====
		if s.EmailService != nil && s.EmailService.IsConfigured() && peminjam.Email != "" {
			alasan := req.CatatanVerifikasi
			if alasan == "" {
				alasan = "Tidak ada alasan yang diberikan."
			}
			emailBody := EmailTemplateRejected(
				peminjam.Nama,
				namaKegiatan,
				namaRuangan,
				alasan,
				peminjaman.TanggalMulai,
				peminjaman.TanggalSelesai,
			)
			s.EmailService.SendEmailHTML(
				peminjam.Email,
				fmt.Sprintf("❌ Peminjaman Ditolak - %s", namaKegiatan),
				emailBody,
			)
			log.Printf("📧 Rejection email sent to: %s", peminjam.Email)
		}
	}

	return nil
}
