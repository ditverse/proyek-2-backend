package repositories

import (
	"backend-sarpras/models"
	"database/sql"
	"fmt"
	"strings"
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
			tanggal_mulai, tanggal_selesai, status, path_surat_digital
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
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
		       status, path_surat_digital, verified_by, verified_at,
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
		       p.status, p.path_surat_digital, p.verified_by, p.verified_at,
		       p.catatan_verifikasi, p.created_at, p.updated_at,
		       r.kode_ruangan, r.nama_ruangan, r.lokasi, r.kapasitas, r.deskripsi,
		       k.kode_kegiatan, k.nama_kegiatan, k.deskripsi, k.tanggal_mulai, k.tanggal_selesai, k.organisasi_kode
		FROM peminjaman p
		LEFT JOIN ruangan r ON p.kode_ruangan = r.kode_ruangan
		LEFT JOIN kegiatan k ON p.kode_kegiatan = k.kode_kegiatan
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

			// Kegiatan fields
			kKode           sql.NullString
			kNama           sql.NullString
			kDeskripsi      sql.NullString
			kTanggalMulai   sql.NullTime
			kTanggalSelesai sql.NullTime
			kOrganisasiKode sql.NullString
		)
		err := rows.Scan(
			&row.KodePeminjaman,
			&row.KodeUser,
			&kodeRuangan,
			&kodeKegiatan,
			&row.TanggalMulai,
			&row.TanggalSelesai,
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
			&kKode,
			&kNama,
			&kDeskripsi,
			&kTanggalMulai,
			&kTanggalSelesai,
			&kOrganisasiKode,
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

		// Populate Kegiatan struct if data exists
		if kKode.Valid {
			row.Kegiatan = &models.Kegiatan{
				KodeKegiatan:   kKode.String,
				NamaKegiatan:   kNama.String,
				Deskripsi:      kDeskripsi.String,
				TanggalMulai:   kTanggalMulai.Time,
				TanggalSelesai: kTanggalSelesai.Time,
				OrganisasiKode: kOrganisasiKode.String,
			}
		}

		result = append(result, row)
	}
	return result, nil
}

func (r *PeminjamanRepository) GetPending() ([]models.Peminjaman, error) {
	query := `
		SELECT p.kode_peminjaman, p.kode_user, p.kode_ruangan, p.kode_kegiatan, p.tanggal_mulai, p.tanggal_selesai,
		       p.status, p.path_surat_digital, p.verified_by, p.verified_at,
		       p.catatan_verifikasi, p.created_at, p.updated_at,
		       r.kode_ruangan, r.nama_ruangan, r.lokasi, r.kapasitas, r.deskripsi,
		       k.kode_kegiatan, k.nama_kegiatan, k.deskripsi, k.tanggal_mulai, k.tanggal_selesai, k.organisasi_kode,
		       u.kode_user, u.nama, u.email, u.role, u.organisasi_kode
		FROM peminjaman p
		LEFT JOIN ruangan r ON p.kode_ruangan = r.kode_ruangan
		LEFT JOIN kegiatan k ON p.kode_kegiatan = k.kode_kegiatan
		LEFT JOIN users u ON p.kode_user = u.kode_user
		WHERE p.status = 'PENDING'
		ORDER BY p.created_at ASC
	`
	rows, err := r.DB.Query(query)
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

			// Kegiatan fields
			kKode           sql.NullString
			kNama           sql.NullString
			kDeskripsi      sql.NullString
			kTanggalMulai   sql.NullTime
			kTanggalSelesai sql.NullTime
			kOrganisasiKode sql.NullString

			// User fields
			uKode           sql.NullString
			uNama           sql.NullString
			uEmail          sql.NullString
			uRole           sql.NullString
			uOrganisasiKode sql.NullString
		)
		err := rows.Scan(
			&row.KodePeminjaman,
			&row.KodeUser,
			&kodeRuangan,
			&kodeKegiatan,
			&row.TanggalMulai,
			&row.TanggalSelesai,
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
			&kKode,
			&kNama,
			&kDeskripsi,
			&kTanggalMulai,
			&kTanggalSelesai,
			&kOrganisasiKode,
			&uKode,
			&uNama,
			&uEmail,
			&uRole,
			&uOrganisasiKode,
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

		// Populate Kegiatan struct if data exists
		if kKode.Valid {
			row.Kegiatan = &models.Kegiatan{
				KodeKegiatan:   kKode.String,
				NamaKegiatan:   kNama.String,
				Deskripsi:      kDeskripsi.String,
				TanggalMulai:   kTanggalMulai.Time,
				TanggalSelesai: kTanggalSelesai.Time,
				OrganisasiKode: kOrganisasiKode.String,
			}
		}

		// Populate User (Peminjam) struct if data exists
		if uKode.Valid {
			peminjam := &models.User{
				KodeUser: uKode.String,
				Nama:     uNama.String,
				Email:    uEmail.String,
				Role:     models.RoleEnum(uRole.String),
			}
			if uOrganisasiKode.Valid {
				peminjam.OrganisasiKode = &uOrganisasiKode.String
			}
			row.Peminjam = peminjam
		}

		result = append(result, row)
	}
	return result, nil
}

