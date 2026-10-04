import { useState, useRef } from 'react'
import { UploadCloud, CheckCircle, AlertCircle, Loader2 } from 'lucide-react'
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
        description: `Logo baru untuk aplikasi "${app.nama}" telah aktif.`,
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
      <DialogContent className='max-w-md'>
        <DialogHeader>
          <DialogTitle>Upload Logo Aplikasi</DialogTitle>
          <DialogDescription>
            Unggah logo representatif untuk aplikasi{' '}
            <span className='font-semibold text-foreground'>{app?.nama}</span>.
            Backend akan otomatis me-resize (maks 512x512) dan mengompresi gambar.
          </DialogDescription>
        </DialogHeader>

        <div className='space-y-4 py-2'>
          {/* Current Logo Display */}
          {app?.ikon_url && !previewUrl && (
            <div className='flex items-center gap-3 p-3 rounded-lg bg-muted/50 border'>
              <img
                src={app.ikon_url}
                alt={app.nama}
                className='size-12 rounded-lg object-contain bg-white p-1 border shadow-xs'
              />
              <div className='text-xs space-y-0.5 overflow-hidden'>
                <p className='font-semibold text-foreground'>Logo Saat Ini</p>
                <p className='text-muted-foreground truncate'>{app.ikon_url}</p>
              </div>
            </div>
          )}

          {/* Dropzone */}
          <div
            onDrop={handleDrop}
            onDragOver={handleDragOver}
            onClick={() => fileInputRef.current?.click()}
            className={`flex flex-col items-center justify-center p-6 border-2 border-dashed rounded-xl cursor-pointer transition-colors ${
              previewUrl
                ? 'border-primary/50 bg-primary/5'
                : 'border-muted-foreground/30 hover:border-primary/50 hover:bg-muted/30'
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
              <div className='flex flex-col items-center gap-2'>
                <img
                  src={previewUrl}
                  alt='Preview'
                  className='size-24 object-contain rounded-lg bg-white p-2 border shadow-sm'
                />
                <div className='text-center'>
                  <p className='text-xs font-semibold text-foreground'>{file?.name}</p>
                  <p className='text-[11px] text-muted-foreground'>
                    {file ? formatBytes(file.size) : ''}
                  </p>
                </div>
                <Badge variant='outline' className='text-[10px] bg-background'>
                  Klik untuk mengganti gambar
                </Badge>
              </div>
            ) : (
              <div className='flex flex-col items-center gap-2 text-center'>
                <div className='p-3 rounded-full bg-primary/10 text-primary'>
                  <UploadCloud className='size-6' />
                </div>
                <div>
                  <p className='text-sm font-medium text-foreground'>
                    Tarik & lepaskan file logo di sini
                  </p>
                  <p className='text-xs text-muted-foreground'>
                    atau klik untuk memilih file dari komputer
                  </p>
                </div>
                <div className='flex flex-wrap gap-1 justify-center pt-2'>
                  <Badge variant='secondary' className='text-[10px]'>WebP</Badge>
                  <Badge variant='secondary' className='text-[10px]'>SVG</Badge>
                  <Badge variant='secondary' className='text-[10px]'>PNG</Badge>
                  <Badge variant='secondary' className='text-[10px]'>JPG/JPEG</Badge>
                  <Badge variant='outline' className='text-[10px]'>Maks 2MB</Badge>
                </div>
              </div>
            )}
          </div>

          {errorMsg && (
            <div className='flex items-center gap-2 p-3 text-xs text-destructive bg-destructive/10 rounded-lg'>
              <AlertCircle className='size-4 shrink-0' />
              <span>{errorMsg}</span>
            </div>
          )}
        </div>

        <DialogFooter className='gap-2'>
          <Button
            type='button'
            variant='outline'
            onClick={handleClose}
            disabled={isUploading}
          >
            Batal
          </Button>
          <Button
            type='button'
            onClick={handleUpload}
            disabled={!file || isUploading}
          >
            {isUploading ? (
              <>
                <Loader2 className='mr-2 size-4 animate-spin' />
                Mengompresi & Menyimpan...
              </>
            ) : (
              <>
                <CheckCircle className='mr-2 size-4' />
                Simpan Logo
              </>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
