# Sprint Planning - Sistem Peminjaman Sarana Prasarana Kampus

> **Metodologi:** Scrum  
> **Total Sprint:** 4 Sprint  
> **Durasi per Sprint:** 2-3 Minggu  
> **Total Fitur:** 56 Fitur

---

## Ringkasan Sprint

| Sprint | Durasi | Fokus | Jumlah Fitur | Status |
|:------:|:------:|-------|:------------:|:------:|
| 1 | 2 Minggu | Fondasi & Autentikasi | 11 | ✅ Selesai |
| 2 | 3 Minggu | Core Peminjaman & Dashboard | 16 | ✅ Selesai |
| 3 | 2 Minggu | Administrasi & Kehadiran | 14 | ✅ Selesai |
| 4 | 2 Minggu | Notifikasi & Finalisasi | 15 | ✅ Selesai |

---

## Sprint 1: Fondasi & Autentikasi

**Durasi:** 2 Minggu  
**Sprint Goal:** *Pengguna dapat melakukan autentikasi dan sistem memiliki infrastruktur database yang siap*

### Sprint Backlog

| No | User Story / Task | Prioritas | Status |
|----|-------------------|:---------:|:------:|
| 1 | Membuat struktur database (migrasi) | Tinggi | ✅ |
| 2 | Membuat trigger auto-generate kode | Tinggi | ✅ |
| 3 | Membuat index untuk optimasi query | Sedang | ✅ |
| 4 | Membuat data awal (seed) | Sedang | ✅ |
| 5 | Implementasi fitur Login | Tinggi | ✅ |
| 6 | Implementasi fitur Register | Tinggi | ✅ |
| 7 | Implementasi manajemen token JWT | Tinggi | ✅ |
| 8 | Implementasi enkripsi password (bcrypt) | Tinggi | ✅ |
| 9 | Implementasi middleware autentikasi | Tinggi | ✅ |
| 10 | Implementasi otorisasi berbasis role | Tinggi | ✅ |
| 11 | Implementasi middleware CORS | Sedang | ✅ |

### Deliverables
- ✅ Database PostgreSQL dengan schema lengkap
- ✅ API Login & Register berfungsi
- ✅ Token JWT dapat di-generate dan divalidasi
- ✅ Middleware keamanan aktif

---

## Sprint 2: Core Peminjaman & Dashboard

**Durasi:** 3 Minggu  
**Sprint Goal:** *Mahasiswa dapat mengajukan peminjaman dan semua pengguna dapat melihat jadwal ruangan*

### Sprint Backlog

| No | User Story / Task | Prioritas | Status |
|----|-------------------|:---------:|:------:|
| 1 | Implementasi form pengajuan peminjaman | Tinggi | ✅ |
| 2 | Implementasi upload surat digital (PDF) | Tinggi | ✅ |
| 3 | Integrasi Supabase Storage | Tinggi | ✅ |
| 4 | Implementasi lihat/unduh surat digital | Sedang | ✅ |
| 5 | Implementasi validasi konflik jadwal | Tinggi | ✅ |
| 6 | Implementasi peminjaman barang | Sedang | ✅ |
| 7 | Implementasi auto-create kegiatan | Sedang | ✅ |
| 8 | Membuat halaman riwayat peminjaman | Sedang | ✅ |
| 9 | Membuat halaman detail peminjaman | Sedang | ✅ |
| 10 | Membuat Dashboard Mahasiswa | Tinggi | ✅ |
| 11 | Membuat Dashboard Tamu (Guest) | Sedang | ✅ |
| 12 | Membuat halaman Jadwal Ruangan | Tinggi | ✅ |
| 13 | Implementasi API jadwal aktif harian | Sedang | ✅ |
| 14 | Implementasi API booked dates | Sedang | ✅ |
| 15 | Membuat widget kalender interaktif | Sedang | ✅ |
| 16 | Implementasi redirect berdasarkan role | Sedang | ✅ |

### Deliverables
- ✅ Mahasiswa dapat submit pengajuan peminjaman
- ✅ Surat digital dapat diupload ke cloud storage
- ✅ Jadwal ruangan dapat dilihat oleh publik
- ✅ Kalender interaktif dengan indikator jadwal

---

## Sprint 3: Administrasi & Kehadiran

**Durasi:** 2 Minggu  
**Sprint Goal:** *Admin Sarpras dapat memverifikasi pengajuan dan Security dapat mencatat kehadiran*

### Sprint Backlog

