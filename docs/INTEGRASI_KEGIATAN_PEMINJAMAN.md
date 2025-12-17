# Integrasi Auto-Create Kegiatan pada Peminjaman

**Tanggal:** 17 Desember 2024  
**Status:** ✅ Selesai  
**Last Updated:** 17 Desember 2024 (keperluan dihapus dari response)

## Ringkasan

Form peminjaman sekarang otomatis membuat record di tabel `kegiatan` saat submit. 

**Perubahan Utama:**
- Field `keperluan` ❌ **dihapus dari response JSON**
- Field `kode_kegiatan` ❌ dihapus dari request
- Diganti dengan `nama_kegiatan` + `deskripsi` di request
- Response include nested object `kegiatan` untuk data lengkap

---

## Perubahan API

### POST /api/peminjaman

#### Request Body Lama
```json
{
  "kode_ruangan": "RNG001",
  "kode_kegiatan": "KGT001",
  "keperluan": "Rapat BEM",
  "tanggal_mulai": "2024-12-20T08:00:00Z",
  "tanggal_selesai": "2024-12-20T12:00:00Z",
  "barang": []
}
```

#### Request Body Baru
```json
{
  "kode_ruangan": "RNG001",
  "nama_kegiatan": "Rapat Kerja BEM",
  "deskripsi": "Koordinasi pengurus semester ganjil",
  "tanggal_mulai": "2024-12-20T08:00:00Z",
  "tanggal_selesai": "2024-12-20T12:00:00Z",
  "barang": []
}
```

---

### GET /api/peminjaman/* (Response Update)

> **Update 17 Desember 2024:** Response sekarang include nested `kegiatan` object. Field `keperluan` **tidak lagi dikirim**.

#### Response Baru
```json
{
  "kode_peminjaman": "PJM001",
  "kode_ruangan": "RNG001",
  "kode_kegiatan": "KGT001",
  "kegiatan": {
    "kode_kegiatan": "KGT001",
    "nama_kegiatan": "Rapat Kerja BEM",
    "deskripsi": "Koordinasi pengurus semester ganjil",
    "tanggal_mulai": "2024-12-20T08:00:00Z",
    "tanggal_selesai": "2024-12-20T12:00:00Z",
    "organisasi_kode": "ORG001"
  },
  "tanggal_mulai": "2024-12-20T08:00:00Z",
  "tanggal_selesai": "2024-12-20T12:00:00Z",
  "status": "PENDING"
}
```

> ⚠️ **Perhatian:** Field `keperluan` **tidak ada** di response!

#### Endpoints yang Terpengaruh
| Endpoint | Deskripsi |
|----------|-----------|
| `GET /api/peminjaman/:id` | Detail peminjaman |
| `GET /api/peminjaman/me` | Riwayat peminjaman user |
| `GET /api/peminjaman/pending` | List pending untuk verifikasi |
| `GET /api/jadwal-aktif` | Jadwal aktif untuk security |
| `GET /api/jadwal-aktif-belum-verifikasi` | Jadwal belum verifikasi |
| `GET /api/laporan/peminjaman` | Laporan peminjaman |

#### Frontend Access
```javascript
// ✅ Cara baru (WAJIB)
const namaKegiatan = peminjaman.kegiatan?.nama_kegiatan || '';
const deskripsi = peminjaman.kegiatan?.deskripsi || '';

// ❌ Cara lama (TIDAK BEKERJA)
// peminjaman.keperluan  <-- field ini sudah tidak ada!
```

---

## Field Changes

| Field | Request | Response | Keterangan |
|-------|---------|----------|------------|
| `kode_kegiatan` | ❌ Dihapus | ✅ Ada | Auto-generated, tidak perlu input |
| `keperluan` | ❌ Dihapus | ❌ **Dihapus** | Gunakan `kegiatan.nama_kegiatan` |
| `nama_kegiatan` | ✅ Wajib | via `kegiatan` | Nama kegiatan yang akan dibuat |
| `deskripsi` | ✅ Opsional | via `kegiatan` | Deskripsi kegiatan |
| `kegiatan` | - | ✅ **Baru** | Nested object dengan data lengkap |

---

## File yang Diubah

### New Files
- `repositories/kegiatan_repository.go` - Repository untuk CRUD kegiatan

### Modified Files

| File | Perubahan |
|------|-----------|
| `models/peminjaman.go` | Update `CreatePeminjamanRequest` struct |
| `services/peminjaman_service.go` | Auto-create kegiatan sebelum peminjaman |
| `handlers/peminjaman_handler.go` | Tambah `KegiatanRepo` + enrichment di 6 GET handlers |
| `internal/router/router.go` | Inisialisasi dan inject `KegiatanRepository` |

---

## Flow Baru

```
POST /api/peminjaman {nama_kegiatan, deskripsi, ...}
        │
        ▼
┌─────────────────────────────────────┐
│ 1. Get user untuk organisasi_kode   │
├─────────────────────────────────────┤
│ 2. CREATE kegiatan (auto-generate)  │
├─────────────────────────────────────┤
│ 3. CREATE peminjaman (linked)       │
├─────────────────────────────────────┤
│ 4. CREATE peminjaman_barang (opt)   │
└─────────────────────────────────────┘
```

---

## Breaking Changes

> ⚠️ **Frontend harus diupdate!**

### Request (POST)
```javascript
// Sebelum
{ kode_kegiatan: "...", keperluan: "..." }

// Sesudah
{ nama_kegiatan: "...", deskripsi: "..." }
```

### Response (GET)
```javascript
// Sebelum
peminjaman.keperluan

// Sesudah
peminjaman.kegiatan?.nama_kegiatan
peminjaman.kegiatan?.deskripsi
```

---

## Backward Compatibility

- ✅ Field `keperluan` **sudah dihapus sepenuhnya** dari backend
- ✅ Kolom `keperluan` di database sudah di-DROP
- ✅ Model, repository, dan service sudah tidak menggunakan `keperluan`
- Frontend wajib akses `kegiatan.nama_kegiatan` sebagai gantinya
- Frontend harus akses `kegiatan.nama_kegiatan` sebagai gantinya


