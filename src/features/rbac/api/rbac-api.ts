import apiClient from '@/lib/api-client'
import { StandardApiResponse } from '@/types/auth'
import {
  Permission,
  RoleWithPermissions,
  UserWithRoles,
  UpdateRolePermissionsPayload,
  AssignUserRolesPayload,
  UsersFilterParams,
} from '@/types/rbac'

export const rbacApi = {
  /**
   * Mengambil seluruh daftar permission di sistem
   */
  async getPermissions(): Promise<Permission[]> {
    const response = await apiClient.get<StandardApiResponse<Permission[]>>('/api/v1/admin/rbac/permissions')
    return response.data.data || []
  },

  /**
   * Mengambil seluruh daftar role beserta permission yang dimilikinya
   */
  async getRoles(): Promise<RoleWithPermissions[]> {
    const response = await apiClient.get<StandardApiResponse<RoleWithPermissions[]>>('/api/v1/admin/rbac/roles')
    return response.data.data || []
  },

  /**
   * Memperbarui daftar hak akses (permissions) untuk suatu role
   */
  async updateRolePermissions(
    roleId: number,
    payload: UpdateRolePermissionsPayload
  ): Promise<RoleWithPermissions> {
    const response = await apiClient.put<StandardApiResponse<RoleWithPermissions>>(
      `/api/v1/admin/rbac/roles/${roleId}/permissions`,
      payload
    )
    return response.data.data
  },

  /**
   * Mengambil daftar pengguna berserta role yang dimiliki
   */
  async getUsers(params?: UsersFilterParams): Promise<{ users: UserWithRoles[]; meta?: any }> {
    const response = await apiClient.get<StandardApiResponse<UserWithRoles[]>>('/api/v1/admin/rbac/users', {
      params,
    })
    return {
      users: response.data.data || [],
      meta: response.data.meta,
    }
  },

  /**
   * Menetapkan / memperbarui role untuk seorang pengguna
   */
  async assignUserRoles(
    userId: number,
    payload: AssignUserRolesPayload
  ): Promise<{ user_id: number; role_ids: number[] }> {
    const response = await apiClient.put<StandardApiResponse<{ user_id: number; role_ids: number[] }>>(
      `/api/v1/admin/rbac/users/${userId}/roles`,
      payload
    )
    return response.data.data
  },
}
