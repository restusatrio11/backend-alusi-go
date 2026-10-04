import { Role } from './auth'

export interface Category {
  id: number
  nama: string
  slug: string
  deskripsi?: string | null
  urutan: number
  ikon?: string | null
  created_at: string
  updated_at: string
}

export type StatusLayanan = 'online' | 'pemeliharaan' | 'kendala' | 'offline'

export interface AppGuide {
  id: number
  app_id: number
  judul: string
  konten: string
  urutan: number
  created_at: string
  updated_at: string
}

export interface App {
  id: number
  category_id: number
  category?: Category
  nama: string
  slug: string
  url: string
  deskripsi?: string | null
  ikon_url?: string | null
  target_pengguna: string // "Semua" | "Internal BPS" | "Mitra" | "Pimpinan"
  pemilik: string
  kontak_admin?: string | null
  urutan: number
  status_layanan: StatusLayanan
  is_public: boolean
  aktif: boolean
  roles?: Role[]
  guides?: AppGuide[]
  total_clicks?: number
  is_favorite?: boolean
  created_at: string
  updated_at: string
}

export interface CreateAppPayload {
  category_id: number
  nama: string
  slug?: string
  url: string
  deskripsi?: string
  ikon_url?: string
  target_pengguna?: string
  pemilik?: string
  kontak_admin?: string
  urutan?: number
  status_layanan?: StatusLayanan
  is_public?: boolean
  role_ids?: number[]
}

export interface UpdateAppPayload {
  category_id: number
  nama: string
  slug?: string
  url: string
  deskripsi?: string
  ikon_url?: string
  target_pengguna?: string
  pemilik?: string
  kontak_admin?: string
  urutan?: number
  status_layanan?: StatusLayanan
  is_public?: boolean
  aktif?: boolean
  role_ids?: number[]
}

export interface CategoryPayload {
  nama: string
  slug?: string
  deskripsi?: string
  urutan?: number
  ikon?: string
}

export interface GuidePayload {
  judul: string
  konten: string
  urutan?: number
}

export interface ReorderAppsPayload {
  app_ids: number[]
}

export interface AppsFilterParams {
  category?: string
  target_pengguna?: string
  status?: string
  persona?: string
  page?: number
  per_page?: number
  q?: string
}
