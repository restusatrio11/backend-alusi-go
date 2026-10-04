import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { AuditLog } from '@/types/audit'
import { ScrollArea } from '@/components/ui/scroll-area'
import { ShieldCheck, Laptop, Globe, Calendar, User, FileCode, Layers } from 'lucide-react'

interface AuditDetailDialogProps {
  log: AuditLog | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function AuditDetailDialog({
  log,
  open,
  onOpenChange,
}: AuditDetailDialogProps) {
  if (!log) return null

  const getActionBadgeVariant = (action: string) => {
    switch (action?.toUpperCase()) {
      case 'CREATE':
        return 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/40 dark:text-emerald-300'
      case 'UPDATE':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900/40 dark:text-blue-300'
      case 'DELETE':
        return 'bg-rose-100 text-rose-800 dark:bg-rose-900/40 dark:text-rose-300'
      case 'PROBE':
        return 'bg-purple-100 text-purple-800 dark:bg-purple-900/40 dark:text-purple-300'
      case 'LOGIN':
        return 'bg-amber-100 text-amber-800 dark:bg-amber-900/40 dark:text-amber-300'
      default:
        return 'bg-muted text-muted-foreground'
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-3xl lg:max-w-4xl max-h-[90vh] overflow-y-auto p-6 md:p-8'>
        <DialogHeader className='space-y-1.5 pb-2 border-b'>
          <div className='flex items-center gap-2 text-primary'>
            <ShieldCheck className='size-5' />
            <DialogTitle className='text-xl font-bold tracking-tight'>
              Rekaman Forensik Jejak Audit #{log.id}
            </DialogTitle>
          </div>
          <DialogDescription className='text-xs text-muted-foreground leading-relaxed'>
            Catatan jejak forensik keamanan, manipulasi data entitas, dan riwayat aktivitas administrator.
          </DialogDescription>
        </DialogHeader>

        <div className='space-y-5 py-3 text-xs'>
          {/* Metadata Grid */}
          <div className='grid grid-cols-1 sm:grid-cols-2 gap-4 p-4 bg-muted/30 rounded-xl border'>
            <div className='flex items-start gap-3'>
              <div className='p-2 rounded-lg bg-primary/10 text-primary shrink-0'>
                <User className='size-4' />
              </div>
              <div className='min-w-0'>
                <span className='text-muted-foreground text-[11px] block'>Aktor Pelaksana</span>
                <span className='font-semibold text-foreground text-sm block truncate'>
                  {log.user_name || 'System / Service Worker'}
                </span>
                {log.user_email && (
                  <span className='text-[11px] text-muted-foreground block font-mono'>
                    {log.user_email}
                  </span>
                )}
              </div>
            </div>

            <div className='flex items-start gap-3'>
              <div className='p-2 rounded-lg bg-blue-500/10 text-blue-600 shrink-0'>
                <Calendar className='size-4' />
              </div>
              <div>
                <span className='text-muted-foreground text-[11px] block'>Waktu Kejadian (Timestamp)</span>
                <span className='font-medium font-mono text-foreground text-sm block'>
                  {new Date(log.created_at).toLocaleString('id-ID')}
                </span>
                <span className='text-[10px] text-muted-foreground font-mono'>
                  {new Date(log.created_at).toISOString()}
                </span>
              </div>
            </div>

            <div className='flex items-start gap-3'>
              <div className='p-2 rounded-lg bg-purple-500/10 text-purple-600 shrink-0'>
                <Layers className='size-4' />
              </div>
              <div>
                <span className='text-muted-foreground text-[11px] block'>Jenis Aksi & Target</span>
                <div className='flex items-center gap-2 mt-1'>
                  <Badge
                    variant='outline'
                    className={`font-mono text-xs px-2 py-0.5 border-0 font-bold ${getActionBadgeVariant(
                      log.action
                    )}`}
                  >
                    {log.action}
                  </Badge>
                  <span className='font-bold text-foreground uppercase text-xs'>
                    {log.entity} {log.entity_id ? `(#${log.entity_id})` : ''}
                  </span>
                </div>
              </div>
            </div>

            <div className='flex items-start gap-3'>
              <div className='p-2 rounded-lg bg-emerald-500/10 text-emerald-600 shrink-0'>
                <Globe className='size-4' />
              </div>
              <div>
                <span className='text-muted-foreground text-[11px] block'>Alamat IP Asal</span>
                <span className='font-mono font-semibold text-foreground text-sm block'>
                  {log.ip_address || '127.0.0.1'}
                </span>
              </div>
            </div>
          </div>

          {/* User Agent */}
          {log.user_agent && (
            <div className='flex items-start gap-3 p-3.5 rounded-xl border bg-card text-xs text-muted-foreground font-mono'>
              <Laptop className='size-4 shrink-0 mt-0.5 text-primary' />
              <div className='min-w-0 space-y-0.5'>
                <span className='font-sans font-semibold text-foreground text-[11px] block'>
                  Header Client User-Agent:
                </span>
                <p className='break-all leading-relaxed'>{log.user_agent}</p>
              </div>
            </div>
          )}

          {/* Details / Description if any */}
          {log.details && (
            <div className='space-y-1.5'>
              <span className='font-semibold text-foreground text-xs uppercase tracking-wider block'>
                Keterangan Aktivitas:
              </span>
              <p className='p-3.5 rounded-xl bg-muted/20 text-foreground border leading-relaxed'>
                {log.details}
              </p>
            </div>
          )}

          {/* JSON Payload Viewer */}
          <div className='space-y-1.5'>
            <span className='font-semibold text-foreground text-xs uppercase tracking-wider flex items-center gap-1.5'>
              <FileCode className='size-3.5 text-primary' />
              Muatan Data Mentah (Raw JSON Payload):
            </span>
            <ScrollArea className='h-64 rounded-xl border bg-zinc-950 p-4 font-mono text-xs text-emerald-400 shadow-inner'>
              <pre className='whitespace-pre-wrap break-all leading-relaxed'>
                {log.payload ? JSON.stringify(log.payload, null, 2) : 'Tidak ada payload perubahan data.'}
              </pre>
            </ScrollArea>
          </div>
        </div>

        <DialogFooter className='pt-4 border-t'>
          <Button
            type='button'
            variant='outline'
            onClick={() => onOpenChange(false)}
            className='h-10 px-5 text-xs'
          >
            Tutup
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
