import { useState, useEffect } from 'react'
import {
  Plus,
  RefreshCw,
  Megaphone,
  Edit,
  Trash2,
  Info,
  AlertTriangle,
  AlertOctagon,
  CheckCircle2,
  Calendar,
  ExternalLink,
} from 'lucide-react'
import { Announcement, AnnouncementType } from '@/types/announcements'
import { App } from '@/types/apps'
import { announcementsApi } from './api/announcements-api'
import { appsApi } from '@/features/apps/api/apps-api'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { AnnouncementFormDialog } from './components/announcement-form-dialog'
import { AnnouncementDeleteDialog } from './components/announcement-delete-dialog'

export function Announcements() {
  const [announcements, setAnnouncements] = useState<Announcement[]>([])
  const [apps, setApps] = useState<App[]>([])
  const [isLoading, setIsLoading] = useState(true)

  const [isFormOpen, setIsFormOpen] = useState(false)
  const [announcementToEdit, setAnnouncementToEdit] = useState<Announcement | null>(null)

  const [isDeleteOpen, setIsDeleteOpen] = useState(false)
  const [announcementToDelete, setAnnouncementToDelete] = useState<Announcement | null>(null)

  const fetchData = async () => {
    setIsLoading(true)
    try {
      const [annRes, appsRes] = await Promise.all([
        announcementsApi.getAdminAnnouncements(1, 50),
        appsApi.getApps({ per_page: 100 }),
      ])
      setAnnouncements(annRes.announcements)
      setApps(appsRes.apps)
    } catch {
      setAnnouncements([])
      setApps([])
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  const handleOpenCreate = () => {
    setAnnouncementToEdit(null)
    setIsFormOpen(true)
  }

  const handleOpenEdit = (ann: Announcement) => {
    setAnnouncementToEdit(ann)
    setIsFormOpen(true)
  }

  const handleOpenDelete = (ann: Announcement) => {
    setAnnouncementToDelete(ann)
    setIsDeleteOpen(true)
  }

  const getTipeBadge = (tipe: AnnouncementType) => {
    switch (tipe) {
      case 'danger':
        return (
          <Badge variant='outline' className='bg-destructive/10 text-destructive border-destructive/30'>
            <AlertOctagon className='mr-1 size-3' /> Gangguan
          </Badge>
        )
      case 'warning':
        return (
          <Badge variant='outline' className='bg-amber-500/10 text-amber-600 border-amber-500/30'>
            <AlertTriangle className='mr-1 size-3' /> Maintenance
          </Badge>
        )
      case 'success':
        return (
          <Badge variant='outline' className='bg-emerald-500/10 text-emerald-600 border-emerald-500/30'>
            <CheckCircle2 className='mr-1 size-3' /> Sukses
          </Badge>
        )
      case 'info':
      default:
        return (
          <Badge variant='outline' className='bg-blue-500/10 text-blue-600 border-blue-500/30'>
            <Info className='mr-1 size-3' /> Informasi
          </Badge>
        )
    }
  }

  return (
    <>
      <Header>
        <div className='flex items-center gap-2'>
          <Megaphone className='size-5 text-primary' />
          <h2 className='text-sm font-semibold tracking-tight'>Broadcast Pengumuman</h2>
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
              <h1 className='text-2xl font-bold tracking-tight'>Pengumuman & Banner Broadcast</h1>
              <p className='text-sm text-muted-foreground'>
                Kelola pesan peringatan pemeliharaan sistem, informasi rilis fitur, dan pengumuman bagi pengguna portal.
              </p>
            </div>

            <div className='flex items-center gap-2'>
              <Button
                variant='outline'
                size='sm'
                onClick={fetchData}
                disabled={isLoading}
              >
                <RefreshCw className={`mr-1.5 size-3.5 ${isLoading ? 'animate-spin' : ''}`} />
                Refresh
              </Button>
              <Button size='sm' onClick={handleOpenCreate}>
                <Plus className='mr-1.5 size-4' />
                Buat Pengumuman
              </Button>
            </div>
          </div>

          {/* Announcements Grid */}
          <div className='grid gap-4'>
            {isLoading ? (
              <div className='p-12 text-center text-sm text-muted-foreground border rounded-xl border-dashed'>
                Memuat daftar pengumuman...
              </div>
            ) : announcements.length === 0 ? (
              <div className='p-12 text-center text-sm text-muted-foreground border rounded-xl border-dashed'>
                Belum ada pengumuman broadcast yang dibuat.
              </div>
            ) : (
              announcements.map((ann) => (
                <Card
                  key={ann.id}
                  className={`border shadow-xs transition-colors ${
                    !ann.aktif ? 'opacity-60 bg-muted/20' : ''
                  }`}
                >
                  <CardHeader className='p-4 pb-2'>
                    <div className='flex flex-col sm:flex-row sm:items-center justify-between gap-2'>
                      <div className='flex items-center gap-2'>
                        {getTipeBadge(ann.tipe)}
                        <CardTitle className='text-base font-semibold'>
                          {ann.judul}
                        </CardTitle>
                        {!ann.aktif && (
                          <Badge variant='secondary' className='text-[10px]'>
                            Nonaktif
                          </Badge>
                        )}
                        {ann.app ? (
                          <Badge variant='outline' className='text-[10px]'>
                            Aplikasi: {ann.app.nama}
                          </Badge>
                        ) : (
                          <Badge variant='outline' className='text-[10px] bg-primary/5 text-primary'>
                            🌐 Global Portal
                          </Badge>
                        )}
                      </div>

                      <div className='flex items-center gap-1 self-end sm:self-auto'>
                        <Button
                          variant='ghost'
                          size='icon'
                          className='size-7'
                          onClick={() => handleOpenEdit(ann)}
                        >
                          <Edit className='size-3.5' />
                        </Button>
                        <Button
                          variant='ghost'
                          size='icon'
                          className='size-7 text-destructive hover:text-destructive'
                          onClick={() => handleOpenDelete(ann)}
                        >
                          <Trash2 className='size-3.5' />
                        </Button>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent className='p-4 pt-2 space-y-2'>
                    <p className='text-sm text-muted-foreground leading-relaxed whitespace-pre-line'>
                      {ann.pesan}
                    </p>

                    <div className='flex flex-wrap items-center gap-4 text-xs text-muted-foreground pt-1 border-t'>
                      {ann.tautan_url && (
                        <a
                          href={ann.tautan_url}
                          target='_blank'
                          rel='noopener noreferrer'
                          className='text-primary hover:underline flex items-center gap-1 font-medium'
                        >
                          {ann.tautan_teks || 'Buka Tautan'}
                          <ExternalLink className='size-3' />
                        </a>
                      )}
                      {(ann.mulai_pada || ann.berakhir_pada) && (
                        <div className='flex items-center gap-1 text-[11px]'>
                          <Calendar className='size-3' />
                          <span>
                            Jadwal:{' '}
                            {ann.mulai_pada
                              ? new Date(ann.mulai_pada).toLocaleDateString()
                              : 'Sekarang'}{' '}
                            -{' '}
                            {ann.berakhir_pada
                              ? new Date(ann.berakhir_pada).toLocaleDateString()
                              : 'Seterusnya'}
                          </span>
                        </div>
                      )}
                    </div>
                  </CardContent>
                </Card>
              ))
            )}
          </div>
        </div>

        <AnnouncementFormDialog
          open={isFormOpen}
          onOpenChange={setIsFormOpen}
          announcementToEdit={announcementToEdit}
          apps={apps}
          onSuccess={fetchData}
        />

        <AnnouncementDeleteDialog
          open={isDeleteOpen}
          onOpenChange={setIsDeleteOpen}
          announcement={announcementToDelete}
          onSuccess={fetchData}
        />
      </Main>
    </>
  )
}
