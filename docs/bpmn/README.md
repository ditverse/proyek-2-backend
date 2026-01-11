# Dokumentasi BPMN - Sistem Peminjaman Sarpras

> Kumpulan diagram BPMN untuk memvisualisasikan proses bisnis aplikasi

## 📁 Daftar Diagram

| No | File | Deskripsi |
|----|------|-----------|
| 1 | [01_main_workflow.md](./01_main_workflow.md) | Alur peminjaman end-to-end |
| 2 | [02_swimlane_collaboration.md](./02_swimlane_collaboration.md) | Kolaborasi antar aktor |
| 3 | [03_state_machine.md](./03_state_machine.md) | State machine status peminjaman |
| 4 | [04_pengajuan_peminjaman.md](./04_pengajuan_peminjaman.md) | Detail proses pengajuan |
| 5 | [05_verifikasi_sarpras.md](./05_verifikasi_sarpras.md) | Proses verifikasi SARPRAS |
| 6 | [06_verifikasi_kehadiran.md](./06_verifikasi_kehadiran.md) | Proses verifikasi kehadiran |
| 7 | [07_notification_system.md](./07_notification_system.md) | Sistem notifikasi & reminder |
| 8 | [08_cancellation_process.md](./08_cancellation_process.md) | Proses pembatalan |
| 9 | [09_use_case.md](./09_use_case.md) | Use case diagram |

## 🛠 Cara Menggunakan

### Melihat Diagram
1. **VS Code**: Install ekstensi "Markdown Preview Mermaid Support"
2. **Online**: Copy kode Mermaid ke [Mermaid Live Editor](https://mermaid.live)

### Export ke Gambar
1. Buka [Mermaid Live Editor](https://mermaid.live)
2. Paste kode Mermaid
3. Klik tombol Download (PNG/SVG)

## 👥 Aktor Utama

| Aktor | Emoji | Deskripsi |
|-------|-------|-----------|
| MAHASISWA | 🎓 | Peminjam ruangan/barang |
| SARPRAS | 🏢 | Verifikator peminjaman |
| SECURITY | 🔒 | Verifikator kehadiran |
| ADMIN | 👑 | Full access |

## 📊 Status Peminjaman

```
PENDING → APPROVED → ONGOING → FINISHED
    ↓         ↓         ↓
 REJECTED  CANCELLED  CANCELLED
```
