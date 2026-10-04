import { useState, useEffect, useMemo } from 'react'
import {
  ShieldCheck,
  KeyRound,
  CheckCircle2,
  Layers,
  FolderTree,
  Megaphone,
  MessageSquareWarning,
  Activity,
  LineChart,
  Lock,
  RotateCcw,
  Save,
  CheckSquare,
  Square,
  Info,
} from 'lucide-react'
import { toast } from 'sonner'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { ThemeSwitch } from '@/components/theme-switch'
import { PermissionGuard } from '@/components/permission-guard'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Checkbox } from '@/components/ui/checkbox'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { rbacApi } from '../api/rbac-api'
import { Permission, RoleWithPermissions } from '@/types/rbac'

const CATEGORY_META: Record<
  string,
  { label: string; icon: React.ElementType; description: string; color: string }
> = {
  apps: {
    label: 'Katalog Aplikasi',
    icon: Layers,
    description: 'Hak akses melihat, menambah, mengubah, menghapus, menyusun urutan dan mengelola panduan aplikasi.',
    color: 'text-blue-500 bg-blue-500/10 dark:bg-blue-500/20',
  },
  categories: {
    label: 'Kategori Aplikasi',
    icon: FolderTree,
    description: 'Hak akses melihat dan mengelola taksonomi serta kelompok kategori aplikasi.',
    color: 'text-emerald-500 bg-emerald-500/10 dark:bg-emerald-500/20',
  },
  announcements: {
    label: 'Pengumuman Portal',
    icon: Megaphone,
    description: 'Hak akses publikasi banner pengumuman darurat dan informasi penting portal.',
    color: 'text-amber-500 bg-amber-500/10 dark:bg-amber-500/20',
  },
  feedbacks: {
    label: 'Umpan Balik & Kendala',
    icon: MessageSquareWarning,
    description: 'Hak akses monitoring aduan pengguna, tindak lanjut dan perubahan status tiket.',
    color: 'text-rose-500 bg-rose-500/10 dark:bg-rose-500/20',
  },
  monitoring: {
    label: 'Pemantauan Layanan',
    icon: Activity,
    description: 'Hak akses pemantauan uptime server, healthcheck probe, dan pemicu ping berkala.',
    color: 'text-teal-500 bg-teal-500/10 dark:bg-teal-500/20',
  },
  analytics: {
    label: 'Analitik & Laporan',
    icon: LineChart,
    description: 'Hak akses statistik metrik kunjungan dan ekspor laporan berkas resmi (PDF & CSV).',
    color: 'text-purple-500 bg-purple-500/10 dark:bg-purple-500/20',
  },
  audit: {
    label: 'Jejak Audit (Audit Trail)',
    icon: ShieldCheck,
    description: 'Hak akses pengawasan log aktivitas pengguna dan rekam jejak transaksi data.',
    color: 'text-indigo-500 bg-indigo-500/10 dark:bg-indigo-500/20',
  },
  rbac: {
    label: 'Keamanan & RBAC',
    icon: Lock,
    description: 'Hak akses konfigurasi role permission matrix dan pengelolaan peran akun pengguna.',
    color: 'text-red-500 bg-red-500/10 dark:bg-red-500/20',
  },
}

