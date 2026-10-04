package domain

import "time"

type User struct {
	ID          int        `json:"id"`
	SSOSub      string     `json:"sso_sub"`
	NIP         *string    `json:"nip,omitempty"`
	Nama        string     `json:"nama"`
	Email       string     `json:"email"`
	SatkerID    *int       `json:"satker_id,omitempty"`
	Satker      *Satker    `json:"satker,omitempty"`
	Roles       []Role     `json:"roles,omitempty"`
	Status      string     `json:"status"` // active, inactive
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (u *User) HasRole(roleName string) bool {
	for _, r := range u.Roles {
		if r.Nama == roleName {
			return true
		}
	}
	return false
}

func (u *User) IsAdmin() bool {
	return u.HasRole("admin")
}
