package exporter

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/repository/postgres"
)

type ReportExporter struct{}

func NewReportExporter() *ReportExporter {
	return &ReportExporter{}
}

// GenerateCatalogCSV exports full application catalog into CSV format
func (e *ReportExporter) GenerateCatalogCSV(apps []domain.App) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// CSV Header
	header := []string{
		"ID", "Nama Aplikasi", "Slug", "Kategori", "URL", "Target Pengguna",
		"Pemilik", "Kontak Admin", "Status Layanan", "Publik", "Aktif", "Total Klik", "Tanggal Dibuat",
	}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	for _, app := range apps {
		catName := ""
		if app.Category != nil {
			catName = app.Category.Nama
		}

		kontak := ""
		if app.KontakAdmin != nil {
			kontak = *app.KontakAdmin
		}

		row := []string{
			fmt.Sprintf("%d", app.ID),
			app.Nama,
			app.Slug,
			catName,
			app.URL,
			app.TargetPengguna,
			app.Pemilik,
			kontak,
			app.StatusLayanan,
			fmt.Sprintf("%t", app.IsPublic),
			fmt.Sprintf("%t", app.Aktif),
			fmt.Sprintf("%d", app.TotalClicks),
			app.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), writer.Error()
}

// GenerateAnalyticsCSV exports executive usage and uptime rollup to CSV format
func (e *ReportExporter) GenerateAnalyticsCSV(
	summary *postgres.DashboardSummary,
	topApps []postgres.TopAppMetric,
	disruptions []postgres.DisruptionSummary,
) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// Section 1: Executive KPI Summary
	_ = writer.Write([]string{"=== RINGKASAN EKSEKUTIF PENGGUNAAN PORTAL ALUSI BPS SUMUT ==="})
	_ = writer.Write([]string{"Tanggal Laporan", time.Now().Format("2006-01-02 15:04:05")})
	_ = writer.Write([]string{""})
	_ = writer.Write([]string{"Metrik", "Nilai"})
	if summary != nil {
		_ = writer.Write([]string{"Total Aplikasi Aktif", fmt.Sprintf("%d", summary.TotalApps)})
		_ = writer.Write([]string{"Total Kategori", fmt.Sprintf("%d", summary.TotalCategories)})
		_ = writer.Write([]string{"Total Pengguna Terdaftar", fmt.Sprintf("%d", summary.TotalUsers)})
		_ = writer.Write([]string{"Total Klik Hari Ini", fmt.Sprintf("%d", summary.TotalClicksToday)})
		_ = writer.Write([]string{"Total Klik Bulan Ini", fmt.Sprintf("%d", summary.TotalClicksThisMonth)})
		_ = writer.Write([]string{"Daily Active Users (DAU)", fmt.Sprintf("%d", summary.DAU)})
		_ = writer.Write([]string{"Monthly Active Users (MAU)", fmt.Sprintf("%d", summary.MAU)})
		_ = writer.Write([]string{"Laporan Kendala Pending", fmt.Sprintf("%d", summary.PendingFeedbacksCount)})
	}
	_ = writer.Write([]string{""})

	// Section 2: Top Applications
	_ = writer.Write([]string{"=== RANKING APLIKASI TERPOPULER (30 HARI TERAKHIR) ==="})
	_ = writer.Write([]string{"Peringkat", "Nama Aplikasi", "Kategori", "Total Klik", "Pengguna Unik"})
	for i, app := range topApps {
		_ = writer.Write([]string{
			fmt.Sprintf("%d", i+1),
			app.Nama,
			app.Category,
			fmt.Sprintf("%d", app.TotalClicks),
			fmt.Sprintf("%d", app.UniqueUsers),
		})
	}
	_ = writer.Write([]string{""})

	// Section 3: SLA Uptime & Disruptions
	_ = writer.Write([]string{"=== MONITORING SLA & RIWAYAT GANGGUAN (30 HARI TERAKHIR) ==="})
	_ = writer.Write([]string{"Nama Aplikasi", "Total Pengecekan", "Jumlah Down", "Jumlah Degraded", "Persentase Uptime SLA (%)", "Terakhir Down"})
	for _, d := range disruptions {
		lastDown := "-"
		if d.LastDownAt != nil {
			lastDown = d.LastDownAt.Format("2006-01-02 15:04:05")
		}
		_ = writer.Write([]string{
			d.AppNama,
			fmt.Sprintf("%d", d.TotalChecks),
			fmt.Sprintf("%d", d.DownChecks),
			fmt.Sprintf("%d", d.DegradedChecks),
			fmt.Sprintf("%.2f%%", d.UptimePercentage),
			lastDown,
		})
	}

	writer.Flush()
	return buf.Bytes(), writer.Error()
}
