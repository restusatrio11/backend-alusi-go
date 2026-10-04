import { useState, useEffect } from 'react'
import {
  RefreshCw,
  MessageSquareWarning,
  Eye,
  CheckCircle2,
  Clock,
  AlertCircle,
  XCircle,
  Search,
  Layers,
} from 'lucide-react'
import { Feedback, FeedbackStatus } from '@/types/feedbacks'
import { feedbacksApi } from './api/feedbacks-api'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { FeedbackDetailDialog } from './components/feedback-detail-dialog'

export function Feedbacks() {
  const [feedbacks, setFeedbacks] = useState<Feedback[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState('')
  const [activeTab, setActiveTab] = useState('all')

  const [isDetailOpen, setIsDetailOpen] = useState(false)
  const [selectedFeedback, setSelectedFeedback] = useState<Feedback | null>(null)

  const fetchFeedbacks = async () => {
    setIsLoading(true)
    try {
      const res = await feedbacksApi.getFeedbacks({
        status: activeTab === 'all' ? undefined : activeTab,
        per_page: 50,
      })
      setFeedbacks(res.feedbacks)
    } catch {
      setFeedbacks([])
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchFeedbacks()
  }, [activeTab])

  const handleOpenDetail = (fb: Feedback) => {
    setSelectedFeedback(fb)
    setIsDetailOpen(true)
  }

  const filteredFeedbacks = feedbacks.filter((fb) => {
    const matchesSearch =
      fb.judul.toLowerCase().includes(searchTerm.toLowerCase()) ||
      fb.isi_laporan.toLowerCase().includes(searchTerm.toLowerCase()) ||
      fb.nama_pelapor.toLowerCase().includes(searchTerm.toLowerCase())
    return matchesSearch
  })

  const getStatusBadge = (st: FeedbackStatus) => {
    switch (st) {
      case 'resolved':
        return (
          <Badge variant='outline' className='bg-emerald-500/10 text-emerald-600 border-emerald-500/30 text-xs'>
            <CheckCircle2 className='mr-1 size-3' /> Selesai
          </Badge>
        )
      case 'in_progress':
        return (
          <Badge variant='outline' className='bg-blue-500/10 text-blue-600 border-blue-500/30 text-xs'>
            <Clock className='mr-1 size-3' /> Diproses
          </Badge>
        )
      case 'closed':
        return (
          <Badge variant='outline' className='bg-muted text-muted-foreground text-xs'>
            <XCircle className='mr-1 size-3' /> Ditutup
          </Badge>
        )
      case 'pending':
      default:
        return (
          <Badge variant='outline' className='bg-amber-500/10 text-amber-600 border-amber-500/30 text-xs'>
            <AlertCircle className='mr-1 size-3' /> Pending
          </Badge>
        )
    }
  }

  return (
    <>
      <Header>
        <div className='flex items-center gap-2'>
          <MessageSquareWarning className='size-5 text-primary' />
          <h2 className='text-sm font-semibold tracking-tight'>Umpan Balik & Masalah</h2>
        </div>
        <div className='ml-auto flex items-center space-x-2'>
          <ThemeSwitch />
          <ProfileDropdown />
        </div>
      </Header>

      <Main>
        <div className='space-y-6'>
          <div className='flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4'>
            <div>
              <h1 className='text-2xl font-bold tracking-tight'>Kotak Masukan & Laporan Kendala</h1>
              <p className='text-sm text-muted-foreground'>
                Pantau saran, kritik, dan laporan bug dari pengguna aplikasi untuk ditindaklanjuti.
              </p>
            </div>

            <Button
              variant='outline'
              size='sm'
              onClick={fetchFeedbacks}
              disabled={isLoading}
            >
              <RefreshCw className={`mr-1.5 size-3.5 ${isLoading ? 'animate-spin' : ''}`} />
              Refresh
            </Button>
          </div>

          {/* Status Tabs and Search */}
          <div className='flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3'>
            <Tabs value={activeTab} onValueChange={setActiveTab}>
              <TabsList>
                <TabsTrigger value='all'>Semua</TabsTrigger>
                <TabsTrigger value='pending'>Pending</TabsTrigger>
                <TabsTrigger value='in_progress'>Diproses</TabsTrigger>
                <TabsTrigger value='resolved'>Selesai</TabsTrigger>
                <TabsTrigger value='closed'>Ditutup</TabsTrigger>
              </TabsList>
            </Tabs>

            <div className='relative w-full sm:w-72'>
              <Search className='absolute left-2.5 top-2.5 size-4 text-muted-foreground' />
              <Input
                placeholder='Cari laporan atau pelapor...'
                className='pl-8 h-9 text-xs'
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
            </div>
          </div>

          {/* Table */}
          <div className='rounded-lg border bg-card shadow-xs overflow-hidden'>
            <Table>
              <TableHeader>
                <TableRow className='bg-muted/40'>
                  <TableHead className='w-14 text-center'>Tiket</TableHead>
                  <TableHead>Pelapor</TableHead>
                  <TableHead>Kategori</TableHead>
                  <TableHead>Aplikasi Terkait</TableHead>
                  <TableHead>Judul Masalah</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Waktu</TableHead>
                  <TableHead className='text-end w-20'>Aksi</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={8} className='h-32 text-center text-muted-foreground'>
                      Memuat daftar masukan...
                    </TableCell>
                  </TableRow>
                ) : filteredFeedbacks.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={8} className='h-32 text-center text-muted-foreground'>
                      Tidak ada laporan tiket masukan pada kategori ini.
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredFeedbacks.map((fb) => (
                    <TableRow key={fb.id} className='hover:bg-muted/30'>
                      <TableCell className='text-center font-mono text-xs font-bold'>
                        #{fb.id}
                      </TableCell>
                      <TableCell>
                        <div>
                          <p className='font-semibold text-xs text-foreground'>{fb.nama_pelapor}</p>
                          <p className='text-[11px] text-muted-foreground truncate'>{fb.email_pelapor || '-'}</p>
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant='secondary' className='text-[10px] uppercase font-mono'>
                          {fb.kategori}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        {fb.app ? (
                          <div className='flex items-center gap-1 text-xs font-medium'>
                            <Layers className='size-3 text-muted-foreground' />
                            <span>{fb.app.nama}</span>
                          </div>
                        ) : (
                          <span className='text-xs text-muted-foreground'>- (Portal Umum)</span>
                        )}
                      </TableCell>
                      <TableCell className='max-w-xs'>
                        <p className='text-xs font-semibold text-foreground truncate'>{fb.judul}</p>
                        <p className='text-[11px] text-muted-foreground truncate'>{fb.isi_laporan}</p>
                      </TableCell>
                      <TableCell>{getStatusBadge(fb.status)}</TableCell>
                      <TableCell className='text-xs text-muted-foreground font-mono'>
                        {new Date(fb.created_at).toLocaleDateString('id-ID')}
                      </TableCell>
                      <TableCell className='text-end'>
                        <Button
                          variant='ghost'
                          size='sm'
                          className='text-xs h-7 gap-1'
                          onClick={() => handleOpenDetail(fb)}
                        >
                          <Eye className='size-3.5' /> Respon
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </div>

        <FeedbackDetailDialog
          open={isDetailOpen}
          onOpenChange={setIsDetailOpen}
          feedback={selectedFeedback}
          onSuccess={fetchFeedbacks}
        />
      </Main>
    </>
  )
}
