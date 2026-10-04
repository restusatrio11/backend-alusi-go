export interface Permission {
  id: number
  kode: string
  nama: string
  kategori: string
  deskripsi?: string
  created_at: string
}

export interface RoleWithPermissions {
  id: number
  nama: string
  deskripsi?: string
  permissions: Permission[]
  created_at: string
}

export interface UserWithRoles {
  id: number
  sso_sub?: string
  username?: string
  user_type: string
  nip?: string
  nama: string
  email: string
  foto_url?: string
  satker_id?: number
  satker_nama?: string
  status: string
  roles: {
    id: number
    nama: string
    deskripsi?: string
  }[]
  created_at: string
  updated_at: string
}

export interface UpdateRolePermissionsPayload {
  permission_codes: string[]
}

export interface AssignUserRolesPayload {
  role_ids: number[]
}

export interface UsersFilterParams {
  page?: number
  limit?: number
  search?: string
  role_id?: number
}
