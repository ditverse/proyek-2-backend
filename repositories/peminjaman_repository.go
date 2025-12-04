package repositories

import (
	"backend-sarpras/models"
	"database/sql"
	"time"
)

type PeminjamanRepository struct {
	DB *sql.DB
}

func NewPeminjamanRepository(db *sql.DB) *PeminjamanRepository {
	return &PeminjamanRepository{DB: db}
}

func (r *PeminjamanRepository) Create(peminjaman *models.Peminjaman) error {
	query := `
		INSERT INTO peminjaman (
			kode_user, kode_ruangan, kode_kegiatan,
			tanggal_mulai, tanggal_selesai, keperluan, status, path_surat_digital
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING kode_peminjaman, created_at
	`
	var kodeRuangan interface{}
	if peminjaman.KodeRuangan != nil {
		kodeRuangan = *peminjaman.KodeRuangan
	}
	var kodeKegiatan interface{}
	if peminjaman.KodeKegiatan != nil {
		kodeKegiatan = *peminjaman.KodeKegiatan
	}
	return r.DB.QueryRow(
		query,
		peminjaman.KodeUser,
		kodeRuangan,
		kodeKegiatan,
		peminjaman.TanggalMulai,
		peminjaman.TanggalSelesai,
		peminjaman.Keperluan,
		peminjaman.Status,
		peminjaman.PathSuratDigital,
	).Scan(&peminjaman.KodePeminjaman, &peminjaman.CreatedAt)
}

