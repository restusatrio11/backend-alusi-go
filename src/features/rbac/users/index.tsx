import { useState, useEffect } from 'react'
import {
  Users2,
  Search,
  RotateCcw,
  UserCheck,
  UserX,
  Filter,
  KeyRound,
} from 'lucide-react'
import { toast } from 'sonner'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { ThemeSwitch } from '@/components/theme-switch'
import { PermissionGuard } from '@/components/permission-guard'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Card, CardContent } from '@/components/ui/card'
import { UserRoleDialog } from './components/user-role-dialog'
import { rbacApi } from '../api/rbac-api'
import { UserWithRoles, RoleWithPermissions } from '@/types/rbac'

export function UserRoleManagement() {
  const [users, setUsers] = useState<UserWithRoles[]>([])
  const [roles, setRoles] = useState<RoleWithPermissions[]>([])
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedRoleId, setSelectedRoleId] = useState<string>('all')
  const [isLoading, setIsLoading] = useState(true)

  // Dialog state
  const [selectedUser, setSelectedUser] = useState<UserWithRoles | null>(null)
  const [isDialogOpen, setIsDialogOpen] = useState(false)

  const loadData = async () => {
    try {
      setIsLoading(true)
      const [usersData, rolesData] = await Promise.all([
        rbacApi.getUsers({
          search: searchQuery || undefined,
          role_id: selectedRoleId !== 'all' ? Number(selectedRoleId) : undefined,
          limit: 100,
        }),
        rbacApi.getRoles(),
      ])

      setUsers(usersData.users)
      setRoles(rolesData)
    } catch (err: any) {
      toast.error('Gagal memuat daftar pengguna: ' + (err.message || 'Terjadi kesalahan'))
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [selectedRoleId])

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    loadData()
  }

  const handleOpenManageRole = (user: UserWithRoles) => {
    setSelectedUser(user)
    setIsDialogOpen(true)
  }

  // Helper for role badge colors
  const getRoleBadge = (roleName: string) => {
    switch (roleName.toLowerCase()) {
      case 'admin':
        return (
          <Badge
            key={roleName}
            className='bg-red-500/15 text-red-700 dark:text-red-400 border-red-200 dark:border-red-900/50 hover:bg-red-500/20 text-[11px]'
          >
            Admin
          </Badge>
        )
      case 'pimpinan':
        return (
          <Badge
            key={roleName}
            className='bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-200 dark:border-amber-900/50 hover:bg-amber-500/20 text-[11px]'
          >
            Pimpinan
          </Badge>
        )
      case 'pegawai':
        return (
          <Badge
            key={roleName}
            className='bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-200 dark:border-emerald-900/50 hover:bg-emerald-500/20 text-[11px]'
          >
            Pegawai
          </Badge>
        )
      default:
        return (
          <Badge key={roleName} variant='secondary' className='text-[11px] capitalize'>
            {roleName}
          </Badge>
        )
    }
  }

  return (
    <>
      <Header fixed>
        <div className='flex items-center gap-2 font-semibold'>
          <Users2 className='h-5 w-5 text-primary' />
          <span>Pengguna & Peran (User Roles)</span>
        </div>
        <div className='ms-auto flex items-center space-x-4'>
          <ThemeSwitch />
          <ProfileDropdown />
        </div>
      </Header>

      <Main className='flex flex-1 flex-col gap-6 p-4 sm:p-6'>
        {/* Page Title & Stats */}
        <div className='flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between'>
          <div>
            <h1 className='text-2xl font-bold tracking-tight'>Manajemen Peran Pengguna</h1>
            <p className='text-sm text-muted-foreground'>
              Kelola penetapan peran (role assignment) untuk seluruh akun pegawai dan administrator portal ALUSI.
            </p>
          </div>
          <Button
            variant='outline'
            size='sm'
            onClick={loadData}
            disabled={isLoading}
          >
            <RotateCcw className={`mr-2 h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
            Segarkan Data
          </Button>
        </div>

        {/* Filter & Search Toolbar */}
        <Card className='shadow-xs'>
          <CardContent className='p-4'>
            <form onSubmit={handleSearchSubmit} className='flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between'>
              <div className='flex flex-1 items-center gap-2'>
                <div className='relative flex-1 max-w-sm'>
                  <Search className='absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground' />
                  <Input
                    type='search'
                    placeholder='Cari nama, email, atau NIP...'
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className='pl-8 h-9 text-xs'
                  />
                </div>
                <Button type='submit' size='sm' variant='secondary' className='h-9 text-xs'>
                  Cari
                </Button>
              </div>

              <div className='flex items-center gap-2'>
                <Filter className='h-4 w-4 text-muted-foreground' />
                <span className='text-xs text-muted-foreground whitespace-nowrap'>Filter Peran:</span>
                <Select value={selectedRoleId} onValueChange={setSelectedRoleId}>
                  <SelectTrigger className='w-[160px] h-9 text-xs'>
                    <SelectValue placeholder='Semua Peran' />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value='all'>Semua Peran</SelectItem>
                    {roles.map((r) => (
                      <SelectItem key={r.id} value={r.id.toString()} className='capitalize text-xs'>
                        {r.nama}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </form>
          </CardContent>
        </Card>

        {/* Users Table */}
        <Card className='shadow-xs overflow-hidden border'>
          <div className='overflow-x-auto'>
            <Table>
              <TableHeader className='bg-muted/40'>
                <TableRow>
                  <TableHead className='w-[220px] font-semibold text-xs'>Nama & Identitas</TableHead>
                  <TableHead className='font-semibold text-xs'>Email</TableHead>
                  <TableHead className='font-semibold text-xs'>Satuan Kerja</TableHead>
                  <TableHead className='font-semibold text-xs'>Tipe Akun</TableHead>
                  <TableHead className='font-semibold text-xs'>Peran (Roles)</TableHead>
                  <TableHead className='font-semibold text-xs'>Status</TableHead>
                  <TableHead className='text-right font-semibold text-xs'>Aksi</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={7} className='h-32 text-center text-muted-foreground text-xs'>
                      <div className='flex items-center justify-center gap-2'>
                        <RotateCcw className='h-4 w-4 animate-spin text-primary' />
                        <span>Memuat data pengguna...</span>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : users.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className='h-32 text-center text-muted-foreground text-xs'>
                      Tidak ada pengguna yang sesuai dengan kriteria pencarian.
                    </TableCell>
                  </TableRow>
                ) : (
                  users.map((user) => (
                    <TableRow key={user.id} className='hover:bg-muted/30 transition-colors'>
                      <TableCell className='font-medium text-xs'>
                        <div className='font-semibold text-foreground'>{user.nama}</div>
                        {(user.nip || user.username || user.sso_sub) && (
                          <div className='text-[11px] text-muted-foreground font-mono'>
                            {user.nip || user.username || user.sso_sub}
                          </div>
                        )}
                      </TableCell>

                      <TableCell className='text-xs text-muted-foreground font-mono'>
                        {user.email}
                      </TableCell>

                      <TableCell className='text-xs text-muted-foreground'>
                        {user.satker_nama || '-'}
                      </TableCell>

                      <TableCell className='text-xs'>
                        <Badge
                          variant='outline'
                          className={
                            user.user_type === 'internal'
                              ? 'border-blue-300 text-blue-700 dark:text-blue-400 bg-blue-50/50 dark:bg-blue-950/20 text-[10px]'
                              : 'text-muted-foreground text-[10px]'
                          }
                        >
                          {user.user_type === 'internal' ? 'BPS SSO' : 'Manual'}
                        </Badge>
                      </TableCell>

                      <TableCell className='text-xs'>
                        <div className='flex flex-wrap gap-1'>
                          {user.roles && user.roles.length > 0 ? (
                            user.roles.map((r) => getRoleBadge(r.nama))
                          ) : (
                            <span className='text-xs text-muted-foreground italic'>
                              Tanpa Peran
                            </span>
                          )}
                        </div>
                      </TableCell>

                      <TableCell className='text-xs'>
                        {user.status === 'active' ? (
                          <span className='flex items-center text-emerald-600 dark:text-emerald-400 text-xs font-medium'>
                            <UserCheck className='mr-1 h-3.5 w-3.5' /> Aktif
                          </span>
                        ) : (
                          <span className='flex items-center text-muted-foreground text-xs'>
                            <UserX className='mr-1 h-3.5 w-3.5' /> Nonaktif
                          </span>
                        )}
                      </TableCell>

                      <TableCell className='text-right'>
                        <PermissionGuard permission='users:manage'>
                          <Button
                            variant='outline'
                            size='sm'
                            className='h-8 text-xs gap-1.5'
                            onClick={() => handleOpenManageRole(user)}
                          >
                            <KeyRound className='h-3.5 w-3.5 text-primary' />
                            Kelola Role
                          </Button>
                        </PermissionGuard>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </Card>

        {/* User Role Assignment Modal Dialog */}
        <UserRoleDialog
          user={selectedUser}
          open={isDialogOpen}
          onOpenChange={setIsDialogOpen}
          roles={roles}
          onSuccess={loadData}
        />
      </Main>
    </>
  )
}
