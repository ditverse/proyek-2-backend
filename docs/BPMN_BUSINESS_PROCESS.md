# Diagram BPMN - Sistem Peminjaman Sarana dan Prasarana Kampus

> **Visualisasi lengkap proses bisnis aplikasi peminjaman ruangan dan barang kampus**

---

## 1. Diagram Proses Utama (Main Business Process)

### 1.1 Alur Peminjaman End-to-End

Diagram ini menunjukkan alur lengkap dari pengajuan hingga selesai:

```mermaid
flowchart TD
    subgraph MAHASISWA["🎓 MAHASISWA/GUEST"]
        A[("Start")] --> B["Mengakses Sistem"]
        B --> C{"Sudah Login?"}
        C -->|Tidak| D["Login/Register"]
        D --> E["Memilih Organisasi"]
        E --> F["Dashboard Mahasiswa"]
        C -->|Ya| F
        F --> G["Membuat Pengajuan Peminjaman"]
        G --> H["Mengisi Form:<br/>- Ruangan<br/>- Tanggal & Waktu<br/>- Nama Kegiatan<br/>- Upload Surat"]
        H --> I["Submit Pengajuan"]
        I --> J["Status: PENDING"]
    end

    subgraph SARPRAS["🏢 SARPRAS"]
        J --> K["Notifikasi Pengajuan Baru"]
        K --> L["Review Pengajuan"]
        L --> M{"Verifikasi"}
        M -->|Approve| N["Status: APPROVED"]
        M -->|Reject| O["Status: REJECTED"]
        O --> P["Notifikasi Ditolak"]
        N --> Q["Notifikasi Disetujui"]
        Q --> R["Notifikasi ke Security"]
    end

    subgraph SECURITY["🔒 SECURITY"]
        R --> S["Melihat Jadwal Aktif"]
        S --> T["Menunggu Hari-H"]
        T --> U["Peminjam Hadir"]
        U --> V["Verifikasi Kehadiran"]
        V --> W{"Status Kehadiran?"}
        W -->|HADIR| X["Status: ONGOING"]
        W -->|TIDAK_HADIR| Y["Status: FINISHED<br/>(Tidak Hadir)"]
        W -->|BATAL| Z["Status: CANCELLED"]
    end

    subgraph COMPLETION["✅ PENYELESAIAN"]
        X --> AA["Kegiatan Berlangsung"]
        AA --> AB["Waktu Selesai"]
        AB --> AC["Status: FINISHED"]
        AC --> AD[("End")]
        Y --> AD
        Z --> AD
        P --> AE[("End - Ditolak")]
    end

    style A fill:#22c55e,color:#fff
    style AD fill:#22c55e,color:#fff
    style AE fill:#ef4444,color:#fff
    style J fill:#f59e0b,color:#fff
    style N fill:#3b82f6,color:#fff
    style O fill:#ef4444,color:#fff
    style X fill:#8b5cf6,color:#fff
    style AC fill:#22c55e,color:#fff
```

---

## 2. Diagram BPMN Swimlane (Per Aktor)

### 2.1 Collaboration Diagram

```mermaid
flowchart TB
    subgraph System["⚙️ SISTEM"]
        SYS1["Auto-generate Kode<br/>PMJ-YYMMDD-XXXX"]
        SYS2["Kirim Email & Notifikasi"]
        SYS3["Reminder 1 Jam Sebelum"]
        SYS4["Log Aktivitas"]
    end

    subgraph Mahasiswa["🎓 MAHASISWA"]
        M1["Buat Pengajuan"] --> M2["Upload Surat Digital"]
        M2 --> M3["Konfirmasi"]
        M3 --> SYS1
        SYS1 --> M4["Menunggu Verifikasi"]
        M4 --> M5["Terima Notifikasi"]
    end

    subgraph Sarpras["🏢 SARPRAS"]
        S1["Terima Notifikasi"] --> S2["Cek Ketersediaan"]
        S2 --> S3["Review Dokumen"]
        S3 --> S4{"Keputusan"}
        S4 -->|OK| S5["Approve"]
        S4 -->|Tolak| S6["Reject + Alasan"]
        S5 --> SYS2
        S6 --> SYS2
    end

    subgraph Security["🔒 SECURITY"]
        SC1["Lihat Jadwal Hari Ini"] --> SC2["Verifikasi Kedatangan"]
        SC2 --> SC3{"Hadir?"}
        SC3 -->|Ya| SC4["Update: HADIR"]
        SC3 -->|Tidak| SC5["Update: TIDAK_HADIR"]
        SC4 --> SYS4
        SC5 --> SYS4
    end

    SYS2 --> S1
    SYS2 --> M5
    SYS2 --> SC1
    SYS3 --> M5
```