func (r *PeminjamanRepository) CreatePeminjamanBarang(kodePeminjamanBarang, kodePeminjaman, kodeBarang string, jumlah int) error {
	query := `
		INSERT INTO peminjaman_barang (kode_peminjaman_barang, kode_peminjaman, kode_barang, jumlah)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.DB.Exec(query, kodePeminjamanBarang, kodePeminjaman, kodeBarang, jumlah)
	return err
}

func (r *PeminjamanRepository) GetByID(kodePeminjaman string) (*models.Peminjaman, error) {
	query := `
		SELECT kode_peminjaman, kode_user, kode_ruangan, kode_kegiatan, tanggal_mulai, tanggal_selesai,
		       keperluan, status, path_surat_digital, verified_by, verified_at,
		       catatan_verifikasi, created_at, updated_at
		FROM peminjaman
		WHERE kode_peminjaman = $1
	`
	var (
		kodeRuangan  sql.NullString
		kodeKegiatan sql.NullString
		verifiedBy   sql.NullString
		verifiedAt   sql.NullTime
		catatan      sql.NullString
		updatedAt    sql.NullTime
	)

	p := &models.Peminjaman{}
	err := r.DB.QueryRow(query, kodePeminjaman).Scan(
		&p.KodePeminjaman,
		&p.KodeUser,
		&kodeRuangan,
		&kodeKegiatan,
		&p.TanggalMulai,
		&p.TanggalSelesai,
		&p.Keperluan,
		&p.Status,
		&p.PathSuratDigital,
		&verifiedBy,
		&verifiedAt,
		&catatan,
		&p.CreatedAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if kodeRuangan.Valid {
		p.KodeRuangan = &kodeRuangan.String
	}
	if kodeKegiatan.Valid {
		p.KodeKegiatan = &kodeKegiatan.String
	}
	if verifiedBy.Valid {
		p.VerifiedBy = &verifiedBy.String
	}
	if verifiedAt.Valid {
		t := verifiedAt.Time
		p.VerifiedAt = &t
	}
	if catatan.Valid {
		p.CatatanVerifikasi = catatan.String
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		p.UpdatedAt = &t
	}
	return p, nil
}

func (r *PeminjamanRepository) scanRows(rows *sql.Rows) ([]models.Peminjaman, error) {
	var result []models.Peminjaman
	for rows.Next() {
		var (
			row          models.Peminjaman
			kodeRuangan  sql.NullString
			kodeKegiatan sql.NullString
			verifiedBy   sql.NullString
			verifiedAt   sql.NullTime
			catatan      sql.NullString
			updatedAt    sql.NullTime
		)
		err := rows.Scan(
			&row.KodePeminjaman,
			&row.KodeUser,
			&kodeRuangan,
			&kodeKegiatan,
			&row.TanggalMulai,
			&row.TanggalSelesai,
			&row.Keperluan,
			&row.Status,
			&row.PathSuratDigital,
			&verifiedBy,
			&verifiedAt,
			&catatan,
			&row.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		if kodeRuangan.Valid {
			row.KodeRuangan = &kodeRuangan.String
		}
		if kodeKegiatan.Valid {
			row.KodeKegiatan = &kodeKegiatan.String
		}
		if verifiedBy.Valid {
			row.VerifiedBy = &verifiedBy.String
		}
		if verifiedAt.Valid {
			t := verifiedAt.Time
			row.VerifiedAt = &t
		}
		if catatan.Valid {
			row.CatatanVerifikasi = catatan.String
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			row.UpdatedAt = &t
		}
		result = append(result, row)
	}
	return result, nil
}

func (r *PeminjamanRepository) GetByPeminjamID(kodeUser string) ([]models.Peminjaman, error) {
	query := `
		SELECT p.kode_peminjaman, p.kode_user, p.kode_ruangan, p.kode_kegiatan, p.tanggal_mulai, p.tanggal_selesai,
		       p.keperluan, p.status, p.path_surat_digital, p.verified_by, p.verified_at,
		       p.catatan_verifikasi, p.created_at, p.updated_at,
		       r.kode_ruangan, r.nama_ruangan, r.lokasi, r.kapasitas, r.deskripsi
		FROM peminjaman p
		LEFT JOIN ruangan r ON p.kode_ruangan = r.kode_ruangan
		WHERE p.kode_user = $1
		ORDER BY p.created_at DESC
	`
	rows, err := r.DB.Query(query, kodeUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Peminjaman
	for rows.Next() {
		var (
			row          models.Peminjaman
			kodeRuangan  sql.NullString
			kodeKegiatan sql.NullString
			verifiedBy   sql.NullString
			verifiedAt   sql.NullTime
			catatan      sql.NullString
			updatedAt    sql.NullTime

			// Ruangan fields
			rKode      sql.NullString
			rNama      sql.NullString
			rLokasi    sql.NullString
			rKapasitas sql.NullInt32
			rDeskripsi sql.NullString
		)
		err := rows.Scan(
			&row.KodePeminjaman,
			&row.KodeUser,
			&kodeRuangan,
			&kodeKegiatan,
			&row.TanggalMulai,
			&row.TanggalSelesai,
			&row.Keperluan,
			&row.Status,
			&row.PathSuratDigital,
			&verifiedBy,
			&verifiedAt,
			&catatan,
			&row.CreatedAt,
			&updatedAt,
			&rKode,
			&rNama,
			&rLokasi,
			&rKapasitas,
			&rDeskripsi,
		)
		if err != nil {
			return nil, err
		}

		if kodeRuangan.Valid {
			row.KodeRuangan = &kodeRuangan.String
		}
		if kodeKegiatan.Valid {
			row.KodeKegiatan = &kodeKegiatan.String
		}
		if verifiedBy.Valid {
			row.VerifiedBy = &verifiedBy.String
		}
		if verifiedAt.Valid {
			t := verifiedAt.Time
			row.VerifiedAt = &t
		}
		if catatan.Valid {
			row.CatatanVerifikasi = catatan.String
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			row.UpdatedAt = &t
		}

		// Populate Ruangan struct if data exists
		if rKode.Valid {
			row.Ruangan = &models.Ruangan{
				KodeRuangan: rKode.String,
				NamaRuangan: rNama.String,
				Lokasi:      rLokasi.String,
				Kapasitas:   int(rKapasitas.Int32),
				Deskripsi:   rDeskripsi.String,
			}
		}

		result = append(result, row)
	}
	return result, nil
}

func (r *PeminjamanRepository) GetPending() ([]models.Peminjaman, error) {
	query := `
		SELECT kode_peminjaman, kode_user, kode_ruangan, kode_kegiatan, tanggal_mulai, tanggal_selesai,
		       keperluan, status, path_surat_digital, verified_by, verified_at,
		       catatan_verifikasi, created_at, updated_at
		FROM peminjaman
		WHERE status = 'PENDING'
		ORDER BY created_at ASC
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *PeminjamanRepository) GetJadwalRuangan(start, end time.Time) ([]models.JadwalRuanganResponse, error) {
	query := `
		SELECT p.kode_peminjaman, p.kode_ruangan, r.nama_ruangan, p.tanggal_mulai, p.tanggal_selesai,
		       p.status, u.nama
		FROM peminjaman p
		JOIN ruangan r ON p.kode_ruangan = r.kode_ruangan
		JOIN users u ON p.kode_user = u.kode_user
		WHERE p.kode_ruangan IS NOT NULL
		  AND p.status = 'APPROVED'
		  AND p.tanggal_mulai <= $2
		  AND p.tanggal_selesai >= $1
		ORDER BY p.tanggal_mulai
	`
	rows, err := r.DB.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jadwal []models.JadwalRuanganResponse
	for rows.Next() {
		var j models.JadwalRuanganResponse
		err := rows.Scan(
			&j.KodePeminjaman,
			&j.KodeRuangan,
			&j.NamaRuangan,
			&j.TanggalMulai,
			&j.TanggalSelesai,
			&j.Status,
			&j.Peminjam,
		)
		if err != nil {
			return nil, err
		}
		jadwal = append(jadwal, j)
	}
	if jadwal == nil {
		jadwal = []models.JadwalRuanganResponse{}
	}
	return jadwal, nil
}

