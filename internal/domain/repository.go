package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int) (*User, error)
	GetBySSOSub(ctx context.Context, ssoSub string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	UpdateLastLogin(ctx context.Context, userID int, loginTime time.Time) error
	AssignRoles(ctx context.Context, userID int, roleIDs []int) error
	List(ctx context.Context, offset, limit int) ([]User, int64, error)
}

type CategoryRepository interface {
	List(ctx context.Context) ([]Category, error)
	GetByID(ctx context.Context, id int) (*Category, error)
	GetBySlug(ctx context.Context, slug string) (*Category, error)
	Create(ctx context.Context, category *Category) error
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id int) error
}

type AppFilter struct {
	CategoryID     *int
	TargetPengguna *string
	RoleIDs        []int
	StatusLayanan  *string
	AktifOnly      bool
	IsPublicOnly   bool
	UserID         *int // For determining is_favorite
	Offset         int
	Limit          int
}

type AppRepository interface {
	List(ctx context.Context, filter AppFilter) ([]App, int64, error)
	GetByID(ctx context.Context, id int, userID *int) (*App, error)
	GetBySlug(ctx context.Context, slug string, userID *int) (*App, error)
	Search(ctx context.Context, query string, categoryID *int, roleIDs []int, isPublicOnly bool, userID *int, limit int) ([]App, error)
	Create(ctx context.Context, app *App, roleIDs []int) error
	Update(ctx context.Context, app *App, roleIDs []int) error
	UpdateStatus(ctx context.Context, id int, status string) error
	Delete(ctx context.Context, id int) error // Soft delete
	Reorder(ctx context.Context, appIDs []int) error
}

type FavoriteRepository interface {
	Add(ctx context.Context, userID, appID int) error
	Remove(ctx context.Context, userID, appID int) error
	ListByUser(ctx context.Context, userID int) ([]App, error)
	IsFavorite(ctx context.Context, userID, appID int) (bool, error)
}

type ClickLogRepository interface {
	Record(ctx context.Context, log *ClickLog) error
	GetTopApps(ctx context.Context, startDate, endDate time.Time, limit int) ([]map[string]interface{}, error)
	GetRecentByUser(ctx context.Context, userID int, limit int) ([]App, error)
	GetSummary(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error)
}

type StatusCheckRepository interface {
	Record(ctx context.Context, check *StatusCheck) error
	GetRecentByApp(ctx context.Context, appID int, limit int) ([]StatusCheck, error)
	GetUptimeSummary(ctx context.Context, days int) ([]map[string]interface{}, error)
}

type AnnouncementRepository interface {
	ListActive(ctx context.Context) ([]Announcement, error)
	ListAll(ctx context.Context, offset, limit int) ([]Announcement, int64, error)
	Create(ctx context.Context, announcement *Announcement) error
	Update(ctx context.Context, announcement *Announcement) error
	Delete(ctx context.Context, id int) error
}

type FeedbackRepository interface {
	Create(ctx context.Context, feedback *Feedback) error
	List(ctx context.Context, status *string, offset, limit int) ([]Feedback, int64, error)
	UpdateStatus(ctx context.Context, id int, status string, adminNotes *string) error
}

type AuditLogRepository interface {
	Record(ctx context.Context, log *AuditLog) error
	List(ctx context.Context, entity *string, offset, limit int) ([]AuditLog, int64, error)
}
