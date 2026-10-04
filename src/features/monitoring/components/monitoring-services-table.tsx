import { useState } from 'react'
import {
  CheckCircle2,
  Clock,
  AlertTriangle,
  XCircle,
  Activity,
  Zap,
} from 'lucide-react'
import { toast } from 'sonner'
import { AppUptimeSummary, StatusLayanan } from '@/types/monitoring'
import { monitoringApi } from '../api/monitoring-api'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'

interface MonitoringServicesTableProps {
  services: AppUptimeSummary[]
  onRefresh: () => void
}

export function MonitoringServicesTable({
  services,
  onRefresh,
}: MonitoringServicesTableProps) {
  const [probingId, setProbingId] = useState<number | null>(null)

  const handleProbe = async (appId: number, appName: string) => {
    setProbingId(appId)
    try {
      await monitoringApi.probeApp(appId)
      toast.success(`Probe kesehatan "${appName}" selesai dilakukan.`)
      onRefresh()
    } catch (err: any) {
      toast.error('Gagal Melakukan Probe', {
        description: err?.message || 'Tidak dapat menghubungi server.',
      })
    } finally {
      setProbingId(null)
    }
  }

  const getStatusBadge = (status: StatusLayanan) => {
    switch (status) {
      case 'online':
        return (
          <Badge variant='outline' className='bg-emerald-500/10 text-emerald-600 border-emerald-500/30 text-xs'>
            <CheckCircle2 className='mr-1 size-3' /> Online
          </Badge>
        )
      case 'pemeliharaan':
        return (
          <Badge variant='outline' className='bg-amber-500/10 text-amber-600 border-amber-500/30 text-xs'>
            <Clock className='mr-1 size-3' /> Maintenance
          </Badge>
        )
      case 'kendala':
        return (
          <Badge variant='outline' className='bg-orange-500/10 text-orange-600 border-orange-500/30 text-xs'>
            <AlertTriangle className='mr-1 size-3' /> Kendala
          </Badge>
        )
      case 'offline':
      default:
        return (
          <Badge variant='outline' className='bg-destructive/10 text-destructive border-destructive/30 text-xs'>
            <XCircle className='mr-1 size-3' /> Offline
          </Badge>
        )
    }
  }

  const getUptimeColor = (pct: number) => {
    if (pct >= 99) return 'text-emerald-600'
    if (pct >= 95) return 'text-amber-600'
    return 'text-destructive'
  }

  return (
    <div className='rounded-lg border bg-card shadow-xs overflow-hidden'>
      <Table>
        <TableHeader>
          <TableRow className='bg-muted/40'>
            <TableHead className='w-12 text-center'>#</TableHead>
            <TableHead>Aplikasi / Layanan</TableHead>
            <TableHead>Status Live</TableHead>
            <TableHead>Uptime (30 Hari)</TableHead>
            <TableHead>Respon Ping</TableHead>
            <TableHead>Pemeriksaan</TableHead>
            <TableHead className='text-end w-28'>Uji Manual</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {services.length === 0 ? (
            <TableRow>
              <TableCell colSpan={7} className='h-32 text-center text-muted-foreground'>
                Belum ada data monitoring layanan yang tercatat.
              </TableCell>
            </TableRow>
          ) : (
            services.map((svc, index) => (
              <TableRow key={svc.app_id} className='hover:bg-muted/30'>
                <TableCell className='text-center text-xs font-mono text-muted-foreground'>
                  {index + 1}
                </TableCell>
                <TableCell>
                  <div className='flex items-center gap-3'>
                    {svc.app_ikon_url ? (
                      <img
                        src={svc.app_ikon_url}
                        alt={svc.app_nama}
                        className='size-8 rounded-lg object-contain bg-white p-1 border shadow-xs shrink-0'
                      />
                    ) : (
                      <div className='size-8 rounded-lg bg-primary/10 text-primary flex items-center justify-center font-bold text-xs shrink-0'>
                        {svc.app_nama.substring(0, 2).toUpperCase()}
                      </div>
                    )}
                    <div>
                      <p className='font-semibold text-sm text-foreground'>{svc.app_nama}</p>
                      <p className='text-xs text-muted-foreground font-mono'>/{svc.app_slug}</p>
                    </div>
                  </div>
                </TableCell>
                <TableCell>{getStatusBadge(svc.current_status)}</TableCell>
                <TableCell>
                  <div className='space-y-1 max-w-[140px]'>
                    <div className='flex items-center justify-between text-xs'>
                      <span className={`font-semibold ${getUptimeColor(svc.uptime_percentage)}`}>
                        {svc.uptime_percentage.toFixed(1)}%
                      </span>
                    </div>
                    <div className='h-1.5 w-full rounded-full bg-muted overflow-hidden'>
                      <div
                        className={`h-full rounded-full ${
                          svc.uptime_percentage >= 99
                            ? 'bg-emerald-500'
                            : svc.uptime_percentage >= 95
                              ? 'bg-amber-500'
                              : 'bg-destructive'
                        }`}
                        style={{ width: `${Math.min(100, Math.max(0, svc.uptime_percentage))}%` }}
                      />
                    </div>
                  </div>
                </TableCell>
                <TableCell>
                  <div className='flex items-center gap-1 text-xs font-mono'>
                    <Zap className='size-3 text-blue-500' />
                    <span>{svc.avg_latency_ms} ms</span>
                  </div>
                </TableCell>
                <TableCell>
                  <div className='text-xs text-muted-foreground'>
                    <span className='text-emerald-600 font-medium'>{svc.success_checks} OK</span>
                    {svc.failed_checks > 0 && (
                      <span className='text-destructive font-medium ml-1'>
                        / {svc.failed_checks} Fail
                      </span>
                    )}
                  </div>
                </TableCell>
                <TableCell className='text-end'>
                  <Button
                    variant='outline'
                    size='sm'
                    className='text-xs h-8'
                    onClick={() => handleProbe(svc.app_id, svc.app_nama)}
                    disabled={probingId === svc.app_id}
                  >
                    <Activity className={`mr-1 size-3.5 ${probingId === svc.app_id ? 'animate-spin' : ''}`} />
                    {probingId === svc.app_id ? 'Probing...' : 'Ping'}
                  </Button>
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  )
}
