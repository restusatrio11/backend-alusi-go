import { StatusLayanan } from './apps'
export type { StatusLayanan }

export interface AppUptimeSummary {
  app_id: number
  app_nama: string
  app_slug: string
  app_ikon_url?: string | null
  current_status: StatusLayanan
  uptime_percentage: number
  avg_latency_ms: number
  last_check_at?: string
  total_checks: number
  success_checks: number
  failed_checks: number
}

export interface ServiceUptimeSummary {
  total_apps: number
  online_apps: number
  maintenance_apps: number
  degraded_apps: number
  offline_apps: number
  overall_uptime_percentage: number
  avg_response_time_ms: number
  services: AppUptimeSummary[]
}

export interface AppStatusHistoryItem {
  id: number
  app_id: number
  status: StatusLayanan
  latency_ms: number
  status_code?: number
  pesan_error?: string
  checked_at: string
}

export interface SSEStatusEvent {
  type: 'snapshot' | 'status_change' | 'ping'
  timestamp: string
  app_id?: number
  status?: StatusLayanan
  latency_ms?: number
  services?: ServiceUptimeSummary
}