---

## 3. State Diagram - Status Peminjaman

### 3.1 State Machine Peminjaman

```mermaid
stateDiagram-v2
    [*] --> PENDING: Pengajuan Dibuat

    PENDING --> APPROVED: Disetujui SARPRAS
    PENDING --> REJECTED: Ditolak SARPRAS
    
    APPROVED --> ONGOING: Kehadiran Terverifikasi
    APPROVED --> CANCELLED: Dibatalkan
    
    ONGOING --> FINISHED: Waktu Selesai
    ONGOING --> CANCELLED: Dibatalkan
    
    REJECTED --> [*]
    FINISHED --> [*]
    CANCELLED --> [*]

    note right of PENDING
        Status awal setelah
        mahasiswa submit
    end note

    note right of APPROVED
        Menunggu hari-H
        dan verifikasi kehadiran
    end note

    note right of ONGOING
        Kegiatan sedang
        berlangsung
    end note
```

---

## 4. Diagram Proses Detail

### 4.1 Proses Pengajuan Peminjaman (Create Peminjaman)

```mermaid
flowchart TD
    A[("Mulai")] --> B["User mengakses form pengajuan"]
    B --> C["Pilih Ruangan"]
    C --> D["Set Tanggal & Waktu"]
    D --> E{"Cek Konflik Jadwal"}
    E -->|Konflik| F["Tampilkan Error"]
    F --> D
    E -->|OK| G["Input Nama Kegiatan"]
    G --> H["Input Deskripsi"]
    H --> I["Upload Surat Digital"]
    I --> J{"File Valid?"}
    J -->|Tidak| K["Error: Format/Size"]
    K --> I
    J -->|Ya| L["Pilih Barang (Opsional)"]
    L --> M["Set Jumlah Barang"]
    M --> N{"Stok Tersedia?"}
    N -->|Tidak| O["Error: Stok Kurang"]
    O --> L
    N -->|Ya| P["Konfirmasi Data"]
    P --> Q["Submit"]
    Q --> R["Generate Kode PMJ"]
    R --> S["Simpan ke Database"]
    S --> T["Buat Notifikasi SARPRAS"]
    T --> U["Kirim Email SARPRAS"]
    U --> V["Tampilkan Sukses"]
    V --> W[("Selesai")]

    style A fill:#22c55e,color:#fff
    style W fill:#22c55e,color:#fff
    style F fill:#ef4444,color:#fff
    style K fill:#ef4444,color:#fff
    style O fill:#ef4444,color:#fff
```

### 4.2 Proses Verifikasi oleh SARPRAS

```mermaid
flowchart TD
    A[("Mulai")] --> B["SARPRAS menerima notifikasi"]
    B --> C["Buka halaman verifikasi"]
    C --> D["Filter: PENDING"]
    D --> E["Pilih pengajuan"]
    E --> F["Lihat detail pengajuan"]
    F --> G["Download surat digital"]
    G --> H["Cek ketersediaan ruangan"]
    H --> I["Cek dokumen kelengkapan"]
    I --> J{"Keputusan"}
    
    J -->|Setuju| K["Klik Approve"]
    K --> L["Input catatan (opsional)"]
    L --> M["Konfirmasi"]
    M --> N["Update Status: APPROVED"]
    N --> O["Notifikasi Mahasiswa"]
    O --> P["Notifikasi Security"]
    P --> Q["Kirim Email Mahasiswa"]
    Q --> R[("Selesai - Approved")]
    
    J -->|Tolak| S["Klik Reject"]
    S --> T["Input alasan penolakan"]
    T --> U["Konfirmasi"]
    U --> V["Update Status: REJECTED"]
    V --> W["Notifikasi Mahasiswa"]
    W --> X["Kirim Email Penolakan"]
    X --> Y[("Selesai - Rejected")]

    style A fill:#22c55e,color:#fff
    style R fill:#3b82f6,color:#fff
    style Y fill:#ef4444,color:#fff
```

