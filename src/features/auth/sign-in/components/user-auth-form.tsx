import { useState } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from '@tanstack/react-router'
import { Loader2, LogIn, ShieldCheck } from 'lucide-react'
import { toast } from 'sonner'
import { useAuthStore } from '@/stores/auth-store'
import { authApi } from '@/features/auth/api/auth-api'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { PasswordInput } from '@/components/password-input'

const formSchema = z.object({
  username: z
    .string()
    .min(1, 'Username / Email / NIP wajib diisi.'),
  password: z
    .string()
    .min(1, 'Kata sandi wajib diisi.'),
})

interface UserAuthFormProps extends React.HTMLAttributes<HTMLFormElement> {
  redirectTo?: string
}

export function UserAuthForm({
  className,
  redirectTo,
  ...props
}: UserAuthFormProps) {
  const [isLoading, setIsLoading] = useState(false)
  const [isSSOLoading, setIsSSOLoading] = useState(false)
  const navigate = useNavigate()
  const { auth } = useAuthStore()

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: '',
      password: '',
    },
  })

  async function onSubmit(data: z.infer<typeof formSchema>) {
    setIsLoading(true)
    try {
      const response = await authApi.manualLogin({
        username: data.username,
        password: data.password,
      })

      // Store authenticated user and token
      auth.setAuth(response.token, response.user)

      toast.success(`Selamat datang kembali, ${response.user.nama}!`, {
        description: 'Autentikasi admin berhasil.',
      })

      const targetPath = redirectTo || '/'
      navigate({ to: targetPath, replace: true })
    } catch (err: any) {
      toast.error('Gagal Masuk', {
        description: err?.message || 'Kredensial login tidak valid.',
      })
    } finally {
      setIsLoading(false)
    }
  }

  async function handleSSOLogin() {
    setIsSSOLoading(true)
    try {
      const loginURL = await authApi.getSSOLoginURL()
      window.location.href = loginURL
    } catch (err: any) {
      toast.error('Gagal Menghubungi SSO BPS Sumut', {
        description: err?.message || 'Server SSO tidak dapat dijangkau.',
      })
      setIsSSOLoading(false)
    }
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className={cn('grid gap-4', className)}
        {...props}
      >
        <FormField
          control={form.control}
          name='username'
          render={({ field }) => (
            <FormItem>
              <FormLabel>Username / Email / NIP</FormLabel>
              <FormControl>
                <Input
                  placeholder='admin / nama@bps.go.id / 1995xxxx'
                  autoComplete='username'
                  disabled={isLoading || isSSOLoading}
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name='password'
          render={({ field }) => (
            <FormItem>
              <FormLabel>Kata Sandi</FormLabel>
              <FormControl>
                <PasswordInput
                  placeholder='Masukkan kata sandi'
                  autoComplete='current-password'
                  disabled={isLoading || isSSOLoading}
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button className='mt-2 w-full' type='submit' disabled={isLoading || isSSOLoading}>
          {isLoading ? <Loader2 className='mr-2 size-4 animate-spin' /> : <LogIn className='mr-2 size-4' />}
          Masuk ke Portal Admin
        </Button>

        <div className='relative my-2'>
          <div className='absolute inset-0 flex items-center'>
            <span className='w-full border-t border-muted' />
          </div>
          <div className='relative flex justify-center text-xs uppercase'>
            <span className='bg-card px-2 text-muted-foreground'>
              Atau masuk menggunakan
            </span>
          </div>
        </div>

        <Button
          variant='outline'
          type='button'
          className='w-full border-primary/30 hover:bg-primary/5'
          disabled={isLoading || isSSOLoading}
          onClick={handleSSOLogin}
        >
          {isSSOLoading ? (
            <Loader2 className='mr-2 size-4 animate-spin' />
          ) : (
            <ShieldCheck className='mr-2 size-4 text-primary' />
          )}
          Masuk dengan SSO BPS Sumut
        </Button>
      </form>
    </Form>
  )
}
