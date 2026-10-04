import apiClient from '@/lib/api-client'
import { StandardApiResponse } from '@/types/auth'
import { Category, CategoryPayload } from '@/types/apps'

export const categoriesApi = {
  /**
   * Mengambil daftar seluruh kategori aplikasi
   */
  async getCategories(): Promise<Category[]> {
    const response = await apiClient.get<StandardApiResponse<Category[]>>('/api/v1/categories')
    return response.data.data || []
  },

  /**
   * Mengambil detail kategori berdasarkan slug
   */
  async getCategoryBySlug(slug: string): Promise<Category> {
    const response = await apiClient.get<StandardApiResponse<Category>>(`/api/v1/categories/${slug}`)
    return response.data.data
  },

  /**
   * Menambahkan kategori baru (Admin)
   */
  async createCategory(payload: CategoryPayload): Promise<Category> {
    const response = await apiClient.post<StandardApiResponse<Category>>('/api/v1/admin/categories', payload)
    return response.data.data
  },

  /**
   * Memperbarui data kategori (Admin)
   */
  async updateCategory(id: number, payload: CategoryPayload): Promise<Category> {
    const response = await apiClient.put<StandardApiResponse<Category>>(`/api/v1/admin/categories/${id}`, payload)
    return response.data.data
  },

  /**
   * Menghapus kategori (Admin)
   */
  async deleteCategory(id: number): Promise<boolean> {
    const response = await apiClient.delete<StandardApiResponse<null>>(`/api/v1/admin/categories/${id}`)
    return response.data.success
  },
}
