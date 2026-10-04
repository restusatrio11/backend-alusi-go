import { useState, useEffect } from 'react'
import { Plus, BookOpen, Trash2, Edit2, Loader2, FileText } from 'lucide-react'
import { toast } from 'sonner'
import { App, AppGuide } from '@/types/apps'
import { guidesApi } from '@/features/guides/api/guides-api'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'

interface AppGuidesDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  app: App | null
}

export function AppGuidesDialog({
  open,
  onOpenChange,
  app,
}: AppGuidesDialogProps) {
  const [guides, setGuides] = useState<AppGuide[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [activeGuideId, setActiveGuideId] = useState<number | null>(null)
  const [judul, setJudul] = useState('')
  const [konten, setKonten] = useState('')
  const [isSaving, setIsSaving] = useState(false)

  const fetchGuides = async () => {
    if (!app) return
    setIsLoading(true)
    try {
      const data = await guidesApi.getGuidesByAppSlug(app.slug)
      setGuides(data)
    } catch {
      setGuides([])
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    if (open && app) {
      fetchGuides()
      resetForm()
    }
  }, [open, app])

  const resetForm = () => {
    setActiveGuideId(null)
    setJudul('')
    setKonten('')
  }

  const handleStartEdit = (guide: AppGuide) => {
    setActiveGuideId(guide.id)
    setJudul(guide.judul)
    setKonten(guide.konten)
  }

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!app || !judul.trim() || !konten.trim()) return

    setIsSaving(true)
    try {
      if (activeGuideId) {
        await guidesApi.updateGuide(activeGuideId, {
          judul,
          konten,
        })
        toast.success('Panduan berhasil diperbarui!')
      } else {
        await guidesApi.createGuide(app.id, {
          judul,
          konten,
          urutan: guides.length + 1,
        })
        toast.success('Panduan baru berhasil ditambahkan!')
      }
      resetForm()
      fetchGuides()
    } catch (err: any) {
      toast.error('Gagal Menyimpan Panduan', {
        description: err?.message || 'Terjadi kesalahan.',
      })
    } finally {
      setIsSaving(false)
    }
  }

  const handleDelete = async (guideId: number) => {
    if (!confirm('Apakah Anda yakin ingin menghapus panduan ini?')) return

    try {
      await guidesApi.deleteGuide(guideId)
      toast.success('Panduan berhasil dihapus.')
      fetchGuides()
      if (activeGuideId === guideId) {
        resetForm()
      }
    } catch (err: any) {
      toast.error('Gagal Menghapus Panduan', {
        description: err?.message || 'Terjadi kesalahan.',
      })
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='max-w-3xl max-h-[85vh] overflow-y-auto'>
        <DialogHeader>
          <div className='flex items-center gap-2 text-primary'>
            <BookOpen className='size-5' />
            <DialogTitle>Panduan & FAQ: {app?.nama}</DialogTitle>
          </div>
          <DialogDescription>
            Kelola petunjuk teknis, FAQ, dan manual penggunaan untuk aplikasi ini.
          </DialogDescription>
        </DialogHeader>

        <div className='grid grid-cols-1 md:grid-cols-2 gap-6 py-2'>
          {/* List of Guides */}
          <div className='space-y-3'>
            <div className='flex items-center justify-between'>
              <h4 className='text-sm font-semibold text-foreground flex items-center gap-1.5'>
                <FileText className='size-4 text-muted-foreground' />
                Daftar Panduan ({guides.length})
              </h4>
              <Button
                variant='outline'
                size='sm'
                onClick={resetForm}
                className='text-xs h-7 px-2'
              >
                <Plus className='size-3 mr-1' /> Tambah Baru
              </Button>
            </div>

            {isLoading ? (
              <div className='flex items-center justify-center p-8'>
                <Loader2 className='size-6 animate-spin text-primary' />
              </div>
            ) : guides.length === 0 ? (
              <div className='p-6 text-center rounded-lg border border-dashed text-xs text-muted-foreground'>
                Belum ada panduan atau FAQ untuk aplikasi ini.
              </div>
            ) : (
              <div className='space-y-2 max-h-[50vh] overflow-y-auto pr-1'>
                {guides.map((guide, idx) => (
                  <Card
                    key={guide.id}
                    className={`transition-colors cursor-pointer border ${
                      activeGuideId === guide.id ? 'border-primary bg-primary/5' : 'hover:border-muted-foreground/40'
                    }`}
                    onClick={() => handleStartEdit(guide)}
                  >
                    <CardHeader className='p-3 pb-1'>
                      <div className='flex items-start justify-between gap-2'>
                        <CardTitle className='text-xs font-semibold leading-snug line-clamp-1'>
                          {idx + 1}. {guide.judul}
                        </CardTitle>
                        <div className='flex items-center gap-1 shrink-0'>
                          <Button
                            variant='ghost'
                            size='icon'
                            className='size-6 text-muted-foreground hover:text-foreground'
                            onClick={(e) => {
                              e.stopPropagation()
                              handleStartEdit(guide)
                            }}
                          >
                            <Edit2 className='size-3' />
                          </Button>
                          <Button
                            variant='ghost'
                            size='icon'
                            className='size-6 text-destructive/80 hover:text-destructive'
                            onClick={(e) => {
                              e.stopPropagation()
                              handleDelete(guide.id)
                            }}
                          >
                            <Trash2 className='size-3' />
                          </Button>
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent className='p-3 pt-0'>
                      <p className='text-[11px] text-muted-foreground line-clamp-2'>
                        {guide.konten}
                      </p>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </div>

          {/* Form Editor */}
          <div className='space-y-3 border-l pl-4'>
            <div className='flex items-center justify-between'>
              <h4 className='text-sm font-semibold text-foreground'>
                {activeGuideId ? 'Edit Panduan' : 'Form Panduan Baru'}
              </h4>
              {activeGuideId && (
                <Badge variant='outline' className='text-[10px]'>
                  Mode Edit
                </Badge>
              )}
            </div>

            <form onSubmit={handleSave} className='space-y-3'>
              <div>
                <label className='text-xs font-medium mb-1 block'>
                  Judul Panduan / Pertanyaan FAQ *
                </label>
                <Input
                  value={judul}
                  onChange={(e) => setJudul(e.target.value)}
                  placeholder='Contoh: Cara Login dengan Akun BPS'
                  required
                />
              </div>

              <div>
                <label className='text-xs font-medium mb-1 block'>
                  Isi Panduan / Jawaban (Markdown) *
                </label>
                <Textarea
                  value={konten}
                  onChange={(e) => setKonten(e.target.value)}
                  placeholder='Tuliskan langkah-langkah penggunaan atau penjelasan detail...'
                  rows={8}
                  required
                />
                <p className='text-[10px] text-muted-foreground mt-1'>
                  Mendukung sintaks Markdown (heading, list, bold, link).
                </p>
              </div>

              <div className='flex items-center justify-end gap-2 pt-2'>
                {activeGuideId && (
                  <Button
                    type='button'
                    variant='outline'
                    size='sm'
                    onClick={resetForm}
                  >
                    Batal
                  </Button>
                )}
                <Button type='submit' size='sm' disabled={isSaving}>
                  {isSaving && <Loader2 className='mr-1.5 size-3.5 animate-spin' />}
                  {activeGuideId ? 'Simpan Perubahan' : 'Tambah Panduan'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
