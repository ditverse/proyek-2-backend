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

	// 1. Get optional date filters
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	
	var startPtr, endPtr *string
	if start != "" {
		startPtr = &start
	}
	if end != "" {
		endPtr = &end
	}

	// 2. Fetch main list (Kehadiran)
	kehadiranList, err := h.KehadiranRepo.GetRiwayat(startPtr, endPtr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if len(kehadiranList) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(kehadiranList)
		return
	}

	// 3. Collect IDs for batch fetching
	peminjamanIDs := make([]string, 0)
	verifierIDs := make([]string, 0)
	verifierSet := make(map[string]bool)

	for _, k := range kehadiranList {
		peminjamanIDs = append(peminjamanIDs, k.KodePeminjaman)
		
		if k.DiverifikasiOleh != nil && !verifierSet[*k.DiverifikasiOleh] {
			verifierIDs = append(verifierIDs, *k.DiverifikasiOleh)
			verifierSet[*k.DiverifikasiOleh] = true
		}
	}

	// 4. Batch Load: Peminjaman & Verifiers
	peminjamanMap, err := h.PeminjamanRepo.GetByIDs(peminjamanIDs)
	if err != nil {
		http.Error(w, "Error loading peminjaman data", http.StatusInternalServerError)
		return
	}
	
	verifierMap := h.UserRepo.GetByIDs(verifierIDs)

	// 5. Collect related IDs from fetched Peminjaman
	userIDs := make([]string, 0)      // Peminjam
	ruanganIDs := make([]string, 0)
	kegiatanIDs := make([]string, 0)
	
	userSet := make(map[string]bool)
	ruanganSet := make(map[string]bool)
	kegiatanSet := make(map[string]bool)

	for _, p := range peminjamanMap {
		if !userSet[p.KodeUser] {
			userIDs = append(userIDs, p.KodeUser)
			userSet[p.KodeUser] = true
		}
		if p.KodeRuangan != nil && !ruanganSet[*p.KodeRuangan] {
			ruanganIDs = append(ruanganIDs, *p.KodeRuangan)
			ruanganSet[*p.KodeRuangan] = true
		}
		if p.KodeKegiatan != nil && !kegiatanSet[*p.KodeKegiatan] {
			kegiatanIDs = append(kegiatanIDs, *p.KodeKegiatan)
			kegiatanSet[*p.KodeKegiatan] = true
		}
	}

	// 6. Batch Load: User (Peminjam), Ruangan, Kegiatan
	userMap := h.UserRepo.GetByIDs(userIDs)
	ruanganMap := h.RuanganRepo.GetByIDs(ruanganIDs)
	kegiatanMap := h.KegiatanRepo.GetByIDs(kegiatanIDs)

	// 7. Assemble Data (In-Memory Association)
	for i := range kehadiranList {
		// Link Peminjaman
		if peminjaman, ok := peminjamanMap[kehadiranList[i].KodePeminjaman]; ok {
			// Deep copy to avoid sharing references if needed, but here we just need to link
			// We need to construct a new object or modify the map value (which is a pointer)
			// But since we want to attach related objects to Peminjaman struct, better to use value copy if modifying
			
			// Link Peminjam (User)
			if peminjam, ok := userMap[peminjaman.KodeUser]; ok {
				peminjam.PasswordHash = "" // Safety
				peminjaman.Peminjam = peminjam
			}
			
			// Link Ruangan
			if peminjaman.KodeRuangan != nil {
				if r, ok := ruanganMap[*peminjaman.KodeRuangan]; ok {
					peminjaman.Ruangan = r
				}
			}
			
			// Link Kegiatan
			if peminjaman.KodeKegiatan != nil {
				if k, ok := kegiatanMap[*peminjaman.KodeKegiatan]; ok {
					peminjaman.Kegiatan = k
				}
			}
			
			kehadiranList[i].Peminjaman = peminjaman
		}
		
		// Link Verifier
		if kehadiranList[i].DiverifikasiOleh != nil {
			if v, ok := verifierMap[*kehadiranList[i].DiverifikasiOleh]; ok {
				v.PasswordHash = ""
				kehadiranList[i].Verifier = v
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kehadiranList)
}

