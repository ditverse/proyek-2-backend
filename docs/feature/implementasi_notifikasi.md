Tentu, ini adalah **Project Requirement Document (PRD)** yang disusun secara teknis dan terstruktur. Dokumen ini dirancang khusus agar mudah dipahami oleh AI Agent (seperti coding assistant) untuk langsung mengimplementasikan kode.

Anda bisa menyalin teks di bawah ini dan memberikannya kepada AI Anda.

---

# Project Requirement Document (PRD)

**Fitur:** Sistem Notifikasi Multi-Channel (Email & WhatsApp)
**Aplikasi:** Peminjaman Sarana Prasarana Kampus (Sarpras)
**Tech Stack:** Golang (Backend), Supabase (PostgreSQL), Vanilla JS (Frontend)

## 1. Overview

Implementasi sistem notifikasi otomatis untuk tiga aktor (Mahasiswa, Sarpras, Security) menggunakan dua saluran komunikasi:

1. **Email (SMTP):** Untuk notifikasi formal (Bukti Approval, Rejection, Pengajuan Baru).
2. **WhatsApp (Fonnte API):** Untuk notifikasi urgensi tinggi (Reminder 1 Jam, Info Security).

Sistem harus berjalan secara *event-driven* (berdasarkan aksi user) dan *time-driven* (scheduler/cron job).

## 2. Database Schema Changes

Instruksi modifikasi database untuk mendukung fitur ini.

### A. Tabel `users`

Tambahkan kolom untuk menyimpan nomor WhatsApp.

```sql
ALTER TABLE users ADD COLUMN no_hp VARCHAR(20);

```

*Note: Pastikan format penyimpanan `no_hp` dinormalisasi (contoh: 628xxx) atau divalidasi di level aplikasi.*

### B. Update Model `User` (`models/user.go`)

Update struct user untuk mapping kolom baru.

```go
type User struct {
    // fields existing...
    NoHP string `json:"no_hp"`
}

```

### C. Enum Jenis Notifikasi (`models/enums.go`)

Tambahkan konstanta berikut untuk standarisasi jenis notifikasi.

```go
const (
    NotifJenisPengajuanMasuk      = "PENGAJUAN_MASUK"
    NotifJenisPeminjamanDisetujui = "PEMINJAMAN_DISETUJUI"
    NotifJenisPeminjamanDitolak   = "PEMINJAMAN_DITOLAK"
    NotifJenisReminder1Jam        = "REMINDER_1JAM"
)

```

## 3. Environment Variables (`.env`)

Tambahkan variabel berikut untuk konfigurasi layanan pihak ketiga.

```env
# Email Configuration (Google SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your_email@gmail.com
SMTP_PASSWORD=your_app_password

# WhatsApp Gateway (Fonnte)
FONNTE_TOKEN=your_fonnte_token
FONNTE_URL=https://api.fonnte.com/send

```

## 4. Functional Requirements & Implementation Specs

### Feature 1: Email Service (`services/email_service.go`)

**Library:** `gopkg.in/gomail.v2`
**Functionality:**

* Membuat struct `EmailService`.
* Method `SendEmail(to string, subject string, body string) error`.
* **Constraint:** Eksekusi pengiriman harus menggunakan Goroutine (`go func`) agar tidak memblokir *main thread*.

### Feature 2: WhatsApp Service (`services/whatsapp_service.go`)

**Provider:** Fonnte (HTTP POST Request)
**Functionality:**

* Membuat struct `WhatsappService`.
* Method `SendMessage(target string, message string) error`.
* Logic: Kirim HTTP POST ke `FONNTE_URL` dengan header `Authorization: FONNTE_TOKEN`.
* **Constraint:** Gunakan Goroutine. Handle validasi nomor (jika user input `08xx`, ubah otomatis jadi `628xx` sebelum kirim).

### Feature 3: Event-Driven Notifications (`services/peminjaman_service.go`)

Integrasikan service notifikasi ke dalam *business logic* peminjaman.

#### Scenario A: Create Peminjaman (Mahasiswa Submit)

