import { Link } from '@tanstack/react-router'
import {
  BadgeCheck,
  ChevronsUpDown,
  LogOut,
} from 'lucide-react'
import useDialogState from '@/hooks/use-dialog-state'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '@/components/ui/sidebar'
import { SignOutDialog } from '@/components/sign-out-dialog'
import { useAuthStore } from '@/stores/auth-store'

type NavUserProps = {
  user?: {
    name: string
    email: string
    avatar: string
  }
}

export function NavUser({ user: defaultUser }: NavUserProps) {
  const { isMobile } = useSidebar()
  const [open, setOpen] = useDialogState()
  const { auth } = useAuthStore()

  const displayName = auth.user?.nama || defaultUser?.name || 'User ALUSI'
  const displayEmail = auth.user?.email || auth.user?.username || defaultUser?.email || 'admin@bps.go.id'
  const displayAvatar = auth.user?.foto_url || defaultUser?.avatar || ''
  const roleName = auth.user?.roles?.[0]?.nama?.toUpperCase() || 'USER'
  const satkerName = auth.user?.satker?.nama || ''

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .slice(0, 2)
      .map((part) => part[0])
      .join('')
      .toUpperCase() || 'AL'
  }

  return (
    <>
      <SidebarMenu>
        <SidebarMenuItem>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <SidebarMenuButton
                size='lg'
                className='data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground'
              >
                <Avatar className='h-8 w-8 rounded-lg'>
                  <AvatarImage src={displayAvatar} alt={displayName} />
                  <AvatarFallback className='rounded-lg bg-primary/10 text-primary font-semibold text-xs'>
                    {getInitials(displayName)}
                  </AvatarFallback>
                </Avatar>
                <div className='grid flex-1 text-start text-sm leading-tight'>
                  <span className='truncate font-semibold'>{displayName}</span>
                  <span className='truncate text-xs text-muted-foreground'>{displayEmail}</span>
                </div>
                <ChevronsUpDown className='ms-auto size-4' />
              </SidebarMenuButton>
            </DropdownMenuTrigger>
            <DropdownMenuContent
              className='w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg'
              side={isMobile ? 'bottom' : 'right'}
              align='end'
              sideOffset={4}
            >
              <DropdownMenuLabel className='p-0 font-normal'>
                <div className='flex items-center gap-2 px-1 py-1.5 text-start text-sm'>
                  <Avatar className='h-8 w-8 rounded-lg'>
                    <AvatarImage src={displayAvatar} alt={displayName} />
                    <AvatarFallback className='rounded-lg bg-primary/10 text-primary font-semibold text-xs'>
                      {getInitials(displayName)}
                    </AvatarFallback>
                  </Avatar>
                  <div className='grid flex-1 text-start text-sm leading-tight'>
                    <span className='truncate font-semibold'>{displayName}</span>
                    <span className='truncate text-xs text-muted-foreground'>{displayEmail}</span>
                    {satkerName && (
                      <span className='truncate text-[10px] text-muted-foreground/80 mt-0.5'>{satkerName}</span>
                    )}
                  </div>
                  <span className='ml-auto text-[10px] bg-primary/10 text-primary px-1.5 py-0.5 rounded font-mono font-bold'>
                    {roleName}
                  </span>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuGroup>
                <DropdownMenuItem asChild>
                  <Link to='/settings'>
                    <BadgeCheck className='size-4' />
                    Profil Pengguna
                  </Link>
                </DropdownMenuItem>
              </DropdownMenuGroup>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                variant='destructive'
                onClick={() => setOpen(true)}
              >
                <LogOut className='size-4' />
                Keluar Sesi
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </SidebarMenuItem>
      </SidebarMenu>

      <SignOutDialog open={!!open} onOpenChange={setOpen} />
    </>
  )
}
