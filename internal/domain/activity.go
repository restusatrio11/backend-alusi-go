package domain

import "time"

type ClickLog struct {
	ID        int64     `json:"id"`
	UserID    *int      `json:"user_id,omitempty"`
	AppID     int       `json:"app_id"`
	SatkerID  *int      `json:"satker_id,omitempty"`
	IPAddress *string   `json:"ip_address,omitempty"`
	UserAgent *string   `json:"user_agent,omitempty"`
	ClickedAt time.Time `json:"clicked_at"`
}

type StatusCheck struct {
	ID             int64     `json:"id"`
	AppID          int       `json:"app_id"`
	Status         string    `json:"status"`
	HTTPStatusCode *int      `json:"http_status_code,omitempty"`
	LatencyMS      int       `json:"latency_ms"`
	ErrorMessage   *string   `json:"error_message,omitempty"`
	CheckedAt      time.Time `json:"checked_at"`
}

type Announcement struct {
	ID        int        `json:"id"`
	Judul     string     `json:"judul"`
	Konten    string     `json:"konten"`
	Tipe      string     `json:"tipe"` // info, warning, maintenance, urgent
	AppID     *int       `json:"app_id,omitempty"`
	App       *App       `json:"app,omitempty"`
	IsActive  bool       `json:"is_active"`
	StartsAt  time.Time  `json:"starts_at"`
	EndsAt    *time.Time `json:"ends_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Feedback struct {
	ID           int       `json:"id"`
	UserID       *int      `json:"user_id,omitempty"`
	User         *User     `json:"user,omitempty"`
	AppID        *int      `json:"app_id,omitempty"`
	App          *App      `json:"app,omitempty"`
	Kategori     string    `json:"kategori"` // kendala, saran, data
	Pesan        string    `json:"pesan"`
	Status       string    `json:"status"` // pending, in_progress, resolved, closed
	CatatanAdmin *string   `json:"catatan_admin,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AuditLog struct {
	ID            int64                  `json:"id"`
	UserID        *int                   `json:"user_id,omitempty"`
	User          *User                  `json:"user,omitempty"`
	Action        string                 `json:"action"` // CREATE, UPDATE, DELETE, LOGIN, LOGOUT
	Entity        string                 `json:"entity"`
	EntityID      *string                `json:"entity_id,omitempty"`
	PayloadBefore map[string]interface{} `json:"payload_before,omitempty"`
	PayloadAfter  map[string]interface{} `json:"payload_after,omitempty"`
	IPAddress     *string                `json:"ip_address,omitempty"`
	UserAgent     *string                `json:"user_agent,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
}

type AuditLogFilter struct {
	UserID    *int
	Action    string
	Entity    string
	EntityID  *string
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	PerPage   int
}
