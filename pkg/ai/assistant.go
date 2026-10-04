package ai

import (
	"context"
	"fmt"
	"strings"

	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/repository/postgres"
)

type Recommendation struct {
	App             domain.App `json:"app"`
	RelevanceScore  float64    `json:"relevance_score"`
	Recommendation  string     `json:"recommendation"`
	SuggestedAction string     `json:"suggested_action"`
}

type AIResponse struct {
	Question        string           `json:"question"`
	Answer          string           `json:"answer"`
	Recommendations []Recommendation `json:"recommendations"`
	DetectedIntent  string           `json:"detected_intent"`
}

type AssistantService struct {
	appRepo      *postgres.AppRepo
	categoryRepo *postgres.CategoryRepo
	guideRepo    *postgres.GuideRepo
}

func NewAssistantService(
	appRepo *postgres.AppRepo,
	categoryRepo *postgres.CategoryRepo,
	guideRepo *postgres.GuideRepo,
) *AssistantService {
	return &AssistantService{
		appRepo:      appRepo,
		categoryRepo: categoryRepo,
		guideRepo:    guideRepo,
	}
}

// Ask processes user question and returns smart app recommendation
func (s *AssistantService) Ask(
	ctx context.Context,
	question string,
	userID *int,
	roleIDs []int,
	isPublicOnly bool,
) (*AIResponse, error) {
	q := strings.TrimSpace(question)
	if q == "" {
		return nil, fmt.Errorf("pertanyaan tidak boleh kosong")
	}

	qLower := strings.ToLower(q)

	// Intent detection
	intent := "general_inquiry"
	if strings.Contains(qLower, "cuti") || strings.Contains(qLower, "presensi") || strings.Contains(qLower, "absen") || strings.Contains(qLower, "kepegawaian") || strings.Contains(qLower, "kinerja") || strings.Contains(qLower, "skp") {
		intent = "kepegawaian_sdm"
	} else if strings.Contains(qLower, "survei") || strings.Contains(qLower, "sensus") || strings.Contains(qLower, "lapangan") || strings.Contains(qLower, "pendataan") || strings.Contains(qLower, "entri") || strings.Contains(qLower, "kuesioner") {
		intent = "survei_sensus"
	} else if strings.Contains(qLower, "data") || strings.Contains(qLower, "publikasi") || strings.Contains(qLower, "tabel") || strings.Contains(qLower, "rilis") || strings.Contains(qLower, "inflasi") || strings.Contains(qLower, "pdrb") {
		intent = "diseminasi_publikasi"
	} else if strings.Contains(qLower, "anggaran") || strings.Contains(qLower, "keuangan") || strings.Contains(qLower, "spk") || strings.Contains(qLower, "honor") || strings.Contains(qLower, "pembayaran") || strings.Contains(qLower, "spd") || strings.Contains(qLower, "perjadin") {
		intent = "keuangan_pengadaan"
	} else if strings.Contains(qLower, "surat") || strings.Contains(qLower, "naskah") || strings.Contains(qLower, "arsip") || strings.Contains(qLower, "internal") || strings.Contains(qLower, "rapat") {
		intent = "tata_kelola_internal"
	}

	var matchingApps []domain.App
	if s.appRepo != nil {
		// 1. Search by full text & trigram similarity
		apps, err := s.appRepo.Search(ctx, q, nil, roleIDs, isPublicOnly, userID, 5)
		if err == nil {
			matchingApps = apps
		}
	}

	// Prepare recommendations
	var recs []Recommendation
	for i, app := range matchingApps {
		score := 0.95 - (float64(i) * 0.1)
		if score < 0.5 {
			score = 0.5
		}

		explanation := fmt.Sprintf("Aplikasi %s dapat digunakan untuk kebutuhan terkait %s.", app.Nama, app.TargetPengguna)
		if app.Deskripsi != nil && *app.Deskripsi != "" {
			explanation = *app.Deskripsi
		}

		action := fmt.Sprintf("Buka aplikasi %s via %s", app.Nama, app.URL)

		recs = append(recs, Recommendation{
			App:             app,
			RelevanceScore:  score,
			Recommendation:  explanation,
			SuggestedAction: action,
		})
	}

	// Formulate natural language answer
	answer := s.synthesizeAnswer(q, intent, recs)

	return &AIResponse{
		Question:        q,
		Answer:          answer,
		Recommendations: recs,
		DetectedIntent:  intent,
	}, nil
}

func (s *AssistantService) synthesizeAnswer(question string, intent string, recs []Recommendation) string {
	if len(recs) == 0 {
		return fmt.Sprintf("Berdasarkan pertanyaan Anda (\"%s\"), saat ini belum ditemukan aplikasi atau layanan internal yang spesifik sesuai kata kunci tersebut. Silakan periksa katalog lengkap aplikasi atau hubungi tim pengelola IT BPS Provinsi Sumatera Utara.", question)
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Halo! Untuk kebutuhan Anda mengenai \"%s\", kami merekomendasikan layanan berikut:\n\n", question))

	for idx, r := range recs {
		b.WriteString(fmt.Sprintf("%d. **%s** (%s)\n", idx+1, r.App.Nama, r.App.TargetPengguna))
		b.WriteString(fmt.Sprintf("   %s\n", r.Recommendation))
		b.WriteString(fmt.Sprintf("   Tautan: %s\n\n", r.App.URL))
	}

	b.WriteString("Silakan klik aplikasi di atas untuk langsung diarahkan ke halaman layanan terkait.")
	return b.String()
}
