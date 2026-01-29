package services

import (
	"backend-sarpras/models"
	"backend-sarpras/repositories"
	"log"
	"time"
)

// StatusScheduler handles automatic status transitions for peminjaman
type StatusScheduler struct {
	peminjamanRepo *repositories.PeminjamanRepository
	ticker         *time.Ticker
	done           chan bool
}

// NewStatusScheduler creates a new StatusScheduler
func NewStatusScheduler(peminjamanRepo *repositories.PeminjamanRepository) *StatusScheduler {
	return &StatusScheduler{
		peminjamanRepo: peminjamanRepo,
		done:           make(chan bool),
	}
}

// Start begins the scheduler with 1 minute interval
func (s *StatusScheduler) Start() {
	s.ticker = time.NewTicker(1 * time.Minute)
	log.Println("📅 Status scheduler started")

	// Run immediately on start
	go s.updateStatuses()

	go func() {
		for {
			select {
			case <-s.done:
				return
			case <-s.ticker.C:
				s.updateStatuses()
			}
		}
	}()
}

// Stop stops the scheduler
func (s *StatusScheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.done <- true
	log.Println("📅 Status scheduler stopped")
}

// updateStatuses performs the status transition logic
func (s *StatusScheduler) updateStatuses() {
	now := time.Now()

	// 1. Update APPROVED/ONGOING → FINISHED (where tanggal_selesai < now)
	finishedCount, err := s.peminjamanRepo.UpdateStatusToFinished(now)
	if err != nil {
		log.Printf("[ERROR] Error updating status to FINISHED: %v", err)
	} else if finishedCount > 0 {
		log.Printf("[OK] Updated %d peminjaman to FINISHED", finishedCount)
	}

	// 2. Update APPROVED → ONGOING (where tanggal_mulai <= now <= tanggal_selesai)
	ongoingCount, err := s.peminjamanRepo.UpdateStatusToOngoing(now)
	if err != nil {
		log.Printf("[ERROR] Error updating status to ONGOING: %v", err)
	} else if ongoingCount > 0 {
		log.Printf("[OK] Updated %d peminjaman to ONGOING", ongoingCount)
	}
}

// GetPeminjamanRepo returns the repository for external use
func (s *StatusScheduler) GetPeminjamanRepo() *repositories.PeminjamanRepository {
	return s.peminjamanRepo
}

// ForceUpdate forces an immediate status update (useful for testing)
func (s *StatusScheduler) ForceUpdate() {
	s.updateStatuses()
}

// UpdateStatusFromApprovedToOngoing updates a specific peminjaman status
// This can be called when a peminjaman is approved and its start time has already passed
func UpdateStatusFromApprovedToOngoing(repo *repositories.PeminjamanRepository, kodePeminjaman string, tanggalMulai, tanggalSelesai time.Time) error {
	now := time.Now()
	if tanggalMulai.Before(now) && tanggalSelesai.After(now) {
		return repo.UpdateStatusOnly(kodePeminjaman, models.StatusPeminjamanOngoing)
	}
	return nil
}
