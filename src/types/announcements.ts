import { App } from './apps'

export type AnnouncementType = 'info' | 'warning' | 'danger' | 'success'

export interface Announcement {
  id: number
  app_id?: number | null
  app?: App | null
  judul: string
  pesan: string
  tipe: AnnouncementType
  tautan_url?: string | null
  tautan_teks?: string | null
  mulai_pada?: string | null
  berakhir_pada?: string | null
  aktif: boolean
  created_at: string
  updated_at: string
}

export interface CreateAnnouncementPayload {
  app_id?: number | null
  judul: string
  pesan: string
  tipe: AnnouncementType
  tautan_url?: string
  tautan_teks?: string
  mulai_pada?: string
  berakhir_pada?: string
  aktif?: boolean
}

export interface UpdateAnnouncementPayload {
  app_id?: number | null
  judul: string
  pesan: string
  tipe: AnnouncementType
  tautan_url?: string
  tautan_teks?: string
  mulai_pada?: string
  berakhir_pada?: string
  aktif?: boolean
}
