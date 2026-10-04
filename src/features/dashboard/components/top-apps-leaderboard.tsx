import { Trophy, Layers } from 'lucide-react'
import { TopAppMetric } from '@/types/analytics'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Progress } from '@/components/ui/progress'
import { Badge } from '@/components/ui/badge'

interface TopAppsLeaderboardProps {
  apps: TopAppMetric[]
  isLoading: boolean
}

export function TopAppsLeaderboard({ apps, isLoading }: TopAppsLeaderboardProps) {
  return (
    <Card className='col-span-1 lg:col-span-3 shadow-xs'>
      <CardHeader className='pb-3'>
        <div className='flex items-center justify-between'>
          <div>
            <CardTitle className='text-base font-semibold flex items-center gap-1.5'>
              <Trophy className='size-4 text-amber-500' />
              Aplikasi Terpopuler
            </CardTitle>
            <CardDescription className='text-xs'>
              Top 10 aplikasi paling sering diluncurkan
            </CardDescription>
          </div>
          <Badge variant='outline' className='text-[10px] font-mono'>
            30 Hari Terakhir
          </Badge>
        </div>
      </CardHeader>

      <CardContent>
        {isLoading ? (
          <div className='h-64 flex items-center justify-center text-xs text-muted-foreground'>
            Memuat peringkat aplikasi...
          </div>
        ) : apps.length === 0 ? (
          <div className='h-64 flex flex-col items-center justify-center text-xs text-muted-foreground gap-1.5'>
            <Layers className='size-8 text-muted-foreground/40' />
            <span>Belum ada catatan peluncuran aplikasi.</span>
          </div>
        ) : (
          <div className='space-y-3.5 max-h-[360px] overflow-y-auto pr-1'>
            {apps.map((app, index) => {
              const rank = index + 1
              return (
                <div
                  key={app.app_id}
                  className='flex items-center gap-3 p-2 rounded-lg hover:bg-muted/50 transition-colors'
                >
                  {/* Rank badge */}
                  <div
                    className={`size-6 rounded-full flex items-center justify-center text-[11px] font-bold shrink-0 ${
                      rank === 1
                        ? 'bg-amber-100 text-amber-800 dark:bg-amber-900/50 dark:text-amber-300'
                        : rank === 2
                        ? 'bg-slate-200 text-slate-700 dark:bg-slate-800 dark:text-slate-300'
                        : rank === 3
                        ? 'bg-amber-700/20 text-amber-900 dark:bg-amber-900/30 dark:text-amber-400'
                        : 'bg-muted text-muted-foreground'
                    }`}
                  >
                    {rank}
                  </div>

                  {/* App Icon */}
                  <div className='size-8 rounded-md bg-muted flex items-center justify-center shrink-0 overflow-hidden border'>
                    {app.ikon_url ? (
                      <img
                        src={app.ikon_url}
                        alt={app.nama}
                        className='size-full object-contain p-1'
                      />
                    ) : (
                      <Layers className='size-4 text-muted-foreground' />
                    )}
                  </div>

                  {/* App Info & Progress */}
                  <div className='flex-1 min-w-0'>
                    <div className='flex items-center justify-between mb-1'>
                      <span className='text-xs font-semibold text-foreground truncate'>
                        {app.nama}
                      </span>
                      <span className='text-xs font-mono font-medium text-foreground shrink-0 ml-2'>
                        {app.total_clicks.toLocaleString('id-ID')} klik
                      </span>
                    </div>

                    <div className='flex items-center gap-2'>
                      <Progress
                        value={app.click_percentage}
                        className='h-1.5 flex-1'
                      />
                      <span className='text-[10px] text-muted-foreground font-mono shrink-0 w-8 text-right'>
                        {app.click_percentage.toFixed(0)}%
                      </span>
                    </div>
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
