package handlers

import (
	storagesvc "backend-sarpras/internal/services"
	"backend-sarpras/middleware"
	"backend-sarpras/models"
	"backend-sarpras/repositories"
	"backend-sarpras/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type PeminjamanHandler struct {
	PeminjamanService *services.PeminjamanService
	PeminjamanRepo    *repositories.PeminjamanRepository
	RuanganRepo       *repositories.RuanganRepository
	UserRepo          *repositories.UserRepository
	OrganisasiRepo    *repositories.OrganisasiRepository
}

func NewPeminjamanHandler(
	peminjamanService *services.PeminjamanService,
	peminjamanRepo *repositories.PeminjamanRepository,
	ruanganRepo *repositories.RuanganRepository,
	userRepo *repositories.UserRepository,
	organisasiRepo *repositories.OrganisasiRepository,
) *PeminjamanHandler {
	return &PeminjamanHandler{
		PeminjamanService: peminjamanService,
		PeminjamanRepo:    peminjamanRepo,
		RuanganRepo:       ruanganRepo,
		UserRepo:          userRepo,
		OrganisasiRepo:    organisasiRepo,
	}
}

func (h *PeminjamanHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreatePeminjamanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	peminjaman, err := h.PeminjamanService.CreatePeminjaman(&req, user.KodeUser)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(peminjaman)
}

