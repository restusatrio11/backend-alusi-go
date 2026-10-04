import apiClient from '@/lib/api-client'
import { StandardApiResponse } from '@/types/auth'
import {
  Announcement,
  CreateAnnouncementPayload,
  UpdateAnnouncementPayload,
} from '@/types/announcements'

export const announcementsApi = {
  /**
   * Mengambil seluruh pengumuman (Admin)
   */
  async getAdminAnnouncements(page = 1, perPage = 50): Promise<{ announcements: Announcement[]; meta?: any }> {
    const response = await apiClient.get<StandardApiResponse<Announcement[]>>(
      '/api/v1/admin/announcements',
      { params: { page, per_page: perPage } }
    )
    return {
      announcements: response.data.data || [],
      meta: response.data.meta,
    }
  },

  /**
   * Mengambil pengumuman yang sedang aktif
   */
  async getActiveAnnouncements(appId?: number): Promise<Announcement[]> {
    const response = await apiClient.get<StandardApiResponse<Announcement[]>>(
      '/api/v1/announcements',
      { params: { app_id: appId } }
    )
    return response.data.data || []
  },

  /**
   * Membuat pengumuman baru (Admin)
   */
  async createAnnouncement(payload: CreateAnnouncementPayload): Promise<Announcement> {
    const response = await apiClient.post<StandardApiResponse<Announcement>>(
      '/api/v1/admin/announcements',
      payload
    )
    return response.data.data
  },

  /**
   * Memperbarui pengumuman (Admin)
   */
  async updateAnnouncement(id: number, payload: UpdateAnnouncementPayload): Promise<Announcement> {
    const response = await apiClient.put<StandardApiResponse<Announcement>>(
      `/api/v1/admin/announcements/${id}`,
      payload
    )
    return response.data.data
  },

  /**
   * Menghapus pengumuman (Admin)
   */
  async deleteAnnouncement(id: number): Promise<boolean> {
    const response = await apiClient.delete<StandardApiResponse<null>>(
      `/api/v1/admin/announcements/${id}`
    )
    return response.data.success
  },
}
