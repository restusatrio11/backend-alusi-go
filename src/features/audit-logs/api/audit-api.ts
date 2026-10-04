import apiClient from '@/lib/api-client'
import { StandardApiResponse } from '@/types/auth'
import { AuditLog, AuditLogFilterParams } from '@/types/audit'

export const auditApi = {
  /**
   * Mengambil log jejak audit aktivitas admin
   */
  async getAuditLogs(params?: AuditLogFilterParams): Promise<{ logs: AuditLog[]; meta?: any }> {
    const response = await apiClient.get<StandardApiResponse<AuditLog[]>>(
      '/api/v1/admin/audit-logs',
      { params }
    )
    return {
      logs: response.data.data || [],
      meta: response.data.meta,
    }
  },
}
