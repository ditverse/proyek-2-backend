package services

import (
	"backend-sarpras/models"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExportService struct{}

func NewExportService() *ExportService {
	return &ExportService{}
}

// GeneratePeminjamanExcel creates an Excel file from peminjaman data
func (s *ExportService) GeneratePeminjamanExcel(data []models.Peminjaman, start, end time.Time, status string) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "Laporan Peminjaman"
	
	// Create or rename default sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1") // Remove default sheet

	// Define styles
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   11,
			Color:  "FFFFFF",
			Family: "Calibri",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"FF6B35"}, // Orange color
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	cellStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical: "top",
			WrapText: true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "CCCCCC", Style: 1},
			{Type: "right", Color: "CCCCCC", Style: 1},
			{Type: "top", Color: "CCCCCC", Style: 1},
			{Type: "bottom", Color: "CCCCCC", Style: 1},
		},
	})

	// Status styles
	approvedStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "006400", // Dark green
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"D4EDDA"}, // Light green background
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "CCCCCC", Style: 1},
			{Type: "right", Color: "CCCCCC", Style: 1},
			{Type: "top", Color: "CCCCCC", Style: 1},
			{Type: "bottom", Color: "CCCCCC", Style: 1},
		},
	})

	rejectedStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "8B0000", // Dark red
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"F8D7DA"}, // Light red background
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "CCCCCC", Style: 1},
			{Type: "right", Color: "CCCCCC", Style: 1},
			{Type: "top", Color: "CCCCCC", Style: 1},
			{Type: "bottom", Color: "CCCCCC", Style: 1},
		},
	})

	pendingStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "856404", // Dark yellow
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"FFF3CD"}, // Light yellow background
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "CCCCCC", Style: 1},
			{Type: "right", Color: "CCCCCC", Style: 1},
			{Type: "top", Color: "CCCCCC", Style: 1},
			{Type: "bottom", Color: "CCCCCC", Style: 1},
		},
	})

	// Title section
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   14,
			Family: "Calibri",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	
	f.SetCellValue(sheetName, "A1", "LAPORAN PEMINJAMAN SARANA PRASARANA")
	f.SetCellStyle(sheetName, "A1", "A1", titleStyle)
	f.MergeCell(sheetName, "A1", "M1")
	
	// Info section
	infoStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   10,
			Family: "Calibri",
		},
	})
	
	filterInfo := fmt.Sprintf("Periode: %s s/d %s", 
		start.Format("02 January 2006"), 
		end.Format("02 January 2006"))
	if status != "" {
		filterInfo += fmt.Sprintf(" | Status: %s", strings.ToUpper(status))
	}
	
	f.SetCellValue(sheetName, "A2", filterInfo)
	f.SetCellStyle(sheetName, "A2", "A2", infoStyle)
	f.MergeCell(sheetName, "A2", "M2")
	
	generatedAt := fmt.Sprintf("Dibuat pada: %s", time.Now().Format("02 January 2006 15:04:05"))
	f.SetCellValue(sheetName, "A3", generatedAt)
	f.SetCellStyle(sheetName, "A3", "A3", infoStyle)
	f.MergeCell(sheetName, "A3", "M3")

	// Headers
	headers := []string{
		"No",
		"Kode Peminjaman",
		"Nama Peminjam",
		"Organisasi",
		"Nama Kegiatan",
		"Ruangan",
		"Barang",
		"Tanggal Mulai",
		"Tanggal Selesai",
		"Status",
		"Verifikator",
		"Tanggal Verifikasi",
		"Catatan Verifikasi",
	}

	headerRow := 5
	for i, header := range headers {
		col := string(rune('A' + i))
		if i >= 26 {
			col = string(rune('A'+(i/26-1))) + string(rune('A'+(i%26)))
		}
		cell := col + strconv.Itoa(headerRow)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Set row height for header
	f.SetRowHeight(sheetName, headerRow, 25)

	// Data rows
	row := headerRow + 1
	for i, p := range data {
		// No
		f.SetCellValue(sheetName, "A"+strconv.Itoa(row), i+1)
		f.SetCellStyle(sheetName, "A"+strconv.Itoa(row), "A"+strconv.Itoa(row), cellStyle)

		// Kode Peminjaman
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row), p.KodePeminjaman)
		f.SetCellStyle(sheetName, "B"+strconv.Itoa(row), "B"+strconv.Itoa(row), cellStyle)

		// Nama Peminjam
		namaPeminjam := ""
		if p.Peminjam != nil {
			namaPeminjam = p.Peminjam.Nama
		}
		f.SetCellValue(sheetName, "C"+strconv.Itoa(row), namaPeminjam)
		f.SetCellStyle(sheetName, "C"+strconv.Itoa(row), "C"+strconv.Itoa(row), cellStyle)

		// Organisasi
		organisasi := ""
		if p.Peminjam != nil && p.Peminjam.Organisasi != nil {
			organisasi = p.Peminjam.Organisasi.NamaOrganisasi
		}
		f.SetCellValue(sheetName, "D"+strconv.Itoa(row), organisasi)
		f.SetCellStyle(sheetName, "D"+strconv.Itoa(row), "D"+strconv.Itoa(row), cellStyle)

		// Nama Kegiatan
		namaKegiatan := ""
		if p.Kegiatan != nil {
			namaKegiatan = p.Kegiatan.NamaKegiatan
		}
		f.SetCellValue(sheetName, "E"+strconv.Itoa(row), namaKegiatan)
		f.SetCellStyle(sheetName, "E"+strconv.Itoa(row), "E"+strconv.Itoa(row), cellStyle)

		// Ruangan
		namaRuangan := "-"
		if p.Ruangan != nil {
			namaRuangan = p.Ruangan.NamaRuangan
		}
		f.SetCellValue(sheetName, "F"+strconv.Itoa(row), namaRuangan)
		f.SetCellStyle(sheetName, "F"+strconv.Itoa(row), "F"+strconv.Itoa(row), cellStyle)

		// Barang (list items)
		barangStr := "-"
		if len(p.Barang) > 0 {
			var barangList []string
			for _, item := range p.Barang {
				namaBarang := "Unknown"
				if item.Barang != nil {
					namaBarang = item.Barang.NamaBarang
				}
				barangList = append(barangList, fmt.Sprintf("%s (x%d)", namaBarang, item.Jumlah))
			}
			barangStr = strings.Join(barangList, "\n")
		}
		f.SetCellValue(sheetName, "G"+strconv.Itoa(row), barangStr)
		f.SetCellStyle(sheetName, "G"+strconv.Itoa(row), "G"+strconv.Itoa(row), cellStyle)

		// Tanggal Mulai
		tanggalMulai := p.TanggalMulai.Format("02/01/2006 15:04")
		f.SetCellValue(sheetName, "H"+strconv.Itoa(row), tanggalMulai)
		f.SetCellStyle(sheetName, "H"+strconv.Itoa(row), "H"+strconv.Itoa(row), cellStyle)

		// Tanggal Selesai
		tanggalSelesai := p.TanggalSelesai.Format("02/01/2006 15:04")
		f.SetCellValue(sheetName, "I"+strconv.Itoa(row), tanggalSelesai)
		f.SetCellStyle(sheetName, "I"+strconv.Itoa(row), "I"+strconv.Itoa(row), cellStyle)

		// Status - with colored styling
		statusCell := "J" + strconv.Itoa(row)
		f.SetCellValue(sheetName, statusCell, string(p.Status))
		
		switch p.Status {
		case models.StatusPeminjamanApproved:
			f.SetCellStyle(sheetName, statusCell, statusCell, approvedStyle)
		case models.StatusPeminjamanRejected:
			f.SetCellStyle(sheetName, statusCell, statusCell, rejectedStyle)
		case models.StatusPeminjamanPending:
			f.SetCellStyle(sheetName, statusCell, statusCell, pendingStyle)
		default:
			f.SetCellStyle(sheetName, statusCell, statusCell, cellStyle)
		}

		// Verifikator
		verifikator := "-"
		if p.Verifier != nil {
			verifikator = p.Verifier.Nama
		}
		f.SetCellValue(sheetName, "K"+strconv.Itoa(row), verifikator)
		f.SetCellStyle(sheetName, "K"+strconv.Itoa(row), "K"+strconv.Itoa(row), cellStyle)

		// Tanggal Verifikasi
		tanggalVerifikasi := "-"
		if p.VerifiedAt != nil {
			tanggalVerifikasi = p.VerifiedAt.Format("02/01/2006 15:04")
		}
		f.SetCellValue(sheetName, "L"+strconv.Itoa(row), tanggalVerifikasi)
		f.SetCellStyle(sheetName, "L"+strconv.Itoa(row), "L"+strconv.Itoa(row), cellStyle)

		// Catatan Verifikasi
		catatan := p.CatatanVerifikasi
		if catatan == "" {
			catatan = "-"
		}
		f.SetCellValue(sheetName, "M"+strconv.Itoa(row), catatan)
		f.SetCellStyle(sheetName, "M"+strconv.Itoa(row), "M"+strconv.Itoa(row), cellStyle)

		// Auto-adjust row height based on content
		if len(p.Barang) > 1 || len(catatan) > 50 {
			f.SetRowHeight(sheetName, row, float64(20+(len(p.Barang)*5)))
		}

		row++
	}

	// Auto-fit columns
	columns := []struct {
		col   string
		width float64
	}{
		{"A", 5},   // No
		{"B", 18},  // Kode Peminjaman
		{"C", 25},  // Nama Peminjam
		{"D", 25},  // Organisasi
		{"E", 30},  // Nama Kegiatan
		{"F", 20},  // Ruangan
		{"G", 30},  // Barang
		{"H", 18},  // Tanggal Mulai
		{"I", 18},  // Tanggal Selesai
		{"J", 12},  // Status
		{"K", 20},  // Verifikator
		{"L", 18},  // Tanggal Verifikasi
		{"M", 35},  // Catatan Verifikasi
	}

	for _, col := range columns {
		f.SetColWidth(sheetName, col.col, col.col, col.width)
	}

	// Freeze header row
	f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		XSplit:      0,
		YSplit:      headerRow,
		TopLeftCell: "A" + strconv.Itoa(headerRow+1),
		ActivePane:  "bottomLeft",
	})

	return f, nil
}
