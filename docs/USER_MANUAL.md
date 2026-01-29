# Panduan Pengguna (User Manual)
Sistem Informasi Peminjaman Sarana dan Prasarana Kampus

Dokumen ini berisi panduan lengkap penggunaan sistem untuk semua role pengguna: Mahasiswa, Staff Sarpras, dan Security.

## Daftar Isi
1. [Akses Sistem](#akses-sistem)
2. [Panduan Mahasiswa](#1-panduan-mahasiswa)
3. [Panduan Staff Sarpras](#2-panduan-staff-sarpras)
4. [Panduan Security](#3-panduan-security)
5. [Panduan Admin](#4-panduan-admin)

---

## Akses Sistem

Sistem dapat diakses melalui web browser (disarankan Google Chrome atau Microsoft Edge terbaru).
- **URL**: `http://localhost:8000` (atau domain produksi)
- **Halaman Login**: `index.html`

Masukan Email dan Password yang telah terdaftar untuk masuk ke dalam sistem. Role pengguna akan dideteksi otomatis saat login.

---

## 1. Panduan Mahasiswa

Role Mahasiswa memiliki akses untuk melihat jadwal dan mengajukan peminjaman.

### 1.1 Dashboard Mahasiswa
Setelah login, Anda akan diarahkan ke halaman Dashboard Mahasiswa (`dashboard-mahasiswa.html`).
- **Jadwal Hari Ini**: Menampilkan daftar peminjaman ruangan yang aktif pada hari ini.
- **Menu Navigasi**:
    - **Dashboard**: Kembali ke halaman utama.
    - **Jadwal Ruangan**: Melihat kalender ketersediaan ruangan.
    - **Riwayat Peminjaman**: Melihat status pengajuan peminjaman Anda.
    - **Profil**: Mengubah password (opsional).

### 1.2 Mengajukan Peminjaman
1. Klik menu **"Ajukan Peminjaman"** atau tombol **"+"** di dashboard.
2. Isi formulir **Formulir Peminjaman**:
    - **Nama Kegiatan**: Masukkan nama acara/kegiatan.
    - **Ruangan**: Pilih ruangan yang ingin dipinjam.
    - **Tanggal & Waktu**: Tentukan waktu mulai dan selesai. Pastikan tidak bentrok dengan jadwal lain.
    - **Barang (Opsional)**: Tambahkan barang jika diperlukan (misal: Proyektor, Sound System).
    - **Deskripsi**: Jelaskan detail kegiatan.
    - **Surat Pengantar**: Upload file PDF surat pengantar dari organisasi/jurusan (Dapat diupload nanti di detail).
3. Klik **Kirim Pengajuan**. Status awal adalah `PENDING`.

### 1.3 Upload Surat Digital
Jika Anda belum mengupload surat saat pengajuan:
1. Masuk ke menu **Riwayat Peminjaman**.
2. Pilih item peminjaman dengan status `PENDING` atau `APPROVED`.
3. Klik tombol **Upload Surat** pada halaman detail.
4. Pilih file PDF dan upload.

### 1.4 Memantau Status
Cek menu **Riwayat Peminjaman** secara berkala.
- **PENDING**: Sedang menunggu verifikasi Sarpras.
- **APPROVED**: Disetujui. Silakan gunakan ruangan sesuai jadwal.
- **REJECTED**: Ditolak. Alasan penolakan dapat dilihat di detail.
- **ONGOING**: Sedang berjalan (sudah check-in/waktu mulai tiba).
- **FINISHED**: Peminjaman selesai.

---

## 2. Panduan Staff Sarpras

Role Sarpras bertanggung jawab memverifikasi peminjaman dan mengelola aset.

### 2.1 Dashboard Sarpras
Halaman utama (`sarpras.html`) menampilkan ringkasan:
- Jumlah pengajuan PENDING.
- Jadwal ruangan hari ini.
- Statistik penggunaan aset.

### 2.2 Verifikasi Peminjaman
1. Buka menu **Verifikasi Peminjaman**.
2. Anda akan melihat daftar pengajuan dengan status `PENDING`.
3. Klik tombol **Detail** pada salah satu pengajuan.
4. Periksa detail kegiatan, jadwal, dan ketersediaan barang.
5. Periksa Surat Pengantar (Download/Preview).
6. Aksi:
    - **Terima (Approve)**: Peminjaman disetujui, jadwal akan terkunci, notifikasi dikirim ke mahasiswa.
    - **Tolak (Reject)**: Masukkan alasan penolakan, lalu kirim. Peminjaman dibatalkan.

### 2.3 Kelola Master Data
Gunakan menu sidebar untuk mengelola data:
- **Kelola Ruangan**: Tambah, Edit, atau Hapus data ruangan (Kapasitas, Fasilitas).
- **Kelola Barang**: Update stok barang, tambah jenis barang baru.
- **Registrasi User**: Mendaftarkan akun baru untuk Mahasiswa atau Staff lain.

### 2.4 Laporan
Menu **Laporan** memungkinkan Anda untuk:
- Melihat rekapitulasi peminjaman per periode.
- **Export to Excel**: Mengunduh data laporan untuk arsip fisik.

---

## 3. Panduan Security

Role Security bertugas memverifikasi kehadiran peminjam di lapangan.

### 3.1 Dashboard Security
Halaman utama (`dashboard-security.html`) fokus pada operasional harian:
- **Jadwal Aktif**: Peminjaman yang sedang berlangsung atau akan segera dimulai.
- **Status Kehadiran**: Indikator apakah peminjam sudah datang.

### 3.2 Verifikasi Kehadiran (Check-In)
1. Saat perwakilan mahasiswa datang untuk menggunakan ruangan, minta mereka menunjukkan Bukti Peminjaman (bisa dari HP mereka).
2. Di dashboard Security, cari jadwal yang sesuai.
3. Klik tombol **Verifikasi Hadir**.
4. Sistem akan mencatat waktu kehadiran dan mengubah status menjadi `ONGOING` (jika waktu sesuai) atau `HADIR`.

### 3.3 Penanganan Batal/Tidak Hadir
Jika peminjam tidak hadir hingga batas waktu toleransi:
1. Klik tombol **Tidak Hadir** atau **Batal** pada jadwal tersebut.
2. Ruangan akan tercatat kosong pada laporan.

---

## 4. Panduan Admin

Admin memiliki akses penuh (Superuser). Selain fitur Sarpras dan Security, Admin dapat:
- Melihat **Log Aktivitas** sistem secara keseluruhan (Audit Trail).
- Mengelola konfigurasi sistem (jika ada).
- Melakukan override status peminjaman darurat.

---

## Bantuan & Troubleshooting

- **Lupa Password**: Hubungi Admin atau Staff Sarpras untuk reset password.
- **Error Sistem**: Laporkan pesan error atau screenshot ke tim IT.
- **Notifikasi Email**: Pastikan email yang terdaftar aktif untuk menerima update status.
