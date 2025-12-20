package models

import "time"

type User struct {
	KodeUser       string      `json:"kode_user"`
	Nama           string      `json:"nama"`
	Email          string      `json:"email"`
	PasswordHash   string      `json:"-"` // tidak di-expose ke JSON
	Role           RoleEnum    `json:"role"`
	NoHP           *string     `json:"no_hp"`
	OrganisasiKode *string     `json:"organisasi_kode"`
	Organisasi     *Organisasi `json:"organisasi,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RegisterRequest struct {
	Nama           string   `json:"nama"`
	Email          string   `json:"email"`
	Password       string   `json:"password"`
	Role           RoleEnum `json:"role"`
	OrganisasiKode *string  `json:"organisasi_kode"`
}
