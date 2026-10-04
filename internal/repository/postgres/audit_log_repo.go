package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"backend-alusi-go/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditLogRepo struct {
	pool *pgxpool.Pool
}

func NewAuditLogRepo(pool *pgxpool.Pool) *AuditLogRepo {
	return &AuditLogRepo{pool: pool}
}

func (r *AuditLogRepo) Record(ctx context.Context, log *domain.AuditLog) error {
	var payloadBeforeJSON []byte
	var payloadAfterJSON []byte
	var err error

	if log.PayloadBefore != nil {
		payloadBeforeJSON, err = json.Marshal(log.PayloadBefore)
		if err != nil {
			payloadBeforeJSON = nil
		}
	}

	if log.PayloadAfter != nil {
		payloadAfterJSON, err = json.Marshal(log.PayloadAfter)
		if err != nil {
			payloadAfterJSON = nil
		}
	}

	query := `
	INSERT INTO audit_logs (
		user_id, action, entity, entity_id, payload_before, payload_after, ip_address, user_agent
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, created_at
	`
	return r.pool.QueryRow(
		ctx, query,
		log.UserID, log.Action, log.Entity, log.EntityID,
		payloadBeforeJSON, payloadAfterJSON, log.IPAddress, log.UserAgent,
	).Scan(&log.ID, &log.CreatedAt)
}

func (r *AuditLogRepo) List(ctx context.Context, filter domain.AuditLogFilter) ([]domain.AuditLog, int64, error) {
	var whereClauses []string
	var args []interface{}
	argIdx := 1

	if filter.UserID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("al.user_id = $%d", argIdx))
		args = append(args, *filter.UserID)
		argIdx++
	}

	if filter.Action != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("al.action = $%d", argIdx))
		args = append(args, filter.Action)
		argIdx++
	}

	if filter.Entity != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("al.entity = $%d", argIdx))
		args = append(args, filter.Entity)
		argIdx++
	}

	if filter.EntityID != nil && *filter.EntityID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("al.entity_id = $%d", argIdx))
		args = append(args, *filter.EntityID)
		argIdx++
	}

	if filter.StartDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("al.created_at >= $%d", argIdx))
		args = append(args, *filter.StartDate)
		argIdx++
	}

	if filter.EndDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("al.created_at <= $%d", argIdx))
		args = append(args, *filter.EndDate)
		argIdx++
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM audit_logs al %s", whereSQL)
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	limit := filter.PerPage
	if limit <= 0 {
		limit = 20
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	limitSQL := fmt.Sprintf("LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	query := fmt.Sprintf(`
	SELECT 
		al.id, al.user_id, al.action, al.entity, al.entity_id,
		al.payload_before, al.payload_after, al.ip_address, al.user_agent, al.created_at,
		u.id, u.nama, u.email
	FROM audit_logs al
	LEFT JOIN users u ON al.user_id = u.id
	%s
	ORDER BY al.created_at DESC
	%s
	`, whereSQL, limitSQL)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list audit logs: %w", err)
	}
	defer rows.Close()

	var list []domain.AuditLog
	for rows.Next() {
		log, err := r.scanAuditLogRow(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, *log)
	}

	return list, total, nil
}

func (r *AuditLogRepo) scanAuditLogRow(rows pgx.Rows) (*domain.AuditLog, error) {
	var a domain.AuditLog
	var beforeBytes, afterBytes []byte
	var uID *int
	var uNama, uEmail *string

	err := rows.Scan(
		&a.ID, &a.UserID, &a.Action, &a.Entity, &a.EntityID,
		&beforeBytes, &afterBytes, &a.IPAddress, &a.UserAgent, &a.CreatedAt,
		&uID, &uNama, &uEmail,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan audit log row: %w", err)
	}

	if len(beforeBytes) > 0 {
		var before map[string]interface{}
		if err := json.Unmarshal(beforeBytes, &before); err == nil {
			a.PayloadBefore = before
		}
	}

	if len(afterBytes) > 0 {
		var after map[string]interface{}
		if err := json.Unmarshal(afterBytes, &after); err == nil {
			a.PayloadAfter = after
		}
	}

	if uID != nil && uNama != nil && uEmail != nil {
		a.User = &domain.User{
			ID:    *uID,
			Nama:  *uNama,
			Email: *uEmail,
		}
	}

	return &a, nil
}
