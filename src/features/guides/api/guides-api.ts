import apiClient from '@/lib/api-client'
import { StandardApiResponse } from '@/types/auth'
import { AppGuide, GuidePayload } from '@/types/apps'

export const guidesApi = {
  /**
   * Mengambil panduan dan FAQ untuk aplikasi berdasarkan slug
   */
  async getGuidesByAppSlug(slug: string): Promise<AppGuide[]> {
    const response = await apiClient.get<StandardApiResponse<AppGuide[]>>(`/api/v1/apps/${slug}/guides`)
    return response.data.data || []
  },

  /**
   * Menambahkan panduan baru ke aplikasi (Admin)
   */
  async createGuide(appId: number, payload: GuidePayload): Promise<AppGuide> {
    const response = await apiClient.post<StandardApiResponse<AppGuide>>(
      `/api/v1/admin/apps/${appId}/guides`,
      payload
    )
    return response.data.data
  },

  /**
   * Memperbarui panduan (Admin)
   */
  async updateGuide(guideId: number, payload: GuidePayload): Promise<AppGuide> {
    const response = await apiClient.put<StandardApiResponse<AppGuide>>(
      `/api/v1/admin/guides/${guideId}`,
      payload
    )
    return response.data.data
  },

  /**
   * Menghapus panduan (Admin)
   */
  async deleteGuide(guideId: number): Promise<boolean> {
    const response = await apiClient.delete<StandardApiResponse<null>>(
      `/api/v1/admin/guides/${guideId}`
    )
    return response.data.success
  },
}
