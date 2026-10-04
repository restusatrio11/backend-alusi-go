import { useSearch } from '@tanstack/react-router'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { AuthLayout } from '../auth-layout'
import { UserAuthForm } from './components/user-auth-form'

export function SignIn() {
  const { redirect } = useSearch({ from: '/(auth)/sign-in' })

  return (
    <AuthLayout>
      <Card className='max-w-md w-full gap-4 shadow-lg border-border/60'>
        <CardHeader className='text-center space-y-1.5'>
          <div className='flex justify-center mb-2'>
            <div className='flex h-12 w-12 items-center justify-center rounded-xl bg-primary text-primary-foreground font-bold text-xl shadow-md'>
              A
            </div>
          </div>
          <CardTitle className='text-xl font-bold tracking-tight'>
            ALUSI Admin Panel
          </CardTitle>
          <CardDescription className='text-xs'>
            Badan Pusat Statistik Provinsi Sumatera Utara
          </CardDescription>
        </CardHeader>
        <CardContent>
          <UserAuthForm redirectTo={redirect} />
        </CardContent>
        <CardFooter className='flex justify-center border-t pt-4'>
          <p className='text-center text-xs text-muted-foreground'>
            Akses Terbatas • Khusus Pengelola & Pimpinan BPS
          </p>
        </CardFooter>
      </Card>
    </AuthLayout>
  )
}