func (r *PeminjamanRepository) GetJadwalRuangan(start, end time.Time) ([]models.JadwalRuanganResponse, error) {
	query := `
		SELECT p.kode_peminjaman, p.kode_ruangan, r.nama_ruangan, p.tanggal_mulai, p.tanggal_selesai,
		       p.status, u.nama, COALESCE(u.organisasi_kode, '') as organisasi_kode,
			   COALESCE(k.nama_kegiatan, 'Penggunaan Ruangan') as nama_kegiatan
		FROM peminjaman p
		JOIN ruangan r ON p.kode_ruangan = r.kode_ruangan
		JOIN users u ON p.kode_user = u.kode_user
		LEFT JOIN kegiatan k ON p.kode_kegiatan = k.kode_kegiatan
		WHERE p.kode_ruangan IS NOT NULL
		  AND p.status IN ('PENDING', 'APPROVED')
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
			&j.Organisasi,
			&j.NamaKegiatan,
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
		       status, path_surat_digital, verified_by, verified_at,
		       catatan_verifikasi, created_at, updated_at
		FROM peminjaman
		WHERE status IN ('PENDING', 'APPROVED')
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
		       p.status, p.path_surat_digital, p.verified_by, p.verified_at,
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
		SELECT p.kode_peminjaman, p.kode_user, p.kode_ruangan, p.kode_kegiatan, p.tanggal_mulai, p.tanggal_selesai,
		       p.status, p.path_surat_digital, p.verified_by, p.verified_at,
		       p.catatan_verifikasi, p.created_at, p.updated_at,
		       r.kode_ruangan, r.nama_ruangan, r.lokasi, r.kapasitas, r.deskripsi,
		       k.kode_kegiatan, k.nama_kegiatan, k.deskripsi, k.tanggal_mulai, k.tanggal_selesai, k.organisasi_kode,
		       u.kode_user, u.nama, u.email, u.role, u.organisasi_kode
		FROM peminjaman p
		LEFT JOIN ruangan r ON p.kode_ruangan = r.kode_ruangan
		LEFT JOIN kegiatan k ON p.kode_kegiatan = k.kode_kegiatan
		LEFT JOIN users u ON p.kode_user = u.kode_user
		WHERE p.created_at >= $1 AND p.created_at <= $2
	`
	args := []interface{}{start, end}
	if status != "" {
		query += " AND p.status = $3"
		args = append(args, status)
	}
	query += " ORDER BY p.created_at DESC"

	rows, err := r.DB.Query(query, args...)
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

			// Kegiatan fields
			kKode           sql.NullString
			kNama           sql.NullString
			kDeskripsi      sql.NullString
			kTanggalMulai   sql.NullTime
			kTanggalSelesai sql.NullTime
			kOrganisasiKode sql.NullString

			// User fields
			uKode           sql.NullString
			uNama           sql.NullString
			uEmail          sql.NullString
			uRole           sql.NullString
			uOrganisasiKode sql.NullString
		)
		err := rows.Scan(
			&row.KodePeminjaman,
			&row.KodeUser,
			&kodeRuangan,
			&kodeKegiatan,
			&row.TanggalMulai,
			&row.TanggalSelesai,
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
			&kKode,
			&kNama,
			&kDeskripsi,
			&kTanggalMulai,
			&kTanggalSelesai,
			&kOrganisasiKode,
			&uKode,
			&uNama,
			&uEmail,
			&uRole,
			&uOrganisasiKode,
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

		// Populate Kegiatan struct if data exists
		if kKode.Valid {
			row.Kegiatan = &models.Kegiatan{
				KodeKegiatan:   kKode.String,
				NamaKegiatan:   kNama.String,
				Deskripsi:      kDeskripsi.String,
				TanggalMulai:   kTanggalMulai.Time,
				TanggalSelesai: kTanggalSelesai.Time,
				OrganisasiKode: kOrganisasiKode.String,
			}
		}

		// Populate User (Peminjam) struct if data exists
		if uKode.Valid {
			peminjam := &models.User{
				KodeUser: uKode.String,
				Nama:     uNama.String,
				Email:    uEmail.String,
				Role:     models.RoleEnum(uRole.String),
			}
			if uOrganisasiKode.Valid {
				peminjam.OrganisasiKode = &uOrganisasiKode.String
			}
			row.Peminjam = peminjam
		}

		result = append(result, row)
	}
	return result, nil
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

