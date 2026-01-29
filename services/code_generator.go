package services

import (
	"fmt"
	"sync"
	"time"
)

var (
	codeCounters = make(map[string]int)
	counterMutex sync.Mutex
	lastDate     string
)

// generateCode membuat kode unik berbasis prefix dengan format PREFIX-YYYY-MM-DD-XXXX.
// Format ini selaras dengan trigger database (contoh: KGT-2026-01-30-0001).
// Menggunakan counter harian untuk menghindari duplikasi.
func generateCode(prefix string) string {
	counterMutex.Lock()
	defer counterMutex.Unlock()

	today := time.Now().Format("2006-01-02")

	// Reset counter jika hari berubah
	if lastDate != today {
		codeCounters = make(map[string]int)
		lastDate = today
	}

	// Increment counter untuk prefix ini
	codeCounters[prefix]++
	seq := codeCounters[prefix]

	return fmt.Sprintf("%s-%s-%04d", prefix, today, seq)
}
