import apiClient from '@/lib/api-client'

export const reportsApi = {
  /**
   * Mengunduh file laporan inventaris katalog aplikasi dalam format CSV
   */
  async downloadCatalogCSV(): Promise<void> {
    const response = await apiClient.get('/api/v1/admin/reports/catalog/export', {
      responseType: 'blob',
    })

    const blob = new Blob([response.data], { type: 'text/csv;charset=utf-8;' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.setAttribute(
      'download',
      `katalog_aplikasi_bps_${new Date().toISOString().slice(0, 10)}.csv`
    )
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  },

  /**
   * Mengunduh file rekapitulasi analitik penggunaan dan SLA dalam format CSV
   */
  async downloadAnalyticsCSV(): Promise<void> {
    const response = await apiClient.get('/api/v1/admin/reports/analytics/export', {
      responseType: 'blob',
    })

    const blob = new Blob([response.data], { type: 'text/csv;charset=utf-8;' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.setAttribute(
      'download',
      `laporan_analitik_portal_bps_${new Date().toISOString().slice(0, 10)}.csv`
    )
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  },
}
