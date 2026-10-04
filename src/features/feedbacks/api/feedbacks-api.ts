import apiClient from '@/lib/api-client'
import { StandardApiResponse } from '@/types/auth'
import { Feedback, UpdateFeedbackStatusPayload } from '@/types/feedbacks'

export const feedbacksApi = {
  /**
   * Mengambil daftar tiket feedback / laporan kendala (Admin)
   */
  async getFeedbacks(params?: {
    status?: string
    app_id?: number
    page?: number
    per_page?: number
  }): Promise<{ feedbacks: Feedback[]; meta?: any }> {
    const response = await apiClient.get<StandardApiResponse<Feedback[]>>(
      '/api/v1/admin/feedbacks',
      { params }
    )
    return {
      feedbacks: response.data.data || [],
      meta: response.data.meta,
    }
  },

  /**
   * Memperbarui status tindak lanjut tiket feedback (Admin)
   */
  async updateFeedbackStatus(
    id: number,
    payload: UpdateFeedbackStatusPayload
  ): Promise<Feedback> {
    const response = await apiClient.put<StandardApiResponse<Feedback>>(
      `/api/v1/admin/feedbacks/${id}`,
      payload
    )
    return response.data.data
  },
}
