import {
  CheckCircle2,
  AlertTriangle,
  Activity,
  Zap,
} from 'lucide-react'
import { ServiceUptimeSummary } from '@/types/monitoring'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface MonitoringStatCardsProps {
  summary: ServiceUptimeSummary | null
}

export function MonitoringStatCards({ summary }: MonitoringStatCardsProps) {
  const total = summary?.total_apps || 0
  const online = summary?.online_apps || 0
  const maintenance = summary?.maintenance_apps || 0
  const degraded = summary?.degraded_apps || 0
  const uptimePercent = summary?.overall_uptime_percentage ?? 100
  const avgLatency = summary?.avg_response_time_ms ?? 0

  return (
    <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-4'>
      {/* Online Services */}
      <Card className='border-emerald-500/30 bg-emerald-500/5 shadow-xs'>
        <CardHeader className='flex flex-row items-center justify-between pb-2'>
          <CardTitle className='text-xs font-semibold text-emerald-600 dark:text-emerald-400'>
            Layanan Online
          </CardTitle>
          <CheckCircle2 className='size-4 text-emerald-500' />
        </CardHeader>
        <CardContent>
          <div className='text-2xl font-bold text-foreground'>
            {online} <span className='text-xs font-normal text-muted-foreground'>/ {total} aplikasi</span>
          </div>
          <p className='text-[11px] text-muted-foreground mt-1'>
            Operasional normal dan dapat diakses
          </p>
        </CardContent>
      </Card>

      {/* Maintenance & Degraded */}
      <Card className='border-amber-500/30 bg-amber-500/5 shadow-xs'>
        <CardHeader className='flex flex-row items-center justify-between pb-2'>
          <CardTitle className='text-xs font-semibold text-amber-600 dark:text-amber-400'>
            Pemeliharaan & Kendala
          </CardTitle>
          <AlertTriangle className='size-4 text-amber-500' />
        </CardHeader>
        <CardContent>
          <div className='text-2xl font-bold text-foreground'>
            {maintenance + degraded}{' '}
            <span className='text-xs font-normal text-muted-foreground'>
              ({maintenance} maint, {degraded} kendala)
            </span>
          </div>
          <p className='text-[11px] text-muted-foreground mt-1'>
            Layanan dalam proses perbaikan
          </p>
        </CardContent>
      </Card>

      {/* Overall SLA Uptime % */}
      <Card className='shadow-xs'>
        <CardHeader className='flex flex-row items-center justify-between pb-2'>
          <CardTitle className='text-xs font-medium text-muted-foreground'>
            Uptime SLA (30 Hari)
          </CardTitle>
          <Activity className='size-4 text-primary' />
        </CardHeader>
        <CardContent>
          <div className='text-2xl font-bold text-foreground'>
            {uptimePercent.toFixed(2)}%
          </div>
          <p className='text-[11px] text-emerald-600 dark:text-emerald-400 mt-1 font-medium'>
            Memenuhi target SLA &gt; 99.0%
          </p>
        </CardContent>
      </Card>

      {/* Avg Response Latency */}
      <Card className='shadow-xs'>
        <CardHeader className='flex flex-row items-center justify-between pb-2'>
          <CardTitle className='text-xs font-medium text-muted-foreground'>
            Rata-rata Respon Ping
          </CardTitle>
          <Zap className='size-4 text-blue-500' />
        </CardHeader>
        <CardContent>
          <div className='text-2xl font-bold text-foreground'>
            {avgLatency} <span className='text-xs font-normal text-muted-foreground'>ms</span>
          </div>
          <p className='text-[11px] text-muted-foreground mt-1'>
            Latensi jaringan probe backend
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
