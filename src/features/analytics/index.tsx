import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { LineChart, RefreshCw, Layers, MousePointerClick } from 'lucide-react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { analyticsApi } from './api/analytics-api'
import { TrafficTrendsChart } from '@/features/dashboard/components/traffic-trends-chart'
import { ExportReportButtons } from '@/features/dashboard/components/export-report-buttons'

export function Analytics() {
  const [days, setDays] = useState<number>(30)

  const {
    data: topApps = [],
    isLoading: isLoadingTopApps,
    refetch: refetchTopApps,
  } = useQuery({
    queryKey: ['analytics', 'top-apps', days],
    queryFn: () => analyticsApi.getTopApps(days, 50),
  })

  const {
    data: trends = [],
    isLoading: isLoadingTrends,
    refetch: refetchTrends,
  } = useQuery({
    queryKey: ['analytics', 'trends', days],
    queryFn: () => analyticsApi.getTrends(days),
  })

  const totalInteractions = topApps.reduce((acc, curr) => acc + curr.total_clicks, 0)
  const mostPopularApp = topApps[0]?.nama || '-'

  const handleRefresh = () => {
    refetchTopApps()
    refetchTrends()
  }

  return (
    <>
      <Header>
        <div className='flex items-center gap-2'>
          <LineChart className='size-5 text-primary' />
          <h2 className='text-sm font-semibold tracking-tight'>Analitik & Statistik Portal</h2>
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
                Analitik Akses & Adopsi Aplikasi
              </h1>
              <p className='text-sm text-muted-foreground'>
                Laporan komprehensif tingkat pemanfaatan modul-modul sistem oleh pegawai BPS.
              </p>
            </div>

            <div className='flex items-center gap-2'>
              <ExportReportButtons />
              <Button
                variant='outline'
                size='sm'
                onClick={handleRefresh}
                className='h-8 text-xs'
              >
                <RefreshCw className='mr-1.5 size-3.5' />
                Refresh
              </Button>
            </div>
          </div>

          {/* Highlights Row */}
          <div className='grid gap-4 sm:grid-cols-3'>
            <Card className='shadow-xs'>
              <CardHeader className='pb-2'>
                <CardTitle className='text-xs font-semibold text-muted-foreground'>
                  Total Interaksi ({days} Hari)
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className='text-2xl font-bold text-foreground flex items-center gap-2'>
                  <MousePointerClick className='size-5 text-primary' />
                  {totalInteractions.toLocaleString('id-ID')}
                </div>
                <p className='text-[11px] text-muted-foreground mt-1'>
                  Total akumulasi klik peluncuran aplikasi
                </p>
              </CardContent>
            </Card>

            <Card className='shadow-xs'>
              <CardHeader className='pb-2'>
                <CardTitle className='text-xs font-semibold text-muted-foreground'>
                  Aplikasi Paling Sering Digunakan
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className='text-xl font-bold text-foreground truncate'>
                  {mostPopularApp}
                </div>
                <p className='text-[11px] text-emerald-600 dark:text-emerald-400 mt-1 font-medium'>
                  Peringkat 1 trafik tertinggi
                </p>
              </CardContent>
            </Card>

            <Card className='shadow-xs'>
              <CardHeader className='pb-2'>
                <CardTitle className='text-xs font-semibold text-muted-foreground'>
                  Aplikasi Terindeks
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className='text-2xl font-bold text-foreground flex items-center gap-2'>
                  <Layers className='size-5 text-amber-500' />
                  {topApps.length} Aplikasi
                </div>
                <p className='text-[11px] text-muted-foreground mt-1'>
                  Memiliki rekaman klik aktif
                </p>
              </CardContent>
            </Card>
          </div>

          {/* Main Trends Chart */}
          <TrafficTrendsChart
            trends={trends}
            days={days}
            onDaysChange={setDays}
            isLoading={isLoadingTrends}
          />

          {/* Full Distribution Breakdown Table */}
          <Card className='shadow-xs'>
            <CardHeader className='pb-3'>
              <div className='flex items-center justify-between'>
                <div>
                  <CardTitle className='text-base font-semibold'>
                    Distribusi Peluncuran per Aplikasi
                  </CardTitle>
                  <CardDescription className='text-xs'>
                    Daftar lengkap metrik klik seluruh aplikasi yang terpasang pada portal
                  </CardDescription>
                </div>
                <Badge variant='outline' className='text-xs font-mono'>
                  Rentang {days} Hari
                </Badge>
              </div>
            </CardHeader>

            <CardContent>
              {isLoadingTopApps ? (
                <div className='h-48 flex items-center justify-center text-xs text-muted-foreground'>
                  Memuat rincian statistik...
                </div>
              ) : topApps.length === 0 ? (
                <div className='h-32 flex items-center justify-center text-xs text-muted-foreground'>
                  Belum ada rekaman klik aplikasi.
                </div>
              ) : (
                <div className='overflow-x-auto'>
                  <table className='w-full text-xs'>
                    <thead>
                      <tr className='border-b text-muted-foreground font-medium text-left'>
                        <th className='pb-2 pl-2 w-12 text-center'>Rank</th>
                        <th className='pb-2'>Nama Aplikasi</th>
                        <th className='pb-2'>Kategori</th>
                        <th className='pb-2 text-right'>Total Klik</th>
                        <th className='pb-2 text-right pr-4 w-48'>Pangsa Penggunaan</th>
                      </tr>
                    </thead>
                    <tbody className='divide-y divide-border/50'>
                      {topApps.map((app, idx) => {
                        const rank = idx + 1
                        return (
                          <tr key={app.app_id} className='hover:bg-muted/40 transition-colors'>
                            <td className='py-3 pl-2 text-center font-bold text-muted-foreground'>
                              #{rank}
                            </td>
                            <td className='py-3'>
                              <div className='flex items-center gap-2'>
                                <div className='size-6 rounded bg-muted flex items-center justify-center shrink-0 border overflow-hidden'>
                                  {app.ikon_url ? (
                                    <img
                                      src={app.ikon_url}
                                      alt={app.nama}
                                      className='size-full object-contain p-0.5'
                                    />
                                  ) : (
                                    <Layers className='size-3 text-muted-foreground' />
                                  )}
                                </div>
                                <span className='font-semibold text-foreground'>
                                  {app.nama}
                                </span>
                              </div>
                            </td>
                            <td className='py-3 text-muted-foreground'>
                              {app.category_nama || 'Umum'}
                            </td>
                            <td className='py-3 text-right font-mono font-bold text-foreground'>
                              {app.total_clicks.toLocaleString('id-ID')}
                            </td>
                            <td className='py-3 pr-4 text-right'>
                              <div className='flex items-center justify-end gap-2'>
                                <Progress value={app.click_percentage} className='h-2 w-24' />
                                <span className='font-mono text-xs w-10 text-right'>
                                  {app.click_percentage.toFixed(1)}%
                                </span>
                              </div>
                            </td>
                          </tr>
                        )
                      })}
                    </tbody>
                  </table>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </Main>
    </>
  )
}
