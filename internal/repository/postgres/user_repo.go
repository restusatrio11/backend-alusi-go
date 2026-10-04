package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"backend-alusi-go/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `
	SELECT 
		u.id, COALESCE(u.sso_sub, ''), u.username, u.password_hash, u.user_type, u.nip, u.nik, u.nama, u.email, u.satker_id, 
		u.status, u.metadata, u.last_login_at, u.created_at, u.updated_at,
		s.id, s.kode, s.nama, s.created_at, s.updated_at
	FROM users u
	LEFT JOIN satker s ON u.satker_id = s.id
	WHERE u.id = $1
	`
	return r.querySingleUser(ctx, query, id)
}

func (r *UserRepo) GetBySSOSub(ctx context.Context, ssoSub string) (*domain.User, error) {
	query := `
	SELECT 
		u.id, COALESCE(u.sso_sub, ''), u.username, u.password_hash, u.user_type, u.nip, u.nik, u.nama, u.email, u.satker_id, 
		u.status, u.metadata, u.last_login_at, u.created_at, u.updated_at,
		s.id, s.kode, s.nama, s.created_at, s.updated_at
	FROM users u
	LEFT JOIN satker s ON u.satker_id = s.id
	WHERE u.sso_sub = $1
	`
	return r.querySingleUser(ctx, query, ssoSub)
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
	SELECT 
		u.id, COALESCE(u.sso_sub, ''), u.username, u.password_hash, u.user_type, u.nip, u.nik, u.nama, u.email, u.satker_id, 
		u.status, u.metadata, u.last_login_at, u.created_at, u.updated_at,
		s.id, s.kode, s.nama, s.created_at, s.updated_at
	FROM users u
	LEFT JOIN satker s ON u.satker_id = s.id
	WHERE u.email = $1
	`
	return r.querySingleUser(ctx, query, email)
}

func (r *UserRepo) GetByUsernameOrEmailOrNIP(ctx context.Context, identifier string) (*domain.User, error) {
	query := `
	SELECT 
		u.id, COALESCE(u.sso_sub, ''), u.username, u.password_hash, u.user_type, u.nip, u.nik, u.nama, u.email, u.satker_id, 
		u.status, u.metadata, u.last_login_at, u.created_at, u.updated_at,
		s.id, s.kode, s.nama, s.created_at, s.updated_at
	FROM users u
	LEFT JOIN satker s ON u.satker_id = s.id
	WHERE LOWER(u.username) = LOWER($1) OR LOWER(u.email) = LOWER($1) OR u.nip = $1
	LIMIT 1
	`
	return r.querySingleUser(ctx, query, identifier)
}

func (r *UserRepo) GetSatkerByKode(ctx context.Context, kode string) (*domain.Satker, error) {
	query := `SELECT id, kode, nama, created_at, updated_at FROM satker WHERE kode = $1`
	var s domain.Satker
	err := r.pool.QueryRow(ctx, query, kode).Scan(&s.ID, &s.Kode, &s.Nama, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get satker by kode %s: %w", kode, err)
	}
	return &s, nil
}

