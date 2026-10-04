import { useState } from 'react'
import {
  MoreHorizontal,
  Edit,
  Trash2,
  Upload,
  BookOpen,
  Activity,
  ExternalLink,
  Search,
  CheckCircle2,
  AlertTriangle,
  XCircle,
  Clock,
  Shield,
  Users,
} from 'lucide-react'
import { toast } from 'sonner'
import { App, Category, StatusLayanan } from '@/types/apps'
import { appsApi } from '../api/apps-api'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

interface AppsTableProps {
  apps: App[]
  categories: Category[]
  isLoading: boolean
  onEdit: (app: App) => void
  onDelete: (app: App) => void
  onUploadLogo: (app: App) => void
  onManageGuides: (app: App) => void
  onRefresh: () => void
}

export function AppsTable({
  apps,
  categories,
  isLoading,
  onEdit,
  onDelete,
  onUploadLogo,
  onManageGuides,
  onRefresh,
}: AppsTableProps) {
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedCategory, setSelectedCategory] = useState<string>('all')
  const [selectedStatus, setSelectedStatus] = useState<string>('all')
  const [probingId, setProbingId] = useState<number | null>(null)

  const handleProbe = async (app: App) => {
    setProbingId(app.id)
    try {
      await appsApi.probeApp(app.id)
      toast.success(`Probe status layanan "${app.nama}" berhasil dikirim.`)
      onRefresh()
    } catch (err: any) {
      toast.error('Probe Gagal', {
        description: err?.message || 'Tidak dapat memverifikasi kesehatan URL.',
      })
    } finally {
      setProbingId(null)
    }
  }

  const filteredApps = apps.filter((app) => {
    const matchesSearch =
      app.nama.toLowerCase().includes(searchTerm.toLowerCase()) ||
      app.slug.toLowerCase().includes(searchTerm.toLowerCase()) ||
      (app.deskripsi && app.deskripsi.toLowerCase().includes(searchTerm.toLowerCase()))

    const matchesCategory =
      selectedCategory === 'all' || String(app.category_id) === selectedCategory

    const matchesStatus =
      selectedStatus === 'all' || app.status_layanan === selectedStatus

    return matchesSearch && matchesCategory && matchesStatus
  })

  const getStatusBadge = (status: StatusLayanan) => {
    switch (status) {
      case 'online':
        return (
          <Badge variant='outline' className='bg-emerald-500/10 text-emerald-600 border-emerald-500/30'>
            <CheckCircle2 className='mr-1 size-3' /> Online
          </Badge>
        )
      case 'pemeliharaan':
        return (
          <Badge variant='outline' className='bg-amber-500/10 text-amber-600 border-amber-500/30'>
            <Clock className='mr-1 size-3' /> Maintenance
          </Badge>
        )
      case 'kendala':
        return (
          <Badge variant='outline' className='bg-orange-500/10 text-orange-600 border-orange-500/30'>
            <AlertTriangle className='mr-1 size-3' /> Kendala
          </Badge>
        )
      case 'offline':
      default:
        return (
          <Badge variant='outline' className='bg-destructive/10 text-destructive border-destructive/30'>
            <XCircle className='mr-1 size-3' /> Offline
          </Badge>
        )
    }
  }

  const getTargetBadge = (target: string) => {
    if (target === 'Semua') {
      return (
        <Badge variant='secondary' className='text-[10px]'>
          <Users className='mr-1 size-2.5' /> Semua
        </Badge>
      )
    }
    return (
      <Badge variant='outline' className='text-[10px]'>
        <Shield className='mr-1 size-2.5' /> {target}
      </Badge>
    )
  }

  return (
    <div className='space-y-4'>
      {/* Filters Toolbar */}
      <div className='flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3'>
        <div className='flex flex-1 items-center gap-2'>
          <div className='relative flex-1 max-w-sm'>
            <Search className='absolute left-2.5 top-2.5 size-4 text-muted-foreground' />
            <Input
              placeholder='Cari nama aplikasi atau kata kunci...'
              className='pl-8 h-9 text-sm'
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
          </div>

          <Select value={selectedCategory} onValueChange={setSelectedCategory}>
            <SelectTrigger className='w-44 h-9 text-xs'>
              <SelectValue placeholder='Semua Kategori' />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value='all'>Semua Kategori</SelectItem>
              {categories.map((c) => (
                <SelectItem key={c.id} value={String(c.id)}>
                  {c.nama}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select value={selectedStatus} onValueChange={setSelectedStatus}>
            <SelectTrigger className='w-36 h-9 text-xs'>
              <SelectValue placeholder='Semua Status' />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value='all'>Semua Status</SelectItem>
              <SelectItem value='online'>Online</SelectItem>
              <SelectItem value='pemeliharaan'>Maintenance</SelectItem>
              <SelectItem value='kendala'>Kendala</SelectItem>
              <SelectItem value='offline'>Offline</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div className='text-xs text-muted-foreground self-end sm:self-center'>
          Menampilkan <span className='font-semibold text-foreground'>{filteredApps.length}</span> dari {apps.length} aplikasi
        </div>
      </div>

      {/* Table Content */}
      <div className='rounded-lg border bg-card shadow-xs overflow-hidden'>
        <Table>
          <TableHeader>
            <TableRow className='bg-muted/40'>
              <TableHead className='w-12 text-center'>#</TableHead>
              <TableHead>Aplikasi</TableHead>
              <TableHead>Kategori</TableHead>
              <TableHead>Target Sasaran</TableHead>
              <TableHead>Status Layanan</TableHead>
              <TableHead>Akses</TableHead>
              <TableHead className='text-end w-16'>Aksi</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={7} className='h-32 text-center text-muted-foreground'>
                  Memuat daftar aplikasi...
                </TableCell>
              </TableRow>
            ) : filteredApps.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} className='h-32 text-center text-muted-foreground'>
                  Tidak ada aplikasi yang sesuai dengan kriteria pencarian.
                </TableCell>
              </TableRow>
            ) : (
              filteredApps.map((app, index) => (
                <TableRow key={app.id} className='hover:bg-muted/30'>
                  <TableCell className='text-center text-xs font-mono text-muted-foreground'>
                    {index + 1}
                  </TableCell>
                  <TableCell>
                    <div className='flex items-center gap-3'>
                      {app.ikon_url ? (
                        <img
                          src={app.ikon_url}
                          alt={app.nama}
                          className='size-9 rounded-lg object-contain bg-white p-1 border shadow-xs shrink-0'
                        />
                      ) : (
                        <div className='size-9 rounded-lg bg-primary/10 text-primary flex items-center justify-center font-bold text-xs shrink-0'>
                          {app.nama.substring(0, 2).toUpperCase()}
                        </div>
                      )}
                      <div className='min-w-0'>
                        <div className='flex items-center gap-1.5'>
                          <span className='font-semibold text-sm text-foreground truncate'>
                            {app.nama}
                          </span>
                          {!app.aktif && (
                            <Badge variant='destructive' className='text-[9px] py-0 px-1'>
                              Nonaktif
                            </Badge>
                          )}
                        </div>
                        <a
                          href={app.url}
                          target='_blank'
                          rel='noopener noreferrer'
                          className='text-xs text-muted-foreground hover:text-primary flex items-center gap-1 truncate max-w-xs'
                        >
                          {app.url}
                          <ExternalLink className='size-2.5 shrink-0 opacity-70' />
                        </a>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant='outline' className='text-xs font-normal'>
                      {app.category?.nama || '-'}
                    </Badge>
                  </TableCell>
                  <TableCell>{getTargetBadge(app.target_pengguna)}</TableCell>
                  <TableCell>{getStatusBadge(app.status_layanan)}</TableCell>
                  <TableCell>
                    <div className='flex items-center gap-1'>
                      {app.is_public ? (
                        <Badge variant='secondary' className='text-[10px] bg-emerald-500/10 text-emerald-600'>
                          Publik
                        </Badge>
                      ) : (
                        <Badge variant='secondary' className='text-[10px] bg-blue-500/10 text-blue-600'>
                          SSO Saja
                        </Badge>
                      )}
                    </div>
                  </TableCell>
                  <TableCell className='text-end'>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant='ghost' size='icon' className='size-8'>
                          <MoreHorizontal className='size-4' />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align='end' className='w-48'>
                        <DropdownMenuLabel className='text-xs'>Aksi Aplikasi</DropdownMenuLabel>
                        <DropdownMenuItem onClick={() => onEdit(app)}>
                          <Edit className='mr-2 size-3.5' /> Edit Metadata
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => onUploadLogo(app)}>
                          <Upload className='mr-2 size-3.5' /> Upload Logo
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => onManageGuides(app)}>
                          <BookOpen className='mr-2 size-3.5' /> Panduan & FAQ
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          onClick={() => handleProbe(app)}
                          disabled={probingId === app.id}
                        >
                          <Activity className='mr-2 size-3.5' />
                          {probingId === app.id ? 'Memeriksa...' : 'Uji Probe Status'}
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem
                          onClick={() => window.open(app.url, '_blank')}
                        >
                          <ExternalLink className='mr-2 size-3.5' /> Kunjungi Web
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem
                          variant='destructive'
                          onClick={() => onDelete(app)}
                        >
                          <Trash2 className='mr-2 size-3.5' /> Nonaktifkan
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
