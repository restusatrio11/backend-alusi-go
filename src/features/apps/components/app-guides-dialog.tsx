import { useState, useEffect } from 'react'
import { Plus, BookOpen, Trash2, Edit2, Loader2, FileText, Info, HelpCircle } from 'lucide-react'
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
      <DialogContent className='sm:max-w-4xl lg:max-w-5xl max-h-[90vh] overflow-y-auto p-6 md:p-8'>
        <DialogHeader className='space-y-1.5 pb-2 border-b'>
          <div className='flex items-center gap-2 text-primary'>
            <BookOpen className='size-5' />
            <DialogTitle className='text-xl font-bold tracking-tight'>
              Kelola Panduan Pengguna & FAQ: {app?.nama}
            </DialogTitle>
          </div>
          <DialogDescription className='text-xs text-muted-foreground leading-relaxed'>
            Kelola petunjuk penggunaan teknis, prosedur operasional standar (SOP), dan tanya-jawab umum (FAQ) bagi pengguna aplikasi.
          </DialogDescription>
        </DialogHeader>

        <div className='grid grid-cols-1 lg:grid-cols-12 gap-6 py-3'>
          {/* List of Guides (5 cols) */}
          <div className='lg:col-span-5 space-y-3'>
            <div className='flex items-center justify-between'>
              <h4 className='text-xs font-semibold text-foreground uppercase tracking-wider flex items-center gap-1.5'>
                <FileText className='size-3.5 text-primary' />
                Daftar Dokumen Panduan ({guides.length})
              </h4>
              <Button
                variant='outline'
                size='sm'
                onClick={resetForm}
                className='text-xs h-8 px-2.5 gap-1'
              >
                <Plus className='size-3.5' /> Tambah Baru
              </Button>
            </div>

            {isLoading ? (
              <div className='flex items-center justify-center p-12 border rounded-xl bg-muted/10'>
                <Loader2 className='size-6 animate-spin text-primary' />
              </div>
            ) : guides.length === 0 ? (
              <div className='p-8 text-center rounded-2xl border-2 border-dashed text-xs text-muted-foreground space-y-1.5 bg-muted/10'>
                <HelpCircle className='size-8 text-muted-foreground/40 mx-auto' />
                <p className='font-medium text-foreground'>Belum ada panduan terdaftar</p>
                <p>Klik tombol &quot;Tambah Baru&quot; untuk menyusun panduan pertama.</p>
              </div>
            ) : (
              <div className='space-y-2.5 max-h-[55vh] overflow-y-auto pr-1'>
                {guides.map((guide, idx) => (
                  <Card
                    key={guide.id}
                    className={`transition-all cursor-pointer border shadow-xs ${
                      activeGuideId === guide.id
                        ? 'border-primary ring-1 ring-primary/30 bg-primary/5'
                        : 'hover:border-muted-foreground/40 bg-card hover:bg-muted/30'
                    }`}
                    onClick={() => handleStartEdit(guide)}
                  >
                    <CardHeader className='p-3.5 pb-1'>
                      <div className='flex items-start justify-between gap-2'>
                        <CardTitle className='text-xs font-semibold leading-snug line-clamp-1'>
                          {idx + 1}. {guide.judul}
                        </CardTitle>
                        <div className='flex items-center gap-1 shrink-0'>
                          <Button
                            variant='ghost'
                            size='icon'
                            className='size-7 text-muted-foreground hover:text-foreground'
                            onClick={(e) => {
                              e.stopPropagation()
                              handleStartEdit(guide)
                            }}
                            title='Edit'
                          >
                            <Edit2 className='size-3.5' />
                          </Button>
                          <Button
                            variant='ghost'
                            size='icon'
                            className='size-7 text-destructive/80 hover:text-destructive hover:bg-destructive/10'
                            onClick={(e) => {
                              e.stopPropagation()
                              handleDelete(guide.id)
                            }}
                            title='Hapus'
                          >
                            <Trash2 className='size-3.5' />
                          </Button>
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent className='p-3.5 pt-0'>
                      <p className='text-[11px] text-muted-foreground line-clamp-2 leading-relaxed'>
                        {guide.konten}
                      </p>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </div>

          {/* Form Editor (7 cols) */}
          <div className='lg:col-span-7 space-y-3 lg:border-l lg:pl-6'>
            <div className='flex items-center justify-between pb-1'>
              <h4 className='text-xs font-semibold text-foreground uppercase tracking-wider'>
                {activeGuideId ? '📝 Edit Konten Panduan' : '✍️ Susun Panduan / FAQ Baru'}
              </h4>
              {activeGuideId && (
                <Badge variant='outline' className='text-[10px] bg-primary/10 text-primary border-primary/30'>
                  Mode Mengedit
                </Badge>
              )}
            </div>

            <form onSubmit={handleSave} className='space-y-4'>
              <div className='space-y-1.5'>
                <label className='text-xs font-medium text-foreground block'>
                  Judul Panduan / Pertanyaan FAQ *
                </label>
                <Input
                  value={judul}
                  onChange={(e) => setJudul(e.target.value)}
                  placeholder='Contoh: Petunjuk Pengajuan Cuti di SIMPEG'
                  className='h-10 text-xs'
                  required
                />
              </div>

              <div className='space-y-1.5'>
                <div className='flex items-center justify-between'>
                  <label className='text-xs font-medium text-foreground block'>
                    Isi Panduan / Langkah Operasional (Format Markdown) *
                  </label>
                  <span className='text-[10px] text-muted-foreground flex items-center gap-1'>
                    <Info className='size-3' /> Mendukung Heading, List, Link
                  </span>
                </div>
                <Textarea
                  value={konten}
                  onChange={(e) => setKonten(e.target.value)}
                  placeholder='1. Buka halaman utama aplikasi&#10;2. Masukkan kredensial login SSO&#10;3. Masuk ke menu formulir...'
                  rows={10}
                  className='text-xs leading-relaxed font-mono resize-y min-h-[220px]'
                  required
                />
              </div>

              <div className='flex items-center justify-end gap-2.5 pt-2 border-t'>
                {activeGuideId && (
                  <Button
                    type='button'
                    variant='outline'
                    size='sm'
                    onClick={resetForm}
                    className='h-9 px-3 text-xs'
                  >
                    Batal Edit
                  </Button>
                )}
                <Button
                  type='submit'
                  size='sm'
                  disabled={isSaving}
                  className='h-9 px-4 text-xs font-semibold shadow-xs'
                >
                  {isSaving && <Loader2 className='mr-1.5 size-3.5 animate-spin' />}
                  {activeGuideId ? 'Simpan Perubahan' : 'Terbitkan Panduan'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
