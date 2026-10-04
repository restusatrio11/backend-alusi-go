import { useState, useRef } from 'react'
import { UploadCloud, CheckCircle, AlertCircle, Loader2, Image as ImageIcon, Sparkles, Trash2 } from 'lucide-react'
import { toast } from 'sonner'
import { App } from '@/types/apps'
import { appsApi } from '../api/apps-api'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'

interface AppLogoDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  app: App | null
  onSuccess: () => void
}

const MAX_FILE_SIZE = 2 * 1024 * 1024 // 2MB
const ACCEPTED_TYPES = ['image/png', 'image/jpeg', 'image/jpg', 'image/webp', 'image/svg+xml']

export function AppLogoDialog({
  open,
  onOpenChange,
  app,
  onSuccess,
}: AppLogoDialogProps) {
  const [file, setFile] = useState<File | null>(null)
  const [previewUrl, setPreviewUrl] = useState<string | null>(null)
  const [isUploading, setIsUploading] = useState(false)
  const [errorMsg, setErrorMsg] = useState<string | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileChange = (selectedFile: File) => {
    setErrorMsg(null)
    if (!ACCEPTED_TYPES.includes(selectedFile.type)) {
      setErrorMsg('Format file tidak didukung. Harap gunakan format PNG, JPG, JPEG, WebP, atau SVG.')
      return
    }

    if (selectedFile.size > MAX_FILE_SIZE) {
      setErrorMsg('Ukuran file melebihi batas maksimum 2MB.')
      return
    }

    setFile(selectedFile)
    const reader = new FileReader()
    reader.onload = () => {
      setPreviewUrl(reader.result as string)
    }
    reader.readAsDataURL(selectedFile)
  }

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      handleFileChange(e.dataTransfer.files[0])
    }
  }

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
  }

  const handleUpload = async () => {
    if (!app || !file) return
    setIsUploading(true)
    setErrorMsg(null)

    try {
      await appsApi.uploadAppLogo(app.id, file)
      toast.success('Logo berhasil diunggah dan dikompresi!', {
        description: `Logo baru untuk aplikasi "${app.nama}" telah aktif secara publik.`,
      })
      onSuccess()
      handleClose()
    } catch (err: any) {
      setErrorMsg(err?.message || 'Gagal mengunggah logo.')
      toast.error('Gagal Mengunggah Logo', {
        description: err?.message || 'Terjadi kesalahan saat memproses gambar di backend.',
      })
    } finally {
      setIsUploading(false)
    }
  }

  const handleClose = () => {
    setFile(null)
    setPreviewUrl(null)
    setErrorMsg(null)
    onOpenChange(false)
  }

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className='sm:max-w-xl lg:max-w-2xl p-6 md:p-8'>
        <DialogHeader className='space-y-1.5 pb-2 border-b'>
          <div className='flex items-center gap-2 text-primary'>
            <ImageIcon className='size-5' />
            <DialogTitle className='text-xl font-bold tracking-tight'>
              Unggah & Kompresi Logo Aplikasi
            </DialogTitle>
          </div>
          <DialogDescription className='text-xs text-muted-foreground leading-relaxed'>
            Perbarui logo ikon untuk aplikasi{' '}
            <span className='font-semibold text-foreground'>&quot;{app?.nama}&quot;</span>.
            Sistem backend akan secara otomatis melakukan resize dan kompresi WebP untuk performa maksimal.
          </DialogDescription>
        </DialogHeader>

        <div className='space-y-5 py-3'>
          {/* Current vs Preview Comparison Grid */}
          <div className='grid grid-cols-1 sm:grid-cols-2 gap-4'>
            {/* Current Active Logo */}
            <div className='p-4 rounded-xl border bg-muted/20 space-y-2'>
              <span className='text-xs font-semibold text-muted-foreground block'>
                Logo Aktif Saat Ini
              </span>
              <div className='flex items-center gap-3'>
                <div className='size-14 rounded-xl bg-background flex items-center justify-center p-1.5 border shadow-xs overflow-hidden shrink-0'>
                  {app?.ikon_url ? (
                    <img
                      src={app.ikon_url}
                      alt={app.nama}
                      className='size-full object-contain'
                    />
                  ) : (
                    <ImageIcon className='size-6 text-muted-foreground/40' />
                  )}
                </div>
                <div className='min-w-0 text-xs'>
                  <p className='font-semibold text-foreground truncate'>
                    {app?.nama || 'Aplikasi'}
                  </p>
                  <p className='text-[11px] text-muted-foreground truncate'>
                    {app?.ikon_url || 'Belum ada logo terpasang'}
                  </p>
                </div>
              </div>
            </div>

            {/* Smart Pipeline Features */}
            <div className='p-4 rounded-xl border bg-primary/5 space-y-2'>
              <span className='text-xs font-semibold text-primary flex items-center gap-1.5'>
                <Sparkles className='size-3.5' /> Fitur Pemrosesan Otomatis
              </span>
              <ul className='text-[11px] text-muted-foreground space-y-1'>
                <li>• Konversi otomatis ke format efisien <strong>WebP</strong>.</li>
                <li>• Standardisasi resolusi proporsional (maks 512x512).</li>
                <li>• Dukungan vektor native untuk file <strong>SVG</strong>.</li>
              </ul>
            </div>
          </div>

          {/* Dropzone */}
          <div
            onDrop={handleDrop}
            onDragOver={handleDragOver}
            onClick={() => fileInputRef.current?.click()}
            className={`flex flex-col items-center justify-center p-8 border-2 border-dashed rounded-2xl cursor-pointer transition-all duration-200 ${
              previewUrl
                ? 'border-primary bg-primary/5 shadow-xs'
                : 'border-muted-foreground/30 hover:border-primary hover:bg-muted/30'
            }`}
          >
            <input
              ref={fileInputRef}
              type='file'
              accept='.png,.jpg,.jpeg,.webp,.svg'
              className='hidden'
              onChange={(e) => {
                if (e.target.files && e.target.files[0]) {
                  handleFileChange(e.target.files[0])
                }
              }}
            />

            {previewUrl ? (
              <div className='flex flex-col items-center gap-3'>
                <div className='size-28 rounded-2xl bg-white p-3 border shadow-sm flex items-center justify-center overflow-hidden'>
                  <img
                    src={previewUrl}
                    alt='Preview'
                    className='size-full object-contain'
                  />
                </div>
                <div className='text-center space-y-0.5'>
                  <p className='text-xs font-semibold text-foreground'>{file?.name}</p>
                  <p className='text-[11px] text-muted-foreground font-mono'>
                    {file ? formatBytes(file.size) : ''}
                  </p>
                </div>
                <div className='flex items-center gap-2 pt-1'>
                  <Badge variant='outline' className='text-[10px] bg-background'>
                    Klik area untuk ganti file
                  </Badge>
                  <Button
                    type='button'
                    variant='ghost'
                    size='sm'
                    className='h-6 px-2 text-xs text-destructive hover:bg-destructive/10'
                    onClick={(e) => {
                      e.stopPropagation()
                      setFile(null)
                      setPreviewUrl(null)
                    }}
                  >
                    <Trash2 className='size-3 mr-1' /> Hapus
                  </Button>
                </div>
              </div>
            ) : (
              <div className='flex flex-col items-center gap-3 text-center'>
                <div className='p-3.5 rounded-full bg-primary/10 text-primary'>
                  <UploadCloud className='size-7' />
                </div>
                <div className='space-y-1'>
                  <p className='text-sm font-semibold text-foreground'>
                    Tarik & lepaskan berkas logo ke sini
                  </p>
                  <p className='text-xs text-muted-foreground'>
                    atau klik untuk memilih file gambar dari perangkat Anda
                  </p>
                </div>
                <div className='flex flex-wrap gap-1.5 justify-center pt-2'>
                  <Badge variant='secondary' className='text-[10px]'>WebP</Badge>
                  <Badge variant='secondary' className='text-[10px]'>SVG</Badge>
                  <Badge variant='secondary' className='text-[10px]'>PNG</Badge>
                  <Badge variant='secondary' className='text-[10px]'>JPG / JPEG</Badge>
                  <Badge variant='outline' className='text-[10px] font-mono'>Maks 2MB</Badge>
                </div>
              </div>
            )}
          </div>

          {errorMsg && (
            <div className='flex items-center gap-2.5 p-3 text-xs text-destructive bg-destructive/10 border border-destructive/20 rounded-xl'>
              <AlertCircle className='size-4 shrink-0' />
              <span className='font-medium'>{errorMsg}</span>
            </div>
          )}
        </div>

        <DialogFooter className='gap-2.5 pt-4 border-t'>
          <Button
            type='button'
            variant='outline'
            onClick={handleClose}
            disabled={isUploading}
            className='h-10 px-4 text-xs'
          >
            Batal
          </Button>
          <Button
            type='button'
            onClick={handleUpload}
            disabled={!file || isUploading}
            className='h-10 px-5 text-xs font-semibold shadow-xs'
          >
            {isUploading ? (
              <>
                <Loader2 className='mr-2 size-4 animate-spin' />
                Mengompresi & Menyimpan...
              </>
            ) : (
              <>
                <CheckCircle className='mr-2 size-4' />
                Simpan & Pasang Logo
              </>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
