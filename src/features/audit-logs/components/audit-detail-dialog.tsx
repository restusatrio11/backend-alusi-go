import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { AuditLog } from '@/types/audit'
import { ScrollArea } from '@/components/ui/scroll-area'
import { ShieldCheck, Laptop, Globe, Calendar, User } from 'lucide-react'

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
      <DialogContent className='sm:max-w-2xl'>
        <DialogHeader>
          <div className='flex items-center gap-2'>
            <ShieldCheck className='size-5 text-primary' />
            <DialogTitle className='text-base font-bold'>
              Detail Jejak Audit #{log.id}
            </DialogTitle>
          </div>
          <DialogDescription className='text-xs'>
            Rekaman aktivitas forensik sistem dan manipulasi data.
          </DialogDescription>
        </DialogHeader>

        <div className='space-y-4 py-2 text-xs'>
          {/* Metadata Grid */}
          <div className='grid grid-cols-2 gap-3 p-3 bg-muted/40 rounded-lg border'>
            <div className='flex items-start gap-2'>
              <User className='size-4 text-muted-foreground shrink-0 mt-0.5' />
              <div>
                <span className='text-muted-foreground text-[11px] block'>Aktor / Pengguna</span>
                <span className='font-semibold text-foreground'>
                  {log.user_name || 'System / Anonymous'}
                </span>
                {log.user_email && (
                  <span className='text-[11px] text-muted-foreground block font-mono'>
                    {log.user_email}
                  </span>
                )}
              </div>
            </div>

            <div className='flex items-start gap-2'>
              <Calendar className='size-4 text-muted-foreground shrink-0 mt-0.5' />
              <div>
                <span className='text-muted-foreground text-[11px] block'>Waktu Kejadian</span>
                <span className='font-medium font-mono text-foreground'>
                  {new Date(log.created_at).toLocaleString('id-ID')}
                </span>
              </div>
            </div>

            <div className='flex items-start gap-2'>
              <ShieldCheck className='size-4 text-muted-foreground shrink-0 mt-0.5' />
              <div>
                <span className='text-muted-foreground text-[11px] block'>Aksi & Entitas</span>
                <div className='flex items-center gap-1.5 mt-0.5'>
                  <Badge
                    variant='outline'
                    className={`font-mono text-[10px] px-1.5 py-0 border-0 ${getActionBadgeVariant(
                      log.action
                    )}`}
                  >
                    {log.action}
                  </Badge>
                  <span className='font-medium text-foreground uppercase'>
                    {log.entity} {log.entity_id ? `(#${log.entity_id})` : ''}
                  </span>
                </div>
              </div>
            </div>

            <div className='flex items-start gap-2'>
              <Globe className='size-4 text-muted-foreground shrink-0 mt-0.5' />
              <div>
                <span className='text-muted-foreground text-[11px] block'>Alamat IP</span>
                <span className='font-mono text-foreground'>
                  {log.ip_address || '127.0.0.1'}
                </span>
              </div>
            </div>
          </div>

          {/* User Agent */}
          {log.user_agent && (
            <div className='flex items-start gap-2 p-2.5 rounded-lg border bg-card text-[11px] text-muted-foreground font-mono'>
              <Laptop className='size-4 shrink-0 mt-0.5' />
              <span className='break-all'>{log.user_agent}</span>
            </div>
          )}

          {/* Details / Description if any */}
          {log.details && (
            <div>
              <span className='font-semibold text-foreground mb-1 block'>Keterangan:</span>
              <p className='p-2.5 rounded-lg bg-muted text-foreground border'>
                {log.details}
              </p>
            </div>
          )}

          {/* JSON Payload Viewer */}
          <div>
            <span className='font-semibold text-foreground mb-1 block'>
              Muatan Data (JSON Payload):
            </span>
            <ScrollArea className='h-48 rounded-lg border bg-zinc-950 p-3 font-mono text-[11px] text-emerald-400'>
              <pre className='whitespace-pre-wrap break-all'>
                {log.payload ? JSON.stringify(log.payload, null, 2) : 'Tidak ada payload data.'}
              </pre>
            </ScrollArea>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
