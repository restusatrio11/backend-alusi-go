import { App } from './apps'

export type FeedbackStatus = 'pending' | 'in_progress' | 'resolved' | 'closed'
export type FeedbackCategory = 'kendala' | 'saran' | 'pertanyaan' | 'lainnya'

export interface Feedback {
  id: number
  app_id?: number | null
  app?: App | null
  user_id?: number | null
  nama_pelapor: string
  email_pelapor?: string | null
  kontak_pelapor?: string | null
  kategori: FeedbackCategory
  judul: string
  isi_laporan: string
  lampiran_url?: string | null
  status: FeedbackStatus
  tanggapan_admin?: string | null
  ditanggapi_pada?: string | null
  created_at: string
  updated_at: string
}

export interface UpdateFeedbackStatusPayload {
  status: FeedbackStatus
  tanggapan_admin?: string
}
