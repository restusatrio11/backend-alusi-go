package domain

import "time"

type Permission struct {
	ID        int       `json:"id"`
	Kode      string    `json:"kode"`
	Nama      string    `json:"nama"`
	Modul     string    `json:"modul"`
	Deskripsi string    `json:"deskripsi"`
	CreatedAt time.Time `json:"created_at"`
}

type RoleWithPermissions struct {
	ID          int          `json:"id"`
	Nama        string       `json:"nama"`
	Deskripsi   string       `json:"deskripsi"`
	Permissions []Permission `json:"permissions"`
	CreatedAt   time.Time    `json:"created_at"`
}

type UserWithRoles struct {
	ID          int        `json:"id"`
	Username    *string    `json:"username,omitempty"`
	NIP         *string    `json:"nip,omitempty"`
	Nama        string     `json:"nama"`
	Email       string     `json:"email"`
	UserType    string     `json:"user_type"`
	Status      string     `json:"status"`
	Satker      *Satker    `json:"satker,omitempty"`
	Roles       []Role     `json:"roles"`
	Permissions []string   `json:"permissions"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
