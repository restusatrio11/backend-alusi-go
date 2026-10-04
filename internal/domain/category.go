package domain

import "time"

type Category struct {
	ID        int       `json:"id"`
	Nama      string    `json:"nama"`
	Slug      string    `json:"slug"`
	Deskripsi *string   `json:"deskripsi,omitempty"`
	Urutan    int       `json:"urutan"`
	Ikon      *string   `json:"ikon,omitempty"`
	AppCount  int       `json:"app_count,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
