package handlers

import (
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

// GetAll - GET /api/organisasi - Public endpoint to get list of all organizations
// Returns a simplified list suitable for dropdowns
func (h *OrganisasiHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
			NamaOrganisasi: org.Nama,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