func (r *UserRepo) GetRoleByNama(ctx context.Context, nama string) (*domain.Role, error) {
	query := `SELECT id, nama, deskripsi, created_at FROM roles WHERE nama = $1`
	var role domain.Role
	err := r.pool.QueryRow(ctx, query, nama).Scan(&role.ID, &role.Nama, &role.Deskripsi, &role.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get role by nama %s: %w", nama, err)
	}
	return &role, nil
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	query := `
	INSERT INTO users (sso_sub, username, password_hash, user_type, nip, nik, nama, email, satker_id, status, metadata, last_login_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING id, created_at, updated_at
	`
	metadataBytes, _ := json.Marshal(user.Metadata)

	var ssoSubVal *string
	if user.SSOSub != "" {
		ssoSubVal = &user.SSOSub
	}

	err := r.pool.QueryRow(
		ctx,
		query,
		ssoSubVal,
		user.Username,
		user.PasswordHash,
		user.UserType,
		user.NIP,
		user.NIK,
		user.Nama,
		user.Email,
		user.SatkerID,
		user.Status,
		metadataBytes,
		user.LastLoginAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	query := `
	UPDATE users
	SET 
		username = $1,
		password_hash = COALESCE($2, password_hash),
		user_type = $3,
		nip = $4,
		nik = $5,
		nama = $6,
		email = $7,
		satker_id = $8,
		status = $9,
		metadata = $10,
		updated_at = NOW()
	WHERE id = $11
	RETURNING updated_at
	`
	metadataBytes, _ := json.Marshal(user.Metadata)

	err := r.pool.QueryRow(
		ctx,
		query,
		user.Username,
		user.PasswordHash,
		user.UserType,
		user.NIP,
		user.NIK,
		user.Nama,
		user.Email,
		user.SatkerID,
		user.Status,
		metadataBytes,
		user.ID,
	).Scan(&user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepo) UpdateLastLogin(ctx context.Context, userID int, loginTime time.Time) error {
	query := `UPDATE users SET last_login_at = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, loginTime, userID)
	return err
}

func (r *UserRepo) AssignRoles(ctx context.Context, userID int, roleIDs []int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Clear existing roles
	if _, err := tx.Exec(ctx, "DELETE FROM user_roles WHERE user_id = $1", userID); err != nil {
		return err
	}

	// Insert new roles
	for _, roleID := range roleIDs {
		if _, err := tx.Exec(ctx, "INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, roleID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *UserRepo) List(ctx context.Context, offset, limit int) ([]domain.User, int64, error) {
	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
	SELECT 
		u.id, COALESCE(u.sso_sub, ''), u.username, u.password_hash, u.user_type, u.nip, u.nik, u.nama, u.email, u.satker_id, 
		u.status, u.metadata, u.last_login_at, u.created_at, u.updated_at,
		s.id, s.kode, s.nama, s.created_at, s.updated_at
	FROM users u
	LEFT JOIN satker s ON u.satker_id = s.id
	ORDER BY u.id DESC
	LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		u, err := r.scanUserWithSatker(rows)
		if err != nil {
			return nil, 0, err
		}
		// Load roles
		roles, err := r.getUserRoles(ctx, u.ID)
		if err == nil {
			u.Roles = roles
		}
		users = append(users, *u)
	}

	return users, total, nil
}

func (r *UserRepo) querySingleUser(ctx context.Context, query string, arg interface{}) (*domain.User, error) {
	rows, err := r.pool.Query(ctx, query, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	u, err := r.scanUserWithSatker(rows)
	if err != nil {
		return nil, err
	}

	roles, err := r.getUserRoles(ctx, u.ID)
	if err == nil {
		u.Roles = roles
	}

	return u, nil
}

func (r *UserRepo) scanUserWithSatker(rows pgx.Rows) (*domain.User, error) {
	var u domain.User
	var metadataBytes []byte
	var satkerID *int
	var satkerKode, satkerNama *string
	var satkerCreated, satkerUpdated *time.Time

	err := rows.Scan(
		&u.ID, &u.SSOSub, &u.Username, &u.PasswordHash, &u.UserType, &u.NIP, &u.NIK, &u.Nama, &u.Email, &u.SatkerID,
		&u.Status, &metadataBytes, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
		&satkerID, &satkerKode, &satkerNama, &satkerCreated, &satkerUpdated,
	)
	if err != nil {
		return nil, err
	}

	if len(metadataBytes) > 0 {
		_ = json.Unmarshal(metadataBytes, &u.Metadata)
	}

	if satkerID != nil && satkerKode != nil && satkerNama != nil {
		u.Satker = &domain.Satker{
			ID:        *satkerID,
			Kode:      *satkerKode,
			Nama:      *satkerNama,
			CreatedAt: *satkerCreated,
			UpdatedAt: *satkerUpdated,
		}
	}

	return &u, nil
}

func (r *UserRepo) getUserRoles(ctx context.Context, userID int) ([]domain.Role, error) {
	query := `
	SELECT r.id, r.nama, r.deskripsi, r.created_at
	FROM roles r
	INNER JOIN user_roles ur ON r.id = ur.role_id
	WHERE ur.user_id = $1
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.ID, &role.Nama, &role.Deskripsi, &role.CreatedAt); err == nil {
			roles = append(roles, role)
		}
	}
	return roles, nil
}
