package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	internalsvc "backend-sarpras/internal/services"
	"backend-sarpras/models"
	"backend-sarpras/repositories"
)

type PeminjamanService struct {
	PeminjamanRepo *repositories.PeminjamanRepository
	BarangRepo     *repositories.BarangRepository
	LogRepo        *repositories.LogAktivitasRepository
	UserRepo       *repositories.UserRepository
	KegiatanRepo   *repositories.KegiatanRepository
	OrganisasiRepo *repositories.OrganisasiRepository
	RuanganRepo    *repositories.RuanganRepository
	MailboxRepo    *repositories.MailboxRepository
	EmailService   *internalsvc.EmailService
}

func NewPeminjamanService(
	peminjamanRepo *repositories.PeminjamanRepository,
	barangRepo *repositories.BarangRepository,
	logRepo *repositories.LogAktivitasRepository,
	userRepo *repositories.UserRepository,
	kegiatanRepo *repositories.KegiatanRepository,
	organisasiRepo *repositories.OrganisasiRepository,
	ruanganRepo *repositories.RuanganRepository,
	mailboxRepo *repositories.MailboxRepository,
	emailService *internalsvc.EmailService,
) *PeminjamanService {
	return &PeminjamanService{
		PeminjamanRepo: peminjamanRepo,
		BarangRepo:     barangRepo,
		LogRepo:        logRepo,
		UserRepo:       userRepo,
		KegiatanRepo:   kegiatanRepo,
		OrganisasiRepo: organisasiRepo,
		RuanganRepo:    ruanganRepo,
		MailboxRepo:    mailboxRepo,
		EmailService:   emailService,
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

	// Validasi: Booking tidak boleh lebih dari 2 minggu ke depan dari hari ini
	maxBookingDate := time.Now().AddDate(0, 0, 14) // 14 hari = 2 minggu
	if tanggalMulai.After(maxBookingDate) {
		return nil, errors.New("booking hanya bisa dilakukan maksimal 2 minggu sebelum tanggal acara")
	}

	// Validasi: Cek apakah ruangan sudah dibooking pada waktu tersebut (mencegah bentrok)
	if req.KodeRuangan != nil && *req.KodeRuangan != "" {
		isAvailable, err := s.PeminjamanRepo.IsRoomAvailable(*req.KodeRuangan, tanggalMulai, tanggalSelesai)
		if err != nil {
			return nil, errors.New("gagal mengecek ketersediaan ruangan")
		}
		if !isAvailable {
			return nil, errors.New("ruangan sudah dibooking pada waktu tersebut, silakan pilih waktu lain")
		}
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
		if err := internalsvc.MoveFile(suratPath, uniquePath); err != nil {
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

	// Send notification to Sarpras (async)
	if s.EmailService != nil && s.EmailService.IsEnabled() {
		go s.sendNewSubmissionEmail(peminjaman.KodePeminjaman)
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

	// Send email notification asynchronously (non-blocking)
	if s.EmailService != nil && s.EmailService.IsEnabled() {
		go s.sendVerificationEmails(kodePeminjaman, verifierKode, req.Status, req.CatatanVerifikasi)
	}

	return nil
}

// sendVerificationEmails sends email notifications after verification (runs async)
func (s *PeminjamanService) sendVerificationEmails(kodePeminjaman, verifierKode string, status models.PeminjamanStatusEnum, catatan string) {
	// 1. Determine JenisPesan based on status
	var jenisPesan string
	if status == models.StatusPeminjamanApproved {
		jenisPesan = models.JenisPesanApproved
	} else {
		jenisPesan = models.JenisPesanRejected
	}

	// 2. Insert log to mailbox
	// We need peminjam's kodeUser, so we fetch minimal data first
	peminjaman, err := s.PeminjamanRepo.GetByID(kodePeminjaman)
	if err != nil || peminjaman == nil {
		log.Printf("[ERROR] Email: failed to get peminjaman %s: %v", kodePeminjaman, err)
		return
	}

	mailbox := &models.Mailbox{
		KodeUser:       peminjaman.KodeUser,
		KodePeminjaman: kodePeminjaman,
		JenisPesan:     jenisPesan,
	}

	if err := s.MailboxRepo.Create(mailbox); err != nil {
		log.Printf("[ERROR] Email: failed to insert mailbox: %v", err)
		// Continue anyway to send email, even if logging fails?
		// Or return? Let's log error but try to send email if we can get data.
	} else {
		log.Printf("[OK] Mailbox: created log %s for %s", mailbox.KodeMailbox, kodePeminjaman)
	}

	// 3. Helper function to fetch data and send email
	// We use the MailboxRepo to get full joined data (normalized approach)
	// If mailbox insert failed, we might not have a KodeMailbox, so we fallback to manual fetching or retry insert.
	// But let's assume insert usually works. If insert failed, we can't use GetFullDataByID comfortably without an ID.
	// Strategy: If mailbox insert succeeded, use GetFullDataByID.

	// Data preparation
	var templateData internalsvc.EmailTemplateData
	var emailTo string

	if mailbox.KodeMailbox != "" {
		// Efficient way: Get from mailbox view
		details, err := s.MailboxRepo.GetFullDataByID(mailbox.KodeMailbox)
		if err != nil || details == nil {
			log.Printf("[ERROR] Email: failed to get mailbox details: %v", err)
			return
		}

		// Map MailboxWithDetails to EmailTemplateData
		// Note: We need to fetch Barang separately as it's a 1-to-Many relation not fully joined in a single row usually (or array agg).
		// Our GetFullDataByID doesn't return barang list in current implementation.
		// So we fetch barang manually.
		barangList, _ := s.PeminjamanRepo.GetPeminjamanBarang(kodePeminjaman)
		var barangItems []internalsvc.BarangItem
		for _, b := range barangList {
			namaBarang := b.KodeBarang
			if barang, err := s.BarangRepo.GetByID(b.KodeBarang); err == nil && barang != nil {
				namaBarang = barang.NamaBarang
			}
			barangItems = append(barangItems, internalsvc.BarangItem{
				NamaBarang: namaBarang,
				Jumlah:     b.Jumlah,
			})
		}

		templateData = internalsvc.EmailTemplateData{
			KodePeminjaman:    details.KodePeminjaman,
			Status:            details.Status,
			CatatanVerifikasi: catatan, // Use current note, not from DB yet if async delay
			NamaPeminjam:      details.NamaPeminjam,
			EmailPeminjam:     details.EmailPeminjam,
			NoHPPeminjam:      details.NoHPPeminjam,
			NamaOrganisasi:    details.NamaOrganisasi,
			NamaKegiatan:      details.NamaKegiatan,
			DeskripsiKegiatan: details.DeskripsiKegiatan,
			NamaRuangan:       details.NamaRuangan,
			LokasiRuangan:     details.LokasiRuangan,
			Kapasitas:         details.Kapasitas,
			TanggalMulai:      details.TanggalMulai,
			TanggalSelesai:    details.TanggalSelesai,
			TanggalVerifikasi: time.Now(),
			NamaVerifikator:   details.NamaVerifikator,
			Barang:            barangItems,
		}
		emailTo = details.EmailTujuan
	} else {
		// Fallback: If mailbox insert failed, we skip sending or handle manual fetch?
		// For now, let's just return to avoid complexity duplication
		return
	}

	// 4. Send Email based on status
	if status == models.StatusPeminjamanApproved {
		// Approved Email
		subject := internalsvc.GetApprovedEmailSubject(templateData.NamaKegiatan)
		htmlBody := internalsvc.BuildApprovedEmailHTML(templateData)
		if err := s.EmailService.SendEmail(emailTo, subject, htmlBody); err != nil {
			log.Printf("[ERROR] Email: failed to send approval to %s: %v", emailTo, err)
		} else {
			log.Printf("[OK] Email: approval sent to %s", emailTo)
		}

		// Security Notification (Only for Approved)
		// We can also insert into mailbox for security users if we want to track them
		securityUsers, err := s.UserRepo.GetByRole(models.RoleSecurity)
		if err == nil && len(securityUsers) > 0 {
			securitySubject := internalsvc.GetSecurityEmailSubject(templateData.NamaKegiatan, templateData.TanggalMulai)
			securityHTML := internalsvc.BuildSecurityNotificationHTML(templateData)

			for _, sec := range securityUsers {
				// Optional: Log to mailbox for security
				secMailbox := &models.Mailbox{
					KodeUser:       sec.KodeUser,
					KodePeminjaman: kodePeminjaman,
					JenisPesan:     models.JenisPesanSecurityNotify,
				}
				s.MailboxRepo.Create(secMailbox)

				// Send email
				if err := s.EmailService.SendEmail(sec.Email, securitySubject, securityHTML); err != nil {
					log.Printf("[ERROR] Email: failed to send to security %s: %v", sec.Email, err)
				} else {
					log.Printf("[OK] Email: security notification sent to %s", sec.Email)
				}
			}
		}

	} else {
		// Rejected Email
		subject := internalsvc.GetRejectedEmailSubject(templateData.NamaKegiatan)
		htmlBody := internalsvc.BuildRejectedEmailHTML(templateData)
		if err := s.EmailService.SendEmail(emailTo, subject, htmlBody); err != nil {
			log.Printf("[ERROR] Email: failed to send rejection to %s: %v", emailTo, err)
		} else {
			log.Printf("[OK] Email: rejection sent to %s", emailTo)
		}
	}
}

// CancelPeminjaman cancels an APPROVED or ONGOING peminjaman
// Only SARPRAS/ADMIN can cancel bookings
func (s *PeminjamanService) CancelPeminjaman(kodePeminjaman, cancellerKode, alasan string) error {
	peminjaman, err := s.PeminjamanRepo.GetByID(kodePeminjaman)
	if err != nil {
		return err
	}
	if peminjaman == nil {
		return errors.New("peminjaman tidak ditemukan")
	}

	// Validate status - can only cancel APPROVED or ONGOING
	if peminjaman.Status != models.StatusPeminjamanApproved && peminjaman.Status != models.StatusPeminjamanOngoing {
		return errors.New("hanya peminjaman dengan status APPROVED atau ONGOING yang dapat dibatalkan")
	}

	// Update status to CANCELLED
	if err := s.PeminjamanRepo.UpdateStatus(kodePeminjaman, models.StatusPeminjamanCanceled, &cancellerKode, alasan); err != nil {
		return err
	}

	// Log activity
	s.LogRepo.Create(&models.LogAktivitas{
		KodeLog:        generateCode("LOG"),
		KodeUser:       &cancellerKode,
		KodePeminjaman: &kodePeminjaman,
		Aksi:           "CANCEL_PEMINJAMAN",
		Keterangan:     fmt.Sprintf("Peminjaman dibatalkan oleh SARPRAS. Alasan: %s", alasan),
	})

	// Send email notification asynchronously (non-blocking)
	if s.EmailService != nil && s.EmailService.IsEnabled() {
		go s.sendCancellationEmail(kodePeminjaman, cancellerKode, alasan)
	}

	return nil
}

// sendCancellationEmail sends email notification after cancellation (runs async)
func (s *PeminjamanService) sendCancellationEmail(kodePeminjaman, cancellerKode, alasan string) {
	// 1. Insert log to mailbox
	// Fetch minimal peminjaman data for KodeUser
	peminjaman, err := s.PeminjamanRepo.GetByID(kodePeminjaman)
	if err != nil || peminjaman == nil {
		log.Printf("[ERROR] Email: failed to get peminjaman %s: %v", kodePeminjaman, err)
		return
	}

	mailbox := &models.Mailbox{
		KodeUser:       peminjaman.KodeUser,
		KodePeminjaman: kodePeminjaman,
		JenisPesan:     models.JenisPesanCancelled,
	}

	if err := s.MailboxRepo.Create(mailbox); err != nil {
		log.Printf("[ERROR] Email: failed to insert mailbox: %v", err)
	} else {
		log.Printf("[OK] Mailbox: created log %s for cancellation %s", mailbox.KodeMailbox, kodePeminjaman)
	}

	// 2. Fetch full data via MailboxRepo
	// Strategy: If mailbox insert succeeded, use GetFullDataByID.

	// Prepare template data
	var templateData internalsvc.EmailTemplateData
	var emailTo string

	if mailbox.KodeMailbox != "" {
		details, err := s.MailboxRepo.GetFullDataByID(mailbox.KodeMailbox)
		if err != nil || details == nil {
			log.Printf("[ERROR] Email: failed to get mailbox details: %v", err)
			return
		}

		// Fetch barang manually
		barangList, _ := s.PeminjamanRepo.GetPeminjamanBarang(kodePeminjaman)
		var barangItems []internalsvc.BarangItem
		for _, b := range barangList {
			namaBarang := b.KodeBarang
			if barang, err := s.BarangRepo.GetByID(b.KodeBarang); err == nil && barang != nil {
				namaBarang = barang.NamaBarang
			}
			barangItems = append(barangItems, internalsvc.BarangItem{
				NamaBarang: namaBarang,
				Jumlah:     b.Jumlah,
			})
		}

		templateData = internalsvc.EmailTemplateData{
			KodePeminjaman:    details.KodePeminjaman,
			Status:            "CANCELLED",
			CatatanVerifikasi: alasan, // Reuse catatan field for cancellation reason
			NamaPeminjam:      details.NamaPeminjam,
			EmailPeminjam:     details.EmailPeminjam,
			NoHPPeminjam:      details.NoHPPeminjam,
			NamaOrganisasi:    details.NamaOrganisasi,
			NamaKegiatan:      details.NamaKegiatan,
			DeskripsiKegiatan: details.DeskripsiKegiatan,
			NamaRuangan:       details.NamaRuangan,
			LokasiRuangan:     details.LokasiRuangan,
			Kapasitas:         details.Kapasitas,
			TanggalMulai:      details.TanggalMulai,
			TanggalSelesai:    details.TanggalSelesai,
			TanggalVerifikasi: time.Now(),
			NamaVerifikator:   "SARPRAS",
			Barang:            barangItems,
		}
		emailTo = details.EmailTujuan

	} else {
		// Fallback not implemented to keep it clean
		return
	}

	// 3. Send Cancellation Email
	subject := internalsvc.GetCancelledEmailSubject(templateData.NamaKegiatan)
	htmlBody := internalsvc.BuildCancelledEmailHTML(templateData)

	if err := s.EmailService.SendEmail(emailTo, subject, htmlBody); err != nil {
		log.Printf("[ERROR] Email: failed to send cancellation to %s: %v", emailTo, err)
	} else {
		log.Printf("[OK] Email: cancellation sent to %s", emailTo)
	}

	// 4. Notify Security
	securityUsers, err := s.UserRepo.GetByRole(models.RoleSecurity)
	if err == nil && len(securityUsers) > 0 {
		securitySubject := fmt.Sprintf("⚠️ Peminjaman Dibatalkan: %s (%s)", templateData.NamaKegiatan, templateData.TanggalMulai.Format("02 Jan 2006"))
		for _, sec := range securityUsers {
			// Optional: Log to mailbox for security
			secMailbox := &models.Mailbox{
				KodeUser:       sec.KodeUser,
				KodePeminjaman: kodePeminjaman,
				JenisPesan:     models.JenisPesanSecurityNotify,
			}
			s.MailboxRepo.Create(secMailbox)

			if err := s.EmailService.SendEmail(sec.Email, securitySubject, htmlBody); err != nil {
				log.Printf("[ERROR] Email: failed to send cancellation to security %s: %v", sec.Email, err)
			} else {
				log.Printf("[OK] Email: cancellation notification sent to security %s", sec.Email)
			}
		}
	}
}

// sendNewSubmissionEmail sends email notification to all Sarpras staff when new submission created
func (s *PeminjamanService) sendNewSubmissionEmail(kodePeminjaman string) {
	// 1. Get Sarpras Users
	sarprasUsers, err := s.UserRepo.GetByRole(models.RoleSarpras)
	if err != nil || len(sarprasUsers) == 0 {
		log.Printf("⚠️ Email: No sarpras found or error: %v", err)
		return
	}

	// 2. Fetch Peminjaman Data (once)
	// We need fetching minimal data first for constructing Mailbox,
	// but we can optimize by fetching full data later inside loop or once here.
	// Since MailboxRepo.GetFullDataByID is good for template, let's stick to the pattern:
	// Insert Mailbox -> GetFullData (this ensures data consistency with mailbox view)

	// Fetch Barang Data once (to reuse)
	barangList, _ := s.PeminjamanRepo.GetPeminjamanBarang(kodePeminjaman)
	var barangItems []internalsvc.BarangItem
	for _, b := range barangList {
		namaBarang := b.KodeBarang
		if barang, err := s.BarangRepo.GetByID(b.KodeBarang); err == nil && barang != nil {
			namaBarang = barang.NamaBarang
		}
		barangItems = append(barangItems, internalsvc.BarangItem{
			NamaBarang: namaBarang,
			Jumlah:     b.Jumlah,
		})
	}

	// 3. Loop Sarpras Users
	for _, staff := range sarprasUsers {
		// Log to Mailbox
		mailbox := &models.Mailbox{
			KodeUser:       staff.KodeUser,
			KodePeminjaman: kodePeminjaman,
			JenisPesan:     models.JenisPesanNewSubmission,
		}

		if err := s.MailboxRepo.Create(mailbox); err != nil {
			log.Printf("[ERROR] Email: failed to insert mailbox for sarpras %s: %v", staff.Email, err)
			continue
		}

		// Get Data for Email Template
		details, err := s.MailboxRepo.GetFullDataByID(mailbox.KodeMailbox)
		if err != nil || details == nil {
			log.Printf("[ERROR] Email: failed to get mailbox details for sarpras %s: %v", staff.Email, err)
			continue
		}

		// Prepare Template Data
		templateData := internalsvc.EmailTemplateData{
			KodePeminjaman: details.KodePeminjaman,
			Status:         "NEW_SUBMISSION",
			NamaPeminjam:   details.NamaPeminjam,
			// No Need EmailPeminjam as recipient is Sarpras
			NamaOrganisasi:    details.NamaOrganisasi,
			NoHPPeminjam:      details.NoHPPeminjam,
			NamaKegiatan:      details.NamaKegiatan,
			DeskripsiKegiatan: details.DeskripsiKegiatan, // Added field to struct if needed
			NamaRuangan:       details.NamaRuangan,
			LokasiRuangan:     details.LokasiRuangan,
			TanggalMulai:      details.TanggalMulai,
			TanggalSelesai:    details.TanggalSelesai,
			Barang:            barangItems,
		}

		// Build and Send Email
		subject := internalsvc.GetNewSubmissionEmailSubject(templateData.NamaKegiatan)
		htmlBody := internalsvc.BuildNewSubmissionEmailHTML(templateData)

		if err := s.EmailService.SendEmail(staff.Email, subject, htmlBody); err != nil {
			log.Printf("[ERROR] Email: failed to send new submission notify to %s: %v", staff.Email, err)
		} else {
			log.Printf("[OK] Email: new submission notify sent to %s", staff.Email)
		}
	}
}
