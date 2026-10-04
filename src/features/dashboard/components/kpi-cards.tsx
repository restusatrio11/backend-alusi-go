import { Layers, MousePointerClick, Activity, UserCheck } from 'lucide-react'
import { DashboardSummary } from '@/types/analytics'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface KPICardsProps {
  summary: DashboardSummary | null
  isLoading: boolean
}

export function KPICards({ summary, isLoading }: KPICardsProps) {
  const totalUsers = summary?.total_users ?? 0
  const totalApps = summary?.total_apps ?? 0
  const onlineApps = summary?.online_apps ?? 0
  const totalClicks = summary?.total_clicks ?? 0
  const dau = summary?.dau ?? 0
  const mau = summary?.mau ?? 0
  const uptime = summary?.overall_uptime_percentage ?? 100

  return (
    <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-4'>
      {/* Total Apps */}
      <Card className='shadow-xs'>
        <CardHeader className='flex flex-row items-center justify-between pb-2'>
          <CardTitle className='text-xs font-semibold text-muted-foreground'>
            Katalog Aplikasi
          </CardTitle>
          <Layers className='size-4 text-primary' />
        </CardHeader>
        <CardContent>
          <div className='text-2xl font-bold text-foreground'>
            {isLoading ? '...' : totalApps}
          </div>
          <p className='text-[11px] text-emerald-600 dark:text-emerald-400 mt-1 font-medium'>
            {onlineApps} layanan operasional online
          </p>
        </CardContent>
      </Card>

      {/* Total Interactions / Clicks */}
      <Card className='shadow-xs'>
        <CardHeader className='flex flex-row items-center justify-between pb-2'>
          <CardTitle className='text-xs font-semibold text-muted-foreground'>
            Total Akses & Klik
          </CardTitle>
          <MousePointerClick className='size-4 text-blue-500' />
        </CardHeader>
        <CardContent>
          <div className='text-2xl font-bold text-foreground'>
            {isLoading ? '...' : totalClicks.toLocaleString('id-ID')}
          </div>
          <p className='text-[11px] text-muted-foreground mt-1'>
            Akumulasi peluncuran modul
          </p>
        </CardContent>
      </Card>

      {/* Active Users (DAU & MAU) */}
      <Card className='shadow-xs'>
        <CardHeader className='flex flex-row items-center justify-between pb-2'>
          <CardTitle className='text-xs font-semibold text-muted-foreground'>
            Pengguna Aktif (DAU / MAU)
          </CardTitle>
          <UserCheck className='size-4 text-emerald-500' />
        </CardHeader>
        <CardContent>
          <div className='text-2xl font-bold text-foreground'>
            {isLoading ? '...' : `${dau} / ${mau}`}
          </div>
          <p className='text-[11px] text-muted-foreground mt-1'>
            {totalUsers} total akun pegawai terdaftar
          </p>
        </CardContent>
      </Card>

      {/* SLA Uptime */}
      <Card className='shadow-xs'>
        <CardHeader className='flex flex-row items-center justify-between pb-2'>
          <CardTitle className='text-xs font-semibold text-muted-foreground'>
            Ketersediaan Sistem (SLA)
          </CardTitle>
          <Activity className='size-4 text-amber-500' />
        </CardHeader>
        <CardContent>
          <div className='text-2xl font-bold text-foreground'>
            {isLoading ? '...' : `${uptime.toFixed(1)}%`}
          </div>
          <p className='text-[11px] text-emerald-600 dark:text-emerald-400 mt-1 font-medium'>
            Uptime 30 hari terakhir
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
