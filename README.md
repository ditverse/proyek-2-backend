# Sistem Informasi Peminjaman Sarana dan Prasarana Kampus

> Backend API untuk sistem peminjaman ruangan dan barang kampus dengan role-based access control

[![Go Version](https://img.shields.io/badge/Go-1.25.3-00ADD8?style=flat&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supabase-336791?style=flat&logo=postgresql)](https://supabase.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 📋 Daftar Isi

- [Tentang Project](#-tentang-project)
- [Fitur Utama](#-fitur-utama)
- [Teknologi](#-teknologi)
- [Struktur Project](#-struktur-project)
- [Setup & Installation](#-setup--installation)
- [API Documentation](#-api-documentation)
- [Database Schema](#-database-schema)
- [Role & Permissions](#-role--permissions)
- [Development](#-development)
- [Troubleshooting](#-troubleshooting)

---

## 🎯 Tentang Project

Sistem web-based untuk mengelola peminjaman ruangan dan barang di lingkungan kampus. Sistem ini mendukung workflow lengkap dari pengajuan peminjaman oleh mahasiswa, verifikasi oleh petugas sarpras, hingga pencatatan kehadiran oleh petugas security.

### Workflow Peminjaman

```
┌─────────────┐      ┌──────────────┐      ┌─────────────┐      ┌──────────────┐
│  MAHASISWA  │─────▶│   SARPRAS    │─────▶│  SECURITY   │─────▶│   FINISHED   │
│  Mengajukan │      │ Verifikasi   │      │  Verifikasi │      │   Selesai    │
│  Peminjaman │      │ APPROVED/    │      │  Kehadiran  │      │              │
│             │      │  REJECTED    │      │             │      │              │
└─────────────┘      └──────────────┘      └─────────────┘      └──────────────┘
    PENDING              APPROVED              ONGOING              FINISHED
```

---

## ✨ Fitur Utama

### 🔐 Authentication & Authorization
- JWT-based authentication
- Role-based access control (RBAC)
- 4 Role: MAHASISWA, SARPRAS, SECURITY, ADMIN

### 📦 Master Data Management
- CRUD Ruangan (kapasitas, lokasi, deskripsi)
- CRUD Barang (stok, lokasi penyimpanan)
- Manajemen Organisasi (HMJ, UKM, BEM, MPM)

### 📝 Peminjaman
- Pengajuan peminjaman ruangan/barang
- Upload surat digital ke Supabase Storage
- Multi-item peminjaman (ruangan + barang)
- Verifikasi oleh petugas sarpras
- Status tracking (PENDING → APPROVED → ONGOING → FINISHED)

### 👥 Kehadiran
- Verifikasi kehadiran oleh security
- Riwayat kehadiran peminjam
- Status kehadiran (HADIR, TIDAK_HADIR, BATAL)

### 🔔 Notifikasi
- Auto-notifikasi saat pengajuan dibuat
- Notifikasi status approved/rejected
- Reminder kehadiran
- Real-time notification count

### 📊 Laporan
- Laporan peminjaman (filter by date, status)
- Laporan kehadiran
- Log aktivitas sistem

---

## 🛠 Teknologi

### Backend
- **Language**: Go 1.25.3 (Native, no framework)
- **Database**: PostgreSQL (Supabase)
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **Database Driver**: pgx/v5

### Storage
- **File Storage**: Supabase Storage
- **Supported Files**: PDF (surat digital)

### Development Tools
- **Hot Reload**: Air
- **Environment**: godotenv

---

## 📁 Struktur Project

```
new-backend/
├── cmd/
│   └── server/
│       └── main.go              # Entry point aplikasi
├── internal/
│   ├── config/
│   │   ├── config.go            # Environment configuration
│   │   └── supabase.go          # Supabase storage config
│   ├── db/
│   │   └── db.go                # Database connection
│   ├── router/
│   │   └── router.go            # HTTP routing & middleware setup
│   └── services/
│       └── storage_service.go   # Supabase storage operations
├── models/                       # Domain models & DTOs
│   ├── user.go
│   ├── peminjaman.go
│   ├── ruangan.go
│   ├── barang.go
│   ├── kehadiran.go
│   ├── notifikasi.go
│   └── enums.go                 # Enum definitions
├── repositories/                 # Data access layer (CRUD)
│   ├── user_repository.go
│   ├── peminjaman_repository.go
│   ├── ruangan_repository.go
│   ├── barang_repository.go
│   ├── kehadiran_repository.go
│   └── notifikasi_repository.go
├── services/                     # Business logic layer
│   ├── auth_service.go
│   ├── peminjaman_service.go
│   ├── kehadiran_service.go
│   └── code_generator.go
├── handlers/                     # HTTP handlers (controllers)
│   ├── auth_handler.go
│   ├── peminjaman_handler.go
│   ├── ruangan_handler.go
│   ├── barang_handler.go
│   ├── kehadiran_handler.go
│   └── notifikasi_handler.go
├── middleware/                   # HTTP middleware
│   ├── auth.go                  # JWT validation & role checking
│   └── cors.go                  # CORS configuration
├── migrations/                   # SQL migration files
│   ├── 001_init_schema.sql
│   ├── erd_new_proyek_2.sql     # Current schema (with triggers)
│   └── 002_auto_generate_codes.sql
├── docs/                         # Documentation
│   └── FIX_KODE_USER_FORMAT.md
├── .air.toml                     # Air configuration
├── .env                          # Environment variables (gitignored)
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

---

## 🚀 Setup & Installation

### Prerequisites

- Go 1.25.3 or higher
- PostgreSQL database (Supabase account)
- Git

### 1. Clone Repository

```bash
git clone <repository-url>
cd new-backend
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

### 4. Setup Supabase Storage

1. Buka **Storage** di Supabase Dashboard
2. Buat bucket baru: `surat-digital`
3. Set policy untuk bucket (public read, authenticated write)

### 5. Environment Variables

Buat file `.env` di root project:

```env
# Database
DATABASE_URL=postgresql://postgres:[PASSWORD]@[HOST]:[PORT]/postgres

# Server
PORT=8000

# JWT Secret (GANTI DI PRODUCTION!)
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
```

**Cara mendapatkan credentials:**
- `DATABASE_URL`: Supabase Dashboard → Settings → Database → Connection String (URI)
- `SUPABASE_SERVICE_KEY`: Supabase Dashboard → Settings → API → service_role key

### 6. Run Server

#### Development (with hot reload)

```bash
# Install Air (jika belum)
go install github.com/cosmtrek/air@latest

# Run with Air
air
```

#### Production

```bash
go run cmd/server/main.go
```

Server akan berjalan di `http://localhost:8000`

### 7. Test API

```bash
# Health check
curl http://localhost:8000/api/health

# Register user
curl -X POST http://localhost:8000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "role": "MAHASISWA"
  }'

# Login
curl -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

---

## 📚 API Documentation

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

#### 🔓 Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/login` | Login user |
| POST | `/auth/register` | Register user baru |
| GET | `/health` | Health check |
| GET | `/info` | Info umum sistem |

#### 🔐 Protected Endpoints

##### Master Data - Ruangan

| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| GET | `/ruangan` | All | List semua ruangan |
| GET | `/ruangan/{id}` | All | Detail ruangan |
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
| POST | `/peminjaman/{id}/upload-surat` | Authenticated | Upload surat digital |
| GET | `/peminjaman/{id}/surat` | Authenticated | Get signed URL surat |
| GET | `/jadwal-ruangan` | All | Jadwal ruangan (calendar) |
| GET | `/jadwal-aktif` | SECURITY, ADMIN | Jadwal aktif untuk security |
| GET | `/jadwal-aktif-belum-verifikasi` | SECURITY, ADMIN | Jadwal belum verifikasi kehadiran |
| GET | `/laporan/peminjaman` | SARPRAS, ADMIN | Laporan peminjaman |

##### Kehadiran

| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| POST | `/kehadiran` | SECURITY, ADMIN | Catat kehadiran peminjam |
| GET | `/laporan/kehadiran` | SARPRAS, SECURITY, ADMIN | Laporan kehadiran |
| GET | `/kehadiran-riwayat` | SECURITY, ADMIN | Riwayat kehadiran by security |

##### Notifikasi

| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| GET | `/notifikasi/me` | Authenticated | List notifikasi user |
| GET | `/notifikasi/count` | Authenticated | Jumlah notifikasi belum dibaca |
| PATCH | `/notifikasi/{id}/dibaca` | Authenticated | Tandai notifikasi sebagai dibaca |

##### Log Aktivitas

| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| GET | `/log-aktivitas` | ADMIN | List semua log aktivitas |

### Request/Response Examples

#### Register User

**Request:**
```json
POST /api/auth/register
Content-Type: application/json

{
  "nama": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "role": "MAHASISWA",
  "organisasi_kode": "ORG-0001"
}
```

**Response:**
```json
{
  "kode_user": "USR-251204-0001",
  "nama": "John Doe",
  "email": "john@example.com",
  "role": "MAHASISWA",
  "organisasi_kode": "ORG-0001",
  "created_at": "2025-12-04T08:00:00Z"
}
```

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
  "tanggal_mulai": "2025-12-10T08:00:00Z",
  "tanggal_selesai": "2025-12-10T12:00:00Z",
  "keperluan": "Rapat Organisasi",
  "path_surat_digital": "surat/2025/12/surat-peminjaman.pdf",
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
  "tanggal_mulai": "2025-12-10T08:00:00Z",
  "tanggal_selesai": "2025-12-10T12:00:00Z",
  "keperluan": "Rapat Organisasi",
  "status": "PENDING",
  "path_surat_digital": "surat/2025/12/surat-peminjaman.pdf",
  "created_at": "2025-12-04T08:00:00Z"
}
```

---

## 🗄 Database Schema

### Tabel Utama

#### users
```sql
kode_user       VARCHAR PRIMARY KEY  -- Format: USR-YYMMDD-0001
nama            VARCHAR
email           VARCHAR UNIQUE
password_hash   VARCHAR
role            role_enum            -- MAHASISWA, SARPRAS, SECURITY, ADMIN
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
keperluan            TEXT
status               peminjaman_status_enum
path_surat_digital   TEXT
verified_by          VARCHAR FK
verified_at          TIMESTAMP
catatan_verifikasi   TEXT
created_at           TIMESTAMP
updated_at           TIMESTAMP
```

#### ruangan
```sql
kode_ruangan    VARCHAR PRIMARY KEY  -- Format: RNG-0001
nama_ruangan    VARCHAR
lokasi          VARCHAR
kapasitas       INT
deskripsi       TEXT
```

#### barang
```sql
kode_barang     VARCHAR PRIMARY KEY  -- Format: BRG-0001
nama_barang     VARCHAR
deskripsi       TEXT
jumlah_total    INT
ruangan_kode    VARCHAR FK
```

### Database Triggers

Sistem menggunakan database triggers untuk auto-generate kode dengan format yang konsisten:

- **`generate_kode_user()`**: Generate `USR-YYMMDD-0001`
- **`generate_kode_peminjaman()`**: Generate `PMJ-YYMMDD-0001`
- **`generate_kode_ruangan()`**: Generate `RNG-0001`
- **`generate_kode_barang()`**: Generate `BRG-0001`

Lihat `migrations/002_auto_generate_codes.sql` untuk detail implementasi.

---

## 👥 Role & Permissions

### MAHASISWA
- ✅ Melihat jadwal ruangan
- ✅ Mengajukan peminjaman
- ✅ Upload surat digital
- ✅ Melihat riwayat peminjaman sendiri
- ✅ Menerima notifikasi status peminjaman

### SARPRAS (Sarana Prasarana)
- ✅ Semua akses MAHASISWA
- ✅ Kelola master data (ruangan, barang)
- ✅ Verifikasi pengajuan peminjaman (approve/reject)
- ✅ Melihat laporan peminjaman
- ✅ Melihat semua pengajuan pending

### SECURITY
- ✅ Melihat jadwal peminjaman aktif
- ✅ Mencatat kehadiran peminjam
- ✅ Melihat riwayat kehadiran
- ✅ Melihat jadwal yang belum diverifikasi kehadirannya

### ADMIN
- ✅ **Semua akses** (full access)
- ✅ Melihat log aktivitas sistem

---

## 💻 Development

### Hot Reload dengan Air

Project ini mendukung Air untuk Windows dan Linux dengan konfigurasi terpisah:

| File | OS | Binary |
|------|-----|--------|
| `.air.windows.toml` | Windows | `tmp\main.exe` |
| `.air.linux.toml` | Linux/macOS | `./tmp/main` |

#### Install Air

```bash
# Windows/Linux/macOS
go install github.com/air-verse/air@latest
```

> **Note untuk Linux**: Pastikan `~/go/bin` sudah ada di PATH. Untuk Fish shell:
> ```bash
> fish_add_path ~/go/bin
> ```

#### Menjalankan Air

**Opsi 1: Menggunakan Script Wrapper (Rekomendasi)**
```bash
# Linux/macOS
./run-air.sh

# Windows
run-air.bat
```

**Opsi 2: Manually dengan config file**
```bash
# Linux/macOS
air -c .air.linux.toml

# Windows
air -c .air.windows.toml
```

**Opsi 3: Default config (sesuaikan dengan OS)**
```bash
air
```

### Code Structure Guidelines

#### Layered Architecture

```
Handler → Service → Repository → Database
```

- **Handlers**: HTTP request/response handling
- **Services**: Business logic
- **Repositories**: Database operations
- **Models**: Data structures

#### Naming Conventions

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

## 🐛 Troubleshooting

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

### Issue: CORS Error

**Problem**: CORS error dari frontend

**Solution**:
1. Set `CORS_ALLOWED_ORIGIN` di `.env`
2. Untuk development: `CORS_ALLOWED_ORIGIN=*`
3. Untuk production: `CORS_ALLOWED_ORIGIN=https://yourdomain.com`

---

## 📝 License

MIT License - feel free to use this project for learning purposes.

---

## 👨‍💻 Contributors

- **Developer**: [Your Name]
- **Project**: Proyek 2 - Sistem Informasi Peminjaman Sarpras

---

## 📞 Support

Jika ada pertanyaan atau issue:
1. Buka issue di repository
2. Lihat dokumentasi di folder `docs/`
3. Contact: [your-email@example.com]

---

**Happy Coding! 🚀**
