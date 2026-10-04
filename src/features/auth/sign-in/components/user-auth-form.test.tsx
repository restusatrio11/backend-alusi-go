import { beforeEach, describe, expect, it, vi } from 'vitest'
import { render, type RenderResult } from 'vitest-browser-react'
import { type Locator, userEvent } from 'vitest/browser'
import { UserAuthForm } from './user-auth-form'

const FORM_MESSAGES = {
  usernameEmpty: 'Username / Email / NIP wajib diisi.',
  passwordEmpty: 'Kata sandi wajib diisi.',
} as const

const navigate = vi.fn()
const setAuthMock = vi.fn()

vi.mock('@/stores/auth-store', () => ({
  useAuthStore: () => ({
    auth: {
      setAuth: setAuthMock,
    },
  }),
}))

vi.mock('@/features/auth/api/auth-api', () => ({
  authApi: {
    manualLogin: vi.fn().mockResolvedValue({
      token: 'mock-access-token',
      user: {
        id: 1,
        nama: 'Admin BPS',
        email: 'admin@bps.go.id',
        username: 'admin',
        user_type: 'internal',
        status: 'active',
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
        roles: [{ id: 1, nama: 'admin', created_at: '2026-01-01T00:00:00Z' }],
      },
    }),
    getSSOLoginURL: vi.fn().mockResolvedValue('https://sso.bps.go.id/login'),
  },
}))

vi.mock('@tanstack/react-router', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@tanstack/react-router')>()
  return {
    ...actual,
    useNavigate: () => navigate,
    Link: ({
      children,
      to,
      className,
      ...rest
    }: {
      children?: React.ReactNode
      to: string
      className?: string
    }) => (
      <a href={to} className={className} {...rest}>
        {children}
      </a>
    ),
  }
})

describe('UserAuthForm', () => {
  describe('Rendering without redirectTo', () => {
    let screen: RenderResult
    let usernameInput: Locator
    let passwordInput: Locator
    let signInButton: Locator
    let ssoButton: Locator

    beforeEach(async () => {
      vi.clearAllMocks()
      screen = await render(<UserAuthForm />)
      usernameInput = screen.getByPlaceholder('admin / nama@bps.go.id / 1995xxxx')
      passwordInput = screen.getByPlaceholder('Masukkan kata sandi')
      signInButton = screen.getByRole('button', { name: /Masuk ke Portal Admin/i })
      ssoButton = screen.getByRole('button', { name: /Masuk dengan SSO BPS Sumut/i })
    })

    it('renders fields and buttons', async () => {
      await expect.element(usernameInput).toBeInTheDocument()
      await expect.element(passwordInput).toBeInTheDocument()
      await expect.element(signInButton).toBeInTheDocument()
      await expect.element(ssoButton).toBeInTheDocument()
    })

    it('shows validation messages when submitting empty form', async () => {
      await userEvent.click(signInButton)

      await expect
        .element(screen.getByText(FORM_MESSAGES.usernameEmpty))
        .toBeInTheDocument()
      await expect
        .element(screen.getByText(FORM_MESSAGES.passwordEmpty))
        .toBeInTheDocument()
    })

    it('authenticates and navigates to default route on success', async () => {
      await userEvent.fill(usernameInput, 'admin')
      await userEvent.fill(passwordInput, 'AdminBPS1200!')

      await userEvent.click(signInButton)

      await vi.waitFor(() => expect(setAuthMock).toHaveBeenCalledOnce())
      expect(setAuthMock).toHaveBeenCalledWith(
        'mock-access-token',
        expect.objectContaining({
          username: 'admin',
          email: 'admin@bps.go.id',
        })
      )

      await vi.waitFor(() =>
        expect(navigate).toHaveBeenCalledWith({ to: '/', replace: true })
      )
    })
  })

  it('navigates to redirectTo when provided', async () => {
    vi.clearAllMocks()

    const screen = await render(
      <UserAuthForm redirectTo='/apps' />
    )

    const usernameInput = screen.getByPlaceholder('admin / nama@bps.go.id / 1995xxxx')
    const passwordInput = screen.getByPlaceholder('Masukkan kata sandi')
    const signInButton = screen.getByRole('button', { name: /Masuk ke Portal Admin/i })

    await userEvent.fill(usernameInput, 'admin')
    await userEvent.fill(passwordInput, 'AdminBPS1200!')

    await userEvent.click(signInButton)

    await vi.waitFor(() => expect(setAuthMock).toHaveBeenCalledOnce())
    await vi.waitFor(() =>
      expect(navigate).toHaveBeenCalledWith({
        to: '/apps',
        replace: true,
      })
    )
  })
})
