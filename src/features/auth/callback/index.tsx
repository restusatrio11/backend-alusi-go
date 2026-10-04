import { useEffect, useRef, useState } from 'react'
import { useNavigate, useSearch } from '@tanstack/react-router'
import { Loader2, ShieldAlert, ShieldCheck } from 'lucide-react'
import { toast } from 'sonner'
import { useAuthStore } from '@/stores/auth-store'
import { authApi } from '@/features/auth/api/auth-api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

export function Callback() {
  const search = useSearch({ strict: false }) as { code?: string; state?: string }
  const navigate = useNavigate()
  const { auth } = useAuthStore()
  const [errorMsg, setErrorMsg] = useState<string | null>(null)
  const isProcessed = useRef(false)

  useEffect(() => {
    if (isProcessed.current) return
    isProcessed.current = true

    const code = search?.code
    const state = search?.state

    if (!code) {
      setErrorMsg('Authorization code dari SSO BPS Sumut tidak ditemukan.')
      return
    }

    async function processCallback() {
      try {
        const response = await authApi.handleSSOCallback(code!, state)
        auth.setAuth(response.token, response.user)

        toast.success(`Selamat datang, ${response.user.nama}!`, {
          description: 'Login SSO BPS Sumut berhasil.',
        })

        navigate({ to: '/', replace: true })
      } catch (err: any) {
        setErrorMsg(err?.message || 'Gagal menukarkan kode otorisasi SSO.')
        toast.error('Autentikasi SSO Gagal', {
          description: err?.message || 'Terjadi kesalahan pada server SSO.',
        })
      }
    }

    processCallback()
  }, [search, auth, navigate])

  return (
    <div className='flex min-h-screen items-center justify-center p-4 bg-background'>
      <Card className='max-w-md w-full text-center shadow-lg border-border/60'>
        <CardHeader className='space-y-2'>
          <div className='flex justify-center'>
            {errorMsg ? (
              <div className='flex h-12 w-12 items-center justify-center rounded-full bg-destructive/10 text-destructive'>
                <ShieldAlert className='size-6' />
              </div>
            ) : (
              <div className='flex h-12 w-12 items-center justify-center rounded-full bg-primary/10 text-primary'>
                <ShieldCheck className='size-6 animate-pulse' />
              </div>
            )}
          </div>
          <CardTitle className='text-lg font-bold'>
            {errorMsg ? 'Autentikasi SSO Gagal' : 'Memproses Login SSO BPS Sumut...'}
          </CardTitle>
          <CardDescription>
            {errorMsg
              ? errorMsg
              : 'Mohon tunggu sebentar, kami sedang memverifikasi identitas Anda dengan portal SSO BPS Sumut.'}
          </CardDescription>
        </CardHeader>
        <CardContent className='flex flex-col items-center justify-center pt-2 pb-6'>
          {!errorMsg ? (
            <Loader2 className='size-8 animate-spin text-primary' />
          ) : (
            <Button
              className='mt-2 w-full'
              onClick={() => navigate({ to: '/sign-in', replace: true })}
            >
              Kembali ke Halaman Login
            </Button>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
