# Product Backlog - Sistem Peminjaman Sarana Prasarana Kampus

> **Hasil Code Audit & Feature Extraction**  
> Dokumen ini berisi daftar lengkap semua fitur yang sudah diimplementasikan dalam sistem berdasarkan analisis kode sumber.

---

## Ringkasan Sistem

**Arsitektur:**
- **Backend:** Golang (net/http, JWT, PostgreSQL)
- **Frontend:** HTML5/JavaScript (Vanilla)
- **Storage:** Supabase (untuk file PDF surat digital)
- **Email:** Gmail API (OAuth2)

**Role Pengguna:**
| Role | Keterangan |
|------|------------|
| MAHASISWA | Peminjam fasilitas kampus |
| SARPRAS | Admin sarana prasarana (verifikator) |
| SECURITY | Petugas keamanan (verifikasi kehadiran) |
| ADMIN | Super admin dengan akses penuh |

---

## Product Backlog Lengkap

### 1. Fitur Autentikasi (Auth)

| No | Kategori | Nama Fitur | Deskripsi Singkat | Bukti Kode/File |
|----|----------|------------|-------------------|-----------------|
| 1 | Auth | **Login User** | Autentikasi pengguna dengan email & password, menghasilkan JWT token | `handlers/auth_handler.go:Login`, `services/auth_service.go:Login` |
| 2 | Auth | **Register User** | Pendaftaran akun baru dengan validasi email unik dan hashing password (bcrypt) | `handlers/auth_handler.go:Register`, `services/auth_service.go:Register` |
| 3 | Auth | **JWT Token Management** | Generate & validasi token JWT dengan expiry 24 jam | `services/auth_service.go`, `middleware/auth.go` |
| 4 | Auth | **Password Hashing** | Enkripsi password menggunakan bcrypt | `services/auth_service.go:HashPassword` |
| 5 | Auth | **Role-based Redirect** | Redirect otomatis ke dashboard sesuai role setelah login | `assets/js/auth.js:71-80` |

---

### 2. Fitur Inti Peminjaman (Core Feature)

| No | Kategori | Nama Fitur | Deskripsi Singkat | Bukti Kode/File |
|----|----------|------------|-------------------|-----------------|
| 6 | Core Feature | **Form Pengajuan Peminjaman** | Mahasiswa mengisi form untuk mengajukan peminjaman ruangan/barang dengan pilihan tanggal & kegiatan | `handlers/peminjaman_handler.go:Create`, `pengajuan-peminjaman.html` |
| 7 | Core Feature | **Upload Surat Digital** | Upload dokumen surat peminjaman dalam format PDF ke Supabase Storage | `handlers/peminjaman_handler.go:UploadSurat`, `internal/services/storage_service.go` |
| 8 | Core Feature | **Lihat Surat Digital** | Download/lihat surat digital dengan Signed URL (expiry 1 jam) | `handlers/peminjaman_handler.go:GetSuratDigital`, `internal/services/storage_service.go:GenerateSignedURL` |
| 9 | Core Feature | **Riwayat Peminjaman** | List semua peminjaman milik user yang sedang login | `handlers/peminjaman_handler.go:GetMyPeminjaman`, `riwayat-peminjaman.html` |
| 10 | Core Feature | **Detail Peminjaman** | Lihat detail lengkap satu peminjaman (ruangan, barang, status, verifikator) | `handlers/peminjaman_handler.go:GetByID`, `detail-laporan-peminjaman.html` |
| 11 | Core Feature | **Peminjaman Barang** | Tambah daftar barang yang ingin dipinjam (dengan jumlah) dalam satu pengajuan | `models/peminjaman.go:CreatePeminjamanBarang`, `services/peminjaman_service.go:CreatePeminjaman` |
| 12 | Core Feature | **Auto-Create Kegiatan** | Otomatis membuat record kegiatan baru saat pengajuan jika belum ada | `services/peminjaman_service.go:CreatePeminjaman:90-106` |
| 13 | Core Feature | **Validasi Konflik Jadwal** | Cek apakah ruangan sudah di-booking pada rentang waktu yang sama | `repositories/peminjaman_repository.go:CheckRuanganKonflik` |

---

### 3. Fitur Dashboard & Jadwal

