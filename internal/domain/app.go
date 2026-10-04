package domain

import "time"

type App struct {
	ID             int       `json:"id"`
	CategoryID     int       `json:"category_id"`
	Category       *Category `json:"category,omitempty"`
	Nama           string    `json:"nama"`
	Slug           string    `json:"slug"`
	URL            string    `json:"url"`
	Deskripsi      *string   `json:"deskripsi,omitempty"`
	IkonURL        *string   `json:"ikon_url,omitempty"`
	TargetPengguna string    `json:"target_pengguna"`
	Pemilik        string    `json:"pemilik"`
	KontakAdmin    *string   `json:"kontak_admin,omitempty"`
	Urutan         int       `json:"urutan"`
	StatusLayanan  string    `json:"status_layanan"` // online, degraded, down, maintenance
	IsPublic       bool      `json:"is_public"`
	Aktif          bool      `json:"aktif"`
	IsFavorite     bool      `json:"is_favorite"`
	TotalClicks    int64     `json:"total_clicks,omitempty"`
	AllowedRoles   []Role    `json:"allowed_roles,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Guide struct {
	ID        int       `json:"id"`
	AppID     int       `json:"app_id"`
	Judul     string    `json:"judul"`
	Konten    string    `json:"konten"`
	Urutan    int       `json:"urutan"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
