import apiClient from '@/lib/api-client'
import { StandardApiResponse } from '@/types/auth'
import {
  App,
  AppsFilterParams,
  CreateAppPayload,
  UpdateAppPayload,
  ReorderAppsPayload,
} from '@/types/apps'

export const appsApi = {
  /**
   * Mengambil daftar aplikasi dari katalog (dengan filter dan pagination)
   */
  async getApps(params?: AppsFilterParams): Promise<{ apps: App[]; meta?: any }> {
    const response = await apiClient.get<StandardApiResponse<App[]>>('/api/v1/apps', {
      params,
    })
    return {
      apps: response.data.data || [],
      meta: response.data.meta,
    }
  },

  /**
   * Mengambil detail aplikasi berdasarkan slug
   */
  async getAppBySlug(slug: string): Promise<App> {
    const response = await apiClient.get<StandardApiResponse<App>>(`/api/v1/apps/${slug}`)
    return response.data.data
  },

  /**
   * Pencarian cerdas aplikasi (Trigram pg_trgm)
   */
  async searchApps(query: string, categorySlug?: string, limit = 20): Promise<App[]> {
    const response = await apiClient.get<StandardApiResponse<App[]>>('/api/v1/apps/search', {
      params: { q: query, category: categorySlug, limit },
    })
    return response.data.data || []
  },

  /**
   * Menambahkan aplikasi baru (Admin)
   */
  async createApp(payload: CreateAppPayload): Promise<App> {
    const response = await apiClient.post<StandardApiResponse<App>>('/api/v1/admin/apps', payload)
    return response.data.data
  },

  /**
   * Memperbarui metadata aplikasi (Admin)
   */
  async updateApp(id: number, payload: UpdateAppPayload): Promise<App> {
    const response = await apiClient.put<StandardApiResponse<App>>(`/api/v1/admin/apps/${id}`, payload)
    return response.data.data
  },

  /**
   * Menonaktifkan / menghapus aplikasi (Admin)
   */
  async deleteApp(id: number): Promise<boolean> {
    const response = await apiClient.delete<StandardApiResponse<null>>(`/api/v1/admin/apps/${id}`)
    return response.data.success
  },

  /**
   * Memperbarui urutan sequence aplikasi (Admin Drag-and-Drop)
   */
  async reorderApps(appIds: number[]): Promise<boolean> {
    const payload: ReorderAppsPayload = { app_ids: appIds }
    const response = await apiClient.put<StandardApiResponse<null>>('/api/v1/admin/apps/reorder', payload)
    return response.data.success
  },

  /**
   * Mengunggah & mengompresi logo aplikasi (Admin)
   */
  async uploadAppLogo(id: number, file: File): Promise<{ app: App; logo_url: string }> {
    const formData = new FormData()
    formData.append('logo', file)

    const response = await apiClient.post<StandardApiResponse<{ app: App; logo_url: string }>>(
      `/api/v1/admin/apps/${id}/logo`,
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      }
    )
    return response.data.data
  },

  /**
   * Trigger manual health probe (Admin)
   */
  async probeApp(id: number): Promise<any> {
    const response = await apiClient.post<StandardApiResponse<any>>(`/api/v1/admin/apps/${id}/probe`)
    return response.data.data
  },
}
