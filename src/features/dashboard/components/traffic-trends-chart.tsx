import { useState } from 'react'
import { LineChart, BarChart3 } from 'lucide-react'
import { TrendMetric } from '@/types/analytics'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

interface TrafficTrendsChartProps {
  trends: TrendMetric[]
  days: number
  onDaysChange: (days: number) => void
  isLoading: boolean
}

export function TrafficTrendsChart({
  trends,
  days,
  onDaysChange,
  isLoading,
}: TrafficTrendsChartProps) {
  const [chartType, setChartType] = useState<'bar' | 'line'>('bar')

  const maxClicks = Math.max(...trends.map((t) => t.total_clicks), 1)

  return (
    <Card className='col-span-1 lg:col-span-4 shadow-xs'>
      <CardHeader className='flex flex-col sm:flex-row sm:items-center justify-between pb-4 gap-2'>
        <div>
          <CardTitle className='text-base font-semibold'>Tren Aktivitas Akses Portal</CardTitle>
          <CardDescription className='text-xs'>
            Statistik klik peluncuran aplikasi dan pengguna aktif unik harian
          </CardDescription>
        </div>

        <div className='flex items-center gap-1.5 self-end sm:self-auto'>
          <div className='flex items-center rounded-lg border bg-muted/40 p-0.5'>
            <Button
              variant={days === 7 ? 'secondary' : 'ghost'}
              size='sm'
              className='h-7 text-xs px-2.5'
              onClick={() => onDaysChange(7)}
            >
              7 Hari
            </Button>
            <Button
              variant={days === 14 ? 'secondary' : 'ghost'}
              size='sm'
              className='h-7 text-xs px-2.5'
              onClick={() => onDaysChange(14)}
            >
              14 Hari
            </Button>
            <Button
              variant={days === 30 ? 'secondary' : 'ghost'}
              size='sm'
              className='h-7 text-xs px-2.5'
              onClick={() => onDaysChange(30)}
            >
              30 Hari
            </Button>
          </div>

          <Button
            variant='outline'
            size='icon'
            className='size-7'
            onClick={() => setChartType(chartType === 'bar' ? 'line' : 'bar')}
            title='Ubah format grafik'
          >
            {chartType === 'bar' ? <LineChart className='size-3.5' /> : <BarChart3 className='size-3.5' />}
          </Button>
        </div>
      </CardHeader>

      <CardContent>
        {isLoading ? (
          <div className='h-64 flex items-center justify-center text-xs text-muted-foreground'>
            Memuat grafik tren...
          </div>
        ) : trends.length === 0 ? (
          <div className='h-64 flex items-center justify-center text-xs text-muted-foreground'>
            Belum ada data aktivitas untuk rentang waktu ini.
          </div>
        ) : (
          <div className='space-y-3'>
            {/* Chart Bars */}
            <div className='h-52 flex items-end gap-1.5 pt-4 pb-2 border-b'>
              {trends.map((item, idx) => {
                const heightPercent = Math.max((item.total_clicks / maxClicks) * 100, 4)
                return (
                  <div
                    key={idx}
                    className='flex-1 flex flex-col items-center gap-1.5 group relative'
                  >
                    {/* Tooltip */}
                    <div className='absolute -top-12 z-20 hidden group-hover:flex flex-col items-center bg-popover text-popover-foreground text-[10px] px-2 py-1 rounded-md shadow-md border whitespace-nowrap'>
                      <span className='font-bold'>{item.date}</span>
                      <span>{item.total_clicks} klik • {item.unique_users} user</span>
                    </div>

                    {/* Bar */}
                    <div className='w-full bg-primary/15 rounded-t group-hover:bg-primary transition-all duration-200 h-full flex items-end justify-center'>
                      <div
                        className='w-full bg-primary rounded-t transition-all duration-300 group-hover:bg-primary/90'
                        style={{ height: `${heightPercent}%` }}
                      />
                    </div>
                  </div>
                )
              })}
            </div>

            {/* X-axis Labels */}
            <div className='flex items-center justify-between text-[10px] text-muted-foreground font-mono px-1'>
              <span>{trends[0]?.date}</span>
              <span className='flex items-center gap-1 text-[11px] font-sans font-medium text-foreground'>
                <span className='size-2 rounded-full bg-primary inline-block' />
                Total Klik Harian
              </span>
              <span>{trends[trends.length - 1]?.date}</span>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