### 4.3 Proses Verifikasi Kehadiran oleh SECURITY

```mermaid
flowchart TD
    A[("Mulai")] --> B["Security login"]
    B --> C["Akses Dashboard Security"]
    C --> D["Lihat jadwal hari ini"]
    D --> E["Filter: Belum Verifikasi"]
    E --> F["Pilih jadwal aktif"]
    F --> G["Lihat detail peminjaman"]
    G --> H{"Peminjam hadir?"}
    
    H -->|Ya, Hadir| I["Pilih status: HADIR"]
    I --> J["Input keterangan (opsional)"]
    J --> K["Konfirmasi"]
    K --> L["Update kehadiran"]
    L --> M["Update status: ONGOING"]
    M --> N["Log aktivitas"]
    N --> O[("Selesai - ONGOING")]
    
    H -->|Tidak Hadir| P["Pilih status: TIDAK_HADIR"]
    P --> Q["Input keterangan"]
    Q --> R["Konfirmasi"]
    R --> S["Update kehadiran"]
    S --> T["Log aktivitas"]
    T --> U[("Selesai - Tidak Hadir")]
    
    H -->|Batal| V["Pilih status: BATAL"]
    V --> W["Input alasan"]
    W --> X["Konfirmasi"]
    X --> Y["Update kehadiran"]
    Y --> Z["Log aktivitas"]
    Z --> AA[("Selesai - Batal")]

    style A fill:#22c55e,color:#fff
    style O fill:#8b5cf6,color:#fff
    style U fill:#f97316,color:#fff
    style AA fill:#ef4444,color:#fff
```

---

## 5. Proses Notifikasi & Reminder

### 5.1 Sistem Notifikasi

```mermaid
flowchart LR
    subgraph Trigger["🎯 Trigger Events"]
        T1["Pengajuan Dibuat"]
        T2["Status Approved"]
        T3["Status Rejected"]
        T4["Status Cancelled"]
        T5["1 Jam Sebelum Acara"]
        T6["Kehadiran Dicatat"]
    end

    subgraph Process["⚙️ Mailbox Service"]
        P1["Generate Notifikasi"]
        P2["Simpan ke Mailbox"]
        P3["Kirim Email (Async)"]
    end

    subgraph Recipients["👥 Penerima"]
        R1["Mahasiswa/Peminjam"]
        R2["SARPRAS/Admin"]
        R3["Security"]
    end

    T1 --> P1
    T2 --> P1
    T3 --> P1
    T4 --> P1
    T5 --> P1
    T6 --> P1

    P1 --> P2
    P2 --> P3

    P3 --> R1
    P3 --> R2
    P3 --> R3
```

### 5.2 Reminder Scheduler Flow

```mermaid
sequenceDiagram
    participant Scheduler as ⏰ Scheduler
    participant MailboxSvc as 📬 Mailbox Service
    participant DB as 🗄️ Database
    participant Email as 📧 Email Service
    participant User as 👤 User

    loop Every 5 Minutes
        Scheduler->>MailboxSvc: ProcessReminders()
        MailboxSvc->>DB: Get peminjaman T-1 jam
        DB-->>MailboxSvc: List peminjaman
        
        loop Each Peminjaman
            MailboxSvc->>DB: Create mailbox notification
            MailboxSvc->>Email: SendReminderEmail (async)
            Email-->>User: Email reminder
        end
    end
```

---

## 6. Proses Pembatalan (Cancellation)

