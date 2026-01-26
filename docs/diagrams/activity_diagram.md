# Activity Diagram - Sistem Peminjaman Sarana Prasarana

Activity diagram menggambarkan alur aktivitas dalam sistem secara detail.

---

## 1. Proses Pengajuan Peminjaman

```mermaid
flowchart TD
    subgraph Mahasiswa
        A1[Buka Halaman Peminjaman]
        A2[Pilih Ruangan]
        A3[Cek Ketersediaan Ruangan]
        A4{Ruangan Tersedia?}
        A5[Pilih Tanggal dan Waktu]
        A6[Isi Nama Kegiatan]
        A7[Pilih Barang Tambahan]
        A8[Upload Surat Digital]
        A9[Submit Pengajuan]
        A10[Terima Notifikasi Pending]
        A11[Pilih Ruangan/Tanggal Lain]
    end

    subgraph Sistem
        S1[Validasi Data Peminjaman]
        S2{Data Valid?}
        S3[Cek Bentrok Jadwal]
        S4{Tidak Bentrok?}
        S5[Generate Kode Peminjaman]
        S6[Simpan ke Database]
        S7[Buat Record Kegiatan]
        S8[Kirim Email ke Sarpras]
        S9[Update Status: PENDING]
        S10[Tampilkan Error Message]
    end

    A1 --> A2
    A2 --> A3
    A3 --> A4
    A4 -->|Ya| A5
    A4 -->|Tidak| A11
    A11 --> A2
    A5 --> A6
    A6 --> A7
    A7 --> A8
    A8 --> A9
    A9 --> S1
    S1 --> S2
    S2 -->|Tidak| S10
    S10 --> A9
    S2 -->|Ya| S3
    S3 --> S4
    S4 -->|Ya| S5
    S4 -->|Tidak| S10
    S5 --> S6
    S6 --> S7
    S7 --> S8
    S8 --> S9
    S9 --> A10
```

---

## 2. Proses Verifikasi oleh Sarpras

```mermaid
flowchart TD
    subgraph Sarpras
        B1[Login ke Sistem]
        B2[Buka Halaman Pending]
        B3[Pilih Peminjaman untuk Review]
        B4[Lihat Detail Peminjaman]
        B5[Download Surat Digital]
        B6{Keputusan?}
        B7[Isi Catatan Verifikasi]
        B8[Klik Approve]
        B9[Klik Reject]
    end

    subgraph Sistem
        S1[Update Status: APPROVED]
        S2[Update Status: REJECTED]
        S3[Set Verified By dan Verified At]
        S4[Simpan Catatan Verifikasi]
        S5[Kirim Email ke Mahasiswa]
        S6[Kirim Notifikasi ke Security]
        S7[Catat Log Aktivitas]
    end

    subgraph Mahasiswa
        M1[Terima Email Notifikasi]
    end

    B1 --> B2
    B2 --> B3
    B3 --> B4
    B4 --> B5
    B5 --> B6
    B6 -->|Approve| B7
    B6 -->|Reject| B7
    B7 --> B6
    B6 -->|Approve| B8
    B6 -->|Reject| B9
    B8 --> S1
    B9 --> S2
    S1 --> S3
    S2 --> S3
    S3 --> S4
    S4 --> S5
    S5 --> S7
    S1 --> S6
    S6 --> S7
    S5 --> M1
```

---

## 3. Proses Verifikasi Kehadiran oleh Security

