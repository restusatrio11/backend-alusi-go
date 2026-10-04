import apiClient from '@/lib/api-client'
import { StandardApiResponse } from '@/types/auth'
import {
  DashboardSummary,
  TopAppMetric,
  TrendMetric,
  DisruptionItem,
} from '@/types/analytics'

export const analyticsApi = {
  /**
   * Mengambil ringkasan metrik eksekutif (KPI Rollup)
   */
  async getDashboardSummary(): Promise<DashboardSummary> {
    const response = await apiClient.get<StandardApiResponse<DashboardSummary>>(
      '/api/v1/admin/analytics/summary'
    )
    return response.data.data
  },

  /**
   * Mengambil daftar ranking aplikasi paling populer berdasarkan klik
   */
  async getTopApps(days = 30, limit = 10): Promise<TopAppMetric[]> {
    const response = await apiClient.get<StandardApiResponse<TopAppMetric[]>>(
      '/api/v1/admin/analytics/top-apps',
      { params: { days, limit } }
    )
    return response.data.data || []
  },

  /**
   * Mengambil data tren grafik klik harian dan pengguna unik
   */
  async getTrends(days = 30): Promise<TrendMetric[]> {
    const response = await apiClient.get<StandardApiResponse<TrendMetric[]>>(
      '/api/v1/admin/analytics/trends',
      { params: { days } }
    )
    return response.data.data || []
  },

  /**
   * Mengambil rekapitulasi gangguan layanan dan SLA
   */
  async getDisruptions(days = 30): Promise<DisruptionItem[]> {
    const response = await apiClient.get<StandardApiResponse<DisruptionItem[]>>(
      '/api/v1/admin/analytics/disruptions',
      { params: { days } }
    )
    return response.data.data || []
  },
}
