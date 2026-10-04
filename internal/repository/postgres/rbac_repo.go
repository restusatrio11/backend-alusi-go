package postgres

import (
	"context"
	"fmt"
	"strings"

	"backend-alusi-go/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RBACRepo struct {
	pool *pgxpool.Pool
}

func NewRBACRepo(pool *pgxpool.Pool) *RBACRepo {
	return &RBACRepo{pool: pool}
}

// ListPermissions returns all registered system permissions grouped by module
func (r *RBACRepo) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	query := `
		SELECT id, kode, nama, modul, COALESCE(deskripsi, ''), created_at
		FROM permissions
		ORDER BY modul ASC, id ASC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}
	defer rows.Close()

	perms := make([]domain.Permission, 0)
	for rows.Next() {
		var p domain.Permission
		if err := rows.Scan(&p.ID, &p.Kode, &p.Nama, &p.Modul, &p.Deskripsi, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		perms = append(perms, p)
	}
	return perms, nil
}

// ListRolesWithPermissions returns all roles with their assigned permissions
func (r *RBACRepo) ListRolesWithPermissions(ctx context.Context) ([]domain.RoleWithPermissions, error) {
	rolesQuery := `
		SELECT id, nama, COALESCE(deskripsi, ''), created_at
		FROM roles
		ORDER BY id ASC
	`
	rows, err := r.pool.Query(ctx, rolesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	roles := make([]domain.RoleWithPermissions, 0)
	for rows.Next() {
		var role domain.RoleWithPermissions
		if err := rows.Scan(&role.ID, &role.Nama, &role.Deskripsi, &role.CreatedAt); err != nil {
			rows.Close()
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		role.Permissions = make([]domain.Permission, 0)
		roles = append(roles, role)
	}
	rows.Close()

	for i := range roles {
		permQuery := `
			SELECT p.id, p.kode, p.nama, p.modul, COALESCE(p.deskripsi, ''), p.created_at
			FROM permissions p
			JOIN role_permissions rp ON p.id = rp.permission_id
			WHERE rp.role_id = $1
			ORDER BY p.modul ASC, p.id ASC
		`
		pRows, err := r.pool.Query(ctx, permQuery, roles[i].ID)
		if err == nil {
			rolePerms := make([]domain.Permission, 0)
			for pRows.Next() {
				var p domain.Permission
				if err := pRows.Scan(&p.ID, &p.Kode, &p.Nama, &p.Modul, &p.Deskripsi, &p.CreatedAt); err == nil {
					rolePerms = append(rolePerms, p)
				}
			}
			pRows.Close()
			roles[i].Permissions = rolePerms
		}
	}

	return roles, nil
}

// GetUserPermissions returns all unique permission codes assigned to a user via their roles
func (r *RBACRepo) GetUserPermissions(ctx context.Context, userID int) ([]string, error) {
	query := `
		SELECT DISTINCT p.kode
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY p.kode ASC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		perms = append(perms, code)
	}
	return perms, nil
}

// UpdateRolePermissions sets the list of permissions assigned to a role
func (r *RBACRepo) UpdateRolePermissions(ctx context.Context, roleID int, permissionCodes []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Delete existing permissions for the role
	_, err = tx.Exec(ctx, "DELETE FROM role_permissions WHERE role_id = $1", roleID)
	if err != nil {
		return fmt.Errorf("failed to clear role permissions: %w", err)
	}

	// Insert new permissions
	if len(permissionCodes) > 0 {
		insertQuery := `
			INSERT INTO role_permissions (role_id, permission_id)
			SELECT $1, p.id
			FROM permissions p
			WHERE p.kode = ANY($2)
			ON CONFLICT DO NOTHING
		`
		_, err = tx.Exec(ctx, insertQuery, roleID, permissionCodes)
		if err != nil {
			return fmt.Errorf("failed to insert role permissions: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// ListUsersWithRoles returns users with pagination and search filter
func (r *RBACRepo) ListUsersWithRoles(ctx context.Context, search string, page, perPage int) ([]domain.UserWithRoles, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 15
	}
	offset := (page - 1) * perPage

	whereClause := "1=1"
	args := []interface{}{}
	if strings.TrimSpace(search) != "" {
		args = append(args, "%"+strings.ToLower(strings.TrimSpace(search))+"%")
		whereClause += fmt.Sprintf(" AND (LOWER(u.nama) LIKE $%d OR LOWER(u.email) LIKE $%d OR LOWER(COALESCE(u.username, '')) LIKE $%d OR COALESCE(u.nip, '') LIKE $%d)", len(args), len(args), len(args), len(args))
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users u WHERE %s", whereClause)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Select users
	queryArgs := append(args, perPage, offset)
	query := fmt.Sprintf(`
		SELECT u.id, u.username, u.nip, u.nama, u.email, u.user_type, u.status,
		       s.id, s.kode, s.nama,
		       u.last_login_at, u.created_at
		FROM users u
		LEFT JOIN satker s ON u.satker_id = s.id
		WHERE %s
		ORDER BY u.id ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, len(queryArgs)-1, len(queryArgs))

	rows, err := r.pool.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []domain.UserWithRoles
	for rows.Next() {
		var u domain.UserWithRoles
		var satkerID *int
		var satkerKode, satkerNama *string

		err := rows.Scan(
			&u.ID, &u.Username, &u.NIP, &u.Nama, &u.Email, &u.UserType, &u.Status,
			&satkerID, &satkerKode, &satkerNama,
			&u.LastLoginAt, &u.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}

		if satkerID != nil && satkerKode != nil && satkerNama != nil {
			u.Satker = &domain.Satker{
				ID:   *satkerID,
				Kode: *satkerKode,
				Nama: *satkerNama,
			}
		}

		users = append(users, u)
	}

	// Fetch roles and permissions for each user
	for i := range users {
		rolesQuery := `
			SELECT r.id, r.nama, COALESCE(r.deskripsi, '')
			FROM roles r
			JOIN user_roles ur ON r.id = ur.role_id
			WHERE ur.user_id = $1
		`
		rRows, err := r.pool.Query(ctx, rolesQuery, users[i].ID)
		if err == nil {
			var roles []domain.Role
			for rRows.Next() {
				var r domain.Role
				if err := rRows.Scan(&r.ID, &r.Nama, &r.Deskripsi); err == nil {
					roles = append(roles, r)
				}
			}
			rRows.Close()
			users[i].Roles = roles
		}

		perms, err := r.GetUserPermissions(ctx, users[i].ID)
		if err == nil {
			users[i].Permissions = perms
		}
	}

	return users, total, nil
}

// AssignUserRoles assigns multiple roles to a specific user
func (r *RBACRepo) AssignUserRoles(ctx context.Context, userID int, roleIDs []int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "DELETE FROM user_roles WHERE user_id = $1", userID)
	if err != nil {
		return fmt.Errorf("failed to clear user roles: %w", err)
	}

	if len(roleIDs) > 0 {
		for _, rID := range roleIDs {
			_, err = tx.Exec(ctx, "INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, rID)
			if err != nil {
				return fmt.Errorf("failed to assign role %d: %w", rID, err)
			}
		}
	}

	return tx.Commit(ctx)
}
