package handlers

import (
	"backend-sarpras/models"
	"backend-sarpras/repositories"
	"backend-sarpras/services"
	"fmt"
	"net/http"
	"time"
)

type ExportHandler struct {
	PeminjamanRepo *repositories.PeminjamanRepository
	RuanganRepo    *repositories.RuanganRepository
	UserRepo       *repositories.UserRepository
	OrganisasiRepo *repositories.OrganisasiRepository
	KegiatanRepo   *repositories.KegiatanRepository
	ExportService  *services.ExportService
}

func NewExportHandler(
	peminjamanRepo *repositories.PeminjamanRepository,
	ruanganRepo *repositories.RuanganRepository,
	userRepo *repositories.UserRepository,
	organisasiRepo *repositories.OrganisasiRepository,
	kegiatanRepo *repositories.KegiatanRepository,
	exportService *services.ExportService,
) *ExportHandler {
	return &ExportHandler{
		PeminjamanRepo: peminjamanRepo,
		RuanganRepo:    ruanganRepo,
		UserRepo:       userRepo,
		OrganisasiRepo: organisasiRepo,
		KegiatanRepo:   kegiatanRepo,
		ExportService:  exportService,
	}
}

// ExportPeminjamanToExcel handles the export endpoint
func (h *ExportHandler) ExportPeminjamanToExcel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters (same as GetLaporan)
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	statusStr := r.URL.Query().Get("status")

	var start, end time.Time
	var err error

	if startStr == "" {
		start = time.Now().AddDate(0, -1, 0) // default 1 bulan lalu
	} else {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			http.Error(w, "Invalid start date", http.StatusBadRequest)
			return
		}
	}

	if endStr == "" {
		end = time.Now()
	} else {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			http.Error(w, "Invalid end date", http.StatusBadRequest)
			return
		}
	}

	var status models.PeminjamanStatusEnum
	if statusStr != "" {
		status = models.PeminjamanStatusEnum(statusStr)
	}

	// Get data from repository (same logic as GetLaporan)
	peminjaman, err := h.PeminjamanRepo.GetLaporan(start, end, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enrich data with relations
	for i := range peminjaman {
		if peminjaman[i].KodeRuangan != nil {
			ruangan, _ := h.RuanganRepo.GetByID(*peminjaman[i].KodeRuangan)
			peminjaman[i].Ruangan = ruangan
		}
		if peminjaman[i].KodeKegiatan != nil {
			kegiatan, _ := h.KegiatanRepo.GetByID(*peminjaman[i].KodeKegiatan)
			peminjaman[i].Kegiatan = kegiatan
		}
		user, _ := h.UserRepo.GetByID(peminjaman[i].KodeUser)
		if user != nil {
			user.PasswordHash = ""
			// Load organisasi data
			if user.OrganisasiKode != nil {
				organisasi, _ := h.OrganisasiRepo.GetByID(*user.OrganisasiKode)
				user.Organisasi = organisasi
			}
			peminjaman[i].Peminjam = user
		}
		// Load verifier
		if peminjaman[i].VerifiedBy != nil {
			verifier, _ := h.UserRepo.GetByID(*peminjaman[i].VerifiedBy)
			if verifier != nil {
				verifier.PasswordHash = ""
				peminjaman[i].Verifier = verifier
			}
		}
		// Load barang
		items, _ := h.PeminjamanRepo.GetPeminjamanBarang(peminjaman[i].KodePeminjaman)
		peminjaman[i].Barang = items
	}

	// Generate Excel file
	file, err := h.ExportService.GeneratePeminjamanExcel(peminjaman, start, end, statusStr)
	if err != nil {
		http.Error(w, "Failed to generate Excel: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Generate filename with timestamp
	filename := fmt.Sprintf("Laporan_Peminjaman_%s.xlsx", time.Now().Format("2006-01-02"))

	// Set headers for file download
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Transfer-Encoding", "binary")

	// Write Excel file to response
	if err := file.Write(w); err != nil {
		http.Error(w, "Failed to write Excel file", http.StatusInternalServerError)
		return
	}
}
