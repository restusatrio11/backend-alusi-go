package postgres

import (
	"context"
	"fmt"
	"strings"

	"backend-alusi-go/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppRepo struct {
	pool *pgxpool.Pool
}

func NewAppRepo(pool *pgxpool.Pool) *AppRepo {
	return &AppRepo{pool: pool}
}

func (r *AppRepo) List(ctx context.Context, filter domain.AppFilter) ([]domain.App, int64, error) {
	var whereClauses []string
	var args []interface{}
	argIdx := 1

	if filter.AktifOnly {
		whereClauses = append(whereClauses, fmt.Sprintf("a.aktif = $%d", argIdx))
		args = append(args, true)
		argIdx++
	}

	if filter.IsPublicOnly {
		whereClauses = append(whereClauses, fmt.Sprintf("a.is_public = $%d", argIdx))
		args = append(args, true)
		argIdx++
	}

	if filter.CategoryID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("a.category_id = $%d", argIdx))
		args = append(args, *filter.CategoryID)
		argIdx++
	}

	if filter.TargetPengguna != nil && *filter.TargetPengguna != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(a.target_pengguna = $%d OR a.target_pengguna = 'Semua')", argIdx))
		args = append(args, *filter.TargetPengguna)
		argIdx++
	}

	if filter.StatusLayanan != nil && *filter.StatusLayanan != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("a.status_layanan = $%d", argIdx))
		args = append(args, *filter.StatusLayanan)
		argIdx++
	}

	// Role-based visibility check: App is public OR has no role restrictions OR user has an allowed role
	if len(filter.RoleIDs) > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf(`(
			a.is_public = true 
			OR NOT EXISTS (SELECT 1 FROM app_access WHERE app_id = a.id)
			OR EXISTS (SELECT 1 FROM app_access aa WHERE aa.app_id = a.id AND aa.role_id = ANY($%d))
		)`, argIdx))
		args = append(args, filter.RoleIDs)
		argIdx++
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM apps a %s", whereSQL)
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count apps: %w", err)
	}

	// Fetch apps with Category and is_favorite
	var favSelect string
	var favJoin string
	if filter.UserID != nil {
		favSelect = fmt.Sprintf(", (f.user_id IS NOT NULL) as is_favorite")
		favJoin = fmt.Sprintf("LEFT JOIN favorites f ON a.id = f.app_id AND f.user_id = $%d", argIdx)
		args = append(args, *filter.UserID)
		argIdx++
	} else {
		favSelect = ", false as is_favorite"
		favJoin = ""
	}

	limitSQL := fmt.Sprintf("LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, filter.Limit, filter.Offset)

	query := fmt.Sprintf(`
	SELECT 
		a.id, a.category_id, a.nama, a.slug, a.url, a.deskripsi, a.ikon_url,
		a.target_pengguna, a.pemilik, a.kontak_admin, a.urutan, a.status_layanan,
		a.is_public, a.aktif, a.created_at, a.updated_at,
		c.id, c.nama, c.slug, c.deskripsi, c.urutan, c.ikon, c.created_at, c.updated_at
		%s
	FROM apps a
	INNER JOIN categories c ON a.category_id = c.id
	%s
	%s
	ORDER BY a.urutan ASC, a.nama ASC
	%s
	`, favSelect, favJoin, whereSQL, limitSQL)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list apps: %w", err)
	}
	defer rows.Close()

	var apps []domain.App
	for rows.Next() {
		app, err := r.scanAppRow(rows)
		if err != nil {
			return nil, 0, err
		}
		apps = append(apps, *app)
	}

	return apps, total, nil
}

func (r *AppRepo) GetByID(ctx context.Context, id int, userID *int) (*domain.App, error) {
	favSelect, favJoin, argIdx, args := r.buildFavSnippet(userID, 1)
	args = append(args, id)

	query := fmt.Sprintf(`
	SELECT 
		a.id, a.category_id, a.nama, a.slug, a.url, a.deskripsi, a.ikon_url,
		a.target_pengguna, a.pemilik, a.kontak_admin, a.urutan, a.status_layanan,
		a.is_public, a.aktif, a.created_at, a.updated_at,
		c.id, c.nama, c.slug, c.deskripsi, c.urutan, c.ikon, c.created_at, c.updated_at
		%s
	FROM apps a
	INNER JOIN categories c ON a.category_id = c.id
	%s
	WHERE a.id = $%d
	`, favSelect, favJoin, argIdx)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	app, err := r.scanAppRow(rows)
	if err != nil {
		return nil, err
	}

	// Load allowed roles
	roles, err := r.getAppRoles(ctx, app.ID)
	if err == nil {
		app.AllowedRoles = roles
	}

	return app, nil
}