| No | Kategori | Nama Fitur | Deskripsi Singkat | Bukti Kode/File |
|----|----------|------------|-------------------|-----------------|
| 14 | Core Feature | **Dashboard Mahasiswa** | Halaman utama mahasiswa dengan statistik peminjaman, jadwal, dan shortcut menu | `dashboard-mahasiswa.html` |
| 15 | Core Feature | **Dashboard Sarpras** | Halaman admin dengan statistik pengajuan pending, approved, rejected | `sarpras.html` |
| 16 | Core Feature | **Dashboard Security** | Halaman security dengan jadwal kegiatan dan verifikasi kehadiran | `dashboard-security.html` |
| 17 | Core Feature | **Dashboard Guest** | Halaman publik untuk melihat jadwal ruangan tanpa login | `dashboard-guest.html`, `assets/js/dashboard-guest.js` |
| 18 | Core Feature | **Jadwal Ruangan** | Tampilan kalender interaktif dengan jadwal peminjaman semua ruangan | `handlers/peminjaman_handler.go:GetJadwalRuangan`, `jadwal-ruangan.html` |
| 19 | Core Feature | **Jadwal Aktif Harian** | List jadwal hari ini yang statusnya APPROVED atau ONGOING | `handlers/peminjaman_handler.go:GetJadwalAktif` |
| 20 | Core Feature | **Widget Kalender Interaktif** | Komponen kalender dengan dot indicator untuk hari yang ada jadwal | `assets/js/calendar-widget.js` |
| 21 | Core Feature | **Booked Dates API** | Mendapatkan tanggal-tanggal yang sudah di-booking untuk ruangan tertentu (untuk disable di calendar picker) | `handlers/peminjaman_handler.go:GetBookedDates` |

---

### 4. Fitur Admin/Sarpras

| No | Kategori | Nama Fitur | Deskripsi Singkat | Bukti Kode/File |
|----|----------|------------|-------------------|-----------------|
| 22 | Admin | **Verifikasi Peminjaman** | Approve/Reject pengajuan peminjaman dengan catatan | `handlers/peminjaman_handler.go:Verifikasi`, `verifikasi-peminjaman.html` |
| 23 | Admin | **List Pengajuan Pending** | Melihat semua pengajuan dengan status PENDING | `handlers/peminjaman_handler.go:GetPending` |
| 24 | Admin | **Pembatalan Peminjaman** | Membatalkan peminjaman yang sudah APPROVED/ONGOING dengan alasan | `handlers/peminjaman_handler.go:CancelPeminjaman`, `services/peminjaman_service.go:CancelPeminjaman` |
| 25 | Admin | **Laporan Peminjaman** | Rekap peminjaman dengan filter tanggal dan status | `handlers/peminjaman_handler.go:GetLaporan`, `laporan-peminjaman.html` |
| 26 | Admin | **Export Excel** | Export laporan peminjaman ke file Excel (.xlsx) dengan styling | `handlers/export_handler.go:ExportPeminjamanToExcel`, `services/export_service.go:GeneratePeminjamanExcel` |
| 27 | Admin | **Kelola Ruangan (CRUD)** | Tambah, edit, hapus, lihat data ruangan | `handlers/ruangan_handler.go`, `kelola-ruangan.html` |
| 28 | Admin | **Kelola Barang (CRUD)** | Tambah, edit, hapus, lihat data inventaris barang | `handlers/barang_handler.go`, `kelola-barang.html` |
| 29 | Admin | **Lihat Log Aktivitas** | Melihat semua log aktivitas sistem (audit trail) | `handlers/log_aktivitas_handler.go:GetAll` |

---

### 5. Fitur Security

| No | Kategori | Nama Fitur | Deskripsi Singkat | Bukti Kode/File |
|----|----------|------------|-------------------|-----------------|
| 30 | Core Feature | **Verifikasi Kehadiran** | Security mencatat kehadiran peminjam (HADIR/TIDAK_HADIR/BATAL) | `handlers/kehadiran_handler.go:Create`, `services/kehadiran_service.go:CreateKehadiran` |
| 31 | Core Feature | **Riwayat Kehadiran** | List riwayat kehadiran yang dicatat oleh security | `handlers/kehadiran_handler.go:GetRiwayatBySecurity`, `riwayat-kehadiran.html` |
| 32 | Core Feature | **Jadwal Aktif Belum Diverifikasi** | List jadwal hari ini yang belum diverifikasi kehadirannya | `handlers/peminjaman_handler.go:GetJadwalAktifBelumVerifikasi` |

---

### 6. Fitur Email Notification

| No | Kategori | Nama Fitur | Deskripsi Singkat | Bukti Kode/File |
|----|----------|------------|-------------------|-----------------|
| 33 | Support | **Email Approved Notification** | Email ke mahasiswa saat peminjaman disetujui | `internal/services/email_templates.go:BuildApprovedEmailHTML` |
| 34 | Support | **Email Rejected Notification** | Email ke mahasiswa saat peminjaman ditolak | `internal/services/email_templates.go:BuildRejectedEmailHTML` |
| 35 | Support | **Email Security Notification** | Email ke security tentang jadwal kegiatan yang disetujui | `internal/services/email_templates.go:BuildSecurityNotificationHTML` |
| 36 | Support | **Email Cancellation Notification** | Email ke mahasiswa saat peminjaman dibatalkan | `internal/services/email_templates.go:BuildCancelledEmailHTML` |
| 37 | Support | **Async Email Sending** | Pengiriman email secara asynchronous (non-blocking) | `services/peminjaman_service.go:sendVerificationEmails`, `services/peminjaman_service.go:sendCancellationEmail` |

---

### 7. Fitur Pendukung (Support)

