export interface Role {
  id: number
  nama: string
  deskripsi?: string
  created_at: string
}

export interface Satker {
  id: number
  kode: string
  nama: string
  created_at: string
  updated_at: string
}

export interface User {
  id: number
  sso_sub?: string
  username?: string
  user_type: string // "internal" | "external"
  nip?: string
  nik?: string
  nama: string
  email: string
  foto_url?: string
  satker_id?: number
  satker?: Satker
  roles?: Role[]
  status: string // "active" | "inactive"
  metadata?: Record<string, any>
  last_login_at?: string
  created_at: string
  updated_at: string
}

export interface StandardApiResponse<T = any> {
  success: boolean
  message: string
  data: T
  meta?: {
    page?: number
    limit?: number
    total_pages?: number
    total_records?: number
  }
  error?: {
    code?: string
    details?: any
  }
}

export interface LoginResponse {
  token: string
  user: User
}
