export interface DashboardSummary {
  total_users: number
  total_apps: number
  online_apps: number
  total_clicks: number
  dau: number
  mau: number
  overall_uptime_percentage: number
}

export interface TopAppMetric {
  app_id: number
  nama: string
  slug: string
  ikon_url?: string | null
  category_nama?: string
  total_clicks: number
  click_percentage: number
}

export interface TrendMetric {
  date: string
  total_clicks: number
  unique_users: number
}

export interface DisruptionItem {
  app_id: number
  nama: string
  slug: string
  ikon_url?: string | null
  total_incidents: number
  total_downtime_minutes: number
  uptime_percentage: number
  last_incident_at?: string
}