| No | Kategori | Nama Fitur | Deskripsi Singkat | Bukti Kode/File |
|----|----------|------------|-------------------|-----------------|
| 38 | Support | **Info Instansi** | Endpoint publik untuk informasi umum instansi | `handlers/info_handler.go:InfoUmumHandler` |
| 39 | Support | **List Organisasi** | Dropdown list organisasi (ORMAWA/UKM) untuk registrasi | `handlers/organisasi_handler.go:GetAll` |
| 40 | Support | **Pendaftaran Organisasi** | Halaman untuk mendaftarkan organisasi baru | `pendaftaran-organisasi.html` |
| 41 | Support | **Health Check** | Endpoint untuk monitoring status server | `router/router.go:100-103` |

---

### 8. Fitur Teknis (Technical Enablers)

| No | Kategori | Nama Fitur | Deskripsi Singkat | Bukti Kode/File |
|----|----------|------------|-------------------|-----------------|
| 42 | Technical | **JWT Authentication Middleware** | Middleware untuk validasi JWT token pada protected routes | `middleware/auth.go:AuthMiddleware` |
| 43 | Technical | **Role-based Authorization** | Middleware untuk membatasi akses berdasarkan role user | `middleware/auth.go:RequireRole` |
| 44 | Technical | **CORS Middleware** | Middleware untuk mengatur Cross-Origin Resource Sharing | `middleware/cors.go:CORSMiddleware` |
| 45 | Technical | **Status Scheduler (Cron Job)** | Background job untuk auto-update status APPROVED→ONGOING→FINISHED | `internal/services/status_scheduler.go` |
| 46 | Technical | **Supabase Storage Integration** | Integrasi upload/download file ke Supabase Storage | `internal/services/storage_service.go` |
| 47 | Technical | **Gmail API Integration** | Integrasi pengiriman email via Gmail API (OAuth2) | `internal/services/email_service.go`, `internal/config/gmail.go` |
| 48 | Technical | **Database Migration** | Script SQL untuk inisialisasi schema database | `migrations/001_init_schema.sql` |
| 49 | Technical | **Auto-Generate Kode** | Trigger database untuk auto-generate kode (USR-xxx, PMJ-xxx, dll) | `migrations/002_auto_generate_codes.sql` |
| 50 | Technical | **Database Indexing** | Index pada kolom-kolom penting untuk optimasi query | `migrations/003_add_indexes.sql` |
| 51 | Technical | **Seed Data** | Data awal untuk testing (ruangan, organisasi, user) | `migrations/seed_initial_data.sql` |
| 52 | Technical | **Log Aktivitas (Audit Trail)** | Pencatatan semua aksi penting dalam sistem | `repositories/log_aktivitas_repository.go`, `models/log_aktivitas.go` |
| 53 | Technical | **API Wrapper Module** | Modul JavaScript untuk API calls dengan auto-auth | `assets/js/api.js` |
| 54 | Technical | **Config Management** | Konfigurasi environment variables (JWT, DB, Supabase, Gmail) | `internal/config/config.go`, `internal/config/supabase.go`, `internal/config/gmail.go` |
| 55 | Technical | **File Move (Rename)** | Fungsi untuk memindahkan file di storage saat status berubah | `internal/services/storage_service.go:MoveFile` |
| 56 | Technical | **Hot Reload Development** | Konfigurasi Air untuk hot reload saat development | `.air.toml`, `.air.windows.toml`, `.air.linux.toml` |

---

## Ringkasan Jumlah Fitur

| Kategori | Jumlah |
|----------|--------|
| Auth | 5 |
| Core Feature | 21 |
| Admin | 8 |
| Support | 6 |
| Technical | 15 |
| **Total** | **56** |

---

## Struktur Database (ERD Summary)

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  organisasi │◄────│    users    │────►│  kehadiran  │
└─────────────┘     └──────┬──────┘     └─────────────┘
                          │
                          ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   ruangan   │◄────│ peminjaman  │────►│  kegiatan   │
└──────┬──────┘     └──────┬──────┘     └─────────────┘
       │                   │
       ▼                   ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   barang    │◄────│ peminjaman_ │     │  notifikasi │
└─────────────┘     │    barang   │     └─────────────┘
                    └─────────────┘     
                                        ┌─────────────┐
                                        │log_aktivitas│
                                        └─────────────┘
```

---

## Catatan Penting untuk Bab 3 Laporan

1. **Status Peminjaman** memiliki lifecycle: `PENDING` → `APPROVED`/`REJECTED` → `ONGOING` → `FINISHED`/`CANCELLED`

2. **Transisi Status Otomatis** dihandle oleh `StatusScheduler` yang berjalan setiap 1 menit

3. **Trigger Notifikasi:**
   - Pengajuan dibuat → Notifikasi ke SARPRAS
   - Diverifikasi → Email + Notifikasi ke Mahasiswa
   - Approved → Email ke Security
   - Dibatalkan → Email ke Mahasiswa

4. **Storage:** Semua file surat digital disimpan di Supabase Storage dengan path: `surat-digital/{status}/{kode_peminjaman}.pdf`

5. **Excel Export** menggunakan library `excelize/v2` dengan styling profesional (warna sesuai status)
