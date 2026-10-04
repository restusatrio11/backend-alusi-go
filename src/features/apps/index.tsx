import { useState, useEffect } from 'react'
import { Plus, RefreshCw, Layers, ArrowUpDown, List } from 'lucide-react'
import { App, Category } from '@/types/apps'
import { appsApi } from './api/apps-api'
import { categoriesApi } from '@/features/categories/api/categories-api'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { AppsTable } from './components/apps-table'
import { AppFormDialog } from './components/app-form-dialog'
import { AppLogoDialog } from './components/app-logo-dialog'
import { AppDeleteDialog } from './components/app-delete-dialog'
import { AppGuidesDialog } from './components/app-guides-dialog'
import { AppReorderView } from './components/app-reorder-view'

export function Apps() {
  const [apps, setApps] = useState<App[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [activeTab, setActiveTab] = useState('table')

  // Dialog States
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [appToEdit, setAppToEdit] = useState<App | null>(null)

  const [isLogoOpen, setIsLogoOpen] = useState(false)
  const [appForLogo, setAppForLogo] = useState<App | null>(null)

  const [isDeleteOpen, setIsDeleteOpen] = useState(false)
  const [appToDelete, setAppToDelete] = useState<App | null>(null)

  const [isGuidesOpen, setIsGuidesOpen] = useState(false)
  const [appForGuides, setAppForGuides] = useState<App | null>(null)

  const fetchData = async () => {
    setIsLoading(true)
    try {
      const [appsRes, catRes] = await Promise.all([
        appsApi.getApps({ per_page: 100 }),
        categoriesApi.getCategories(),
      ])
      setApps(appsRes.apps)
      setCategories(catRes)
    } catch {
      setApps([])
      setCategories([])
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  const handleOpenCreate = () => {
    setAppToEdit(null)
    setIsFormOpen(true)
  }

  const handleOpenEdit = (app: App) => {
    setAppToEdit(app)
    setIsFormOpen(true)
  }

  const handleOpenLogo = (app: App) => {
    setAppForLogo(app)
    setIsLogoOpen(true)
  }

  const handleOpenDelete = (app: App) => {
    setAppToDelete(app)
    setIsDeleteOpen(true)
  }

  const handleOpenGuides = (app: App) => {
    setAppForGuides(app)
    setIsGuidesOpen(true)
  }

  return (
    <>
      {/* Top Header */}
      <Header>
        <div className='flex items-center gap-2'>
          <Layers className='size-5 text-primary' />
          <h2 className='text-sm font-semibold tracking-tight'>Katalog Aplikasi ALUSI</h2>
        </div>
        <div className='ml-auto flex items-center space-x-2'>
          <ThemeSwitch />
          <ProfileDropdown />
        </div>
      </Header>

      {/* Main Content */}
      <Main>
        <div className='space-y-6'>
          {/* Page Heading */}
          <div className='flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4'>
            <div>
              <h1 className='text-2xl font-bold tracking-tight'>Manajemen Aplikasi</h1>
              <p className='text-sm text-muted-foreground'>
                Kelola metadata aplikasi, kategori, logo, hak akses SSO, dan panduan teknis portal ALUSI BPS.
              </p>
            </div>

            <div className='flex items-center gap-2'>
              <Button
                variant='outline'
                size='sm'
                onClick={fetchData}
                disabled={isLoading}
              >
                <RefreshCw className={`mr-1.5 size-3.5 ${isLoading ? 'animate-spin' : ''}`} />
                Refresh
              </Button>
              <Button size='sm' onClick={handleOpenCreate}>
                <Plus className='mr-1.5 size-4' />
                Tambah Aplikasi
              </Button>
            </div>
          </div>

          {/* View Tabs */}
          <Tabs value={activeTab} onValueChange={setActiveTab} className='space-y-4'>
            <div className='flex items-center justify-between'>
              <TabsList>
                <TabsTrigger value='table' className='flex items-center gap-1.5'>
                  <List className='size-3.5' />
                  Daftar Tabel
                </TabsTrigger>
                <TabsTrigger value='reorder' className='flex items-center gap-1.5'>
                  <ArrowUpDown className='size-3.5' />
                  Urutan Tampilan
                </TabsTrigger>
              </TabsList>
            </div>

            <TabsContent value='table' className='m-0'>
              <AppsTable
                apps={apps}
                categories={categories}
                isLoading={isLoading}
                onEdit={handleOpenEdit}
                onDelete={handleOpenDelete}
                onUploadLogo={handleOpenLogo}
                onManageGuides={handleOpenGuides}
                onRefresh={fetchData}
              />
            </TabsContent>

            <TabsContent value='reorder' className='m-0'>
              <AppReorderView
                initialApps={apps}
                onSuccess={fetchData}
              />
            </TabsContent>
          </Tabs>
        </div>

        {/* Dialogs */}
        <AppFormDialog
          open={isFormOpen}
          onOpenChange={setIsFormOpen}
          appToEdit={appToEdit}
          categories={categories}
          onSuccess={fetchData}
        />

        <AppLogoDialog
          open={isLogoOpen}
          onOpenChange={setIsLogoOpen}
          app={appForLogo}
          onSuccess={fetchData}
        />

        <AppDeleteDialog
          open={isDeleteOpen}
          onOpenChange={setIsDeleteOpen}
          app={appToDelete}
          onSuccess={fetchData}
        />

        <AppGuidesDialog
          open={isGuidesOpen}
          onOpenChange={setIsGuidesOpen}
          app={appForGuides}
        />
      </Main>
    </>
  )
}
