export interface AuditLog {
  id: number
  user_id?: number | null
  user_name?: string | null
  user_email?: string | null
  action: string // "CREATE" | "UPDATE" | "DELETE" | "PROBE" | "LOGIN"
  entity: string // "app" | "category" | "guide" | "announcement" | "feedback"
  entity_id?: string | null
  ip_address?: string | null
  user_agent?: string | null
  payload?: Record<string, any> | null
  details?: string | null
  created_at: string
}

export interface AuditLogFilterParams {
  user_id?: number
  action?: string
  entity?: string
  entity_id?: string
  start_date?: string
  end_date?: string
  page?: number
  per_page?: number
}
