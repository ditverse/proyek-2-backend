# Dokumentasi Proyek Sistem Informasi Peminjaman Sarana Prasarana

Selamat datang di pusat dokumentasi untuk **Sistem Informasi Peminjaman Sarana dan Prasarana Kampus**.
Dokumen ini disusun untuk membantu pengembang, penguji, dan pengguna akhir dalam memahami, menggunakan, dan mengembangkan sistem.

## 📚 Kategori Dokumentasi

### 1. Dokumentasi Utama (General)
- **[README Utama](../README.md)**: Gambaran umum proyek, fitur, setup, dan instalasi.
- **[User Manual (Panduan Pengguna)](USER_MANUAL.md)**: Panduan penggunaan aplikasi untuk role Mahasiswa, Sarpras, dan Security.

### 2. Dokumentasi Teknis (Technical)
Dokumentasi ini ditujukan untuk pengembang (Developers).

- **API Specification**: [OpenAPI / Swagger Definition](api/openapi.yaml)
- **Database Schema**: [ERD & Schema](../README.md#database-schema)
- **Architecture**:
  - [Architecture Overview](diagrams/architecture.md)
  - [Class Diagram](diagrams/class_diagram.md)
- **Business Process (BPMN)**:
  - [Business Process Overview](BPMN_BUSINESS_PROCESS.md)

### 3. Diagram Perancangan (Design Diagrams)
Kumpulan diagram UML dan Workflow visual.

- **[Diagrams Folder](diagrams/)**: Berisi Use Case, Activity, Sequence, dan Class diagrams.
- **[BPMN Folder](bpmn/)**: Diagram business process model and notation.

### 4. Laporan & Pengujian (Reports & Testing)
Dokumen hasil pengujian dan perencanaan fitur.

- **[Product Backlog](product_backlog.md)**: Daftar fitur dan user stories.
- **[Laporan Black Box Testing](LAPORAN_PENGUJIAN_BLACK_BOX.md)**: Hasil pengujian fungsionalitas.
- **[Laporan White Box Testing](LAPORAN_PENGUJIAN_WHITE_BOX.md)**: Hasil pengujian logika kode.
- **[Sprint Planning](sprint_planning.md)**: Perencanaan iterasi pengembangan.

### 5. Dokumentasi Fitur Spesifik
Penjelasan mendalam mengenai implementasi fitur tertentu.

- [Integrasi Kegiatan & Peminjaman](INTEGRASI_KEGIATAN_PEMINJAMAN.md)
- [Fix Kode User Format](FIX_KODE_USER_FORMAT.md)
- [Workflow Upload Surat](UPLOAD_SURAT_WORKFLOW.md)
- [Supabase Storage Fix](FIX_SUPABASE_STORAGE.md)

---

## Struktur Folder Dokumentasi

```
docs/
├── api/             # API Definitions (OpenAPI/Swagger)
├── bpmn/            # Business Process Diagrams
├── diagrams/        # UML Diagrams (Sequence, Class, Use Case)
├── testing/         # Testing scripts or additional logs
├── USER_MANUAL.md   # Panduan Pengguna Akhir
└── README.md        # File ini
```
