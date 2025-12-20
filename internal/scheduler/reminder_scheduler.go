package scheduler

import (
	"database/sql"
	"log"
	"time"

	"backend-sarpras/models"
	"backend-sarpras/repositories"
	"backend-sarpras/services"
)

// ReminderScheduler handles automatic reminder notifications
type ReminderScheduler struct {
	db              *sql.DB
	peminjamanRepo  *repositories.PeminjamanRepository
	userRepo        *repositories.UserRepository
	notifikasiRepo  *repositories.NotifikasiRepository
	ruanganRepo     *repositories.RuanganRepository
	whatsappService *services.WhatsappService
	ticker          *time.Ticker
	stopChan        chan bool
}

// NewReminderScheduler creates a new reminder scheduler
func NewReminderScheduler(
	db *sql.DB,
	peminjamanRepo *repositories.PeminjamanRepository,
	userRepo *repositories.UserRepository,
	notifikasiRepo *repositories.NotifikasiRepository,
	ruanganRepo *repositories.RuanganRepository,
	whatsappService *services.WhatsappService,
) *ReminderScheduler {
	return &ReminderScheduler{
		db:              db,
		peminjamanRepo:  peminjamanRepo,
		userRepo:        userRepo,
		notifikasiRepo:  notifikasiRepo,
		ruanganRepo:     ruanganRepo,
		whatsappService: whatsappService,
		stopChan:        make(chan bool),
	}
}

// Start begins the scheduler with the specified interval
func (s *ReminderScheduler) Start(interval time.Duration) {
	s.ticker = time.NewTicker(interval)
	log.Printf("⏰ Reminder scheduler started with interval: %v", interval)

	go func() {
		// Run immediately on start
		s.checkAndSendReminders()

		for {
			select {
			case <-s.ticker.C:
				s.checkAndSendReminders()
			case <-s.stopChan:
				s.ticker.Stop()
				log.Println("⏰ Reminder scheduler stopped")
				return
			}
		}
	}()
}

// Stop halts the scheduler
func (s *ReminderScheduler) Stop() {
	s.stopChan <- true
}

// checkAndSendReminders queries for peminjaman ending within 1 hour and sends reminders
func (s *ReminderScheduler) checkAndSendReminders() {
	now := time.Now()
	oneHourLater := now.Add(1 * time.Hour)

	log.Printf("🔍 Checking for peminjaman ending between %s and %s", now.Format("15:04"), oneHourLater.Format("15:04"))

	// Query for approved peminjaman that will end within 1 hour
	// and haven't received a REMINDER_1JAM notification yet
	query := `
		SELECT p.kode_peminjaman, p.kode_user, p.kode_ruangan, p.tanggal_selesai
		FROM peminjaman p
		WHERE p.status = 'APPROVED'
		  AND p.tanggal_selesai <= $1
		  AND p.tanggal_selesai > $2
		  AND NOT EXISTS (
			SELECT 1 FROM notifikasi n 
			WHERE n.kode_peminjaman = p.kode_peminjaman 
			  AND n.jenis_notifikasi = 'REMINDER_1JAM'
		  )
	`

	rows, err := s.db.Query(query, oneHourLater, now)
	if err != nil {
		log.Printf("❌ Error querying peminjaman for reminder: %v", err)
		return
	}
	defer rows.Close()

	reminderCount := 0
	for rows.Next() {
		var (
			kodePeminjaman string
			kodeUser       string
			kodeRuangan    sql.NullString
			tanggalSelesai time.Time
		)

		if err := rows.Scan(&kodePeminjaman, &kodeUser, &kodeRuangan, &tanggalSelesai); err != nil {
			log.Printf("❌ Error scanning peminjaman row: %v", err)
			continue
		}

		// Get user data for WhatsApp number
		user, err := s.userRepo.GetByID(kodeUser)
		if err != nil || user == nil {
			log.Printf("⚠️ Could not get user %s: %v", kodeUser, err)
			continue
		}

		// Get ruangan name
		namaRuangan := "Ruangan"
		if kodeRuangan.Valid {
			if ruangan, err := s.ruanganRepo.GetByID(kodeRuangan.String); err == nil && ruangan != nil {
				namaRuangan = ruangan.NamaRuangan
			}
		}

		// Send WhatsApp reminder
		if s.whatsappService != nil && s.whatsappService.IsConfigured() && user.NoHP != nil && *user.NoHP != "" {
			waMessage := services.WATemplateReminder1Hour(namaRuangan)
			s.whatsappService.SendMessage(*user.NoHP, waMessage)
			log.Printf("💬 Reminder WhatsApp sent to %s for peminjaman %s", *user.NoHP, kodePeminjaman)
		}

		// Create notification record to prevent duplicate reminders
		s.notifikasiRepo.Create(&models.Notifikasi{
			KodeNotifikasi:  generateReminderCode(),
			KodeUser:        kodeUser,
			KodePeminjaman:  &kodePeminjaman,
			JenisNotifikasi: models.NotifReminder1Jam,
			Pesan:           "Reminder: Waktu peminjaman ruangan " + namaRuangan + " tinggal 1 jam lagi",
			Status:          models.NotifikasiTerkirim,
		})

		reminderCount++
	}

	if reminderCount > 0 {
		log.Printf("✅ Sent %d reminder(s)", reminderCount)
	}
}

// generateReminderCode generates a unique code for reminder notifications
func generateReminderCode() string {
	return "NTF" + time.Now().Format("20060102150405") + randomString(4)
}

// randomString generates a random alphanumeric string
func randomString(n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(1 * time.Nanosecond) // Ensure different values
	}
	return string(b)
}