export function RolesPermissionMatrix() {
  const [roles, setRoles] = useState<RoleWithPermissions[]>([])
  const [allPermissions, setAllPermissions] = useState<Permission[]>([])
  const [selectedRoleId, setSelectedRoleId] = useState<number | null>(null)
  const [selectedCodes, setSelectedCodes] = useState<Set<string>>(new Set())
  const [initialCodes, setInitialCodes] = useState<Set<string>>(new Set())
  const [isLoading, setIsLoading] = useState(true)
  const [isSaving, setIsSaving] = useState(false)

  // Load roles and permissions
  const fetchData = async () => {
    try {
      setIsLoading(true)
      const [fetchedRoles, fetchedPerms] = await Promise.all([
        rbacApi.getRoles(),
        rbacApi.getPermissions(),
      ])

      setRoles(fetchedRoles)
      setAllPermissions(fetchedPerms)

      if (fetchedRoles.length > 0) {
        const defaultRole = fetchedRoles[0]
        setSelectedRoleId(defaultRole.id)
        const currentPermCodes = new Set(defaultRole.permissions.map((p) => p.kode))
        setSelectedCodes(new Set(currentPermCodes))
        setInitialCodes(new Set(currentPermCodes))
      }
    } catch (err: any) {
      toast.error('Gagal memuat data role dan permission: ' + (err.message || 'Terjadi kesalahan'))
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  // Handle role tab change
  const handleSelectRole = (roleIdStr: string) => {
    const roleId = Number(roleIdStr)
    setSelectedRoleId(roleId)
    const targetRole = roles.find((r) => r.id === roleId)
    if (targetRole) {
      const codes = new Set(targetRole.permissions.map((p) => p.kode))
      setSelectedCodes(new Set(codes))
      setInitialCodes(new Set(codes))
    }
  }

  // Active selected role object
  const activeRole = useMemo(() => {
    return roles.find((r) => r.id === selectedRoleId)
  }, [roles, selectedRoleId])

  // Group permissions by category
  const permissionsByCategory = useMemo(() => {
    const map: Record<string, Permission[]> = {}
    allPermissions.forEach((perm) => {
      const cat = perm.kategori || 'other'
      if (!map[cat]) map[cat] = []
      map[cat].push(perm)
    })
    return map
  }, [allPermissions])

  // Toggle single permission
  const handleTogglePermission = (code: string) => {
    if (activeRole?.nama === 'admin') {
      toast.info('Role Admin memiliki akses penuh (Full Superadmin Bypass)')
      return
    }

    setSelectedCodes((prev) => {
      const next = new Set(prev)
      if (next.has(code)) {
        next.delete(code)
      } else {
        next.add(code)
      }
      return next
    })
  }

  // Toggle category all
  const handleToggleCategory = (categoryKey: string) => {
    if (activeRole?.nama === 'admin') return

    const categoryPerms = permissionsByCategory[categoryKey] || []
    const categoryCodes = categoryPerms.map((p) => p.kode)
    const allSelected = categoryCodes.every((c) => selectedCodes.has(c))

    setSelectedCodes((prev) => {
      const next = new Set(prev)
      if (allSelected) {
        categoryCodes.forEach((c) => next.delete(c))
      } else {
        categoryCodes.forEach((c) => next.add(c))
      }
      return next
    })
  }

  // Check if dirty (changes exist)
  const isDirty = useMemo(() => {
    if (selectedCodes.size !== initialCodes.size) return true
    for (const code of selectedCodes) {
      if (!initialCodes.has(code)) return true
    }
    return false
  }, [selectedCodes, initialCodes])

  // Save changes
  const handleSaveChanges = async () => {
    if (!selectedRoleId) return

    try {
      setIsSaving(true)
      const payload = {
        permission_codes: Array.from(selectedCodes),
      }
      const updated = await rbacApi.updateRolePermissions(selectedRoleId, payload)

      // Update local state
      setRoles((prev) =>
        prev.map((r) => (r.id === selectedRoleId ? updated : r))
      )
      const nextCodes = new Set(updated.permissions.map((p) => p.kode))
      setSelectedCodes(new Set(nextCodes))
      setInitialCodes(new Set(nextCodes))

      toast.success(`Hak akses untuk role "${updated.nama}" berhasil diperbarui!`)
    } catch (err: any) {
      toast.error('Gagal menyimpan perubahan hak akses: ' + (err.response?.data?.message || err.message))
    } finally {
      setIsSaving(false)
    }
  }

  // Reset to initial
  const handleReset = () => {
    setSelectedCodes(new Set(initialCodes))
  }

  return (
    <>
      <Header fixed>
        <div className='flex items-center gap-2 font-semibold'>
          <KeyRound className='h-5 w-5 text-primary' />
          <span>Manajemen Role & Hak Akses (RBAC Matrix)</span>
        </div>
        <div className='ms-auto flex items-center space-x-4'>
          <ThemeSwitch />
          <ProfileDropdown />
        </div>
      </Header>

      <Main className='flex flex-1 flex-col gap-6 p-4 sm:p-6'>
        {/* Header Title Bar */}
        <div className='flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between'>
          <div>
            <h1 className='text-2xl font-bold tracking-tight'>Matriks Hak Akses & Role</h1>
            <p className='text-sm text-muted-foreground'>
              Konfigurasikan izin granular (view, create, edit, delete, manage) untuk setiap tingkat peran sistem.
            </p>
          </div>
          <div className='flex items-center gap-2'>
            <Button
              variant='outline'
              size='sm'
              onClick={fetchData}
              disabled={isLoading || isSaving}
            >
              <RotateCcw className={`mr-2 h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
              Muat Ulang
            </Button>
            {isDirty && (
              <Button
                variant='ghost'
                size='sm'
                onClick={handleReset}
                disabled={isSaving}
              >
                Batalkan
              </Button>
            )}
            <PermissionGuard permission='rbac:manage'>
              <Button
                size='sm'
                onClick={handleSaveChanges}
                disabled={!isDirty || isSaving || activeRole?.nama === 'admin'}
                className='bg-primary text-primary-foreground shadow-sm'
              >
                <Save className='mr-2 h-4 w-4' />
                {isSaving ? 'Menyimpan...' : 'Simpan Perubahan'}
              </Button>
            </PermissionGuard>
          </div>
        </div>

        {isLoading ? (
          <div className='flex h-64 items-center justify-center rounded-xl border border-dashed bg-muted/20'>
            <div className='flex flex-col items-center gap-2 text-muted-foreground'>
              <KeyRound className='h-8 w-8 animate-pulse text-primary' />
              <span>Memuat matriks RBAC dan daftar hak akses...</span>
            </div>
          </div>
        ) : (
          <>
            {/* Role Tabs */}
            <div className='flex flex-col gap-4'>
              <Tabs
                value={selectedRoleId?.toString()}
                onValueChange={handleSelectRole}
                className='w-full'
              >
                <div className='flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between border-b pb-4'>
                  <TabsList className='h-11 p-1 bg-muted/60'>
                    {roles.map((role) => {
                      const isRoleAdmin = role.nama === 'admin'
                      const permCount = isRoleAdmin
                        ? allPermissions.length
                        : role.permissions.length

                      return (
                        <TabsTrigger
                          key={role.id}
                          value={role.id.toString()}
                          className='relative gap-2 px-4 py-2 font-medium capitalize data-[state=active]:bg-background data-[state=active]:shadow-sm'
                        >
                          <span>{role.nama}</span>
                          <Badge
                            variant={role.id === selectedRoleId ? 'default' : 'secondary'}
                            className='ml-1 text-[11px] px-1.5 py-0 rounded-full'
                          >
                            {permCount} izin
                          </Badge>
                        </TabsTrigger>
                      )
                    })}
                  </TabsList>

                  {/* Active Role Summary */}
                  {activeRole && (
                    <div className='flex items-center gap-2 text-xs text-muted-foreground'>
                      <Info className='h-4 w-4 text-primary' />
                      <span>
                        {activeRole.deskripsi || `Role ${activeRole.nama} pada sistem ALUSI`}
                      </span>
                      {activeRole.nama === 'admin' && (
                        <Badge variant='outline' className='border-amber-500/50 bg-amber-500/10 text-amber-700 dark:text-amber-400 text-[11px]'>
                          Full Bypass
                        </Badge>
                      )}
                    </div>
                  )}
                </div>
              </Tabs>
            </div>

            {/* Permission Categories Grid */}
            <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5'>
              {Object.entries(permissionsByCategory).map(([categoryKey, perms]) => {
                const meta = CATEGORY_META[categoryKey] || {
                  label: categoryKey.toUpperCase(),
                  icon: KeyRound,
                  description: 'Hak akses modul ' + categoryKey,
                  color: 'text-primary bg-primary/10',
                }
                const IconComponent = meta.icon
                const catCodes = perms.map((p) => p.kode)
                const isAllSelected = catCodes.every((c) => selectedCodes.has(c))
                const selectedInCat = catCodes.filter((c) => selectedCodes.has(c)).length

                return (
                  <Card key={categoryKey} className='border shadow-sm flex flex-col justify-between overflow-hidden'>
                    <div>
                      <CardHeader className='p-4 pb-3 bg-muted/30 border-b'>
                        <div className='flex items-center justify-between'>
                          <div className='flex items-center gap-2.5'>
                            <div className={`p-2 rounded-lg ${meta.color}`}>
                              <IconComponent className='h-4 w-4' />
                            </div>
                            <div>
                              <CardTitle className='text-sm font-semibold'>{meta.label}</CardTitle>
                              <CardDescription className='text-xs line-clamp-1'>
                                {selectedInCat} dari {perms.length} izin aktif
                              </CardDescription>
                            </div>
                          </div>

                          <Button
                            variant='ghost'
                            size='sm'
                            className='h-7 px-2 text-xs text-muted-foreground hover:text-foreground'
                            onClick={() => handleToggleCategory(categoryKey)}
                            disabled={activeRole?.nama === 'admin'}
                          >
                            {isAllSelected ? (
                              <>
                                <Square className='mr-1 h-3.5 w-3.5' /> Lepas Semua
                              </>
                            ) : (
                              <>
                                <CheckSquare className='mr-1 h-3.5 w-3.5' /> Pilih Semua
                              </>
                            )}
                          </Button>
                        </div>
                      </CardHeader>

                      <CardContent className='p-4 space-y-3.5'>
                        {perms.map((perm) => {
                          const isChecked =
                            activeRole?.nama === 'admin' || selectedCodes.has(perm.kode)

                          return (
                            <label
                              key={perm.id}
                              className={`flex items-start gap-3 p-2.5 rounded-lg border transition-colors cursor-pointer ${
                                isChecked
                                  ? 'bg-primary/5 border-primary/20 dark:bg-primary/10'
                                  : 'bg-background hover:bg-muted/50 border-border'
                              } ${activeRole?.nama === 'admin' ? 'cursor-not-allowed opacity-90' : ''}`}
                            >
                              <Checkbox
                                id={`perm-${perm.id}`}
                                checked={isChecked}
                                onCheckedChange={() => handleTogglePermission(perm.kode)}
                                disabled={activeRole?.nama === 'admin'}
                                className='mt-0.5'
                              />
                              <div className='flex-1 space-y-1 select-none'>
                                <div className='flex items-center justify-between gap-2'>
                                  <span className='text-xs font-semibold leading-none'>
                                    {perm.nama}
                                  </span>
                                  <code className='text-[10px] font-mono px-1.5 py-0.5 rounded bg-muted text-muted-foreground font-semibold'>
                                    {perm.kode}
                                  </code>
                                </div>
                                {perm.deskripsi && (
                                  <p className='text-[11px] text-muted-foreground leading-relaxed'>
                                    {perm.deskripsi}
                                  </p>
                                )}
                              </div>
                            </label>
                          )
                        })}
                      </CardContent>
                    </div>

                    <div className='p-3 bg-muted/15 border-t flex items-center justify-between text-[11px] text-muted-foreground'>
                      <span>Status Modul:</span>
                      {selectedInCat === perms.length ? (
                        <span className='flex items-center text-emerald-600 dark:text-emerald-400 font-medium'>
                          <CheckCircle2 className='mr-1 h-3.5 w-3.5' /> Lengkap (Full Access)
                        </span>
                      ) : selectedInCat > 0 ? (
                        <span className='text-amber-600 dark:text-amber-400 font-medium'>
                          Parsial ({selectedInCat}/{perms.length})
                        </span>
                      ) : (
                        <span className='text-muted-foreground'>Tanpa Izin (Restricted)</span>
                      )}
                    </div>
                  </Card>
                )
              })}
            </div>
          </>
        )}
      </Main>
    </>
  )
}