// IsRoomAvailable mengecek apakah ruangan tersedia pada rentang waktu tertentu
// Ruangan dianggap tidak tersedia jika ada peminjaman dengan status PENDING atau APPROVED
// yang waktunya overlap dengan rentang waktu yang diminta
func (r *PeminjamanRepository) IsRoomAvailable(kodeRuangan string, mulai, selesai time.Time) (bool, error) {
	query := `
		SELECT COUNT(*) FROM peminjaman 
		WHERE kode_ruangan = $1 
		  AND status IN ('PENDING', 'APPROVED')
		  AND tanggal_mulai < $3 
		  AND tanggal_selesai > $2
	`
	var count int
	err := r.DB.QueryRow(query, kodeRuangan, mulai, selesai).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// BookedDateInfo menyimpan informasi tanggal yang sudah dibooking
type BookedDateInfo struct {
	TanggalMulai   time.Time `json:"tanggal_mulai"`
	TanggalSelesai time.Time `json:"tanggal_selesai"`
	Status         string    `json:"status"`
}

// GetBookedDates mendapatkan daftar rentang waktu yang sudah dibooking untuk ruangan tertentu
// dalam rentang waktu yang diminta
func (r *PeminjamanRepository) GetBookedDates(kodeRuangan string, startRange, endRange time.Time) ([]BookedDateInfo, error) {
	query := `
		SELECT tanggal_mulai, tanggal_selesai, status
		FROM peminjaman 
		WHERE kode_ruangan = $1 
		  AND status IN ('PENDING', 'APPROVED')
		  AND tanggal_mulai <= $3 
		  AND tanggal_selesai >= $2
		ORDER BY tanggal_mulai
	`
	rows, err := r.DB.Query(query, kodeRuangan, startRange, endRange)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []BookedDateInfo
	for rows.Next() {
		var info BookedDateInfo
		if err := rows.Scan(&info.TanggalMulai, &info.TanggalSelesai, &info.Status); err != nil {
			return nil, err
		}
		result = append(result, info)
	}
	if result == nil {
		result = []BookedDateInfo{}
	}
	return result, nil
}

// GetPeminjamanBarangByIDs mengambil semua barang untuk multiple peminjaman sekaligus (batch query)
// Returns a map where key is kode_peminjaman and value is slice of PeminjamanBarangDetail
func (r *PeminjamanRepository) GetPeminjamanBarangByIDs(kodePeminjamanList []string) (map[string][]models.PeminjamanBarangDetail, error) {
	result := make(map[string][]models.PeminjamanBarangDetail)
	if len(kodePeminjamanList) == 0 {
		return result, nil
	}

	// Build placeholders ($1, $2, $3, ...)
	placeholders := make([]string, len(kodePeminjamanList))
	args := make([]interface{}, len(kodePeminjamanList))
	for i, kode := range kodePeminjamanList {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = kode
	}

	query := fmt.Sprintf(`
		SELECT pb.kode_peminjaman, pb.kode_peminjaman_barang, pb.kode_barang, pb.jumlah,
		       b.kode_barang, b.nama_barang, b.deskripsi, b.jumlah_total, b.ruangan_kode
		FROM peminjaman_barang pb
		JOIN barang b ON pb.kode_barang = b.kode_barang
		WHERE pb.kode_peminjaman IN (%s)
	`, strings.Join(placeholders, ", "))

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			kodePeminjaman string
			item           models.PeminjamanBarangDetail
			barang         models.Barang
			ruanganKode    sql.NullString
		)
		err := rows.Scan(
			&kodePeminjaman,
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
		result[kodePeminjaman] = append(result[kodePeminjaman], item)
	}

	return result, nil
}

// UpdateStatusToOngoing updates APPROVED peminjaman to ONGOING
// where tanggal_mulai <= now <= tanggal_selesai
func (r *PeminjamanRepository) UpdateStatusToOngoing(now time.Time) (int64, error) {
	query := `
		UPDATE peminjaman
		SET status = 'ONGOING',
		    updated_at = NOW()
		WHERE status = 'APPROVED'
		  AND tanggal_mulai <= $1
		  AND tanggal_selesai > $1
	`
	result, err := r.DB.Exec(query, now)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdateStatusToFinished updates APPROVED/ONGOING peminjaman to FINISHED
// where tanggal_selesai < now
func (r *PeminjamanRepository) UpdateStatusToFinished(now time.Time) (int64, error) {
	query := `
		UPDATE peminjaman
		SET status = 'FINISHED',
		    updated_at = NOW()
		WHERE status IN ('APPROVED', 'ONGOING')
		  AND tanggal_selesai < $1
	`
	result, err := r.DB.Exec(query, now)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdateStatusOnly updates only the status field without changing verifier info
// Used by scheduler for automatic status transitions
func (r *PeminjamanRepository) UpdateStatusOnly(kodePeminjaman string, status models.PeminjamanStatusEnum) error {
	query := `
		UPDATE peminjaman
		SET status = $1,
		    updated_at = NOW()
		WHERE kode_peminjaman = $2
	`
	_, err := r.DB.Exec(query, status, kodePeminjaman)
	return err
}
