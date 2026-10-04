import { AlertTriangle, CheckCircle2, Clock, Layers } from 'lucide-react'
import { DisruptionItem } from '@/types/analytics'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'

interface DisruptionsWidgetProps {
  disruptions: DisruptionItem[]
  isLoading: boolean
}

export function DisruptionsWidget({
  disruptions,
  isLoading,
}: DisruptionsWidgetProps) {
  return (
    <Card className='shadow-xs'>
      <CardHeader className='pb-3'>
        <div className='flex items-center justify-between'>
          <div>
            <CardTitle className='text-base font-semibold flex items-center gap-1.5'>
              <AlertTriangle className='size-4 text-amber-500' />
              Rekapitulasi Gangguan Layanan & SLA
            </CardTitle>
            <CardDescription className='text-xs'>
              Riwayat downtime dan ketersediaan layanan 30 hari terakhir
            </CardDescription>
          </div>
          <Badge variant='outline' className='text-[10px]'>
            {disruptions.length} Sistem Terpantau
          </Badge>
        </div>
      </CardHeader>

      <CardContent>
        {isLoading ? (
          <div className='h-48 flex items-center justify-center text-xs text-muted-foreground'>
            Memuat data histori gangguan...
          </div>
        ) : disruptions.length === 0 ? (
          <div className='h-36 flex flex-col items-center justify-center text-xs text-muted-foreground gap-1.5'>
            <CheckCircle2 className='size-7 text-emerald-500' />
            <span className='font-medium text-foreground'>Seluruh Sistem Berjalan Normal</span>
            <span>Tidak ada insiden downtime dalam 30 hari terakhir.</span>
          </div>
        ) : (
          <div className='overflow-x-auto'>
            <table className='w-full text-xs'>
              <thead>
                <tr className='border-b text-muted-foreground font-medium text-left'>
                  <th className='pb-2 pl-2'>Aplikasi</th>
                  <th className='pb-2 text-center'>Total Insiden</th>
                  <th className='pb-2 text-center'>Total Downtime</th>
                  <th className='pb-2 text-center'>SLA Uptime</th>
                  <th className='pb-2 text-right pr-2'>Insiden Terakhir</th>
                </tr>
              </thead>
              <tbody className='divide-y divide-border/50'>
                {disruptions.map((item) => {
                  const uptimePct = item.uptime_percentage ?? 100
                  const isHealthy = uptimePct >= 99.0
                  return (
                    <tr
                      key={item.app_id}
                      className='hover:bg-muted/40 transition-colors'
                    >
                      <td className='py-2.5 pl-2'>
                        <div className='flex items-center gap-2'>
                          <div className='size-6 rounded bg-muted flex items-center justify-center shrink-0 border overflow-hidden'>
                            {item.ikon_url ? (
                              <img
                                src={item.ikon_url}
                                alt={item.nama}
                                className='size-full object-contain p-0.5'
                              />
                            ) : (
                              <Layers className='size-3 text-muted-foreground' />
                            )}
                          </div>
                          <span className='font-medium text-foreground truncate max-w-[200px]'>
                            {item.nama}
                          </span>
                        </div>
                      </td>

                      <td className='py-2.5 text-center font-mono'>
                        {(item.total_incidents ?? 0) > 0 ? (
                          <Badge
                            variant='secondary'
                            className='text-[10px] bg-amber-100 text-amber-800 dark:bg-amber-900/40 dark:text-amber-300'
                          >
                            {item.total_incidents}x
                          </Badge>
                        ) : (
                          <span className='text-muted-foreground'>0</span>
                        )}
                      </td>

                      <td className='py-2.5 text-center font-mono'>
                        {(item.total_downtime_minutes ?? 0) > 0 ? (
                          <span className='text-rose-600 dark:text-rose-400 font-medium flex items-center justify-center gap-1'>
                            <Clock className='size-3' />
                            {item.total_downtime_minutes} mnt
                          </span>
                        ) : (
                          <span className='text-muted-foreground'>-</span>
                        )}
                      </td>

                      <td className='py-2.5 text-center font-mono font-medium'>
                        <span
                          className={
                            isHealthy
                              ? 'text-emerald-600 dark:text-emerald-400'
                              : 'text-amber-600 dark:text-amber-400'
                          }
                        >
                          {uptimePct.toFixed(1)}%
                        </span>
                      </td>

                      <td className='py-2.5 text-right pr-2 text-muted-foreground font-mono text-[11px]'>
                        {item.last_incident_at
                          ? new Date(item.last_incident_at).toLocaleDateString('id-ID', {
                              day: '2-digit',
                              month: 'short',
                              hour: '2-digit',
                              minute: '2-digit',
                            })
                          : 'Tidak ada'}
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
  )
}
