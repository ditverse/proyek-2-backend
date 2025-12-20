package main

import (
	"backend-sarpras/internal/config"
	"backend-sarpras/internal/db"
	"backend-sarpras/internal/router"
	"backend-sarpras/internal/scheduler"
	"backend-sarpras/middleware"
	"backend-sarpras/repositories"
	"backend-sarpras/services"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load konfigurasi dari environment variable
	cfg := config.Load()

	// Initialize JWT secret
	middleware.InitJWTSecret(cfg.JWTSecret)

	// Buka koneksi database Supabase
	conn := db.Open(cfg.DatabaseURL)
	defer conn.Close()

	// Setup router dengan config
	handler := router.New(conn, cfg)

	// Initialize repositories and services for scheduler
	peminjamanRepo := repositories.NewPeminjamanRepository(conn)
	userRepo := repositories.NewUserRepository(conn)
	notifikasiRepo := repositories.NewNotifikasiRepository(conn)
	ruanganRepo := repositories.NewRuanganRepository(conn)
	whatsappService := services.NewWhatsappService(cfg)

	// Start reminder scheduler (runs every 5 minutes)
	reminderScheduler := scheduler.NewReminderScheduler(
		conn,
		peminjamanRepo,
		userRepo,
		notifikasiRepo,
		ruanganRepo,
		whatsappService,
	)
	reminderScheduler.Start(5 * time.Minute)
	defer reminderScheduler.Stop()

	// Create HTTP server dengan timeout
	server := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        handler,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Jalankan server HTTP di goroutine terpisah
	log.Printf("Server running on http://localhost:%s", cfg.Port)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\n⏱️  Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server shutdown error: %v", err)
	}

	log.Println("✅ Server shutdown complete")
}
