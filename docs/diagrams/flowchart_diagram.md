# Flowchart Diagram - Sistem Peminjaman Sarana Prasarana Kampus

## Flowchart Keseluruhan Sistem

```mermaid
flowchart TD
    subgraph SISTEM["SISTEM PEMINJAMAN SARANA PRASARANA KAMPUS"]
        direction TB
        
        START([Mulai]) --> LOGIN[User Mengakses Sistem]
        LOGIN --> AUTH{Sudah Login?}
        
        AUTH -->|Belum| GUEST_ACTION{Aksi Guest}
        GUEST_ACTION -->|Login| DO_LOGIN[Proses Login]
        GUEST_ACTION -->|Register| DO_REGISTER[Proses Registrasi]
        GUEST_ACTION -->|Lihat Jadwal| VIEW_SCHEDULE[Lihat Jadwal Ruangan]
        VIEW_SCHEDULE --> AUTH
        DO_REGISTER --> DO_LOGIN
        DO_LOGIN --> CHECK_ROLE{Cek Role User}
        
        AUTH -->|Sudah| CHECK_ROLE
        
        CHECK_ROLE -->|Mahasiswa| MAHASISWA_FLOW
        CHECK_ROLE -->|Sarpras| SARPRAS_FLOW
        CHECK_ROLE -->|Security| SECURITY_FLOW
        
        subgraph MAHASISWA_FLOW["Alur Mahasiswa"]
            M_DASHBOARD[Dashboard Mahasiswa]
            M_AJUKAN[Ajukan Peminjaman]
            M_RIWAYAT[Lihat Riwayat]
            M_NOTIF[Cek Notifikasi]
            
            M_DASHBOARD --> M_AJUKAN
            M_DASHBOARD --> M_RIWAYAT
            M_DASHBOARD --> M_NOTIF
            
            M_AJUKAN --> M_FORM[Isi Form + Upload Surat]
            M_FORM --> M_SUBMIT[Submit Pengajuan]
            M_SUBMIT --> M_PENDING[Status: PENDING]
            M_PENDING --> NOTIF_SARPRAS[Notifikasi ke Sarpras]
        end
        
        subgraph SARPRAS_FLOW["Alur Sarpras"]
            S_DASHBOARD[Dashboard Sarpras]
            S_PENDING[Lihat Pengajuan Pending]
            S_VERIF[Verifikasi Peminjaman]
            S_MASTER[Kelola Master Data]
            S_LAPORAN[Export Laporan]
            
            S_DASHBOARD --> S_PENDING
            S_DASHBOARD --> S_MASTER
            S_DASHBOARD --> S_LAPORAN
            
            S_PENDING --> S_REVIEW[Review Pengajuan]
            S_REVIEW --> S_DECISION{Keputusan}
            S_DECISION -->|Approve| S_APPROVED[Status: APPROVED]
            S_DECISION -->|Reject| S_REJECTED[Status: REJECTED]
            S_APPROVED --> NOTIF_MHS[Notifikasi ke Mahasiswa]
            S_REJECTED --> NOTIF_MHS
            S_APPROVED --> NOTIF_SEC[Notifikasi ke Security]
        end
        
        subgraph SECURITY_FLOW["Alur Security"]
            SC_DASHBOARD[Dashboard Security]
            SC_JADWAL[Lihat Jadwal Aktif]
            SC_VERIF[Verifikasi Kehadiran]
            SC_RIWAYAT[Riwayat Kehadiran]
            
            SC_DASHBOARD --> SC_JADWAL
            SC_DASHBOARD --> SC_RIWAYAT
            
            SC_JADWAL --> SC_CHECK{Peminjam Hadir?}
            SC_CHECK -->|Ya| SC_HADIR[Status: HADIR / ONGOING]
            SC_CHECK -->|Tidak| SC_TIDAK[Status: TIDAK_HADIR]
            SC_CHECK -->|Batal| SC_BATAL[Status: BATAL]
            SC_HADIR --> SC_SELESAI[Kegiatan Selesai]
            SC_SELESAI --> SC_FINISHED[Status: FINISHED]
        end
        
        NOTIF_SEC --> SC_JADWAL
        
        M_NOTIF --> END_SESSION
        S_LAPORAN --> END_SESSION
        SC_FINISHED --> END_SESSION
        S_REJECTED --> END_SESSION
        SC_TIDAK --> END_SESSION
        SC_BATAL --> END_SESSION
        
        END_SESSION([Selesai / Logout])
    end
```

## Flowchart Proses Utama Peminjaman

```mermaid
flowchart TD
    subgraph Mahasiswa
        A[Mulai] --> B[Login ke Sistem]
        B --> C{Login Berhasil?}
        C -->|Tidak| B
        C -->|Ya| D[Akses Dashboard Mahasiswa]
        D --> E[Pilih Menu Pengajuan Peminjaman]
        E --> F[Isi Form Peminjaman]
        F --> G[Pilih Ruangan]
        G --> H[Pilih Tanggal dan Waktu]
        H --> I[Tambah Barang Opsional]
        I --> J[Upload Surat Digital PDF]
        J --> K[Submit Pengajuan]
        K --> L[Sistem Membuat Notifikasi ke Sarpras]
    end

    subgraph Sarpras
        L --> M[Sarpras Menerima Notifikasi]
        M --> N[Review Pengajuan Peminjaman]
        N --> O{Keputusan Verifikasi}
        O -->|Approve| P[Status: APPROVED]
        O -->|Reject| Q[Status: REJECTED]
        P --> R[Notifikasi ke Mahasiswa: Disetujui]
        Q --> S[Notifikasi ke Mahasiswa: Ditolak]
        S --> T[Selesai - Ditolak]
    end

    subgraph Security
        R --> U[Hari H Peminjaman]
        U --> V[Security Melihat Jadwal Aktif]
        V --> W[Peminjam Datang]
        W --> X{Verifikasi Kehadiran}
        X -->|Hadir| Y[Status: ONGOING]
        X -->|Tidak Hadir| Z[Status: TIDAK_HADIR]
        X -->|Batal| AA[Status: BATAL]
        Y --> AB[Kegiatan Berlangsung]
        AB --> AC[Kegiatan Selesai]
        AC --> AD[Status: FINISHED]
        AD --> AE[Selesai - Sukses]
        Z --> AF[Selesai - Tidak Hadir]
        AA --> AG[Selesai - Dibatalkan]
    end
```