func (r *PeminjamanRepository) GetJadwalAktif(start, end time.Time) ([]models.Peminjaman, error) {
	query := `
		SELECT kode_peminjaman, kode_user, kode_ruangan, kode_kegiatan, tanggal_mulai, tanggal_selesai,
		       keperluan, status, path_surat_digital, verified_by, verified_at,
		       catatan_verifikasi, created_at, updated_at
		FROM peminjaman
		WHERE status = 'APPROVED'
		  AND kode_ruangan IS NOT NULL
		  AND tanggal_mulai >= $1
		  AND tanggal_selesai <= $2
		ORDER BY tanggal_mulai
	`
	rows, err := r.DB.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *PeminjamanRepository) GetJadwalAktifBelumVerifikasi(start, end time.Time) ([]models.Peminjaman, error) {
	query := `
		SELECT p.kode_peminjaman, p.kode_user, p.kode_ruangan, p.kode_kegiatan, p.tanggal_mulai, p.tanggal_selesai,
		       p.keperluan, p.status, p.path_surat_digital, p.verified_by, p.verified_at,
		       p.catatan_verifikasi, p.created_at, p.updated_at
		FROM peminjaman p
		LEFT JOIN kehadiran_peminjam k ON k.kode_peminjaman = p.kode_peminjaman
		WHERE p.status = 'APPROVED'
		  AND p.kode_ruangan IS NOT NULL
		  AND p.tanggal_mulai >= $1
		  AND p.tanggal_selesai <= $2
		  AND k.kode_kehadiran IS NULL
		ORDER BY p.tanggal_mulai
	`
	rows, err := r.DB.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *PeminjamanRepository) UpdateStatus(kodePeminjaman string, status models.PeminjamanStatusEnum, verifiedBy *string, catatan string) error {
	query := `
		UPDATE peminjaman
		SET status = $1,
		    verified_by = $2,
		    verified_at = NOW(),
		    catatan_verifikasi = $3,
		    updated_at = NOW()
		WHERE kode_peminjaman = $4
	`
	var verifier interface{}
	if verifiedBy != nil {
		verifier = *verifiedBy
	}
	_, err := r.DB.Exec(query, status, verifier, catatan, kodePeminjaman)
	return err
}

func (r *PeminjamanRepository) GetPeminjamanBarang(kodePeminjaman string) ([]models.PeminjamanBarangDetail, error) {
	query := `
		SELECT pb.kode_peminjaman_barang, pb.kode_barang, pb.jumlah,
		       b.kode_barang, b.nama_barang, b.deskripsi, b.jumlah_total, b.ruangan_kode
		FROM peminjaman_barang pb
		JOIN barang b ON pb.kode_barang = b.kode_barang
		WHERE pb.kode_peminjaman = $1
	`
	rows, err := r.DB.Query(query, kodePeminjaman)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.PeminjamanBarangDetail
	for rows.Next() {
		var (
			item        models.PeminjamanBarangDetail
			barang      models.Barang
			ruanganKode sql.NullString
		)
		err := rows.Scan(
			&item.KodePeminjamanBarang,
			&item.KodeBarang,
			&item.Jumlah,
			&barang.KodeBarang,
			&barang.NamaBarang,
			&barang.Deskripsi,
			&barang.JumlahTotal,
			&ruanganKode,
		)
		if err != nil {
			return nil, err
		}
		if ruanganKode.Valid {
			barang.RuanganKode = &ruanganKode.String
		}
		item.Barang = &barang
		items = append(items, item)
	}
	return items, nil
}

func (r *PeminjamanRepository) GetLaporan(start, end time.Time, status models.PeminjamanStatusEnum) ([]models.Peminjaman, error) {
	query := `
		SELECT kode_peminjaman, kode_user, kode_ruangan, kode_kegiatan, tanggal_mulai, tanggal_selesai,
		       keperluan, status, path_surat_digital, verified_by, verified_at,
		       catatan_verifikasi, created_at, updated_at
		FROM peminjaman
		WHERE created_at >= $1 AND created_at <= $2
	`
	args := []interface{}{start, end}
	if status != "" {
		query += " AND status = $3"
		args = append(args, status)
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *PeminjamanRepository) UpdateSuratDigitalURL(kodePeminjaman string, path string) error {
	query := `
		UPDATE peminjaman
		SET path_surat_digital = $1,
		    updated_at = NOW()
		WHERE kode_peminjaman = $2
	`
	_, err := r.DB.Exec(query, path, kodePeminjaman)
	return err
}
