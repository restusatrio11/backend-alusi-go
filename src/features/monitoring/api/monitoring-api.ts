import apiClient from '@/lib/api-client'
import { StandardApiResponse } from '@/types/auth'
import {
  ServiceUptimeSummary,
  AppStatusHistoryItem,
} from '@/types/monitoring'

export const monitoringApi = {
  /**
   * Mengambil ringkasan uptime dan SLA seluruh layanan
   */
  async getServiceStatusSummary(days = 30): Promise<ServiceUptimeSummary> {
    const response = await apiClient.get<StandardApiResponse<ServiceUptimeSummary>>(
      '/api/v1/services/status',
      { params: { days } }
    )
    return response.data.data
  },

  /**
   * Mengambil riwayat status ping aplikasi
   */
  async getAppStatusHistory(slugOrId: string | number, limit = 30): Promise<AppStatusHistoryItem[]> {
    const response = await apiClient.get<StandardApiResponse<AppStatusHistoryItem[]>>(
      `/api/v1/apps/${slugOrId}/status-history`,
      { params: { limit } }
    )
    return response.data.data || []
  },

  /**
   * Trigger manual health probe (Admin)
   */
  async probeApp(appId: number): Promise<any> {
    const response = await apiClient.post<StandardApiResponse<any>>(
      `/api/v1/admin/apps/${appId}/probe`
    )
    return response.data.data
  },

  /**
   * URL endpoint untuk Server-Sent Events (SSE) stream
   */
  getSSEStreamURL(): string {
    const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
    return `${baseURL}/services/realtime-status`
  },
}
