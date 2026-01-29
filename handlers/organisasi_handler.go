package handlers

import (
	"backend-sarpras/models"
	"backend-sarpras/repositories"
	"encoding/json"
	"net/http"
)

type OrganisasiHandler struct {
	OrganisasiRepo *repositories.OrganisasiRepository
}

func NewOrganisasiHandler(organisasiRepo *repositories.OrganisasiRepository) *OrganisasiHandler {
	return &OrganisasiHandler{OrganisasiRepo: organisasiRepo}
}

// OrganisasiListItem is a simplified response structure for dropdown/list
type OrganisasiListItem struct {
	KodeOrganisasi string `json:"kode_organisasi"`
	NamaOrganisasi string `json:"nama_organisasi"`
}

// OrganisasiCreateRequest is the request structure for creating a new organization
type OrganisasiCreateRequest struct {
	KodeOrganisasi  string `json:"kode_organisasi"`
	Nama            string `json:"nama"`
	JenisOrganisasi string `json:"jenis_organisasi"`
	Kontak          string `json:"kontak"`
}

// GetAll - GET /api/organisasi - Public endpoint to get list of all organizations
// Returns a simplified list suitable for dropdowns
func (h *OrganisasiHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	orgs, err := h.OrganisasiRepo.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch organizations", http.StatusInternalServerError)
		return
	}

	// Map to simplified response format expected by frontend
	result := make([]OrganisasiListItem, len(orgs))
	for i, org := range orgs {
		result[i] = OrganisasiListItem{
			KodeOrganisasi: org.KodeOrganisasi,
			NamaOrganisasi: org.NamaOrganisasi,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Create - POST /api/organisasi - Public endpoint to create a new organization
func (h *OrganisasiHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req OrganisasiCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.KodeOrganisasi == "" || req.Nama == "" || req.JenisOrganisasi == "" {
		http.Error(w, "Kode organisasi, nama, dan jenis organisasi wajib diisi", http.StatusBadRequest)
		return
	}

	// Check if organization already exists
	existing, _ := h.OrganisasiRepo.GetByID(req.KodeOrganisasi)
	if existing != nil {
		http.Error(w, "Kode organisasi sudah terdaftar", http.StatusConflict)
		return
	}

	// Create organization model
	org := &models.Organisasi{
		KodeOrganisasi:  req.KodeOrganisasi,
		NamaOrganisasi:  req.Nama,
		JenisOrganisasi: models.JenisOrganisasiEnum(req.JenisOrganisasi),
		Kontak:          req.Kontak,
	}

	if err := h.OrganisasiRepo.Create(org); err != nil {
		http.Error(w, "Gagal mendaftarkan organisasi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":         "Organisasi berhasil didaftarkan",
		"kode_organisasi": org.KodeOrganisasi,
	})
}
