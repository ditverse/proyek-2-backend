package services

import (
	"fmt"
	"strings"
	"time"
)

// EmailTemplateData contains all data needed for email templates
type EmailTemplateData struct {
	// Peminjaman info
	KodePeminjaman    string
	Status            string
	CatatanVerifikasi string

	// Peminjam (Mahasiswa) info
	NamaPeminjam   string
	EmailPeminjam  string
	NoHPPeminjam   string
	NamaOrganisasi string

	// Kegiatan info
	NamaKegiatan      string
	DeskripsiKegiatan string

	// Ruangan info
	NamaRuangan   string
	LokasiRuangan string
	Kapasitas     int

	// Time info
	TanggalMulai      time.Time
	TanggalSelesai    time.Time
	TanggalVerifikasi time.Time

	// Verifier info
	NamaVerifikator string

	// Barang list
	Barang []BarangItem
}

// BarangItem represents borrowed item for template
type BarangItem struct {
	NamaBarang string
	Jumlah     int
}

// BuildApprovedEmailHTML builds HTML email for approved loan
func BuildApprovedEmailHTML(data EmailTemplateData) string {
	barangList := buildBarangListHTML(data.Barang)
	catatan := data.CatatanVerifikasi
	if catatan == "" {
		catatan = "-"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Peminjaman Disetujui</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f5f5f5; }
        .container { background-color: #ffffff; border-radius: 10px; padding: 30px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #28a745, #20c997); color: white; padding: 20px; border-radius: 10px 10px 0 0; margin: -30px -30px 20px -30px; text-align: center; }
        .header h1 { margin: 0; font-size: 24px; }
        .header .icon { font-size: 48px; margin-bottom: 10px; }
        .section { background: #f8f9fa; border-left: 4px solid #28a745; padding: 15px; margin: 15px 0; border-radius: 0 5px 5px 0; }
        .section-title { font-weight: bold; color: #28a745; margin-bottom: 10px; font-size: 14px; text-transform: uppercase; letter-spacing: 1px; }
        .info-row { display: flex; margin: 8px 0; }
        .info-label { color: #666; min-width: 140px; }
        .info-value { color: #333; font-weight: 500; }
        .status-approved { background: #d4edda; color: #155724; padding: 5px 15px; border-radius: 20px; font-weight: bold; display: inline-block; }
        .barang-list { list-style: none; padding: 0; margin: 0; }
        .barang-list li { padding: 8px 0; border-bottom: 1px solid #eee; }
        .barang-list li:last-child { border-bottom: none; }
        .catatan-box { background: #fff3cd; border: 1px solid #ffc107; padding: 15px; border-radius: 5px; margin: 15px 0; }
        .catatan-title { color: #856404; font-weight: bold; margin-bottom: 5px; }
        .warning-box { background: #e7f3ff; border: 1px solid #0066cc; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .warning-title { color: #0066cc; font-weight: bold; margin-bottom: 10px; }
        .warning-list { margin: 0; padding-left: 20px; }
        .warning-list li { margin: 5px 0; color: #333; }
        .footer { text-align: center; margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="icon">✅</div>
            <h1>Peminjaman Disetujui</h1>
        </div>
        
        <p>Kepada Yth. <strong>%s</strong><br>
        %s</p>
        
        <p>Pengajuan peminjaman fasilitas Anda telah <strong>DISETUJUI</strong> dengan detail sebagai berikut:</p>
        
        <div class="section">
            <div class="section-title">📋 Informasi Peminjaman</div>
            <div class="info-row"><span class="info-label">Kode Peminjaman</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Status</span><span class="info-value"><span class="status-approved">✅ DISETUJUI</span></span></div>
            <div class="info-row"><span class="info-label">Diverifikasi oleh</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Tanggal Verifikasi</span><span class="info-value">%s</span></div>
        </div>
        
        <div class="section">
            <div class="section-title">📝 Detail Kegiatan</div>
            <div class="info-row"><span class="info-label">Nama Kegiatan</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Deskripsi</span><span class="info-value">%s</span></div>
        </div>
        
        <div class="section">
            <div class="section-title">📍 Waktu & Tempat</div>
            <div class="info-row"><span class="info-label">Ruangan</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Lokasi</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Kapasitas</span><span class="info-value">%d orang</span></div>
            <div class="info-row"><span class="info-label">Mulai</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Selesai</span><span class="info-value">%s</span></div>
        </div>
        
        %s
        
        <div class="catatan-box">
            <div class="catatan-title">📝 Catatan Verifikasi:</div>
            <p style="margin: 0;">%s</p>
        </div>
        
        <div class="warning-box">
            <div class="warning-title">⚠️ Perhatian:</div>
            <ol class="warning-list">
                <li>Harap hadir tepat waktu sesuai jadwal</li>
                <li>Pastikan membawa KTM/identitas yang valid</li>
                <li>Jaga kebersihan dan fasilitas yang dipinjam</li>
                <li>Kembalikan barang dalam kondisi baik setelah selesai</li>
            </ol>
        </div>
        
        <div class="footer">
            <p>Terima kasih telah menggunakan sistem peminjaman Sarpras.</p>
            <p><strong>Unit Sarana dan Prasarana</strong></p>
        </div>
    </div>
</body>
</html>`,
		EscapeHTML(data.NamaPeminjam),
		EscapeHTML(data.NamaOrganisasi),
		EscapeHTML(data.KodePeminjaman),
		EscapeHTML(data.NamaVerifikator),
		FormatDate(data.TanggalVerifikasi),
		EscapeHTML(data.NamaKegiatan),
		EscapeHTML(data.DeskripsiKegiatan),
		EscapeHTML(data.NamaRuangan),
		EscapeHTML(data.LokasiRuangan),
		data.Kapasitas,
		FormatDate(data.TanggalMulai),
		FormatDate(data.TanggalSelesai),
		barangList,
		EscapeHTML(catatan),
	)

	return html
}

// BuildRejectedEmailHTML builds HTML email for rejected loan
func BuildRejectedEmailHTML(data EmailTemplateData) string {
	alasan := data.CatatanVerifikasi
	if alasan == "" {
		alasan = "Tidak ada keterangan"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Peminjaman Ditolak</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f5f5f5; }
        .container { background-color: #ffffff; border-radius: 10px; padding: 30px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #dc3545, #c82333); color: white; padding: 20px; border-radius: 10px 10px 0 0; margin: -30px -30px 20px -30px; text-align: center; }
        .header h1 { margin: 0; font-size: 24px; }
        .header .icon { font-size: 48px; margin-bottom: 10px; }
        .section { background: #f8f9fa; border-left: 4px solid #dc3545; padding: 15px; margin: 15px 0; border-radius: 0 5px 5px 0; }
        .section-title { font-weight: bold; color: #dc3545; margin-bottom: 10px; font-size: 14px; text-transform: uppercase; letter-spacing: 1px; }
        .info-row { display: flex; margin: 8px 0; }
        .info-label { color: #666; min-width: 140px; }
        .info-value { color: #333; font-weight: 500; }
        .status-rejected { background: #f8d7da; color: #721c24; padding: 5px 15px; border-radius: 20px; font-weight: bold; display: inline-block; }
        .alasan-box { background: #f8d7da; border: 2px solid #dc3545; padding: 20px; border-radius: 5px; margin: 20px 0; }
        .alasan-title { color: #721c24; font-weight: bold; margin-bottom: 10px; font-size: 16px; }
        .alasan-text { color: #721c24; font-size: 15px; margin: 0; }
        .saran-box { background: #d1ecf1; border: 1px solid #17a2b8; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .saran-title { color: #0c5460; font-weight: bold; margin-bottom: 5px; }
        .footer { text-align: center; margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="icon">❌</div>
            <h1>Peminjaman Ditolak</h1>
        </div>
        
        <p>Kepada Yth. <strong>%s</strong><br>
        %s</p>
        
        <p>Mohon maaf, pengajuan peminjaman fasilitas Anda <strong>TIDAK DAPAT DISETUJUI</strong>.</p>
        
        <div class="section">
            <div class="section-title">📋 Informasi Peminjaman</div>
            <div class="info-row"><span class="info-label">Kode Peminjaman</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Status</span><span class="info-value"><span class="status-rejected">❌ DITOLAK</span></span></div>
            <div class="info-row"><span class="info-label">Diverifikasi oleh</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Tanggal Verifikasi</span><span class="info-value">%s</span></div>
        </div>
        
        <div class="section">
            <div class="section-title">📝 Detail Pengajuan yang Ditolak</div>
            <div class="info-row"><span class="info-label">Nama Kegiatan</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Ruangan</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Tanggal</span><span class="info-value">%s - %s</span></div>
        </div>
        
        <div class="alasan-box">
            <div class="alasan-title">🚫 Alasan Penolakan:</div>
            <p class="alasan-text">%s</p>
        </div>
        
        <div class="saran-box">
            <div class="saran-title">💡 Saran:</div>
            <p style="margin: 0; color: #0c5460;">Silakan perbaiki pengajuan sesuai dengan catatan di atas dan ajukan kembali melalui sistem. Jika ada pertanyaan, silakan hubungi Unit Sarpras.</p>
        </div>
        
        <div class="footer">
            <p>Terima kasih atas pengertiannya.</p>
            <p><strong>Unit Sarana dan Prasarana</strong></p>
        </div>
    </div>
</body>
</html>`,
		EscapeHTML(data.NamaPeminjam),
		EscapeHTML(data.NamaOrganisasi),
		EscapeHTML(data.KodePeminjaman),
		EscapeHTML(data.NamaVerifikator),
		FormatDate(data.TanggalVerifikasi),
		EscapeHTML(data.NamaKegiatan),
		EscapeHTML(data.NamaRuangan),
		FormatDateShort(data.TanggalMulai),
		FormatDateShort(data.TanggalSelesai),
		EscapeHTML(alasan),
	)

	return html
}

// BuildSecurityNotificationHTML builds HTML email for security staff
func BuildSecurityNotificationHTML(data EmailTemplateData) string {
	barangList := buildBarangListHTML(data.Barang)
	noHP := data.NoHPPeminjam
	if noHP == "" {
		noHP = "-"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Notifikasi Kegiatan</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f5f5f5; }
        .container { background-color: #ffffff; border-radius: 10px; padding: 30px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #007bff, #0056b3); color: white; padding: 20px; border-radius: 10px 10px 0 0; margin: -30px -30px 20px -30px; text-align: center; }
        .header h1 { margin: 0; font-size: 24px; }
        .header .icon { font-size: 48px; margin-bottom: 10px; }
        .section { background: #f8f9fa; border-left: 4px solid #007bff; padding: 15px; margin: 15px 0; border-radius: 0 5px 5px 0; }
        .section-title { font-weight: bold; color: #007bff; margin-bottom: 10px; font-size: 14px; text-transform: uppercase; letter-spacing: 1px; }
        .info-row { display: flex; margin: 8px 0; }
        .info-label { color: #666; min-width: 140px; }
        .info-value { color: #333; font-weight: 500; }
        .pj-section { background: #e7f3ff; border-left: 4px solid #0066cc; padding: 15px; margin: 15px 0; border-radius: 0 5px 5px 0; }
        .barang-list { list-style: none; padding: 0; margin: 0; }
        .barang-list li { padding: 8px 0; border-bottom: 1px solid #eee; }
        .barang-list li:last-child { border-bottom: none; }
        .tugas-box { background: #fff3cd; border: 1px solid #ffc107; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .tugas-title { color: #856404; font-weight: bold; margin-bottom: 10px; }
        .tugas-list { margin: 0; padding-left: 20px; }
        .tugas-list li { margin: 5px 0; }
        .footer { text-align: center; margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; color: #666; font-size: 12px; }
        .verified-info { background: #d4edda; padding: 10px 15px; border-radius: 5px; margin-top: 20px; font-size: 13px; color: #155724; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="icon">🔔</div>
            <h1>Notifikasi Kegiatan</h1>
        </div>
        
        <p>Kepada Yth. <strong>Petugas Security</strong>,</p>
        
        <p>Berikut informasi kegiatan yang telah diverifikasi dan perlu dimonitor:</p>
        
        <div class="section">
            <div class="section-title">📋 Informasi Kegiatan</div>
            <div class="info-row"><span class="info-label">Nama Kegiatan</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Deskripsi</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Kode Peminjaman</span><span class="info-value">%s</span></div>
        </div>
        
        <div class="pj-section">
            <div class="section-title">👤 Penanggung Jawab</div>
            <div class="info-row"><span class="info-label">Nama</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Organisasi</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">No. HP</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Email</span><span class="info-value">%s</span></div>
        </div>
        
        <div class="section">
            <div class="section-title">📍 Waktu & Lokasi</div>
            <div class="info-row"><span class="info-label">Ruangan</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Lokasi</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Mulai</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Selesai</span><span class="info-value">%s</span></div>
        </div>
        
        %s
        
        <div class="tugas-box">
            <div class="tugas-title">📋 Tugas Monitoring:</div>
            <ol class="tugas-list">
                <li>Verifikasi identitas PJ saat check-in</li>
                <li>Pastikan kegiatan berjalan sesuai jadwal</li>
                <li>Pantau kondisi fasilitas selama kegiatan</li>
                <li>Laporkan jika ada kendala/masalah</li>
            </ol>
        </div>
        
        <div class="verified-info">
            ✅ Diverifikasi oleh: <strong>%s</strong> pada %s
        </div>
        
        <div class="footer">
            <p>Email ini dikirim otomatis oleh Sistem Sarpras</p>
        </div>
    </div>
</body>
</html>`,
		EscapeHTML(data.NamaKegiatan),
		EscapeHTML(data.DeskripsiKegiatan),
		EscapeHTML(data.KodePeminjaman),
		EscapeHTML(data.NamaPeminjam),
		EscapeHTML(data.NamaOrganisasi),
		EscapeHTML(noHP),
		EscapeHTML(data.EmailPeminjam),
		EscapeHTML(data.NamaRuangan),
		EscapeHTML(data.LokasiRuangan),
		FormatDate(data.TanggalMulai),
		FormatDate(data.TanggalSelesai),
		barangList,
		EscapeHTML(data.NamaVerifikator),
		FormatDate(data.TanggalVerifikasi),
	)

	return html
}

// buildBarangListHTML builds HTML list of borrowed items
func buildBarangListHTML(barang []BarangItem) string {
	if len(barang) == 0 {
		return ""
	}

	var items []string
	for _, b := range barang {
		items = append(items, fmt.Sprintf(`<li>📦 <strong>%s</strong> - %d unit</li>`,
			EscapeHTML(b.NamaBarang), b.Jumlah))
	}

	return fmt.Sprintf(`
        <div class="section">
            <div class="section-title">📦 Barang yang Dipinjam</div>
            <ul class="barang-list">
                %s
            </ul>
        </div>`, strings.Join(items, "\n                "))
}

// GetApprovedEmailSubject returns subject for approved email
func GetApprovedEmailSubject(namaKegiatan string) string {
	return fmt.Sprintf("Peminjaman Fasilitas Disetujui - %s", namaKegiatan)
}

// GetRejectedEmailSubject returns subject for rejected email
func GetRejectedEmailSubject(namaKegiatan string) string {
	return fmt.Sprintf("Peminjaman Fasilitas Ditolak - %s", namaKegiatan)
}

// GetSecurityEmailSubject returns subject for security notification
func GetSecurityEmailSubject(namaKegiatan string, tanggal time.Time) string {
	return fmt.Sprintf("Jadwal Kegiatan: %s - %s", namaKegiatan, FormatDateShort(tanggal))
}

// GetCancelledEmailSubject returns subject for cancelled email
func GetCancelledEmailSubject(namaKegiatan string) string {
	return fmt.Sprintf("Peminjaman Dibatalkan: %s", namaKegiatan)
}

// BuildCancelledEmailHTML builds HTML email for cancelled loan
func BuildCancelledEmailHTML(data EmailTemplateData) string {
	alasan := data.CatatanVerifikasi
	if alasan == "" {
		alasan = "Tidak ada keterangan"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Peminjaman Dibatalkan</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f5f5f5; }
        .container { background-color: #ffffff; border-radius: 10px; padding: 30px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #fd7e14, #dc3545); color: white; padding: 20px; border-radius: 10px 10px 0 0; margin: -30px -30px 20px -30px; text-align: center; }
        .header h1 { margin: 0; font-size: 24px; }
        .header .icon { font-size: 48px; margin-bottom: 10px; }
        .section { background: #f8f9fa; border-left: 4px solid #fd7e14; padding: 15px; margin: 15px 0; border-radius: 0 5px 5px 0; }
        .section-title { font-weight: bold; color: #fd7e14; margin-bottom: 10px; font-size: 14px; text-transform: uppercase; letter-spacing: 1px; }
        .info-row { display: flex; margin: 8px 0; }
        .info-label { color: #666; min-width: 140px; }
        .info-value { color: #333; font-weight: 500; }
        .status-cancelled { background: #ffe5d0; color: #c92a2a; padding: 5px 15px; border-radius: 20px; font-weight: bold; display: inline-block; }
        .alasan-box { background: #fff3cd; border: 2px solid #fd7e14; padding: 20px; border-radius: 5px; margin: 20px 0; }
        .alasan-title { color: #856404; font-weight: bold; margin-bottom: 10px; font-size: 16px; }
        .alasan-text { color: #856404; font-size: 15px; margin: 0; }
        .saran-box { background: #d1ecf1; border: 1px solid #17a2b8; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .saran-title { color: #0c5460; font-weight: bold; margin-bottom: 5px; }
        .footer { text-align: center; margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="icon">⚠️</div>
            <h1>Peminjaman Dibatalkan</h1>
        </div>
        
        <p>Kepada Yth. <strong>%s</strong>,</p>
        
        <p>Kami informasikan bahwa peminjaman fasilitas Anda telah <strong>DIBATALKAN</strong> oleh petugas Sarpras.</p>
        
        <div class="section">
            <div class="section-title">📋 Informasi Peminjaman</div>
            <div class="info-row"><span class="info-label">Kode Peminjaman</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Status</span><span class="info-value"><span class="status-cancelled">⚠️ DIBATALKAN</span></span></div>
            <div class="info-row"><span class="info-label">Dibatalkan oleh</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Tanggal Pembatalan</span><span class="info-value">%s</span></div>
        </div>
        
        <div class="section">
            <div class="section-title">📝 Detail Pengajuan yang Dibatalkan</div>
            <div class="info-row"><span class="info-label">Nama Kegiatan</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Ruangan</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Tanggal</span><span class="info-value">%s - %s</span></div>
        </div>
        
        <div class="alasan-box">
            <div class="alasan-title">📝 Alasan Pembatalan:</div>
            <p class="alasan-text">%s</p>
        </div>
        
        <div class="saran-box">
            <div class="saran-title">💡 Informasi:</div>
            <p style="margin: 0; color: #0c5460;">Jika Anda merasa pembatalan ini tidak sesuai atau ada pertanyaan, silakan hubungi Unit Sarpras untuk informasi lebih lanjut.</p>
        </div>
        
        <div class="footer">
            <p>Terima kasih atas pengertiannya.</p>
            <p><strong>Unit Sarana dan Prasarana</strong></p>
        </div>
    </div>
</body>
</html>`,
		EscapeHTML(data.NamaPeminjam),
		EscapeHTML(data.KodePeminjaman),
		EscapeHTML(data.NamaVerifikator),
		FormatDate(data.TanggalVerifikasi),
		EscapeHTML(data.NamaKegiatan),
		EscapeHTML(data.NamaRuangan),
		FormatDateShort(data.TanggalMulai),
		FormatDateShort(data.TanggalSelesai),
		EscapeHTML(alasan),
	)

	return html
}

// GetNewSubmissionEmailSubject returns subject for new submission email
func GetNewSubmissionEmailSubject(namaKegiatan string) string {
	return fmt.Sprintf("Pengajuan Peminjaman Baru: %s", namaKegiatan)
}

// BuildNewSubmissionEmailHTML builds HTML email for Sarpras notification
func BuildNewSubmissionEmailHTML(data EmailTemplateData) string {
	barangList := buildBarangListHTML(data.Barang)

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Pengajuan Peminjaman Baru</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f5f5f5; }
        .container { background-color: #ffffff; border-radius: 10px; padding: 30px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #17a2b8, #138496); color: white; padding: 20px; border-radius: 10px 10px 0 0; margin: -30px -30px 20px -30px; text-align: center; }
        .header h1 { margin: 0; font-size: 24px; }
        .header .icon { font-size: 48px; margin-bottom: 10px; }
        .section { background: #f8f9fa; border-left: 4px solid #17a2b8; padding: 15px; margin: 15px 0; border-radius: 0 5px 5px 0; }
        .section-title { font-weight: bold; color: #17a2b8; margin-bottom: 10px; font-size: 14px; text-transform: uppercase; letter-spacing: 1px; }
        .info-row { display: flex; margin: 8px 0; }
        .info-label { color: #666; min-width: 140px; }
        .info-value { color: #333; font-weight: 500; }
        .status-new { background: #d1ecf1; color: #0c5460; padding: 5px 15px; border-radius: 20px; font-weight: bold; display: inline-block; }
        .peminjam-box { background: #e2e3e5; border-left: 4px solid #6c757d; padding: 15px; margin: 15px 0; border-radius: 0 5px 5px 0; }
        .barang-list { list-style: none; padding: 0; margin: 0; }
        .barang-list li { padding: 8px 0; border-bottom: 1px solid #eee; }
        .barang-list li:last-child { border-bottom: none; }
        .action-box { background: #fff3cd; border: 1px solid #ffc107; padding: 20px; border-radius: 5px; margin: 25px 0; text-align: center; }
        .action-title { color: #856404; font-weight: bold; margin-bottom: 10px; font-size: 16px; }
        .footer { text-align: center; margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="icon">📝</div>
            <h1>Pengajuan Baru</h1>
        </div>
        
        <p>Halo <strong>Tim Sarpras</strong>,</p>
        
        <p>Terdapat pengajuan peminjaman fasilitas baru yang membutuhkan verifikasi Anda.</p>
        
        <div class="section">
            <div class="section-title">📋 Informasi Pengajuan</div>
            <div class="info-row"><span class="info-label">Kode Peminjaman</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Status</span><span class="info-value"><span class="status-new">🆕 MENUNGGU VERIFIKASI</span></span></div>
            <div class="info-row"><span class="info-label">Tanggal Pengajuan</span><span class="info-value">%s</span></div>
        </div>
        
        <div class="peminjam-box">
            <div class="section-title">👤 Data Peminjam</div>
            <div class="info-row"><span class="info-label">Nama</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Organisasi</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">No. HP</span><span class="info-value">%s</span></div>
        </div>

        <div class="section">
            <div class="section-title">📝 Detail Kegiatan</div>
            <div class="info-row"><span class="info-label">Nama Kegiatan</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Ruangan</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Lokasi</span><span class="info-value">%s</span></div>
            <div class="info-row"><span class="info-label">Tanggal</span><span class="info-value">%s s/d %s</span></div>
        </div>
        
        %s
        
        <div class="action-box">
            <div class="action-title">⚡ Tindakan Diperlukan:</div>
            <p style="margin: 0; color: #856404;">Silakan login ke dashboard sistem untuk melihat detail lengkap dan melakukan verifikasi (Setujui/Tolak).</p>
        </div>
        
        <div class="footer">
            <p>Email ini dikirim otomatis oleh Sistem Sarpras</p>
        </div>
    </div>
</body>
</html>`,
		EscapeHTML(data.KodePeminjaman),
		FormatDate(time.Now()), // Waktu pengajuan (saat ini)
		EscapeHTML(data.NamaPeminjam),
		EscapeHTML(data.NamaOrganisasi),
		EscapeHTML(data.NoHPPeminjam),
		EscapeHTML(data.NamaKegiatan),
		EscapeHTML(data.NamaRuangan),
		EscapeHTML(data.LokasiRuangan),
		FormatDate(data.TanggalMulai),
		FormatDate(data.TanggalSelesai),
		barangList,
	)

	return html
}
