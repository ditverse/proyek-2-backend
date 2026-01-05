package handlers

import (
	"backend-sarpras/middleware"
	"backend-sarpras/models"
	"backend-sarpras/repositories"
	"backend-sarpras/services"
	"encoding/json"
	"net/http"
)

type KehadiranHandler struct {
	KehadiranService *services.KehadiranService
	KehadiranRepo    *repositories.KehadiranRepository
	PeminjamanRepo   *repositories.PeminjamanRepository
	RuanganRepo      *repositories.RuanganRepository
	UserRepo         *repositories.UserRepository
	KegiatanRepo     *repositories.KegiatanRepository
}

func NewKehadiranHandler(
	kehadiranService *services.KehadiranService,
	kehadiranRepo *repositories.KehadiranRepository,
	peminjamanRepo *repositories.PeminjamanRepository,
	ruanganRepo *repositories.RuanganRepository,
	userRepo *repositories.UserRepository,
	kegiatanRepo *repositories.KegiatanRepository,
) *KehadiranHandler {
	return &KehadiranHandler{
		KehadiranService: kehadiranService,
		KehadiranRepo:    kehadiranRepo,
		PeminjamanRepo:   peminjamanRepo,
		RuanganRepo:      ruanganRepo,
		UserRepo:         userRepo,
		KegiatanRepo:     kegiatanRepo,
	}
}

func (h *KehadiranHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateKehadiranRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.KehadiranService.CreateKehadiran(&req, user.KodeUser)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Kehadiran berhasil dicatat"})
}

func (h *KehadiranHandler) GetByPeminjamanID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	peminjamanIDStr := r.URL.Query().Get("peminjaman_id")
	if peminjamanIDStr == "" {
		http.Error(w, "peminjaman_id required", http.StatusBadRequest)
		return
	}

	// TODO: parse peminjamanID dari query string
	// Untuk sekarang, return empty atau implementasi sesuai kebutuhan
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]models.KehadiranPeminjam{})
}

func (h *KehadiranHandler) GetRiwayatBySecurity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get optional date filters from query params
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	
	var startPtr, endPtr *string
	if start != "" {
		startPtr = &start
	}
	if end != "" {
		endPtr = &end
	}

	// Get kehadiran data from repository
	kehadiranList, err := h.KehadiranRepo.GetRiwayat(startPtr, endPtr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enrich each kehadiran with Peminjaman, Peminjam, Ruangan, Kegiatan, Verifier
	for i := range kehadiranList {
		// Get Peminjaman data
		peminjaman, err := h.PeminjamanRepo.GetByID(kehadiranList[i].KodePeminjaman)
		if err == nil && peminjaman != nil {
			// Get Peminjam (User)
			if peminjaman.KodeUser != "" {
				peminjam, err := h.UserRepo.GetByID(peminjaman.KodeUser)
				if err == nil && peminjam != nil {
					peminjaman.Peminjam = peminjam
				}
			}
			
			// Get Ruangan
			if peminjaman.KodeRuangan != nil {
				ruangan, err := h.RuanganRepo.GetByID(*peminjaman.KodeRuangan)
				if err == nil && ruangan != nil {
					peminjaman.Ruangan = ruangan
				}
			}
			
			// Get Kegiatan
			if peminjaman.KodeKegiatan != nil {
				kegiatan, err := h.KegiatanRepo.GetByID(*peminjaman.KodeKegiatan)
				if err == nil && kegiatan != nil {
					peminjaman.Kegiatan = kegiatan
				}
			}
			
			kehadiranList[i].Peminjaman = peminjaman
		}
		
		// Get Verifier (Security user who verified)
		if kehadiranList[i].DiverifikasiOleh != nil {
			verifier, err := h.UserRepo.GetByID(*kehadiranList[i].DiverifikasiOleh)
			if err == nil && verifier != nil {
				kehadiranList[i].Verifier = verifier
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kehadiranList)
}

