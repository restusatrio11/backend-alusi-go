import { useState, useEffect } from 'react'
import { Plus, RefreshCw, FolderTree, Edit, Trash2, Hash } from 'lucide-react'
import { Category } from '@/types/apps'
import { categoriesApi } from './api/categories-api'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { PermissionGuard } from '@/components/permission-guard'
import { CategoryFormDialog } from './components/category-form-dialog'
import { CategoryDeleteDialog } from './components/category-delete-dialog'

export function Categories() {
  const [categories, setCategories] = useState<Category[]>([])
  const [isLoading, setIsLoading] = useState(true)

  const [isFormOpen, setIsFormOpen] = useState(false)
  const [categoryToEdit, setCategoryToEdit] = useState<Category | null>(null)

  const [isDeleteOpen, setIsDeleteOpen] = useState(false)
  const [categoryToDelete, setCategoryToDelete] = useState<Category | null>(null)

  const fetchCategories = async () => {
    setIsLoading(true)
    try {
      const data = await categoriesApi.getCategories()
      setCategories(data)
    } catch {
      setCategories([])
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchCategories()
  }, [])

  const handleOpenCreate = () => {
    setCategoryToEdit(null)
    setIsFormOpen(true)
  }

  const handleOpenEdit = (category: Category) => {
    setCategoryToEdit(category)
    setIsFormOpen(true)
  }

  const handleOpenDelete = (category: Category) => {
    setCategoryToDelete(category)
    setIsDeleteOpen(true)
  }

  return (
    <>
      <Header>
        <div className='flex items-center gap-2'>
          <FolderTree className='size-5 text-primary' />
          <h2 className='text-sm font-semibold tracking-tight'>Kategori Aplikasi</h2>
        </div>
        <div className='ml-auto flex items-center space-x-2'>
          <ThemeSwitch />
          <ProfileDropdown />
        </div>
      </Header>

      <Main>
        <div className='space-y-6'>
          <div className='flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4'>
            <div>
              <h1 className='text-2xl font-bold tracking-tight'>Kategori Aplikasi</h1>
              <p className='text-sm text-muted-foreground'>
                Kelola taksonomi dan pengelompokan modul aplikasi pada portal ALUSI.
              </p>
            </div>

            <div className='flex items-center gap-2'>
              <Button
                variant='outline'
                size='sm'
                onClick={fetchCategories}
                disabled={isLoading}
              >
                <RefreshCw className={`mr-1.5 size-3.5 ${isLoading ? 'animate-spin' : ''}`} />
                Refresh
              </Button>
              <PermissionGuard permission='categories:manage'>
                <Button size='sm' onClick={handleOpenCreate}>
                  <Plus className='mr-1.5 size-4' />
                  Tambah Kategori
                </Button>
              </PermissionGuard>
            </div>
          </div>

          <div className='rounded-lg border bg-card shadow-xs overflow-hidden'>
            <Table>
              <TableHeader>
                <TableRow className='bg-muted/40'>
                  <TableHead className='w-16 text-center'>Urutan</TableHead>
                  <TableHead>Nama Kategori</TableHead>
                  <TableHead>Slug Identifier</TableHead>
                  <TableHead>Ikon</TableHead>
                  <TableHead>Deskripsi</TableHead>
                  <TableHead className='text-end w-24'>Aksi</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={6} className='h-32 text-center text-muted-foreground'>
                      Memuat daftar kategori...
                    </TableCell>
                  </TableRow>
                ) : categories.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className='h-32 text-center text-muted-foreground'>
                      Belum ada kategori yang ditambahkan.
                    </TableCell>
                  </TableRow>
                ) : (
                  categories.map((cat) => (
                    <TableRow key={cat.id} className='hover:bg-muted/30'>
                      <TableCell className='text-center font-mono font-bold text-xs'>
                        <Badge variant='outline' className='font-mono'>
                          <Hash className='size-2.5 mr-0.5' />
                          {cat.urutan}
                        </Badge>
                      </TableCell>
                      <TableCell className='font-semibold text-sm'>
                        {cat.nama}
                      </TableCell>
                      <TableCell className='font-mono text-xs text-muted-foreground'>
                        {cat.slug}
                      </TableCell>
                      <TableCell>
                        <Badge variant='secondary' className='text-[10px] font-mono'>
                          {cat.ikon || 'layers'}
                        </Badge>
                      </TableCell>
                      <TableCell className='text-xs text-muted-foreground max-w-xs truncate'>
                        {cat.deskripsi || '-'}
                      </TableCell>
                      <TableCell className='text-end'>
                        <div className='flex items-center justify-end gap-1'>
                          <PermissionGuard permission='categories:manage'>
                            <Button
                              variant='ghost'
                              size='icon'
                              className='size-7'
                              onClick={() => handleOpenEdit(cat)}
                            >
                              <Edit className='size-3.5' />
                            </Button>
                            <Button
                              variant='ghost'
                              size='icon'
                              className='size-7 text-destructive hover:text-destructive'
                              onClick={() => handleOpenDelete(cat)}
                            >
                              <Trash2 className='size-3.5' />
                            </Button>
                          </PermissionGuard>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </div>

        <CategoryFormDialog
          open={isFormOpen}
          onOpenChange={setIsFormOpen}
          categoryToEdit={categoryToEdit}
          onSuccess={fetchCategories}
        />

        <CategoryDeleteDialog
          open={isDeleteOpen}
          onOpenChange={setIsDeleteOpen}
          category={categoryToDelete}
          onSuccess={fetchCategories}
        />
      </Main>
    </>
  )
}