func (r *AppRepo) GetBySlug(ctx context.Context, slug string, userID *int) (*domain.App, error) {
	favSelect, favJoin, argIdx, args := r.buildFavSnippet(userID, 1)
	args = append(args, slug)

	query := fmt.Sprintf(`
	SELECT 
		a.id, a.category_id, a.nama, a.slug, a.url, a.deskripsi, a.ikon_url,
		a.target_pengguna, a.pemilik, a.kontak_admin, a.urutan, a.status_layanan,
		a.is_public, a.aktif, a.created_at, a.updated_at,
		c.id, c.nama, c.slug, c.deskripsi, c.urutan, c.ikon, c.created_at, c.updated_at
		%s
	FROM apps a
	INNER JOIN categories c ON a.category_id = c.id
	%s
	WHERE a.slug = $%d
	`, favSelect, favJoin, argIdx)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	app, err := r.scanAppRow(rows)
	if err != nil {
		return nil, err
	}

	roles, err := r.getAppRoles(ctx, app.ID)
	if err == nil {
		app.AllowedRoles = roles
	}

	return app, nil
}

func (r *AppRepo) Search(
	ctx context.Context,
	searchKeyword string,
	categoryID *int,
	roleIDs []int,
	isPublicOnly bool,
	userID *int,
	limit int,
) ([]domain.App, error) {
	if limit <= 0 {
		limit = 20
	}

	var whereClauses []string
	var args []interface{}
	argIdx := 1

	whereClauses = append(whereClauses, "a.aktif = true")

	if isPublicOnly {
		whereClauses = append(whereClauses, fmt.Sprintf("a.is_public = $%d", argIdx))
		args = append(args, true)
		argIdx++
	}

	if categoryID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("a.category_id = $%d", argIdx))
		args = append(args, *categoryID)
		argIdx++
	}

	if len(roleIDs) > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf(`(
			a.is_public = true 
			OR NOT EXISTS (SELECT 1 FROM app_access WHERE app_id = a.id)
			OR EXISTS (SELECT 1 FROM app_access aa WHERE aa.app_id = a.id AND aa.role_id = ANY($%d))
		)`, argIdx))
		args = append(args, roleIDs)
		argIdx++
	}

	// Trigram and text match filter
	whereClauses = append(whereClauses, fmt.Sprintf(`(
		a.nama ILIKE $%d 
		OR a.deskripsi ILIKE $%d 
		OR similarity(a.nama, $%d) > 0.15
		OR similarity(COALESCE(a.deskripsi, ''), $%d) > 0.15
	)`, argIdx, argIdx, argIdx+1, argIdx+1))
	args = append(args, "%"+searchKeyword+"%", searchKeyword)
	queryArgIdx := argIdx + 1
	argIdx += 2

	favSelect, favJoin, argIdx, favArgs := r.buildFavSnippet(userID, argIdx)
	args = append(args, favArgs...)

	args = append(args, limit)
	limitArgIdx := argIdx

	whereSQL := "WHERE " + strings.Join(whereClauses, " AND ")

	query := fmt.Sprintf(`
	SELECT 
		a.id, a.category_id, a.nama, a.slug, a.url, a.deskripsi, a.ikon_url,
		a.target_pengguna, a.pemilik, a.kontak_admin, a.urutan, a.status_layanan,
		a.is_public, a.aktif, a.created_at, a.updated_at,
		c.id, c.nama, c.slug, c.deskripsi, c.urutan, c.ikon, c.created_at, c.updated_at
		%s
	FROM apps a
	INNER JOIN categories c ON a.category_id = c.id
	%s
	%s
	ORDER BY 
		similarity(a.nama, $%d) DESC,
		a.urutan ASC
	LIMIT $%d
	`, favSelect, favJoin, whereSQL, queryArgIdx, limitArgIdx)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search apps: %w", err)
	}
	defer rows.Close()

	var apps []domain.App
	for rows.Next() {
		app, err := r.scanAppRow(rows)
		if err != nil {
			return nil, err
		}
		apps = append(apps, *app)
	}

	return apps, nil
}