* **Trigger:** `CreatePeminjaman` function.
* **Action:**
1. Ambil data semua user dengan role `SARPRAS`.
2. Kirim **Email** ke semua email Sarpras.
3. **Subject:** "Pengajuan Peminjaman Baru - [Nama Mahasiswa]"



#### Scenario B: Verifikasi - APPROVED (Sarpras Approve)

* **Trigger:** `VerifikasiPeminjaman` (ketika `status == APPROVED`).
* **Action 1 (Ke Mahasiswa):**
* Kirim **Email** "Peminjaman Disetujui" (Lampirkan detail surat/link).
* Kirim **WhatsApp** "Halo, pengajuan peminjaman ruangan [Nama Ruangan] telah DISETUJUI. Cek email untuk surat izin."


* **Action 2 (Ke Security):**
* Ambil data user role `SECURITY` (atau hardcode nomor pos satpam).
* Kirim **WhatsApp**: "👮 INFO: Kegiatan Baru Disetujui.\nKegiatan: [Nama]\nRuangan: [Nama]\nWaktu: [Mulai] - [Selesai]."



#### Scenario C: Verifikasi - REJECTED (Sarpras Reject)

* **Trigger:** `VerifikasiPeminjaman` (ketika `status == REJECTED`).
* **Action (Ke Mahasiswa):**
* Kirim **Email** "Peminjaman Ditolak".
* Body email wajib menyertakan `catatan_verifikasi` (alasan penolakan).



### Feature 4: Time-Driven Scheduler (`cmd/server/main.go` / `internal/scheduler`)

Sistem background worker untuk mengecek waktu peminjaman.

**Logic:**

1. Gunakan `time.NewTicker` (interval 5 menit).
2. **Query Database:** Cari peminjaman dengan kriteria:
* `status` = 'APPROVED'
* `tanggal_selesai` <= Waktu Sekarang + 1 Jam.
* `tanggal_selesai` > Waktu Sekarang.
* Belum pernah dikirim notifikasi tipe `REMINDER_1JAM` (Cek tabel `notifikasi` atau buat mekanisme flagging).


3. **Action:**
* Loop hasil query.
* Kirim **WhatsApp** ke Peminjam: "⚠️ REMINDER: Kegiatan di [Ruangan] berakhir dalam 1 jam. Harap rapikan ruangan sebelum meninggalkan lokasi."
* Simpan log ke tabel `notifikasi` agar reminder tidak dikirim berulang kali.



## 5. Message Templates (Copywriting)

Gunakan template string ini dalam kode:

**A. WhatsApp Template**

1. **Mahasiswa (Approved):**
> "✅ Status Update: DISETUJUI\nKegiatan: %s\nRuangan: %s\n\nSilakan cek email untuk surat izin digital."


2. **Security (Info):**
> "👮 MONITOR KEGIATAN\nJudul: %s\nLokasi: %s\nJam: %s s/d %s\n\nMohon dipantau."


3. **Reminder 1 Jam:**
> "⏳ REMINDER WAKTU\nSisa waktu peminjaman ruangan %s tinggal 1 jam lagi. Mohon persiapan untuk check-out."



**B. Email Template (Simple HTML)**

* Gunakan format HTML sederhana dengan kop "Sistem Sarpras".
* Sertakan variabel: Nama User, Nama Kegiatan, Tanggal, dan Status.

## 6. Acceptance Criteria

1. Service dapat mengirim email via Gmail SMTP tanpa error timeout.
2. Service dapat mengirim WhatsApp via Fonnte dan diterima di device target.
3. Proses pengiriman notifikasi tidak membuat loading website menjadi lama (Asynchronous).
4. Scheduler berjalan otomatis di background tanpa perlu trigger manual.
5. Tidak ada spam notifikasi (reminder hanya dikirim 1 kali per peminjaman).

---

**Instruksi Tambahan untuk AI:**

* Mohon buatkan file baru untuk `services/email_service.go` dan `services/whatsapp_service.go`.
* Inject kedua service tersebut ke dalam struct `PeminjamanService` di `services/peminjaman_service.go`.
* Implementasikan error handling: Jika notifikasi gagal, cukup log error-nya (`log.Println`), jangan batalkan transaksi database utama.