```mermaid
flowchart TD
    subgraph Security
        C1[Login ke Sistem]
        C2[Buka Jadwal Aktif Hari Ini]
        C3[Pilih Peminjaman]
        C4[Verifikasi Identitas Peminjam]
        C5{Peminjam Hadir?}
        C6[Pilih Status: HADIR]
        C7[Pilih Status: TIDAK HADIR]
        C8[Pilih Status: BATAL]
        C9[Isi Keterangan]
        C10[Submit Kehadiran]
    end

    subgraph Sistem
        S1[Validasi Peminjaman APPROVED]
        S2[Simpan Record Kehadiran]
        S3{Status Kehadiran?}
        S4[Update Peminjaman: ONGOING]
        S5[Update Peminjaman: TIDAK_HADIR]
        S6[Update Peminjaman: BATAL]
        S7[Catat Log Aktivitas]
        S8[Kirim Notifikasi ke Sarpras]
    end

    C1 --> C2
    C2 --> C3
    C3 --> C4
    C4 --> C5
    C5 -->|Ya| C6
    C5 -->|Tidak| C7
    C5 -->|Dibatalkan| C8
    C6 --> C9
    C7 --> C9
    C8 --> C9
    C9 --> C10
    C10 --> S1
    S1 --> S2
    S2 --> S3
    S3 -->|HADIR| S4
    S3 -->|TIDAK_HADIR| S5
    S3 -->|BATAL| S6
    S4 --> S7
    S5 --> S7
    S6 --> S7
    S7 --> S8
```

---

## 4. Proses Pembatalan Peminjaman

```mermaid
flowchart TD
    subgraph Sarpras
        D1[Buka Detail Peminjaman]
        D2{Status APPROVED/ONGOING?}
        D3[Klik Batalkan Peminjaman]
        D4[Isi Alasan Pembatalan]
        D5[Konfirmasi Pembatalan]
        D6[Tidak Dapat Dibatalkan]
    end

    subgraph Sistem
        S1[Update Status: CANCELLED]
        S2[Set Canceller dan Timestamp]
        S3[Simpan Alasan Pembatalan]
        S4[Kirim Email ke Mahasiswa]
        S5[Kirim Email ke Security]
        S6[Catat Log Aktivitas]
    end

    subgraph Stakeholder
        M1[Mahasiswa Terima Notifikasi]
        SC1[Security Terima Notifikasi]
    end

    D1 --> D2
    D2 -->|Ya| D3
    D2 -->|Tidak| D6
    D3 --> D4
    D4 --> D5
    D5 --> S1
    S1 --> S2
    S2 --> S3
    S3 --> S4
    S4 --> S5
    S5 --> S6
    S4 --> M1
    S5 --> SC1
```

---

## 5. Proses Export Laporan

```mermaid
flowchart TD
    subgraph Sarpras
        E1[Buka Halaman Laporan]
        E2[Set Filter Tanggal]
        E3[Set Filter Status]
        E4[Klik Export Excel]
    end

    subgraph Sistem
        S1[Query Data Peminjaman]
        S2[Join Data Ruangan dan User]
        S3[Generate File Excel]
        S4[Set Headers dan Formatting]
        S5[Return File Download]
    end

    E1 --> E2
    E2 --> E3
    E3 --> E4
    E4 --> S1
    S1 --> S2
    S2 --> S3
    S3 --> S4
    S4 --> S5
    S5 --> E4
```

---

## 6. State Diagram - Status Peminjaman

```mermaid
stateDiagram-v2
    [*] --> PENDING : Mahasiswa Submit
    
    PENDING --> APPROVED : Sarpras Approve
    PENDING --> REJECTED : Sarpras Reject
    
    APPROVED --> ONGOING : Security: HADIR
    APPROVED --> CANCELLED : Sarpras Batalkan
    APPROVED --> TIDAK_HADIR : Security: TIDAK_HADIR
    APPROVED --> BATAL : Security: BATAL
    
    ONGOING --> FINISHED : Waktu Selesai
    ONGOING --> CANCELLED : Sarpras Batalkan
    
    REJECTED --> [*]
    FINISHED --> [*]
    CANCELLED --> [*]
    TIDAK_HADIR --> [*]
    BATAL --> [*]
```

---

## Ringkasan Aktivitas per Role

| Role | Aktivitas Utama |
|------|----------------|
| **Mahasiswa** | Mengajukan peminjaman, Upload surat, Melihat status |
| **Sarpras** | Verifikasi peminjaman, Kelola master data, Membatalkan peminjaman, Export laporan |
| **Security** | Verifikasi kehadiran, Melihat jadwal aktif |
| **Admin** | Semua aktivitas + Melihat log sistem |