func (r *AppRepo) Create(ctx context.Context, app *domain.App, roleIDs []int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `
	INSERT INTO apps (
		category_id, nama, slug, url, deskripsi, ikon_url, 
		target_pengguna, pemilik, kontak_admin, urutan, status_layanan, is_public, aktif
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(
		ctx, query,
		app.CategoryID, app.Nama, app.Slug, app.URL, app.Deskripsi, app.IkonURL,
		app.TargetPengguna, app.Pemilik, app.KontakAdmin, app.Urutan, app.StatusLayanan, app.IsPublic, app.Aktif,
	).Scan(&app.ID, &app.CreatedAt, &app.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert app: %w", err)
	}

	// Insert role access if not public
	for _, roleID := range roleIDs {
		_, err := tx.Exec(ctx, "INSERT INTO app_access (app_id, role_id) VALUES ($1, $2)", app.ID, roleID)
		if err != nil {
			return fmt.Errorf("failed to insert app_access: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *AppRepo) Update(ctx context.Context, app *domain.App, roleIDs []int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `
	UPDATE apps
	SET 
		category_id = $1, nama = $2, slug = $3, url = $4, deskripsi = $5, ikon_url = $6,
		target_pengguna = $7, pemilik = $8, kontak_admin = $9, urutan = $10, status_layanan = $11,
		is_public = $12, aktif = $13, updated_at = NOW()
	WHERE id = $14
	RETURNING updated_at
	`
	err = tx.QueryRow(
		ctx, query,
		app.CategoryID, app.Nama, app.Slug, app.URL, app.Deskripsi, app.IkonURL,
		app.TargetPengguna, app.Pemilik, app.KontakAdmin, app.Urutan, app.StatusLayanan,
		app.IsPublic, app.Aktif, app.ID,
	).Scan(&app.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update app: %w", err)
	}

	// Refresh role access
	if _, err := tx.Exec(ctx, "DELETE FROM app_access WHERE app_id = $1", app.ID); err != nil {
		return err
	}

	for _, roleID := range roleIDs {
		_, err := tx.Exec(ctx, "INSERT INTO app_access (app_id, role_id) VALUES ($1, $2)", app.ID, roleID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *AppRepo) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE apps SET status_layanan = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, status, id)
	return err
}

func (r *AppRepo) Delete(ctx context.Context, id int) error {
	// Soft delete
	query := `UPDATE apps SET aktif = false, updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *AppRepo) Reorder(ctx context.Context, appIDs []int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for idx, id := range appIDs {
		_, err := tx.Exec(ctx, "UPDATE apps SET urutan = $1 WHERE id = $2", idx+1, id)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *AppRepo) buildFavSnippet(userID *int, startIdx int) (string, string, int, []interface{}) {
	if userID != nil {
		favSelect := ", (f.user_id IS NOT NULL) as is_favorite"
		favJoin := fmt.Sprintf("LEFT JOIN favorites f ON a.id = f.app_id AND f.user_id = $%d", startIdx)
		return favSelect, favJoin, startIdx + 1, []interface{}{*userID}
	}
	return ", false as is_favorite", "", startIdx, []interface{}{}
}

func (r *AppRepo) scanAppRow(rows pgx.Rows) (*domain.App, error) {
	var a domain.App
	var c domain.Category

	err := rows.Scan(
		&a.ID, &a.CategoryID, &a.Nama, &a.Slug, &a.URL, &a.Deskripsi, &a.IkonURL,
		&a.TargetPengguna, &a.Pemilik, &a.KontakAdmin, &a.Urutan, &a.StatusLayanan,
		&a.IsPublic, &a.Aktif, &a.CreatedAt, &a.UpdatedAt,
		&c.ID, &c.Nama, &c.Slug, &c.Deskripsi, &c.Urutan, &c.Ikon, &c.CreatedAt, &c.UpdatedAt,
		&a.IsFavorite,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan app row: %w", err)
	}

	a.Category = &c
	return &a, nil
}

func (r *AppRepo) getAppRoles(ctx context.Context, appID int) ([]domain.Role, error) {
	query := `
	SELECT r.id, r.nama, r.deskripsi, r.created_at
	FROM roles r
	INNER JOIN app_access aa ON r.id = aa.role_id
	WHERE aa.app_id = $1
	`
	rows, err := r.pool.Query(ctx, query, appID)
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
