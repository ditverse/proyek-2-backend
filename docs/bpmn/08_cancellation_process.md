# Proses Pembatalan (Cancellation)

> Detail proses pembatalan peminjaman oleh SARPRAS/Admin

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

## Status yang Dapat Dibatalkan

| Status Awal | Dapat Dibatalkan? | Aktor |
|-------------|-------------------|-------|
| `PENDING` | ❌ Tidak | - |
| `APPROVED` | ✅ Ya | SARPRAS, ADMIN |
| `REJECTED` | ❌ Tidak | - |
| `ONGOING` | ✅ Ya | SARPRAS, ADMIN |
| `FINISHED` | ❌ Tidak | - |
| `CANCELLED` | ❌ Tidak | - |
