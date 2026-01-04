package services

import (
	"errors"
	"time"
)

// ValidateSubmissionTime validates that the submission is made during office hours
// Office hours: Monday-Friday, 07:00-17:00 WIB
func ValidateSubmissionTime() error {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return errors.New("gagal memuat timezone")
	}

	now := time.Now().In(location)

	// Check if today is weekday (Monday-Friday)
	weekday := now.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return errors.New("Pengajuan peminjaman hanya dapat dilakukan pada hari kerja (Senin-Jumat)")
	}

	// Check if current time is within office hours (07:00-17:00)
	hour := now.Hour()
	if hour < 7 || hour >= 17 {
		return errors.New("Pengajuan peminjaman hanya dapat dilakukan pada jam kerja (07:00-17:00 WIB)")
	}

	return nil
}

// ValidateRentalPeriod validates that the rental period is within allowed time
// Allowed: Monday-Saturday, 07:00-17:00 WIB
// Not allowed: Sunday, before 07:00, after 17:00
func ValidateRentalPeriod(tanggalMulai, tanggalSelesai time.Time) error {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return errors.New("gagal memuat timezone")
	}

	// Convert to WIB timezone
	start := tanggalMulai.In(location)
	end := tanggalSelesai.In(location)

	// 1. Check start day (Monday-Saturday only, no Sunday)
	if start.Weekday() == time.Sunday {
		return errors.New("Peminjaman tidak dapat dilakukan pada hari Minggu")
	}

	// 2. Check end day (Monday-Saturday only, no Sunday)
	if end.Weekday() == time.Sunday {
		return errors.New("Peminjaman harus selesai sebelum hari Minggu")
	}

	// 3. Check start hour (must be >= 07:00 and < 17:00)
	if start.Hour() < 7 {
		return errors.New("Waktu mulai peminjaman minimal pukul 07:00 WIB")
	}
	if start.Hour() >= 17 {
		return errors.New("Waktu mulai peminjaman maksimal pukul 16:59 WIB")
	}

	// 4. Check end hour (must be <= 17:00)
	// Allow exactly 17:00:00 but not 17:00:01+
	if end.Hour() > 17 || (end.Hour() == 17 && end.Minute() > 0) {
		return errors.New("Waktu selesai peminjaman maksimal pukul 17:00 WIB")
	}
	if end.Hour() < 7 {
		return errors.New("Waktu selesai peminjaman minimal pukul 07:00 WIB")
	}

	// 5. Check if rental period includes Sunday (for multi-day rentals)
	// Loop through each day from start to end
	current := start
	for current.Before(end) || current.Equal(end) {
		if current.Weekday() == time.Sunday {
			return errors.New("Peminjaman tidak dapat melewati hari Minggu")
		}
		current = current.Add(24 * time.Hour)
	}

	return nil
}
