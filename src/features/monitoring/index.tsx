import { Activity, RefreshCw } from 'lucide-react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { useSSEStatus } from './hooks/use-sse-status'
import { MonitoringStatCards } from './components/monitoring-stat-cards'
import { MonitoringServicesTable } from './components/monitoring-services-table'

export function Monitoring() {
  const { summary, isConnected, lastUpdated, refresh } = useSSEStatus()

  return (
    <>
      <Header>
        <div className='flex items-center gap-2'>
          <Activity className='size-5 text-primary' />
          <h2 className='text-sm font-semibold tracking-tight'>Monitoring Layanan Realtime</h2>
        </div>
        <div className='ml-auto flex items-center space-x-2'>
          <ThemeSwitch />
          <ProfileDropdown />
        </div>
      </Header>

      <Main>
        <div className='space-y-6'>
          {/* Header Title & SSE Connection Status */}
          <div className='flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4'>
            <div>
              <div className='flex items-center gap-2'>
                <h1 className='text-2xl font-bold tracking-tight'>Status & Uptime Layanan</h1>
                {isConnected ? (
                  <Badge variant='outline' className='bg-emerald-500/10 text-emerald-600 border-emerald-500/30 gap-1.5 py-0.5 text-[11px]'>
                    <span className='size-2 rounded-full bg-emerald-500 animate-pulse' />
                    SSE Terhubung
                  </Badge>
                ) : (
                  <Badge variant='outline' className='bg-amber-500/10 text-amber-600 border-amber-500/30 gap-1.5 py-0.5 text-[11px]'>
                    <span className='size-2 rounded-full bg-amber-500' />
                    Menghubungkan SSE...
                  </Badge>
                )}
              </div>
              <p className='text-sm text-muted-foreground'>
                Pemantauan status operasional aplikasi secara realtime melalui protokol Server-Sent Events (SSE).
              </p>
            </div>

            <div className='flex items-center gap-2'>
              {lastUpdated && (
                <span className='text-xs text-muted-foreground hidden sm:inline'>
                  Update: {lastUpdated.toLocaleTimeString()}
                </span>
              )}
              <Button variant='outline' size='sm' onClick={refresh}>
                <RefreshCw className='mr-1.5 size-3.5' />
                Refresh
              </Button>
            </div>
          </div>

          {/* Metric Stats Cards */}
          <MonitoringStatCards summary={summary} />

          {/* Services Table */}
          <div className='space-y-3'>
            <div className='flex items-center justify-between'>
              <h3 className='text-base font-semibold'>Daftar Layanan Aplikasi Terpantau</h3>
              <span className='text-xs text-muted-foreground'>
                Total: {summary?.services?.length || 0} Layanan
              </span>
            </div>
            <MonitoringServicesTable
              services={summary?.services || []}
              onRefresh={refresh}
            />
          </div>
        </div>
      </Main>
    </>
  )
}
