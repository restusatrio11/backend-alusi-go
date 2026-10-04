import { createFileRoute, redirect } from '@tanstack/react-router'
import { AuditLogs } from '@/features/audit-logs'
import { useAuthStore } from '@/stores/auth-store'

export const Route = createFileRoute('/_authenticated/audit-logs/')({
  beforeLoad: () => {
    const auth = useAuthStore.getState().auth
    if (!auth.hasPermission('audit:view')) {
      throw redirect({
        to: '/403',
      })
    }
  },
  component: AuditLogs,
})
