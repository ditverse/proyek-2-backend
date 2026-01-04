package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DatabaseURL       string
	JWTSecret         string
	CORSAllowedOrigin string
}

func Load() *Config {
	// Load .env file jika ada (untuk local development)
	_ = godotenv.Load()

	// Ambil env variables - required
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// Jangan fatal di development — pakai default dev secret dan log warning.
		// Untuk production, selalu set env JWT_SECRET ke nilai kuat.
		jwtSecret = "dev-secret-change-me-in-production"
		log.Println("⚠️  JWT_SECRET not set — using development default (change in production)")
	}

	// Optional env variables - dengan default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	corsOrigin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "*" // default untuk development
	}

	return &Config{
		Port:              port,
		DatabaseURL:       dbURL,
		JWTSecret:         jwtSecret,
		CORSAllowedOrigin: corsOrigin,
	}
}
