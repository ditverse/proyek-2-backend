# Ringkasan Fitur Sistem (Bab 3)

## 3.X Daftar Fitur Sistem

Sistem Informasi Peminjaman Sarana dan Prasarana Kampus dikembangkan dengan total **56 fitur** yang dikelompokkan dalam 7 modul utama sebagai berikut:

### Tabel 3.X Ringkasan Modul dan Fitur Sistem

| No | Modul | Jumlah Fitur | Deskripsi |
|----|-------|:------------:|-----------|
| 1 | **Autentikasi** | 5 | Manajemen akses pengguna meliputi login, registrasi, token JWT, enkripsi password, dan pengarahan berdasarkan role |
| 2 | **Peminjaman** | 8 | Proses inti peminjaman meliputi form pengajuan, upload surat digital, riwayat peminjaman, detail peminjaman, peminjaman barang, dan validasi konflik jadwal |
| 3 | **Jadwal & Dashboard** | 8 | Antarmuka utama setiap role meliputi 4 dashboard (Mahasiswa, Sarpras, Security, Tamu), jadwal ruangan, jadwal aktif harian, dan widget kalender interaktif |
| 4 | **Administrasi** | 8 | Fitur pengelolaan oleh admin meliputi verifikasi peminjaman, pembatalan, laporan, ekspor Excel, pengelolaan ruangan, pengelolaan barang, dan log aktivitas |
| 5 | **Kehadiran** | 3 | Pencatatan kehadiran oleh petugas keamanan meliputi verifikasi kehadiran, riwayat kehadiran, dan daftar jadwal belum diverifikasi |
| 6 | **Notifikasi Email** | 5 | Notifikasi otomatis via Gmail API meliputi email persetujuan, penolakan, pembatalan, notifikasi ke security, dan pengiriman asinkron |
| 7 | **Infrastruktur Teknis** | 15 | Komponen pendukung sistem meliputi middleware autentikasi, otorisasi role, CORS, penjadwal status otomatis, integrasi Supabase Storage, integrasi Gmail API, migrasi database, dan manajemen konfigurasi |
| | **Total** | **56** | |

---

### Penjelasan Setiap Modul

#### 1. Modul Autentikasi
Modul ini menangani proses autentikasi dan otorisasi pengguna. Sistem mendukung 4 role pengguna: Mahasiswa, Sarpras, Security, dan Admin. Setiap pengguna akan diarahkan ke dashboard yang sesuai dengan role-nya setelah berhasil login.

#### 2. Modul Peminjaman
Modul inti sistem yang memungkinkan mahasiswa untuk mengajukan peminjaman ruangan dan/atau barang. Setiap pengajuan wajib menyertakan surat digital dalam format PDF. Sistem secara otomatis melakukan validasi konflik jadwal untuk memastikan ruangan tidak double-booking.

#### 3. Modul Jadwal & Dashboard
Modul ini menyediakan antarmuka visual untuk setiap role pengguna. Terdapat widget kalender interaktif yang menampilkan jadwal peminjaman dengan indikator visual. Dashboard tamu dapat diakses tanpa login untuk melihat jadwal ruangan.

#### 4. Modul Administrasi
Modul khusus untuk admin Sarpras yang mencakup verifikasi pengajuan, pembatalan peminjaman, pembuatan laporan, dan pengelolaan master data (ruangan dan barang). Laporan dapat diekspor ke format Excel dengan styling profesional.

#### 5. Modul Kehadiran
Modul untuk petugas keamanan (Security) yang bertugas memverifikasi kehadiran peminjam pada saat kegiatan berlangsung. Status kehadiran dapat berupa Hadir, Tidak Hadir, atau Batal.

#### 6. Modul Notifikasi Email
Sistem mengirimkan notifikasi email secara otomatis menggunakan Gmail API. Email dikirim secara asinkron agar tidak menghambat respons sistem. Notifikasi dikirim saat pengajuan disetujui, ditolak, atau dibatalkan.

#### 7. Modul Infrastruktur Teknis
Komponen non-fungsional yang mendukung operasional sistem meliputi middleware keamanan, background job untuk transisi status otomatis (APPROVED → ONGOING → FINISHED), integrasi cloud storage, dan manajemen konfigurasi environment.

---

> **Catatan:** Daftar lengkap setiap fitur beserta deskripsi detailnya dapat dilihat pada **Lampiran A**.

---

## Lifecycle Status Peminjaman

```
┌─────────┐     Disetujui     ┌──────────┐     Waktu Mulai     ┌─────────┐     Waktu Selesai     ┌──────────┐
│ PENDING │ ─────────────────►│ APPROVED │ ──────────────────► │ ONGOING │ ─────────────────────► │ FINISHED │
└─────────┘                   └──────────┘                     └─────────┘                       └──────────┘
     │                              │                                │
     │ Ditolak                      │ Dibatalkan                     │ Dibatalkan
     ▼                              ▼                                ▼
┌──────────┐                  ┌───────────┐                   ┌───────────┐
│ REJECTED │                  │ CANCELLED │                   │ CANCELLED │
└──────────┘                  └───────────┘                   └───────────┘
```

Transisi status APPROVED → ONGOING → FINISHED dilakukan secara otomatis oleh sistem setiap 1 menit berdasarkan waktu mulai dan selesai peminjaman.
