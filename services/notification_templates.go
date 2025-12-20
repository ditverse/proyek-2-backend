package services

import (
	"fmt"
	"time"
)

// ===============================
// WhatsApp Message Templates
// ===============================

// WATemplateApproved returns WhatsApp message for approved peminjaman
func WATemplateApproved(namaKegiatan, namaRuangan string) string {
	return fmt.Sprintf(`✅ Status Update: DISETUJUI
Kegiatan: %s
Ruangan: %s

Silakan cek email untuk surat izin digital.`, namaKegiatan, namaRuangan)
}

// WATemplateSecurity returns WhatsApp message for security notification
func WATemplateSecurity(namaKegiatan, namaRuangan string, mulai, selesai time.Time) string {
	return fmt.Sprintf(`👮 MONITOR KEGIATAN
Judul: %s
Lokasi: %s
Jam: %s s/d %s

Mohon dipantau.`, namaKegiatan, namaRuangan,
		mulai.Format("02 Jan 2006 15:04"),
		selesai.Format("02 Jan 2006 15:04"))
}

// WATemplateReminder1Hour returns WhatsApp reminder message (1 hour before end)
func WATemplateReminder1Hour(namaRuangan string) string {
	return fmt.Sprintf(`⏳ REMINDER WAKTU
Sisa waktu peminjaman ruangan %s tinggal 1 jam lagi.
Mohon persiapan untuk check-out.`, namaRuangan)
}

// ===============================
// Email Templates (HTML)
// ===============================

// EmailTemplatePengajuanBaru returns HTML email for new peminjaman notification to Sarpras
func EmailTemplatePengajuanBaru(namaMahasiswa, namaKegiatan, namaRuangan string, mulai, selesai time.Time) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .header { background: #2196F3; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .info-table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
        .info-table td { padding: 10px; border-bottom: 1px solid #eee; }
        .info-table td:first-child { font-weight: bold; width: 150px; }
        .footer { background: #f5f5f5; padding: 15px; text-align: center; font-size: 12px; color: #666; }
        .btn { display: inline-block; background: #4CAF50; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; margin-top: 15px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>📋 Pengajuan Peminjaman Baru</h1>
    </div>
    <div class="content">
        <p>Halo Tim Sarpras,</p>
        <p>Terdapat pengajuan peminjaman baru yang memerlukan verifikasi Anda.</p>
        
        <table class="info-table">
            <tr><td>Pemohon</td><td>%s</td></tr>
            <tr><td>Kegiatan</td><td>%s</td></tr>
            <tr><td>Ruangan</td><td>%s</td></tr>
            <tr><td>Waktu Mulai</td><td>%s</td></tr>
            <tr><td>Waktu Selesai</td><td>%s</td></tr>
        </table>
        
        <p>Silakan login ke sistem untuk melakukan verifikasi.</p>
    </div>
    <div class="footer">
        <p>Sistem Sarpras - Peminjaman Sarana Prasarana Kampus</p>
        <p>Email ini dikirim secara otomatis, mohon tidak membalas.</p>
    </div>
</body>
</html>`, namaMahasiswa, namaKegiatan, namaRuangan,
		mulai.Format("02 January 2006, 15:04 WIB"),
		selesai.Format("02 January 2006, 15:04 WIB"))
}

// EmailTemplateApproved returns HTML email for approved peminjaman
func EmailTemplateApproved(namaMahasiswa, namaKegiatan, namaRuangan string, mulai, selesai time.Time) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .header { background: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .info-table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
        .info-table td { padding: 10px; border-bottom: 1px solid #eee; }
        .info-table td:first-child { font-weight: bold; width: 150px; }
        .footer { background: #f5f5f5; padding: 15px; text-align: center; font-size: 12px; color: #666; }
        .status { background: #E8F5E9; border-left: 4px solid #4CAF50; padding: 15px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="header">
        <h1>✅ Peminjaman Disetujui</h1>
    </div>
    <div class="content">
        <p>Halo <strong>%s</strong>,</p>
        
        <div class="status">
            <strong>Status:</strong> DISETUJUI ✅
        </div>
        
        <p>Pengajuan peminjaman Anda telah <strong>DISETUJUI</strong>. Berikut adalah detail peminjaman:</p>
        
        <table class="info-table">
            <tr><td>Kegiatan</td><td>%s</td></tr>
            <tr><td>Ruangan</td><td>%s</td></tr>
            <tr><td>Waktu Mulai</td><td>%s</td></tr>
            <tr><td>Waktu Selesai</td><td>%s</td></tr>
        </table>
        
        <p>📌 <strong>Catatan Penting:</strong></p>
        <ul>
            <li>Silakan tunjukkan email ini sebagai bukti izin kepada petugas Security.</li>
            <li>Pastikan Anda hadir tepat waktu sesuai jadwal.</li>
            <li>Jaga kebersihan dan kerapian ruangan selama penggunaan.</li>
        </ul>
    </div>
    <div class="footer">
        <p>Sistem Sarpras - Peminjaman Sarana Prasarana Kampus</p>
        <p>Email ini dikirim secara otomatis, mohon tidak membalas.</p>
    </div>
</body>
</html>`, namaMahasiswa, namaKegiatan, namaRuangan,
		mulai.Format("02 January 2006, 15:04 WIB"),
		selesai.Format("02 January 2006, 15:04 WIB"))
}

// EmailTemplateRejected returns HTML email for rejected peminjaman
func EmailTemplateRejected(namaMahasiswa, namaKegiatan, namaRuangan, alasanPenolakan string, mulai, selesai time.Time) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .header { background: #f44336; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .info-table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
        .info-table td { padding: 10px; border-bottom: 1px solid #eee; }
        .info-table td:first-child { font-weight: bold; width: 150px; }
        .footer { background: #f5f5f5; padding: 15px; text-align: center; font-size: 12px; color: #666; }
        .status { background: #FFEBEE; border-left: 4px solid #f44336; padding: 15px; margin: 20px 0; }
        .reason { background: #FFF3E0; border-left: 4px solid #FF9800; padding: 15px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="header">
        <h1>❌ Peminjaman Ditolak</h1>
    </div>
    <div class="content">
        <p>Halo <strong>%s</strong>,</p>
        
        <div class="status">
            <strong>Status:</strong> DITOLAK ❌
        </div>
        
        <p>Mohon maaf, pengajuan peminjaman Anda telah <strong>DITOLAK</strong>. Berikut adalah detail pengajuan:</p>
        
        <table class="info-table">
            <tr><td>Kegiatan</td><td>%s</td></tr>
            <tr><td>Ruangan</td><td>%s</td></tr>
            <tr><td>Waktu Mulai</td><td>%s</td></tr>
            <tr><td>Waktu Selesai</td><td>%s</td></tr>
        </table>
        
        <div class="reason">
            <strong>📝 Alasan Penolakan:</strong><br>
            %s
        </div>
        
        <p>Jika Anda memiliki pertanyaan, silakan hubungi petugas Sarpras.</p>
    </div>
    <div class="footer">
        <p>Sistem Sarpras - Peminjaman Sarana Prasarana Kampus</p>
        <p>Email ini dikirim secara otomatis, mohon tidak membalas.</p>
    </div>
</body>
</html>`, namaMahasiswa, namaKegiatan, namaRuangan,
		mulai.Format("02 January 2006, 15:04 WIB"),
		selesai.Format("02 January 2006, 15:04 WIB"),
		alasanPenolakan)
}