```mermaid
flowchart TD
    A[("Mulai")] --> B{"Status Saat Ini?"}
    
    B -->|APPROVED| C["SARPRAS/Admin bisa batalkan"]
    B -->|ONGOING| D["SARPRAS/Admin bisa batalkan"]
    B -->|Lainnya| E["Tidak bisa dibatalkan"]
    
    C --> F["Input alasan pembatalan"]
    D --> F
    
    F --> G["Konfirmasi pembatalan"]
    G --> H["Update Status: CANCELLED"]
    H --> I["Notifikasi Peminjam"]
    I --> J["Notifikasi Security"]
    J --> K["Kirim Email Pembatalan"]
    K --> L["Log Aktivitas"]
    L --> M[("Selesai")]
    
    E --> N[("Error: Status tidak valid")]

    style A fill:#22c55e,color:#fff
    style M fill:#ef4444,color:#fff
    style N fill:#ef4444,color:#fff
```

---

## 7. Diagram Use Case

```mermaid
flowchart TD
    subgraph Actors["Aktor"]
        MHS["🎓 Mahasiswa/Guest"]
        SRP["🏢 SARPRAS"]
        SEC["🔒 Security"]
        ADM["👑 Admin"]
    end

    subgraph UseCases["Use Cases"]
        UC1["📝 Mengajukan Peminjaman"]
        UC2["📄 Upload Surat Digital"]
        UC3["📅 Melihat Jadwal Ruangan"]
        UC4["📋 Melihat Riwayat Peminjaman"]
        UC5["✅ Verifikasi Peminjaman"]
        UC6["❌ Menolak Peminjaman"]
        UC7["🚫 Membatalkan Peminjaman"]
        UC8["📦 Kelola Barang"]
        UC9["🏠 Kelola Ruangan"]
        UC10["👥 Verifikasi Kehadiran"]
        UC11["📊 Lihat Laporan Peminjaman"]
        UC12["🔔 Kelola Notifikasi"]
        UC13["📜 Lihat Log Aktivitas"]
    end

    MHS --> UC1
    MHS --> UC2
    MHS --> UC3
    MHS --> UC4
    MHS --> UC12

    SRP --> UC5
    SRP --> UC6
    SRP --> UC7
    SRP --> UC8
    SRP --> UC9
    SRP --> UC11
    SRP --> UC12

    SEC --> UC3
    SEC --> UC10
    SEC --> UC12

    ADM --> UC5
    ADM --> UC6
    ADM --> UC7
    ADM --> UC8
    ADM --> UC9
    ADM --> UC10
    ADM --> UC11
    ADM --> UC12
    ADM --> UC13
```

---

## 8. Ringkasan Status dan Transisi

| Status | Deskripsi | Transisi Berikutnya | Aktor |
|--------|-----------|---------------------|-------|
| `PENDING` | Pengajuan baru dibuat | APPROVED, REJECTED | SARPRAS |
| `APPROVED` | Disetujui, menunggu hari-H | ONGOING, CANCELLED | SECURITY/SARPRAS |
| `REJECTED` | Ditolak permanen | (terminal) | - |
| `ONGOING` | Kegiatan sedang berlangsung | FINISHED, CANCELLED | System/SARPRAS |
| `FINISHED` | Kegiatan telah selesai | (terminal) | - |
| `CANCELLED` | Dibatalkan | (terminal) | - |

---

## 9. Catatan Implementasi

> [!NOTE]
> Diagram-diagram di atas merepresentasikan proses bisnis yang diimplementasikan dalam:
> - [peminjaman_service.go](file:///home/versedroid/Documents/project-2-if/proyek-2-backend/services/peminjaman_service.go) - Logika peminjaman
> - [kehadiran_service.go](file:///home/versedroid/Documents/project-2-if/proyek-2-backend/services/kehadiran_service.go) - Verifikasi kehadiran
> - [mailbox_service.go](file:///home/versedroid/Documents/project-2-if/proyek-2-backend/services/mailbox_service.go) - Sistem notifikasi

> [!TIP]
> Untuk melihat diagram secara interaktif, Anda dapat:
> 1. Copy Mermaid code ke [Mermaid Live Editor](https://mermaid.live)
> 2. Menggunakan ekstensi VS Code "Markdown Preview Mermaid Support"
> 3. Export ke format gambar (PNG/SVG) dari Mermaid Live Editor
