package usecase

import (
	"context"
	"fmt"

	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/repository/postgres"
)

type UpdateRolePermissionsInput struct {
	PermissionCodes []string `json:"permission_codes" binding:"required"`
}

type AssignUserRolesInput struct {
	RoleIDs []int `json:"role_ids" binding:"required"`
}

type RBACUsecase struct {
	rbacRepo *postgres.RBACRepo
}

func NewRBACUsecase(rbacRepo *postgres.RBACRepo) *RBACUsecase {
	return &RBACUsecase{rbacRepo: rbacRepo}
}

func (u *RBACUsecase) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	if u.rbacRepo == nil {
		return nil, fmt.Errorf("rbac repository is not initialized")
	}
	return u.rbacRepo.ListPermissions(ctx)
}

func (u *RBACUsecase) ListRolesWithPermissions(ctx context.Context) ([]domain.RoleWithPermissions, error) {
	if u.rbacRepo == nil {
		return nil, fmt.Errorf("rbac repository is not initialized")
	}
	return u.rbacRepo.ListRolesWithPermissions(ctx)
}

func (u *RBACUsecase) UpdateRolePermissions(ctx context.Context, roleID int, input UpdateRolePermissionsInput) error {
	if u.rbacRepo == nil {
		return fmt.Errorf("rbac repository is not initialized")
	}
	return u.rbacRepo.UpdateRolePermissions(ctx, roleID, input.PermissionCodes)
}

func (u *RBACUsecase) ListUsers(ctx context.Context, search string, page, perPage int) ([]domain.UserWithRoles, int, error) {
	if u.rbacRepo == nil {
		return nil, 0, fmt.Errorf("rbac repository is not initialized")
	}
	return u.rbacRepo.ListUsersWithRoles(ctx, search, page, perPage)
}

func (u *RBACUsecase) AssignUserRoles(ctx context.Context, userID int, input AssignUserRolesInput) error {
	if u.rbacRepo == nil {
		return fmt.Errorf("rbac repository is not initialized")
	}
	return u.rbacRepo.AssignUserRoles(ctx, userID, input.RoleIDs)
}
