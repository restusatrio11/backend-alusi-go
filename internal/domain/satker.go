package domain

import "time"

type Satker struct {
	ID        int       `json:"id"`
	Kode      string    `json:"kode"`
	Nama      string    `json:"nama"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Role struct {
	ID        int       `json:"id"`
	Nama      string    `json:"nama"`
	Deskripsi string    `json:"deskripsi,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