## Flowchart Proses Registrasi User

```mermaid
flowchart TD
    A[Mulai] --> B[Akses Halaman Register]
    B --> C[Isi Data Registrasi]
    C --> D[Pilih Role: Mahasiswa]
    D --> E[Pilih Organisasi]
    E --> F[Submit Registrasi]
    F --> G{Validasi Data}
    G -->|Email Sudah Terdaftar| H[Tampilkan Error]
    H --> C
    G -->|Valid| I[Simpan User ke Database]
    I --> J[Generate Kode User: USR-YYMMDD-XXXX]
    J --> K[Redirect ke Halaman Login]
    K --> L[Selesai]
```

## Flowchart Proses Verifikasi Peminjaman oleh Sarpras

```mermaid
flowchart TD
    A[Mulai] --> B[Login sebagai Sarpras]
    B --> C[Akses Dashboard Sarpras]
    C --> D[Lihat Daftar Pengajuan Pending]
    D --> E{Ada Pengajuan Pending?}
    E -->|Tidak| F[Menunggu Pengajuan Baru]
    F --> D
    E -->|Ya| G[Pilih Pengajuan untuk Direview]
    G --> H[Lihat Detail Peminjaman]
    H --> I[Download dan Cek Surat Digital]
    I --> J[Cek Ketersediaan Ruangan]
    J --> K{Ruangan Tersedia?}
    K -->|Tidak| L[Reject dengan Catatan]
    K -->|Ya| M{Surat Valid?}
    M -->|Tidak| N[Reject dengan Catatan]
    M -->|Ya| O[Approve Peminjaman]
    L --> P[Update Status: REJECTED]
    N --> P
    O --> Q[Update Status: APPROVED]
    P --> R[Kirim Notifikasi ke Mahasiswa]
    Q --> R
    R --> S[Log Aktivitas Verifikasi]
    S --> T[Selesai]
```

## Flowchart Proses Verifikasi Kehadiran oleh Security

```mermaid
flowchart TD
    A[Mulai] --> B[Login sebagai Security]
    B --> C[Akses Dashboard Security]
    C --> D[Lihat Jadwal Aktif Hari Ini]
    D --> E{Ada Jadwal Aktif?}
    E -->|Tidak| F[Tidak Ada Kegiatan]
    F --> G[Selesai]
    E -->|Ya| H[Pilih Peminjaman untuk Verifikasi]
    H --> I[Peminjam Datang ke Lokasi]
    I --> J{Peminjam Hadir?}
    J -->|Ya| K[Catat Status: HADIR]
    J -->|Tidak Datang| L[Catat Status: TIDAK_HADIR]
    J -->|Dibatalkan| M[Catat Status: BATAL]
    K --> N[Update Status Peminjaman: ONGOING]
    L --> O[Update Status Peminjaman]
    M --> O
    N --> P[Kirim Notifikasi ke Sarpras]
    O --> P
    P --> Q[Log Aktivitas Kehadiran]
    Q --> R[Selesai]
```

## Flowchart Proses Export Laporan

```mermaid
flowchart TD
    A[Mulai] --> B[Login sebagai Sarpras]
    B --> C[Akses Menu Laporan Peminjaman]
    C --> D[Filter Data Laporan]
    D --> E[Pilih Rentang Tanggal]
    E --> F[Pilih Status Optional]
    F --> G[Klik Tombol Export]
    G --> H{Pilih Format}
    H -->|PDF| I[Generate PDF Report]
    H -->|Excel| J[Generate Excel Report]
    I --> K[Download File PDF]
    J --> L[Download File Excel]
    K --> M[Selesai]
    L --> M
```

## Flowchart Kelola Master Data Ruangan

```mermaid
flowchart TD
    A[Mulai] --> B[Login sebagai Sarpras]
    B --> C[Akses Menu Kelola Ruangan]
    C --> D[Lihat Daftar Ruangan]
    D --> E{Aksi yang Dipilih}
    E -->|Tambah| F[Isi Form Ruangan Baru]
    E -->|Edit| G[Pilih Ruangan]
    E -->|Hapus| H[Pilih Ruangan]
    F --> I[Input Nama, Lokasi, Kapasitas, Deskripsi]
    I --> J[Simpan Ruangan]
    J --> K[Generate Kode: RNG-XXXX]
    G --> L[Edit Data Ruangan]
    L --> M[Update Ruangan]
    H --> N{Konfirmasi Hapus?}
    N -->|Ya| O[Hapus Ruangan]
    N -->|Tidak| D
    K --> P[Tampilkan Pesan Sukses]
    M --> P
    O --> P
    P --> D
```
