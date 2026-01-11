# State Diagram - Status Peminjaman

> State Machine menunjukkan transisi status peminjaman

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

## Ringkasan Status dan Transisi

| Status | Deskripsi | Transisi Berikutnya | Aktor |
|--------|-----------|---------------------|-------|
| `PENDING` | Pengajuan baru dibuat | APPROVED, REJECTED | SARPRAS |
| `APPROVED` | Disetujui, menunggu hari-H | ONGOING, CANCELLED | SECURITY/SARPRAS |
| `REJECTED` | Ditolak permanen | (terminal) | - |
| `ONGOING` | Kegiatan sedang berlangsung | FINISHED, CANCELLED | System/SARPRAS |
| `FINISHED` | Kegiatan telah selesai | (terminal) | - |
| `CANCELLED` | Dibatalkan | (terminal) | - |