| No | User Story / Task | Prioritas | Status |
|----|-------------------|:---------:|:------:|
| 1 | Implementasi verifikasi peminjaman (approve/reject) | Tinggi | ✅ |
| 2 | Membuat halaman daftar pengajuan pending | Tinggi | ✅ |
| 3 | Membuat Dashboard Sarpras | Tinggi | ✅ |
| 4 | Implementasi CRUD ruangan | Sedang | ✅ |
| 5 | Membuat halaman kelola ruangan | Sedang | ✅ |
| 6 | Implementasi CRUD barang | Sedang | ✅ |
| 7 | Membuat halaman kelola barang | Sedang | ✅ |
| 8 | Implementasi laporan peminjaman | Sedang | ✅ |
| 9 | Membuat halaman laporan peminjaman | Sedang | ✅ |
| 10 | Implementasi verifikasi kehadiran | Tinggi | ✅ |
| 11 | Membuat Dashboard Security | Tinggi | ✅ |
| 12 | Membuat halaman riwayat kehadiran | Sedang | ✅ |
| 13 | Implementasi API jadwal belum diverifikasi | Sedang | ✅ |
| 14 | Implementasi pencatatan log aktivitas | Sedang | ✅ |

### Deliverables
- ✅ Admin dapat approve/reject pengajuan
- ✅ Admin dapat mengelola master data (ruangan & barang)
- ✅ Security dapat verifikasi kehadiran peminjam
- ✅ Laporan peminjaman dapat dilihat

---

## Sprint 4: Notifikasi & Finalisasi

**Durasi:** 2 Minggu  
**Sprint Goal:** *Sistem mengirim notifikasi email otomatis dan semua fitur pendukung lengkap*

### Sprint Backlog

| No | User Story / Task | Prioritas | Status |
|----|-------------------|:---------:|:------:|
| 1 | Integrasi Gmail API (OAuth2) | Tinggi | ✅ |
| 2 | Implementasi email notifikasi persetujuan | Tinggi | ✅ |
| 3 | Implementasi email notifikasi penolakan | Tinggi | ✅ |
| 4 | Implementasi email notifikasi ke security | Sedang | ✅ |
| 5 | Implementasi pengiriman email asinkron | Sedang | ✅ |
| 6 | Implementasi penjadwal status otomatis (cron) | Tinggi | ✅ |
| 7 | Implementasi pembatalan peminjaman | Sedang | ✅ |
| 8 | Implementasi email notifikasi pembatalan | Sedang | ✅ |
| 9 | Implementasi ekspor laporan ke Excel | Sedang | ✅ |
| 10 | Membuat halaman detail laporan peminjaman | Sedang | ✅ |
| 11 | Implementasi endpoint info instansi | Rendah | ✅ |
| 12 | Implementasi daftar organisasi | Sedang | ✅ |
| 13 | Membuat halaman pendaftaran organisasi | Rendah | ✅ |
| 14 | Implementasi health check endpoint | Rendah | ✅ |
| 15 | Implementasi manajemen konfigurasi | Sedang | ✅ |

### Deliverables
- ✅ Email notifikasi terkirim otomatis saat verifikasi
- ✅ Status peminjaman berubah otomatis sesuai waktu
- ✅ Laporan dapat diekspor ke Excel
- ✅ Sistem siap untuk production

---

## Timeline Keseluruhan

```
Minggu 1-2     │ Sprint 1: Fondasi & Autentikasi
               │ ████████████████████████████████████████
               │
Minggu 3-5     │ Sprint 2: Core Peminjaman & Dashboard  
               │ ████████████████████████████████████████████████████████████
               │
Minggu 6-7     │ Sprint 3: Administrasi & Kehadiran
               │ ████████████████████████████████████████
               │
Minggu 8-9     │ Sprint 4: Notifikasi & Finalisasi
               │ ████████████████████████████████████████
```

**Total Durasi:** 9 Minggu

---

## Velocity Chart

| Sprint | Planned | Completed | Velocity |
|:------:|:-------:|:---------:|:--------:|
| 1 | 11 | 11 | 100% |
| 2 | 16 | 16 | 100% |
| 3 | 14 | 14 | 100% |
| 4 | 15 | 15 | 100% |

---

## Catatan Implementasi

1. **Sprint 1** menjadi fondasi karena semua fitur lain bergantung pada autentikasi dan database
2. **Sprint 2** memiliki durasi lebih panjang karena mencakup fitur inti dengan kompleksitas tinggi
3. **Sprint 3** dan **Sprint 4** dapat dikerjakan paralel jika tim lebih besar
4. Integrasi Gmail API memerlukan setup OAuth2 credentials terlebih dahulu
5. Status Scheduler berjalan sebagai background job setiap 1 menit
