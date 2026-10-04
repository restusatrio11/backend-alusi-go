import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  ShieldAlert,
  Search,
  RefreshCw,
  Eye,
  Layers,
} from 'lucide-react'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { auditApi } from './api/audit-api'
import { AuditLog, AuditLogFilterParams } from '@/types/audit'
import { AuditDetailDialog } from './components/audit-detail-dialog'

export function AuditLogs() {
  const [selectedLog, setSelectedLog] = useState<AuditLog | null>(null)
  const [isDetailOpen, setIsDetailOpen] = useState(false)

  // Filters
  const [searchQuery, setSearchQuery] = useState('')
  const [actionFilter, setActionFilter] = useState<string>('ALL')
  const [entityFilter, setEntityFilter] = useState<string>('ALL')
  const [page, setPage] = useState<number>(1)
  const perPage = 15

  const filterParams: AuditLogFilterParams = {
    page,
    per_page: perPage,
    action: actionFilter !== 'ALL' ? actionFilter : undefined,
    entity: entityFilter !== 'ALL' ? entityFilter : undefined,
  }

  const { data, isLoading, refetch, isRefetching } = useQuery({
    queryKey: ['audit-logs', page, actionFilter, entityFilter],
    queryFn: () => auditApi.getAuditLogs(filterParams),
  })

  const rawLogs = data?.logs || []

  // Client-side search filter for actor / ip / details
  const logs = rawLogs.filter((log) => {
    if (!searchQuery.trim()) return true
    const q = searchQuery.toLowerCase()
    return (
      (log.user_name && log.user_name.toLowerCase().includes(q)) ||
      (log.user_email && log.user_email.toLowerCase().includes(q)) ||
      (log.ip_address && log.ip_address.includes(q)) ||
      (log.details && log.details.toLowerCase().includes(q)) ||
      (log.entity_id && String(log.entity_id).includes(q))
    )
  })

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

  const handleOpenDetail = (log: AuditLog) => {
    setSelectedLog(log)
    setIsDetailOpen(true)
  }

  return (
    <>
      <Header>
        <div className='flex items-center gap-2'>
          <ShieldAlert className='size-5 text-primary' />
          <h2 className='text-sm font-semibold tracking-tight'>Jejak Audit Sistem</h2>
        </div>
        <div className='ml-auto flex items-center space-x-2'>
          <ThemeSwitch />
          <ProfileDropdown />
        </div>
      </Header>

      <Main>
        <div className='space-y-6'>
          {/* Header Title */}
          <div className='flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4'>
            <div>
              <h1 className='text-2xl font-bold tracking-tight'>Log Forensik & Jejak Audit</h1>
              <p className='text-sm text-muted-foreground'>
                Pencatatan riwayat setiap aksi manipulasi data, perubahan status, dan login administrator.
              </p>
            </div>

            <Button
              variant='outline'
              size='sm'
              onClick={() => refetch()}
              disabled={isRefetching}
              className='h-8 text-xs'
            >
              <RefreshCw
                className={`mr-1.5 size-3.5 ${isRefetching ? 'animate-spin' : ''}`}
              />
              Refresh Log
            </Button>
          </div>

          {/* Filter Bar */}
          <div className='flex flex-col sm:flex-row items-stretch sm:items-center gap-2.5 p-3 rounded-xl border bg-card/60'>
            <div className='relative flex-1'>
              <Search className='absolute left-2.5 top-2.5 size-4 text-muted-foreground' />
              <Input
                placeholder='Cari aktor, IP, atau keterangan...'
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className='pl-8 h-9 text-xs'
              />
            </div>

            <div className='flex items-center gap-2'>
              <Select value={actionFilter} onValueChange={setActionFilter}>
                <SelectTrigger className='w-[130px] h-9 text-xs'>
                  <SelectValue placeholder='Aksi' />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value='ALL'>Semua Aksi</SelectItem>
                  <SelectItem value='CREATE'>CREATE</SelectItem>
                  <SelectItem value='UPDATE'>UPDATE</SelectItem>
                  <SelectItem value='DELETE'>DELETE</SelectItem>
                  <SelectItem value='PROBE'>PROBE</SelectItem>
                  <SelectItem value='LOGIN'>LOGIN</SelectItem>
                </SelectContent>
              </Select>

              <Select value={entityFilter} onValueChange={setEntityFilter}>
                <SelectTrigger className='w-[140px] h-9 text-xs'>
                  <SelectValue placeholder='Entitas' />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value='ALL'>Semua Entitas</SelectItem>
                  <SelectItem value='app'>Aplikasi</SelectItem>
                  <SelectItem value='category'>Kategori</SelectItem>
                  <SelectItem value='guide'>Panduan / FAQ</SelectItem>
                  <SelectItem value='announcement'>Pengumuman</SelectItem>
                  <SelectItem value='feedback'>Umpan Balik</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          {/* Audit Logs Table */}
          <div className='rounded-xl border bg-card overflow-hidden shadow-xs'>
            <div className='overflow-x-auto'>
              <table className='w-full text-xs'>
                <thead>
                  <tr className='border-b bg-muted/40 text-muted-foreground font-medium text-left'>
                    <th className='py-3 pl-4'>Waktu</th>
                    <th className='py-3 px-3'>Aktor</th>
                    <th className='py-3 px-3 text-center'>Aksi</th>
                    <th className='py-3 px-3'>Entitas Target</th>
                    <th className='py-3 px-3'>Alamat IP</th>
                    <th className='py-3 px-3'>Keterangan</th>
                    <th className='py-3 pr-4 text-right'>Opsi</th>
                  </tr>
                </thead>
                <tbody className='divide-y divide-border/60'>
                  {isLoading ? (
                    <tr>
                      <td colSpan={7} className='py-12 text-center text-muted-foreground'>
                        Memuat jejak audit...
                      </td>
                    </tr>
                  ) : logs.length === 0 ? (
                    <tr>
                      <td colSpan={7} className='py-12 text-center text-muted-foreground'>
                        <div className='flex flex-col items-center gap-1.5'>
                          <Layers className='size-8 text-muted-foreground/30' />
                          <span>Tidak ada catatan audit yang cocok dengan kriteria filter.</span>
                        </div>
                      </td>
                    </tr>
                  ) : (
                    logs.map((log) => (
                      <tr
                        key={log.id}
                        className='hover:bg-muted/40 transition-colors'
                      >
                        {/* Waktu */}
                        <td className='py-3 pl-4 font-mono text-muted-foreground whitespace-nowrap'>
                          {new Date(log.created_at).toLocaleString('id-ID', {
                            day: '2-digit',
                            month: 'short',
                            year: 'numeric',
                            hour: '2-digit',
                            minute: '2-digit',
                            second: '2-digit',
                          })}
                        </td>

                        {/* Aktor */}
                        <td className='py-3 px-3 font-medium text-foreground'>
                          <div>
                            <span className='block'>{log.user_name || 'System'}</span>
                            {log.user_email && (
                              <span className='text-[10px] text-muted-foreground block font-mono'>
                                {log.user_email}
                              </span>
                            )}
                          </div>
                        </td>

                        {/* Aksi */}
                        <td className='py-3 px-3 text-center'>
                          <Badge
                            variant='outline'
                            className={`font-mono text-[10px] px-2 py-0.5 border-0 ${getActionBadgeVariant(
                              log.action
                            )}`}
                          >
                            {log.action}
                          </Badge>
                        </td>

                        {/* Entitas */}
                        <td className='py-3 px-3'>
                          <span className='font-mono font-medium text-foreground uppercase'>
                            {log.entity}
                          </span>
                          {log.entity_id && (
                            <span className='text-muted-foreground ml-1 font-mono text-[11px]'>
                              #{log.entity_id}
                            </span>
                          )}
                        </td>

                        {/* IP Address */}
                        <td className='py-3 px-3 font-mono text-muted-foreground'>
                          {log.ip_address || '-'}
                        </td>

                        {/* Keterangan */}
                        <td className='py-3 px-3 max-w-[240px] truncate text-muted-foreground'>
                          {log.details || '-'}
                        </td>

                        {/* Actions */}
                        <td className='py-3 pr-4 text-right'>
                          <Button
                            variant='ghost'
                            size='sm'
                            className='h-7 px-2 text-xs gap-1'
                            onClick={() => handleOpenDetail(log)}
                          >
                            <Eye className='size-3.5' />
                            <span>Detail</span>
                          </Button>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>

            {/* Pagination footer */}
            <div className='flex items-center justify-between px-4 py-3 border-t bg-muted/20 text-xs'>
              <span className='text-muted-foreground'>
                Halaman <span className='font-bold text-foreground'>{page}</span> (Menampilkan {logs.length} data)
              </span>
              <div className='flex items-center gap-2'>
                <Button
                  variant='outline'
                  size='sm'
                  className='h-7 text-xs'
                  disabled={page <= 1 || isLoading}
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                >
                  Sebelumnya
                </Button>
                <Button
                  variant='outline'
                  size='sm'
                  className='h-7 text-xs'
                  disabled={logs.length < perPage || isLoading}
                  onClick={() => setPage((p) => p + 1)}
                >
                  Selanjutnya
                </Button>
              </div>
            </div>
          </div>
        </div>
      </Main>

      {/* Detail Dialog */}
      <AuditDetailDialog
        log={selectedLog}
        open={isDetailOpen}
        onOpenChange={setIsDetailOpen}
      />
    </>
  )
}
