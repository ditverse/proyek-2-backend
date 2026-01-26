# ERD (Entity Relationship Diagram) - Sistem Peminjaman Sarana Prasarana

Diagram ini menggambarkan struktur database dan relasi antar entitas dalam sistem.

## Diagram

```mermaid
erDiagram
    %% ==========================================
    %% ENTITIES
    %% ==========================================
    
    USERS {
        varchar kode_user PK "USR-YYMMDD-XXXX"
        varchar nama
        varchar email UK
        varchar password_hash
        enum role "MAHASISWA|SARPRAS|SECURITY|ADMIN"
        varchar no_hp
        varchar organisasi_kode FK
        timestamp created_at
    }

    ORGANISASI {
        varchar kode_organisasi PK "ORG-XXXX"
        varchar nama
        varchar kontak
        enum jenis_organisasi "ORMAWA|UKM"
        timestamp created_at
    }

    RUANGAN {
        varchar kode_ruangan PK "RNG-XXXX"
        varchar nama_ruangan
        varchar lokasi
        int kapasitas
        text deskripsi
    }

    BARANG {
        varchar kode_barang PK "BRG-XXXX"
        varchar nama_barang
        text deskripsi
        int jumlah_total
        varchar ruangan_kode FK
    }

    KEGIATAN {
        varchar kode_kegiatan PK "KGT-YYMMDD-XXXX"
        varchar nama_kegiatan
        text deskripsi
        timestamp tanggal_mulai
        timestamp tanggal_selesai
        varchar organisasi_kode FK
        timestamp created_at
        timestamp updated_at
    }

    PEMINJAMAN {
        varchar kode_peminjaman PK "PMJ-YYMMDD-XXXX"
        varchar kode_user FK
        varchar kode_ruangan FK
        varchar kode_kegiatan FK
        timestamp tanggal_mulai
        timestamp tanggal_selesai
        enum status "PENDING|APPROVED|REJECTED|ONGOING|FINISHED|CANCELLED"
        text path_surat_digital
        varchar verified_by FK
        timestamp verified_at
        text catatan_verifikasi
        timestamp created_at
        timestamp updated_at
    }

    PEMINJAMAN_BARANG {
        varchar kode_peminjaman_barang PK "PMB-XXXX"
        varchar kode_peminjaman FK
        varchar kode_barang FK
        int jumlah
    }

    KEHADIRAN_PEMINJAM {
        varchar kode_kehadiran PK "KHD-YYMMDD-XXXX"
        varchar kode_peminjaman FK
        enum status_kehadiran "HADIR|TIDAK_HADIR|BATAL"
        timestamp waktu_konfirmasi
        text keterangan
        varchar diverifikasi_oleh FK
        timestamp created_at
        timestamp updated_at
    }

    MAILBOX {
        varchar kode_mailbox PK "MBX-YYMMDD-XXXX"
        varchar kode_user FK
        varchar kode_peminjaman FK
        varchar jenis_pesan "APPROVED|REJECTED|CANCELLED|SECURITY_NOTIFY|NEW_SUBMISSION"
        timestamp created_at
    }

    LOG_AKTIVITAS {
        varchar kode_log PK "LOG-YYMMDD-XXXX"
        varchar kode_user FK
        varchar kode_peminjaman FK
        varchar aksi
        text keterangan
        timestamp waktu
        timestamp updated_at
    }

    %% ==========================================
    %% RELATIONSHIPS
    %% ==========================================

    USERS ||--o{ PEMINJAMAN : "mengajukan"
    USERS ||--o{ PEMINJAMAN : "memverifikasi"
    USERS ||--o{ KEHADIRAN_PEMINJAM : "mencatat"
    USERS ||--o{ LOG_AKTIVITAS : "melakukan"
    USERS ||--o{ MAILBOX : "menerima"
    USERS }o--|| ORGANISASI : "berasal dari"

    ORGANISASI ||--o{ KEGIATAN : "menyelenggarakan"

    RUANGAN ||--o{ PEMINJAMAN : "dipinjam"
    RUANGAN ||--o{ BARANG : "menyimpan"

    BARANG ||--o{ PEMINJAMAN_BARANG : "dipinjam"

    KEGIATAN ||--o{ PEMINJAMAN : "terkait"

    PEMINJAMAN ||--o{ PEMINJAMAN_BARANG : "memiliki"
    PEMINJAMAN ||--o| KEHADIRAN_PEMINJAM : "dicatat kehadirannya"
    PEMINJAMAN ||--o{ LOG_AKTIVITAS : "tercatat"
    PEMINJAMAN ||--o{ MAILBOX : "dikirim notifikasi"
```

## Deskripsi Entitas

### Master Data

| Entitas | Deskripsi |
|---------|-----------|
| **USERS** | Data pengguna sistem (Mahasiswa, Sarpras, Security, Admin) |
| **ORGANISASI** | Organisasi kemahasiswaan (ORMAWA, UKM) |
| **RUANGAN** | Ruangan yang tersedia untuk dipinjam |
| **BARANG** | Barang inventaris yang dapat dipinjam |

### Transaksi

| Entitas | Deskripsi |
|---------|-----------|
| **KEGIATAN** | Kegiatan yang terkait dengan peminjaman |
| **PEMINJAMAN** | Transaksi peminjaman ruangan/barang |
| **PEMINJAMAN_BARANG** | Detail barang yang dipinjam (relasi many-to-many) |
| **KEHADIRAN_PEMINJAM** | Catatan kehadiran peminjam pada hari H |

### Sistem

| Entitas | Deskripsi |
|---------|-----------|
| **MAILBOX** | Log notifikasi email yang terkirim |
| **LOG_AKTIVITAS** | Audit trail aktivitas pengguna |

## Kardinalitas Relasi

| Relasi | Kardinalitas | Keterangan |
|--------|--------------|------------|
| USERS - PEMINJAMAN | 1:N | Satu user dapat memiliki banyak peminjaman |
| USERS - ORGANISASI | N:1 | Banyak user dapat berasal dari satu organisasi |
| PEMINJAMAN - RUANGAN | N:1 | Banyak peminjaman dapat menggunakan satu ruangan |
| PEMINJAMAN - PEMINJAMAN_BARANG | 1:N | Satu peminjaman dapat memiliki banyak item barang |
| PEMINJAMAN_BARANG - BARANG | N:1 | Banyak item peminjaman dapat merujuk ke satu barang |
| PEMINJAMAN - KEHADIRAN_PEMINJAM | 1:0..1 | Satu peminjaman memiliki maksimal satu catatan kehadiran |
| PEMINJAMAN - KEGIATAN | N:1 | Banyak peminjaman dapat terkait dengan satu kegiatan |
| KEGIATAN - ORGANISASI | N:1 | Banyak kegiatan dapat diselenggarakan oleh satu organisasi |

## Format Kode

| Entitas | Format | Contoh |
|---------|--------|--------|
| Users | USR-YYMMDD-XXXX | USR-260126-0001 |
| Peminjaman | PMJ-YYMMDD-XXXX | PMJ-260126-0001 |
| Kegiatan | KGT-YYMMDD-XXXX | KGT-260126-0001 |
| Kehadiran | KHD-YYMMDD-XXXX | KHD-260126-0001 |
| Mailbox | MBX-YYMMDD-XXXX | MBX-260126-0001 |
| Log | LOG-YYMMDD-XXXX | LOG-260126-0001 |
| Ruangan | RNG-XXXX | RNG-0001 |
| Barang | BRG-XXXX | BRG-0001 |
| Organisasi | ORG-XXXX | ORG-0001 |
