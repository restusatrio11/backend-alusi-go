import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { LayoutDashboard, RefreshCw } from 'lucide-react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { analyticsApi } from '@/features/analytics/api/analytics-api'
import { KPICards } from './components/kpi-cards'
import { TrafficTrendsChart } from './components/traffic-trends-chart'
import { TopAppsLeaderboard } from './components/top-apps-leaderboard'
import { DisruptionsWidget } from './components/disruptions-widget'
import { ExportReportButtons } from './components/export-report-buttons'

export function Dashboard() {
  const [trendDays, setTrendDays] = useState<number>(30)

  // Query: Dashboard Summary Rollup
  const {
    data: summary,
    isLoading: isLoadingSummary,
    refetch: refetchSummary,
    isRefetching: isRefetchingSummary,
  } = useQuery({
    queryKey: ['analytics', 'summary'],
    queryFn: () => analyticsApi.getDashboardSummary(),
  })

  // Query: Top Apps Leaderboard
  const {
    data: topApps = [],
    isLoading: isLoadingTopApps,
    refetch: refetchTopApps,
  } = useQuery({
    queryKey: ['analytics', 'top-apps'],
    queryFn: () => analyticsApi.getTopApps(30, 10),
  })

  // Query: Traffic Trends
  const {
    data: trends = [],
    isLoading: isLoadingTrends,
    refetch: refetchTrends,
  } = useQuery({
    queryKey: ['analytics', 'trends', trendDays],
    queryFn: () => analyticsApi.getTrends(trendDays),
  })

  // Query: Disruptions & SLA
  const {
    data: disruptions = [],
    isLoading: isLoadingDisruptions,
    refetch: refetchDisruptions,
  } = useQuery({
    queryKey: ['analytics', 'disruptions'],
    queryFn: () => analyticsApi.getDisruptions(30),
  })

  const handleRefreshAll = () => {
    refetchSummary()
    refetchTopApps()
    refetchTrends()
    refetchDisruptions()
  }

  const isRefreshing = isRefetchingSummary

  return (
    <>
      <Header>
        <div className='flex items-center gap-2'>
          <LayoutDashboard className='size-5 text-primary' />
          <h2 className='text-sm font-semibold tracking-tight'>
            Portal Eksekutif BPS
          </h2>
        </div>
        <div className='ml-auto flex items-center space-x-2'>
          <ThemeSwitch />
          <ProfileDropdown />
        </div>
      </Header>

      <Main>
        <div className='space-y-6'>
          {/* Header Title & Actions */}
          <div className='flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4'>
            <div>
              <h1 className='text-2xl font-bold tracking-tight'>
                Ringkasan Eksekutif ALUSI
              </h1>
              <p className='text-sm text-muted-foreground'>
                Pemantauan metrik agregat penggunaan modul, trafik pengguna harian, dan SLA sistem.
              </p>
            </div>

            <div className='flex items-center gap-2'>
              <ExportReportButtons />
              <Button
                variant='outline'
                size='sm'
                onClick={handleRefreshAll}
                disabled={isRefreshing}
                className='h-8 text-xs'
              >
                <RefreshCw
                  className={`mr-1.5 size-3.5 ${isRefreshing ? 'animate-spin' : ''}`}
                />
                Refresh
              </Button>
            </div>
          </div>

          {/* Metric KPI Cards */}
          <KPICards summary={summary || null} isLoading={isLoadingSummary} />

          {/* Charts & Leaderboard Row */}
          <div className='grid grid-cols-1 gap-4 lg:grid-cols-7'>
            <TrafficTrendsChart
              trends={trends}
              days={trendDays}
              onDaysChange={setTrendDays}
              isLoading={isLoadingTrends}
            />
            <TopAppsLeaderboard
              apps={topApps}
              isLoading={isLoadingTopApps}
            />
          </div>

          {/* Disruptions & SLA Recaps */}
          <DisruptionsWidget
            disruptions={disruptions}
            isLoading={isLoadingDisruptions}
          />
        </div>
      </Main>
    </>
  )
}
