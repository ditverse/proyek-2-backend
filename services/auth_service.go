package services

import (
	"backend-sarpras/models"
	"backend-sarpras/repositories"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo  *repositories.UserRepository
	JWTSecret string
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		JWTSecret: jwtSecret,
	}
}

func (s *AuthService) Login(email, password string) (*models.LoginResponse, error) {
	user, err := s.UserRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("Email tidak ditemukan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("Password salah")
	}

	claims := jwt.MapClaims{
		"kode_user": user.KodeUser,
		"email":     user.Email,
		"role":      user.Role,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return nil, err
	}

	user.PasswordHash = ""

	return &models.LoginResponse{
		Token: tokenString,
		User:  *user,
	}, nil
}

func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	if req.Email == "" || req.Password == "" || req.Nama == "" {
		return nil, errors.New("nama, email, dan password wajib diisi")
	}
	if req.Role == "" {
		req.Role = models.RoleMahasiswa
	}

	existing, err := s.UserRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	passwordHash, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Nama:           req.Nama,
		Email:          req.Email,
		PasswordHash:   passwordHash,
		Role:           req.Role,
		NoHP:           req.NoHP,
		OrganisasiKode: req.OrganisasiKode,
	}

	if err := s.UserRepo.Create(user); err != nil {
		return nil, err
	}

	user.PasswordHash = ""
	return user, nil
}