func (h *PeminjamanHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	kode, err := extractKodePeminjaman(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	peminjaman, err := h.PeminjamanRepo.GetByID(kode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if peminjaman == nil {
		http.Error(w, "Peminjaman not found", http.StatusNotFound)
		return
	}

	// Enrich dengan data relasi
	if peminjaman.KodeRuangan != nil {
		ruangan, _ := h.RuanganRepo.GetByID(*peminjaman.KodeRuangan)
		peminjaman.Ruangan = ruangan
	}
	user, _ := h.UserRepo.GetByID(peminjaman.KodeUser)
	if user != nil {
		user.PasswordHash = ""
		peminjaman.Peminjam = user
	}
	items, _ := h.PeminjamanRepo.GetPeminjamanBarang(peminjaman.KodePeminjaman)
	peminjaman.Barang = items

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peminjaman)
}

func (h *PeminjamanHandler) GetMyPeminjaman(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	peminjaman, err := h.PeminjamanRepo.GetByPeminjamID(user.KodeUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enrich dengan data relasi
	for i := range peminjaman {
		if peminjaman[i].KodeRuangan != nil {
			ruangan, _ := h.RuanganRepo.GetByID(*peminjaman[i].KodeRuangan)
			peminjaman[i].Ruangan = ruangan
		}
		user, _ := h.UserRepo.GetByID(peminjaman[i].KodeUser)
		if user != nil {
			user.PasswordHash = ""
			peminjaman[i].Peminjam = user
		}
		items, _ := h.PeminjamanRepo.GetPeminjamanBarang(peminjaman[i].KodePeminjaman)
		peminjaman[i].Barang = items
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peminjaman)
}

func (h *PeminjamanHandler) UploadSurat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Pastikan user login
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Ambil kode peminjaman dari URL: /api/peminjaman/{kode}/upload-surat
	kode, err := extractKodePeminjaman(strings.TrimSuffix(r.URL.Path, "/upload-surat"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Parse multipart form (max 5MB)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		http.Error(w, "Gagal parsing form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("surat")
	if err != nil {
		http.Error(w, "File surat wajib diupload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validasi ukuran max 2MB
	if header.Size > 2*1024*1024 {
		http.Error(w, "Ukuran file maksimal 2MB", http.StatusBadRequest)
		return
	}

	// Validasi MIME type via sniff
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	contentType := http.DetectContentType(buf[:n])
	if contentType != "application/pdf" {
		http.Error(w, "File harus berupa PDF", http.StatusBadRequest)
		return
	}

	// Reset reader ke awal
	if _, err := file.Seek(0, 0); err != nil {
		http.Error(w, "Gagal membaca file", http.StatusInternalServerError)
		return
	}

	// Baca semua bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Gagal membaca file", http.StatusInternalServerError)
		return
	}

	// Path file di bucket (Supabase Storage)
	objectPath := fmt.Sprintf("peminjaman/%s/surat.pdf", kode)

	// Upload ke Supabase Storage
	if err := storagesvc.UploadPDFToSupabase(objectPath, fileBytes); err != nil {
		http.Error(w, "Gagal upload file ke storage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Simpan path ke database lewat repository
	if err := h.PeminjamanRepo.UpdateSuratDigitalURL(kode, objectPath); err != nil {
		http.Error(w, "Gagal menyimpan path surat ke database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":            "Surat berhasil diupload",
		"path_surat_digital": objectPath, // New field name
		"surat_digital_url":  objectPath, // Old field name (backward compatibility)
	})
}

func (h *PeminjamanHandler) GetSuratDigital(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	kode, err := extractKodePeminjaman(strings.TrimSuffix(r.URL.Path, "/surat"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	peminjaman, err := h.PeminjamanRepo.GetByID(kode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if peminjaman == nil {
		http.Error(w, "Peminjaman tidak ditemukan", http.StatusNotFound)
		return
	}

	if peminjaman.PathSuratDigital == "" {
		http.Error(w, "Surat belum tersedia", http.StatusNotFound)
		return
	}

	// Hanya peminjam atau petugas SARPRAS/ADMIN yang boleh mengakses
	if user.KodeUser != peminjaman.KodeUser && user.Role != models.RoleSarpras && user.Role != models.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	signedURL, err := storagesvc.GenerateSignedURL(peminjaman.PathSuratDigital)
	if err != nil {
		http.Error(w, "Gagal membuat tautan surat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"signed_url": signedURL,
		"path":       peminjaman.PathSuratDigital,
	})
}

func (h *PeminjamanHandler) GetPending(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	peminjaman, err := h.PeminjamanRepo.GetPending()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enrich dengan data relasi
	for i := range peminjaman {
		if peminjaman[i].KodeRuangan != nil {
			ruangan, _ := h.RuanganRepo.GetByID(*peminjaman[i].KodeRuangan)
			peminjaman[i].Ruangan = ruangan
		}
		user, _ := h.UserRepo.GetByID(peminjaman[i].KodeUser)
		if user != nil {
			user.PasswordHash = ""
			peminjaman[i].Peminjam = user
		}
		items, _ := h.PeminjamanRepo.GetPeminjamanBarang(peminjaman[i].KodePeminjaman)
		peminjaman[i].Barang = items
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peminjaman)
}

func (h *PeminjamanHandler) Verifikasi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	kode, err := extractKodePeminjaman(strings.TrimSuffix(r.URL.Path, "/verifikasi"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req models.VerifikasiPeminjamanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.PeminjamanService.VerifikasiPeminjaman(kode, user.KodeUser, &req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Verifikasi berhasil"})
}

func (h *PeminjamanHandler) GetJadwalRuangan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var start, end time.Time
	var err error

	if startStr == "" {
		start = time.Now()
	} else {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			http.Error(w, "Invalid start date", http.StatusBadRequest)
			return
		}
	}

	if endStr == "" {
		end = start.AddDate(0, 1, 0) // default 1 bulan ke depan
	} else {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			http.Error(w, "Invalid end date", http.StatusBadRequest)
			return
		}
	}

	jadwal, err := h.PeminjamanRepo.GetJadwalRuangan(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jadwal)
}

func (h *PeminjamanHandler) GetJadwalAktif(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var start, end time.Time
	var err error

	if startStr == "" {
		start = time.Now()
	} else {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			http.Error(w, "Invalid start date", http.StatusBadRequest)
			return
		}
	}

	if endStr == "" {
		end = start.AddDate(0, 0, 7) // default 7 hari ke depan
	} else {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			http.Error(w, "Invalid end date", http.StatusBadRequest)
			return
		}
	}

	peminjaman, err := h.PeminjamanRepo.GetJadwalAktif(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enrich dengan data relasi
	for i := range peminjaman {
		if peminjaman[i].KodeRuangan != nil {
			ruangan, _ := h.RuanganRepo.GetByID(*peminjaman[i].KodeRuangan)
			peminjaman[i].Ruangan = ruangan
		}
		user, _ := h.UserRepo.GetByID(peminjaman[i].KodeUser)
		if user != nil {
			user.PasswordHash = ""
			peminjaman[i].Peminjam = user
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peminjaman)
}

func (h *PeminjamanHandler) GetJadwalAktifBelumVerifikasi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var start, end time.Time
	var err error

	if startStr == "" {
		start = time.Now()
	} else {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			http.Error(w, "Invalid start date", http.StatusBadRequest)
			return
		}
	}

	if endStr == "" {
		end = start.AddDate(0, 0, 7)
	} else {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			http.Error(w, "Invalid end date", http.StatusBadRequest)
			return
		}
	}

	peminjaman, err := h.PeminjamanRepo.GetJadwalAktifBelumVerifikasi(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enrich dengan data relasi
	for i := range peminjaman {
		if peminjaman[i].KodeRuangan != nil {
			ruangan, _ := h.RuanganRepo.GetByID(*peminjaman[i].KodeRuangan)
			peminjaman[i].Ruangan = ruangan
		}
		user, _ := h.UserRepo.GetByID(peminjaman[i].KodeUser)
		if user != nil {
			user.PasswordHash = ""
			peminjaman[i].Peminjam = user
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peminjaman)
}

func (h *PeminjamanHandler) GetLaporan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	peminjaman, err := h.PeminjamanRepo.GetLaporan(start, end, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range peminjaman {
		if peminjaman[i].KodeRuangan != nil {
			ruangan, _ := h.RuanganRepo.GetByID(*peminjaman[i].KodeRuangan)
			peminjaman[i].Ruangan = ruangan
		}
		user, _ := h.UserRepo.GetByID(peminjaman[i].KodeUser)
		if user != nil {
			user.PasswordHash = ""
			// Load organisasi data jika user memiliki organisasi_kode
			if user.OrganisasiKode != nil {
				organisasi, _ := h.OrganisasiRepo.GetByID(*user.OrganisasiKode)
				user.Organisasi = organisasi
			}
			peminjaman[i].Peminjam = user
		}
		items, _ := h.PeminjamanRepo.GetPeminjamanBarang(peminjaman[i].KodePeminjaman)
		peminjaman[i].Barang = items
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peminjaman)
}

func extractKodePeminjaman(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", http.ErrNoLocation
	}

	segments := strings.Split(strings.Trim(path, "/"), "/")
	for i := 0; i < len(segments); i++ {
		if segments[i] == "peminjaman" && i+1 < len(segments) {
			if segments[i+1] == "" {
				return "", http.ErrNoLocation
			}
			return segments[i+1], nil
		}
	}
	return "", http.ErrNoLocation
}
