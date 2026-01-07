package router

import (
	"backend-sarpras/handlers"
	"backend-sarpras/internal/config"
	internalsvc "backend-sarpras/internal/services"
	"backend-sarpras/middleware"
	"backend-sarpras/repositories"
	"backend-sarpras/services"
	"database/sql"
	"net/http"
	"strings"
)

func New(db *sql.DB, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Initialize JWT secret in middleware
	middleware.InitJWTSecret(cfg.JWTSecret)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	ruanganRepo := repositories.NewRuanganRepository(db)
	barangRepo := repositories.NewBarangRepository(db)
	peminjamanRepo := repositories.NewPeminjamanRepository(db)
	kehadiranRepo := repositories.NewKehadiranRepository(db)
	notifikasiRepo := repositories.NewNotifikasiRepository(db)
	logRepo := repositories.NewLogAktivitasRepository(db)
	organisasiRepo := repositories.NewOrganisasiRepository(db)
	kegiatanRepo := repositories.NewKegiatanRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	emailService := internalsvc.NewEmailService()
	peminjamanService := services.NewPeminjamanService(
		peminjamanRepo,
		barangRepo,
		notifikasiRepo,
		logRepo,
		userRepo,
		kegiatanRepo,
		organisasiRepo,
		ruanganRepo,
		emailService,
	)
	kehadiranService := services.NewKehadiranService(
		kehadiranRepo,
		peminjamanRepo,
		logRepo,
	)
	exportService := services.NewExportService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	ruanganHandler := handlers.NewRuanganHandler(ruanganRepo)
	barangHandler := handlers.NewBarangHandler(barangRepo)
	peminjamanHandler := handlers.NewPeminjamanHandler(
		peminjamanService,
		peminjamanRepo,
		ruanganRepo,
		userRepo,
		organisasiRepo,
		kegiatanRepo,
	)
	kehadiranHandler := handlers.NewKehadiranHandler(
		kehadiranService,
		kehadiranRepo,
		peminjamanRepo,
		ruanganRepo,
		userRepo,
		kegiatanRepo,
	)
	exportHandler := handlers.NewExportHandler(
		peminjamanRepo,
		ruanganRepo,
		userRepo,
		organisasiRepo,
		kegiatanRepo,
		exportService,
	)
	notifikasiHandler := handlers.NewNotifikasiHandler(notifikasiRepo)
	logHandler := handlers.NewLogAktivitasHandler(logRepo)
	infoHandler := handlers.InfoUmumHandler
	organisasiHandler := handlers.NewOrganisasiHandler(organisasiRepo)

	// Use centralized CORS middleware from middleware package
	corsMiddleware := middleware.CORSMiddleware

	withAuth := func(h http.Handler) http.Handler {
		return corsMiddleware(middleware.AuthMiddleware(h))
	}

	withRole := func(h http.Handler, roles ...string) http.Handler {
		return corsMiddleware(middleware.AuthMiddleware(middleware.RequireRole(roles...)(h)))
	}

	// Public routes
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	// Public info umum endpoint
	mux.HandleFunc("/api/info", infoHandler)
	// Public organisasi endpoint - for registration dropdown
	mux.HandleFunc("/api/organisasi", organisasiHandler.GetAll)

	// Protected routes - Ruangan
	mux.Handle("/api/ruangan", corsMiddleware(http.HandlerFunc(ruanganHandler.GetAll)))
	mux.Handle("/api/ruangan/", corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Handle /api/ruangan/{kode}/booked-dates endpoint
		if strings.HasSuffix(path, "/booked-dates") {
			if r.Method == http.MethodGet {
				peminjamanHandler.GetBookedDates(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		switch r.Method {
		case http.MethodGet:
			ruanganHandler.GetByID(w, r)
		case http.MethodPut:
			withRole(http.HandlerFunc(ruanganHandler.Update), "SARPRAS", "ADMIN").ServeHTTP(w, r)
		case http.MethodDelete:
			withRole(http.HandlerFunc(ruanganHandler.Delete), "SARPRAS", "ADMIN").ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
	mux.Handle("/api/ruangan/create", withRole(http.HandlerFunc(ruanganHandler.Create), "SARPRAS", "ADMIN"))

	// Protected routes - Barang
	mux.Handle("/api/barang", corsMiddleware(http.HandlerFunc(barangHandler.GetAll)))
	mux.Handle("/api/barang/", corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			barangHandler.GetByID(w, r)
		case http.MethodPut:
			withRole(http.HandlerFunc(barangHandler.Update), "SARPRAS", "ADMIN").ServeHTTP(w, r)
		case http.MethodDelete:
			withRole(http.HandlerFunc(barangHandler.Delete), "SARPRAS", "ADMIN").ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
	mux.Handle("/api/barang/create", withRole(http.HandlerFunc(barangHandler.Create), "SARPRAS", "ADMIN"))

	// Protected routes - Peminjaman
	mux.Handle("/api/peminjaman", withAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			peminjamanHandler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
	mux.HandleFunc("/api/peminjaman/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > len("/api/peminjaman/") {
			remaining := path[len("/api/peminjaman/"):]

			switch {
			case strings.HasSuffix(remaining, "/verifikasi"):
				if r.Method == http.MethodPost {
					withRole(http.HandlerFunc(peminjamanHandler.Verifikasi), "SARPRAS", "ADMIN").ServeHTTP(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}

			case strings.HasSuffix(remaining, "/upload-surat"):
				if r.Method == http.MethodPost {
					withAuth(http.HandlerFunc(peminjamanHandler.UploadSurat)).ServeHTTP(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}

			case strings.HasSuffix(remaining, "/surat"):
				if r.Method == http.MethodGet {
					withAuth(http.HandlerFunc(peminjamanHandler.GetSuratDigital)).ServeHTTP(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}

			case strings.HasSuffix(remaining, "/cancel"):
				if r.Method == http.MethodPost {
					withRole(http.HandlerFunc(peminjamanHandler.CancelPeminjaman), "SARPRAS", "ADMIN").ServeHTTP(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}

			default:
				// Regular ID lookup
				if r.Method == http.MethodGet {
					corsMiddleware(http.HandlerFunc(peminjamanHandler.GetByID)).ServeHTTP(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			}
		}
	})

	mux.Handle("/api/peminjaman/me", withAuth(http.HandlerFunc(peminjamanHandler.GetMyPeminjaman)))
	mux.Handle("/api/peminjaman/pending", withRole(http.HandlerFunc(peminjamanHandler.GetPending), "SARPRAS", "ADMIN"))
	mux.Handle("/api/jadwal-ruangan", corsMiddleware(http.HandlerFunc(peminjamanHandler.GetJadwalRuangan)))
	mux.Handle("/api/jadwal-aktif", corsMiddleware(http.HandlerFunc(peminjamanHandler.GetJadwalAktif)))
	mux.Handle("/api/jadwal-aktif-belum-verifikasi", corsMiddleware(http.HandlerFunc(peminjamanHandler.GetJadwalAktifBelumVerifikasi)))
	mux.Handle("/api/laporan/peminjaman", withRole(http.HandlerFunc(peminjamanHandler.GetLaporan), "SARPRAS", "ADMIN"))
	mux.Handle("/api/laporan/peminjaman/export", withRole(http.HandlerFunc(exportHandler.ExportPeminjamanToExcel), "SARPRAS", "ADMIN"))

	// Protected routes - Kehadiran
	mux.Handle("/api/kehadiran", withRole(http.HandlerFunc(kehadiranHandler.Create), "SECURITY", "ADMIN"))
	mux.Handle("/api/laporan/kehadiran", withRole(http.HandlerFunc(kehadiranHandler.GetByPeminjamanID), "SARPRAS", "ADMIN", "SECURITY"))
	mux.Handle("/api/kehadiran-riwayat", withRole(http.HandlerFunc(kehadiranHandler.GetRiwayatBySecurity), "SECURITY", "ADMIN"))

	// Protected routes - Notifikasi
	mux.Handle("/api/notifikasi/me", withAuth(http.HandlerFunc(notifikasiHandler.GetMyNotifikasi)))
	mux.Handle("/api/notifikasi/count", withAuth(http.HandlerFunc(notifikasiHandler.GetUnreadCount)))
	mux.HandleFunc("/api/notifikasi/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > len("/api/notifikasi/") {
			idOrAction := path[len("/api/notifikasi/"):]
			if idOrAction == "dibaca" || path[len(path)-len("/dibaca"):] == "/dibaca" {
				withAuth(http.HandlerFunc(notifikasiHandler.MarkAsRead)).ServeHTTP(w, r)
			}
		}
	})

	// Protected routes - Log Aktivitas
	mux.Handle("/api/log-aktivitas", withRole(http.HandlerFunc(logHandler.GetAll), "ADMIN"))

	// NOTE: static file serving removed — this backend is API-only.

	return corsMiddleware(mux)
}
