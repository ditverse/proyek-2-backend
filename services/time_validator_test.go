package services

import (
	"testing"
	"time"
)

// TestValidateSubmissionTime tests the submission time validation
func TestValidateSubmissionTime(t *testing.T) {
	// Note: This test will behave differently based on when it's run
	// To properly test, you'd need to mock time.Now()
	// For now, this is a placeholder for manual testing
	err := ValidateSubmissionTime()
	t.Logf("Current submission validation result: %v", err)
}

// TestValidateRentalPeriod tests various rental period scenarios
func TestValidateRentalPeriod(t *testing.T) {
	location, _ := time.LoadLocation("Asia/Jakarta")

	tests := []struct {
		name        string
		start       time.Time
		end         time.Time
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Valid: Monday 09:00 - Monday 16:00",
			start:       time.Date(2026, 1, 5, 9, 0, 0, 0, location),  // Monday
			end:         time.Date(2026, 1, 5, 16, 0, 0, 0, location), // Same Monday
			shouldError: false,
		},
		{
			name:        "Valid: Wednesday 08:00 - Thursday 15:00",
			start:       time.Date(2026, 1, 7, 8, 0, 0, 0, location),  // Wednesday
			end:         time.Date(2026, 1, 8, 15, 0, 0, 0, location), // Thursday
			shouldError: false,
		},
		{
			name:        "Valid: Saturday 10:00 - Saturday 14:00",
			start:       time.Date(2026, 1, 10, 10, 0, 0, 0, location), // Saturday
			end:         time.Date(2026, 1, 10, 14, 0, 0, 0, location), // Same Saturday
			shouldError: false,
		},
		{
			name:        "Invalid: Sunday 10:00 - Sunday 12:00",
			start:       time.Date(2026, 1, 11, 10, 0, 0, 0, location), // Sunday
			end:         time.Date(2026, 1, 11, 12, 0, 0, 0, location), // Same Sunday
			shouldError: true,
			errorMsg:    "Peminjaman tidak dapat dilakukan pada hari Minggu",
		},
		{
			name:        "Invalid: Monday 06:00 - Monday 10:00 (start too early)",
			start:       time.Date(2026, 1, 5, 6, 0, 0, 0, location),  // Monday 06:00
			end:         time.Date(2026, 1, 5, 10, 0, 0, 0, location), // Monday 10:00
			shouldError: true,
			errorMsg:    "Waktu mulai peminjaman minimal pukul 07:00 WIB",
		},
		{
			name:        "Invalid: Friday 15:00 - Friday 18:00 (end too late)",
			start:       time.Date(2026, 1, 9, 15, 0, 0, 0, location), // Friday 15:00
			end:         time.Date(2026, 1, 9, 18, 0, 0, 0, location), // Friday 18:00
			shouldError: true,
			errorMsg:    "Waktu selesai peminjaman maksimal pukul 17:00 WIB",
		},
		{
			name:        "Invalid: Saturday 16:00 - Sunday 10:00 (cross to Sunday)",
			start:       time.Date(2026, 1, 10, 16, 0, 0, 0, location), // Saturday 16:00
			end:         time.Date(2026, 1, 11, 10, 0, 0, 0, location), // Sunday 10:00
			shouldError: true,
			errorMsg:    "Peminjaman harus selesai sebelum hari Minggu",
		},
		{
			name:        "Invalid: Friday 14:00 - Monday 10:00 (includes Sunday)",
			start:       time.Date(2026, 1, 9, 14, 0, 0, 0, location),  // Friday 14:00
			end:         time.Date(2026, 1, 12, 10, 0, 0, 0, location), // Monday 10:00
			shouldError: true,
			errorMsg:    "Peminjaman tidak dapat melewati hari Minggu",
		},
		{
			name:        "Valid: Monday 07:00 - Wednesday 17:00 (edge case: exact limits)",
			start:       time.Date(2026, 1, 5, 7, 0, 0, 0, location),  // Monday 07:00 (min)
			end:         time.Date(2026, 1, 7, 17, 0, 0, 0, location), // Wednesday 17:00 (max)
			shouldError: false,
		},
		{
			name:        "Invalid: Monday 17:00 - Monday 18:00 (start at 17:00+)",
			start:       time.Date(2026, 1, 5, 17, 0, 0, 0, location), // Monday 17:00
			end:         time.Date(2026, 1, 5, 18, 0, 0, 0, location), // Monday 18:00
			shouldError: true,
			errorMsg:    "Waktu mulai peminjaman maksimal pukul 16:59 WIB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRentalPeriod(tt.start, tt.end)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error '%s', got '%s'", tt.errorMsg, err.Error())
				} else {
					t.Logf("✅ Correctly rejected: %s", err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				} else {
					t.Logf("✅ Correctly accepted")
				}
			}
		})
	}
}
