package usecase

import (
	"context"

	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/repository/postgres"
	"backend-alusi-go/pkg/exporter"
)

type OpenAPIMetadataItem struct {
	Identifier     string `json:"identifier"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	LandingPage    string `json:"landing_page"`
	Theme          string `json:"theme"`
	Publisher      string `json:"publisher"`
	AccessRights   string `json:"access_rights"`
	TargetAudience string `json:"target_audience"`
	ServiceStatus  string `json:"service_status"`
	Issued         string `json:"issued"`
	Modified       string `json:"modified"`
}

type OpenAPICatalogResponse struct {
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Publisher   string                `json:"publisher"`
	License     string                `json:"license"`
	TotalItems  int                   `json:"total_items"`
	Dataset     []OpenAPIMetadataItem `json:"dataset"`
}

type ReportUsecase struct {
	appRepo       *postgres.AppRepo
	analyticsRepo *postgres.AnalyticsRepo
	exporter      *exporter.ReportExporter
}

func NewReportUsecase(
	appRepo *postgres.AppRepo,
	analyticsRepo *postgres.AnalyticsRepo,
	exporter *exporter.ReportExporter,
) *ReportUsecase {
	return &ReportUsecase{
		appRepo:       appRepo,
		analyticsRepo: analyticsRepo,
		exporter:      exporter,
	}
}

func (u *ReportUsecase) GenerateCatalogCSV(ctx context.Context) ([]byte, error) {
	var apps []domain.App
	if u.appRepo != nil {
		list, _, err := u.appRepo.List(ctx, domain.AppFilter{AktifOnly: false, Limit: 1000})
		if err == nil {
			apps = list
		}
	}
	return u.exporter.GenerateCatalogCSV(apps)
}

func (u *ReportUsecase) GenerateAnalyticsCSV(ctx context.Context) ([]byte, error) {
	var summary *postgres.DashboardSummary
	var topApps []postgres.TopAppMetric
	var disruptions []postgres.DisruptionSummary

	if u.analyticsRepo != nil {
		summary, _ = u.analyticsRepo.GetDashboardSummary(ctx)
		topApps, _ = u.analyticsRepo.GetTopApps(ctx, 30, 20)
		disruptions, _ = u.analyticsRepo.GetDisruptionSummary(ctx, 30)
	}

	return u.exporter.GenerateAnalyticsCSV(summary, topApps, disruptions)
}

func (u *ReportUsecase) GetOpenAPICatalog(ctx context.Context) (*OpenAPICatalogResponse, error) {
	var apps []domain.App
	if u.appRepo != nil {
		list, _, err := u.appRepo.List(ctx, domain.AppFilter{AktifOnly: true, IsPublicOnly: true, Limit: 500})
		if err == nil {
			apps = list
		}
	}

	var items []OpenAPIMetadataItem
	for _, a := range apps {
		theme := "Umum"
		if a.Category != nil {
			theme = a.Category.Nama
		}

		desc := ""
		if a.Deskripsi != nil {
			desc = *a.Deskripsi
		}

		items = append(items, OpenAPIMetadataItem{
			Identifier:     a.Slug,
			Title:          a.Nama,
			Description:    desc,
			LandingPage:    a.URL,
			Theme:          theme,
			Publisher:      a.Pemilik,
			AccessRights:   "public",
			TargetAudience: a.TargetPengguna,
			ServiceStatus:  a.StatusLayanan,
			Issued:         a.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Modified:       a.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	res := &OpenAPICatalogResponse{
		Title:       "Katalog Terbuka Layanan Aplikasi BPS Provinsi Sumatera Utara",
		Description: "Metadata portal aplikasi resmi BPS Provinsi Sumatera Utara sesuai standar Satu Data Indonesia",
		Publisher:   "Badan Pusat Statistik Provinsi Sumatera Utara",
		License:     "https://opendatacommons.org/licenses/pddl/1-0/",
		TotalItems:  len(items),
		Dataset:     items,
	}

	return res, nil
}
