# Sistem Informasi Peminjaman Sarana dan Prasarana Kampus

> Backend API untuk sistem peminjaman ruangan dan barang kampus dengan role-based access control

[![Go Version](https://img.shields.io/badge/Go-1.25.3-00ADD8?style=flat&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supabase-336791?style=flat&logo=postgresql)](https://supabase.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## Daftar Isi

- [Tentang Project](#tentang-project)
- [Fitur Utama](#fitur-utama)
- [Teknologi](#teknologi)
- [Struktur Project](#struktur-project)
- [Setup dan Installation](#setup-dan-installation)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [Role dan Permissions](#role-dan-permissions)
- [Development](#development)
- [Troubleshooting](#troubleshooting)

---

## Tentang Project

Sistem web-based untuk mengelola peminjaman ruangan dan barang di lingkungan kampus. Sistem ini mendukung workflow lengkap dari pengajuan peminjaman oleh mahasiswa, verifikasi oleh petugas sarpras, hingga pencatatan kehadiran oleh petugas security.

### Workflow Peminjaman

```
+--------------+      +---------------+      +--------------+      +---------------+
|  MAHASISWA   |----->|    SARPRAS    |----->|   SECURITY   |----->|   FINISHED    |
|  Mengajukan  |      |  Verifikasi   |      |  Verifikasi  |      |    Selesai    |
|  Peminjaman  |      |  APPROVED/    |      |  Kehadiran   |      |               |
|              |      |  REJECTED     |      |              |      |               |
+--------------+      +---------------+      +--------------+      +---------------+
    PENDING              APPROVED              ONGOING              FINISHED
```

---

## Fitur Utama

### Authentication dan Authorization
- JWT-based authentication
- Role-based access control (RBAC)
- 4 Role: MAHASISWA, SARPRAS, SECURITY, ADMIN

### Master Data Management
- CRUD Ruangan (kapasitas, lokasi, deskripsi)
- CRUD Barang (stok, lokasi penyimpanan)
- Manajemen Organisasi (ORMAWA, UKM)

### Peminjaman
- Pengajuan peminjaman ruangan/barang
- Upload surat digital ke Supabase Storage
- Multi-item peminjaman (ruangan + barang)
- Verifikasi oleh petugas sarpras
- Pembatalan peminjaman yang sudah disetujui
- Status tracking (PENDING, APPROVED, REJECTED, ONGOING, FINISHED, CANCELLED)

### Kehadiran
- Verifikasi kehadiran oleh security
- Riwayat kehadiran peminjam
- Status kehadiran (HADIR, TIDAK_HADIR, BATAL)

### Email Notification
- Notifikasi email otomatis via Gmail API
- Notifikasi saat pengajuan baru dibuat (ke Sarpras)
- Notifikasi status approved/rejected (ke Mahasiswa)
- Notifikasi kehadiran (ke Security)
- Notifikasi pembatalan peminjaman

### Laporan dan Export
- Laporan peminjaman (filter by date, status)
- Laporan kehadiran
- Export laporan ke Excel
- Log aktivitas sistem

---

## Teknologi

### Backend
- **Language**: Go 1.25.3 (Native, tanpa framework)
- **Database**: PostgreSQL (Supabase)
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **Database Driver**: pgx/v5
- **Email**: Gmail API (google.golang.org/api/gmail)
- **Excel Export**: excelize

### Storage
- **File Storage**: Supabase Storage
- **Supported Files**: PDF (surat digital)

### Development Tools
- **Hot Reload**: Air
- **Environment**: godotenv

---

## Struktur Project

```
proyek-2-backend/
├── cmd/
│   ├── server/
│   │   └── main.go                    # Entry point aplikasi
│   └── oauth_token/                   # Tool untuk generate Gmail OAuth token
├── internal/
│   ├── config/
│   │   ├── config.go                  # Environment configuration
│   │   ├── gmail.go                   # Gmail API configuration
│   │   └── supabase.go                # Supabase storage config
│   ├── db/
│   │   └── db.go                      # Database connection
│   ├── router/
│   │   └── router.go                  # HTTP routing dan middleware setup
│   └── services/
│       ├── email_service.go           # Gmail API email sending
│       ├── email_templates.go         # HTML email templates
│       ├── status_scheduler.go        # Auto-update status scheduler
│       └── storage_service.go         # Supabase storage operations
├── models/                            # Domain models dan DTOs
│   ├── user.go
│   ├── peminjaman.go
│   ├── ruangan.go
│   ├── barang.go
│   ├── kehadiran.go
│   ├── kegiatan.go
│   ├── organisasi.go
│   ├── mailbox.go                     # Email notification log
│   ├── log_aktivitas.go
│   └── enums.go                       # Enum definitions
├── repositories/                      # Data access layer (CRUD)
│   ├── user_repository.go
│   ├── peminjaman_repository.go
│   ├── ruangan_repository.go
│   ├── barang_repository.go
│   ├── kehadiran_repository.go
│   ├── kegiatan_repository.go
│   ├── organisasi_repository.go
│   ├── mailbox_repository.go
│   └── log_aktivitas_repository.go
├── services/                          # Business logic layer
│   ├── auth_service.go
│   ├── peminjaman_service.go
│   ├── kehadiran_service.go
│   ├── export_service.go
│   └── code_generator.go
├── handlers/                          # HTTP handlers (controllers)
│   ├── auth_handler.go
│   ├── peminjaman_handler.go
│   ├── ruangan_handler.go
│   ├── barang_handler.go
│   ├── kehadiran_handler.go
│   ├── export_handler.go
│   ├── organisasi_handler.go
│   ├── log_aktivitas_handler.go
│   └── info_handler.go
├── middleware/                        # HTTP middleware
│   ├── auth.go                        # JWT validation dan role checking
│   └── cors.go                        # CORS configuration
├── migrations/                        # SQL migration files
│   ├── erd_new_proyek_2.sql           # Current schema (with triggers)
│   ├── 002_auto_generate_codes.sql
│   ├── 004_redesign_mailbox.sql
│   └── schema_database.sql
├── docs/                              # Documentation
│   ├── diagrams/                      # UML diagrams
│   │   ├── use_case_diagram.puml      # PlantUML Use Case
│   │   ├── class_diagram.md           # Mermaid Class Diagram
│   │   ├── sequence_diagram.md        # Mermaid Sequence Diagrams
│   │   └── flowchart_diagram.md       # Mermaid Flowcharts
│   ├── bpmn/                          # BPMN business process diagrams
│   ├── product_backlog.md
│   ├── sprint_planning.md
│   └── ...
├── .air.toml                          # Air configuration (default)
├── .air.windows.toml                  # Air configuration (Windows)
├── .air.linux.toml                    # Air configuration (Linux)
├── run-air.bat                        # Air wrapper script (Windows)
├── run-air.sh                         # Air wrapper script (Linux)
├── .env                               # Environment variables (gitignored)
├── .env.example                       # Environment variables template
├── gmail_token.json                   # Gmail OAuth token (gitignored)
├── go.mod
├── go.sum
└── README.md
```

---

## Setup dan Installation

### Prerequisites

- Go 1.25.3 atau lebih tinggi
- PostgreSQL database (Supabase account)
- Git
- Google Cloud Console account (untuk Gmail API)

### 1. Clone Repository

```bash
git clone <repository-url>
cd proyek-2-backend
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Setup Database

1. Buat project di [Supabase](https://supabase.com/)
2. Buka **SQL Editor** di Supabase Dashboard
3. Jalankan migration file:
   - `migrations/erd_new_proyek_2.sql` (schema utama)
   - `migrations/002_auto_generate_codes.sql` (triggers untuk auto-generate kode)
   - `migrations/004_redesign_mailbox.sql` (mailbox untuk email logging)

### 4. Setup Supabase Storage

1. Buka **Storage** di Supabase Dashboard
2. Buat bucket baru: `surat-digital`
3. Set policy untuk bucket (public read, authenticated write)

### 5. Setup Gmail API (Optional)

Untuk mengaktifkan fitur notifikasi email:

1. Buat project di [Google Cloud Console](https://console.cloud.google.com/)
2. Enable Gmail API
3. Buat OAuth 2.0 credentials (Desktop app)
4. Download `credentials.json` dan letakkan di root project
5. Jalankan token generator:
   ```bash
   go run cmd/oauth_token/main.go
   ```
6. Ikuti instruksi untuk generate `gmail_token.json`

### 6. Environment Variables

Buat file `.env` berdasarkan `.env.example`:

```env
# Database
DATABASE_URL=postgresql://postgres:[PASSWORD]@[HOST]:[PORT]/postgres

# Server
PORT=8000

# JWT Secret
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Supabase Storage
SUPABASE_URL=https://[PROJECT-ID].supabase.co
SUPABASE_SERVICE_KEY=your-service-role-key
SUPABASE_BUCKET_NAME=surat-digital

# Storage Configuration
STORAGE_SIGNED_URL_EXPIRES=600
MAX_UPLOAD_SIZE_MB=2

# CORS
CORS_ALLOWED_ORIGIN=*

# Gmail API (Optional - untuk notifikasi email)
GMAIL_CREDENTIALS_FILE=credentials.json
GMAIL_TOKEN_FILE=gmail_token.json
GMAIL_SENDER_EMAIL=your-email@gmail.com
```

### 7. Run Server

#### Development (with hot reload)

```bash
# Install Air (jika belum)
go install github.com/air-verse/air@latest

# Run with Air
air

# Atau menggunakan script wrapper
# Windows:
run-air.bat

# Linux/macOS:
./run-air.sh
```

#### Production

```bash
go run cmd/server/main.go
```

Server akan berjalan di `http://localhost:8000`

### 8. Test API

```bash
# Health check
curl http://localhost:8000/api/health

# Login
curl -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'
```

---

## API Documentation

### Base URL
```
http://localhost:8000/api
```

### Authentication

Semua endpoint yang memerlukan autentikasi harus menyertakan header:
```
Authorization: Bearer <JWT_TOKEN>
```

### Endpoints

#### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/login` | Login user |
| POST | `/auth/register` | Register user baru (hanya oleh SARPRAS) |
| GET | `/health` | Health check |
| GET | `/info` | Info umum sistem |
| GET | `/organisasi` | List organisasi (untuk dropdown register) |
| GET | `/jadwal-ruangan` | Jadwal ruangan publik |
| GET | `/jadwal-aktif` | Jadwal aktif hari ini |

#### Protected Endpoints

##### Master Data - Ruangan

| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| GET | `/ruangan` | All | List semua ruangan |
| GET | `/ruangan/{id}` | All | Detail ruangan |
| GET | `/ruangan/{id}/booked-dates` | All | Tanggal yang sudah dibooking |
| POST | `/ruangan/create` | SARPRAS, ADMIN | Tambah ruangan |
| PUT | `/ruangan/{id}` | SARPRAS, ADMIN | Update ruangan |
| DELETE | `/ruangan/{id}` | SARPRAS, ADMIN | Hapus ruangan |

##### Master Data - Barang

| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| GET | `/barang` | All | List semua barang |
| GET | `/barang/{id}` | All | Detail barang |
| POST | `/barang/create` | SARPRAS, ADMIN | Tambah barang |
| PUT | `/barang/{id}` | SARPRAS, ADMIN | Update barang |
| DELETE | `/barang/{id}` | SARPRAS, ADMIN | Hapus barang |

##### Peminjaman

| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| POST | `/peminjaman` | Authenticated | Buat pengajuan peminjaman |
| GET | `/peminjaman/me` | Authenticated | List peminjaman milik user |
| GET | `/peminjaman/{id}` | All | Detail peminjaman |
| GET | `/peminjaman/pending` | SARPRAS, ADMIN | List pengajuan pending |
| POST | `/peminjaman/{id}/verifikasi` | SARPRAS, ADMIN | Verifikasi peminjaman |
| POST | `/peminjaman/{id}/cancel` | SARPRAS, ADMIN | Batalkan peminjaman |
| POST | `/peminjaman/{id}/upload-surat` | Authenticated | Upload surat digital |
| GET | `/peminjaman/{id}/surat` | Authenticated | Get signed URL surat |
| GET | `/jadwal-aktif-belum-verifikasi` | SECURITY, ADMIN | Jadwal belum verifikasi kehadiran |
| GET | `/laporan/peminjaman` | SARPRAS, ADMIN | Laporan peminjaman |
| GET | `/laporan/peminjaman/export` | SARPRAS, ADMIN | Export laporan ke Excel |

##### Kehadiran

| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| POST | `/kehadiran` | SECURITY, ADMIN | Catat kehadiran peminjam |
| GET | `/laporan/kehadiran` | SARPRAS, SECURITY, ADMIN | Laporan kehadiran |
| GET | `/kehadiran-riwayat` | SECURITY, ADMIN | Riwayat kehadiran by security |

##### Log Aktivitas

| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| GET | `/log-aktivitas` | ADMIN | List semua log aktivitas |

### Request/Response Examples

#### Login

**Request:**
```json
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "kode_user": "USR-251204-0001",
    "nama": "John Doe",
    "email": "john@example.com",
    "role": "MAHASISWA",
    "created_at": "2025-12-04T08:00:00Z"
  }
}
```

#### Create Peminjaman

**Request:**
```json
POST /api/peminjaman
Authorization: Bearer <token>
Content-Type: application/json

{
  "kode_ruangan": "RNG-0001",
  "nama_kegiatan": "Rapat Organisasi",
  "deskripsi": "Rapat bulanan",
  "tanggal_mulai": "2025-12-10T08:00:00Z",
  "tanggal_selesai": "2025-12-10T12:00:00Z",
  "barang": [
    {
      "kode_barang": "BRG-0001",
      "jumlah": 2
    }
  ]
}
```

**Response:**
```json
{
  "kode_peminjaman": "PMJ-251204-0001",
  "kode_user": "USR-251204-0001",
  "kode_ruangan": "RNG-0001",
  "kode_kegiatan": "KGT-251204-0001",
  "tanggal_mulai": "2025-12-10T08:00:00Z",
  "tanggal_selesai": "2025-12-10T12:00:00Z",
  "status": "PENDING",
  "created_at": "2025-12-04T08:00:00Z"
}
```

---

## Database Schema

### Tabel Utama

#### users
```sql
kode_user       VARCHAR PRIMARY KEY  -- Format: USR-YYMMDD-0001
nama            VARCHAR
email           VARCHAR UNIQUE
password_hash   VARCHAR
role            role_enum            -- MAHASISWA, SARPRAS, SECURITY, ADMIN
no_hp           VARCHAR
organisasi_kode VARCHAR FK
created_at      TIMESTAMP
```

#### peminjaman
```sql
kode_peminjaman      VARCHAR PRIMARY KEY  -- Format: PMJ-YYMMDD-0001
kode_user            VARCHAR FK
kode_ruangan         VARCHAR FK
kode_kegiatan        VARCHAR FK
tanggal_mulai        TIMESTAMP
tanggal_selesai      TIMESTAMP
status               peminjaman_status_enum
path_surat_digital   TEXT
verified_by          VARCHAR FK
verified_at          TIMESTAMP
catatan_verifikasi   TEXT
created_at           TIMESTAMP
updated_at           TIMESTAMP
```

#### mailbox
```sql
kode_mailbox     VARCHAR PRIMARY KEY  -- Format: MBX-YYMMDD-0001
kode_user        VARCHAR FK           -- Penerima notifikasi
kode_peminjaman  VARCHAR FK
jenis_pesan      VARCHAR              -- APPROVED, REJECTED, CANCELLED, etc.
created_at       TIMESTAMP
```

### Database Triggers

Sistem menggunakan database triggers untuk auto-generate kode dengan format yang konsisten:

- **`generate_kode_user()`**: Generate `USR-YYMMDD-0001`
- **`generate_kode_peminjaman()`**: Generate `PMJ-YYMMDD-0001`
- **`generate_kode_ruangan()`**: Generate `RNG-0001`
- **`generate_kode_barang()`**: Generate `BRG-0001`
- **`generate_kode_mailbox()`**: Generate `MBX-YYMMDD-0001`

---

## Role dan Permissions

### MAHASISWA
- Melihat jadwal ruangan
- Mengajukan peminjaman
- Upload surat digital
- Melihat riwayat peminjaman sendiri
- Menerima notifikasi email status peminjaman

### SARPRAS (Sarana Prasarana / Admin)
- Registrasi user baru
- Kelola master data (ruangan, barang, organisasi)
- Verifikasi pengajuan peminjaman (approve/reject)
- Membatalkan peminjaman yang sudah disetujui
- Melihat laporan peminjaman
- Export laporan ke Excel
- Melihat log aktivitas

### SECURITY
- Melihat jadwal peminjaman aktif
- Mencatat kehadiran peminjam
- Melihat riwayat kehadiran
- Melihat jadwal yang belum diverifikasi kehadirannya
- Menerima notifikasi email untuk jadwal aktif

### ADMIN
- Semua akses (full access)
- Melihat log aktivitas sistem

---

## Development

### Hot Reload dengan Air

Project ini mendukung Air untuk Windows dan Linux dengan konfigurasi terpisah:

| File | OS | Binary |
|------|-----|--------|
| `.air.windows.toml` | Windows | `tmp\main.exe` |
| `.air.linux.toml` | Linux/macOS | `./tmp/main` |

#### Install Air

```bash
go install github.com/air-verse/air@latest
```

#### Menjalankan Air

```bash
# Menggunakan Script Wrapper (Rekomendasi)
# Linux/macOS
./run-air.sh

# Windows
run-air.bat

# Atau manually dengan config file
air -c .air.windows.toml  # Windows
air -c .air.linux.toml    # Linux
```

### Layered Architecture

```
Handler -> Service -> Repository -> Database
```

- **Handlers**: HTTP request/response handling
- **Services**: Business logic
- **Repositories**: Database operations
- **Models**: Data structures

### Naming Conventions

- **Kode**: `PREFIX-YYMMDD-0001` atau `PREFIX-0001`
  - User: `USR-251204-0001`
  - Peminjaman: `PMJ-251204-0001`
  - Ruangan: `RNG-0001`
  - Barang: `BRG-0001`

- **Enums**: PascalCase dengan suffix `Enum`
  - `RoleEnum`, `PeminjamanStatusEnum`

- **Functions**: camelCase
  - `GetByID()`, `Create()`, `UpdateStatus()`

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestFunctionName ./path/to/package
```

---

## Troubleshooting

### Issue: Kode User Format Salah

**Problem**: Kode user tersimpan sebagai `USR-1733295825123456789` bukan `USR-251204-0001`

**Solution**: Lihat dokumentasi lengkap di `docs/FIX_KODE_USER_FORMAT.md`

### Issue: Database Connection Failed

**Problem**: `DATABASE_URL is required` atau connection timeout

**Solution**:
1. Pastikan `.env` file ada dan berisi `DATABASE_URL`
2. Cek koneksi internet
3. Verifikasi credentials Supabase
4. Test koneksi manual:
   ```bash
   psql "$DATABASE_URL" -c "SELECT 1"
   ```

### Issue: JWT Token Invalid

**Problem**: `Unauthorized` error meskipun sudah login

**Solution**:
1. Pastikan `JWT_SECRET` sama di `.env` dan saat generate token
2. Cek expiry token (default 24 jam)
3. Format header: `Authorization: Bearer <token>`

### Issue: Upload Surat Gagal

**Problem**: Error saat upload file PDF

**Solution**:
1. Cek `SUPABASE_SERVICE_KEY` di `.env`
2. Pastikan bucket `surat-digital` sudah dibuat
3. Cek policy bucket (authenticated write)
4. Verifikasi file size < `MAX_UPLOAD_SIZE_MB`

### Issue: Email Notification Tidak Terkirim

**Problem**: Email tidak terkirim setelah verifikasi

**Solution**:
1. Pastikan Gmail API sudah di-enable di Google Cloud Console
2. Cek `credentials.json` dan `gmail_token.json` sudah ada
3. Verifikasi `GMAIL_SENDER_EMAIL` di `.env`
4. Cek log server untuk error detail

### Issue: CORS Error

**Problem**: CORS error dari frontend

**Solution**:
1. Set `CORS_ALLOWED_ORIGIN` di `.env`
2. Untuk development: `CORS_ALLOWED_ORIGIN=*`
3. Untuk production: `CORS_ALLOWED_ORIGIN=https://yourdomain.com`

---

## License

MIT License - feel free to use this project for learning purposes.

---

## Contributors

- **Project**: Proyek 2 - Sistem Informasi Peminjaman Sarpras
- **Institution**: Politeknik Negeri Bandung

---

## Support

Jika ada pertanyaan atau issue:
1. Buka issue di repository
2. Lihat dokumentasi di folder `docs/`
