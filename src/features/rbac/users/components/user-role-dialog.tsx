import { useState, useEffect } from 'react'
import { KeyRound, Shield, User as UserIcon, Check } from 'lucide-react'
import { toast } from 'sonner'
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
import { Checkbox } from '@/components/ui/checkbox'
import { rbacApi } from '../../api/rbac-api'
import { UserWithRoles, RoleWithPermissions } from '@/types/rbac'

interface UserRoleDialogProps {
  user: UserWithRoles | null
  open: boolean
  onOpenChange: (open: boolean) => void
  roles: RoleWithPermissions[]
  onSuccess: () => void
}

export function UserRoleDialog({
  user,
  open,
  onOpenChange,
  roles,
  onSuccess,
}: UserRoleDialogProps) {
  const [selectedRoleIds, setSelectedRoleIds] = useState<Set<number>>(new Set())
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    if (user && open) {
      const ids = new Set(user.roles.map((r) => r.id))
      setSelectedRoleIds(ids)
    }
  }, [user, open])

  const handleToggleRole = (roleId: number) => {
    setSelectedRoleIds((prev) => {
      const next = new Set(prev)
      if (next.has(roleId)) {
        next.delete(roleId)
      } else {
        next.add(roleId)
      }
      return next
    })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!user) return

    try {
      setIsSubmitting(true)
      await rbacApi.assignUserRoles(user.id, {
        role_ids: Array.from(selectedRoleIds),
      })
      toast.success(`Role untuk pengguna ${user.nama} berhasil diperbarui!`)
      onOpenChange(false)
      onSuccess()
    } catch (err: any) {
      toast.error(
        'Gagal memperbarui peran pengguna: ' +
          (err.response?.data?.message || err.message)
      )
    } finally {
      setIsSubmitting(false)
    }
  }

  if (!user) return null

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-[560px] max-h-[90vh] overflow-y-auto p-6'>
        <form onSubmit={handleSubmit}>
          <DialogHeader className='space-y-2 pb-4 border-b'>
            <div className='flex items-center gap-2 text-primary'>
              <KeyRound className='h-5 w-5' />
              <DialogTitle className='text-lg font-bold'>
                Kelola Peran & Akses Pengguna
              </DialogTitle>
            </div>
            <DialogDescription className='text-xs text-muted-foreground'>
              Tetapkan satu atau lebih peran (roles) untuk pengguna ini. Peran menentukan izin dan visibilitas menu di seluruh sistem.
            </DialogDescription>
          </DialogHeader>

          {/* User Information Card */}
          <div className='my-5 rounded-lg border bg-muted/30 p-4 space-y-2 text-xs'>
            <div className='flex items-center justify-between'>
              <span className='font-semibold text-foreground flex items-center gap-1.5'>
                <UserIcon className='h-3.5 w-3.5 text-primary' />
                {user.nama}
              </span>
              <Badge variant='outline' className='text-[10px] capitalize'>
                {user.user_type}
              </Badge>
            </div>
            <div className='grid grid-cols-2 gap-2 text-muted-foreground pt-1 border-t'>
              <div>
                <span className='block text-[10px] uppercase tracking-wider'>Email</span>
                <span className='text-foreground font-medium'>{user.email}</span>
              </div>
              <div>
                <span className='block text-[10px] uppercase tracking-wider'>NIP / Username</span>
                <span className='text-foreground font-medium'>
                  {user.nip || user.username || user.sso_sub || '-'}
                </span>
              </div>
            </div>
          </div>

          {/* Role Checkbox Selection List */}
          <div className='space-y-3 mb-6'>
            <label className='text-xs font-semibold text-foreground flex items-center gap-1.5'>
              <Shield className='h-3.5 w-3.5 text-primary' />
              Pilih Peran yang Diberikan:
            </label>

            <div className='space-y-2.5'>
              {roles.map((role) => {
                const isSelected = selectedRoleIds.has(role.id)
                return (
                  <label
                    key={role.id}
                    className={`flex items-start gap-3 p-3 rounded-lg border transition-all cursor-pointer ${
                      isSelected
                        ? 'bg-primary/5 border-primary shadow-xs dark:bg-primary/10'
                        : 'bg-background hover:bg-muted/40 border-border'
                    }`}
                  >
                    <Checkbox
                      checked={isSelected}
                      onCheckedChange={() => handleToggleRole(role.id)}
                      className='mt-0.5'
                    />
                    <div className='flex-1 select-none space-y-1'>
                      <div className='flex items-center justify-between'>
                        <span className='text-xs font-bold capitalize text-foreground flex items-center gap-1.5'>
                          {role.nama}
                          {role.nama === 'admin' && (
                            <Badge className='bg-red-500/10 text-red-600 dark:text-red-400 border-red-200 text-[10px] py-0'>
                              Superadmin
                            </Badge>
                          )}
                          {role.nama === 'pimpinan' && (
                            <Badge className='bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-200 text-[10px] py-0'>
                              Manajemen
                            </Badge>
                          )}
                        </span>
                        <span className='text-[11px] text-muted-foreground'>
                          {role.permissions.length} izin terkait
                        </span>
                      </div>
                      <p className='text-[11px] text-muted-foreground leading-relaxed'>
                        {role.deskripsi || `Hak akses level ${role.nama}`}
                      </p>
                    </div>
                  </label>
                )
              })}
            </div>
          </div>

          <DialogFooter className='gap-2 pt-4 border-t'>
            <Button
              type='button'
              variant='outline'
              onClick={() => onOpenChange(false)}
              disabled={isSubmitting}
            >
              Batal
            </Button>
            <Button
              type='submit'
              disabled={isSubmitting}
              className='bg-primary text-primary-foreground'
            >
              {isSubmitting ? (
                'Menyimpan...'
              ) : (
                <>
                  <Check className='mr-1.5 h-4 w-4' /> Simpan Peran Pengguna
                </>
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
