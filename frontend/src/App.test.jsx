import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import App from './App'
import { AuthProvider } from './auth/AuthContext'
import { ACCESS_TOKEN_KEY } from './auth/session'

// Mencegah Dashboard page memanggil API sungguhan saat test navigasi App.
// Tanpa ini, api.get('/dashboard/...') dipanggil → server mengembalikan 401
// → token dihapus → user terpaksa logout → heading dashboard tidak ditemukan.
vi.mock('./components/NotificationBell', () => ({
  default: () => <button type="button">Notifikasi</button>,
}))
vi.mock('./services/dashboard', () => ({
  fetchAdminDashboard: vi.fn(() => new Promise(() => {})),      // loading selamanya
  fetchUserDashboard: vi.fn(() => new Promise(() => {})),       // loading selamanya
  fetchMonitoringDashboard: vi.fn(() => new Promise(() => {})), // loading selamanya
}))

const adminUser = {
  user_id: '1',
  username: 'admin',
  role: 'admin_gateway',
  app_name: 'API Gateway',
}

function createApiClient() {
  return {
    get: vi.fn(),
    post: vi.fn(),
  }
}

function renderApp({
  path = '/',
  apiClient = createApiClient(),
} = {}) {
  return {
    apiClient,
    ...render(
      <MemoryRouter initialEntries={[path]}>
        <AuthProvider apiClient={apiClient}>
          <App />
        </AuthProvider>
      </MemoryRouter>,
    ),
  }
}

function loginResponse({
  token = 'signed-token',
  role = 'admin_gateway',
  appName = 'API Gateway',
  dashboardUrl = '/dashboard/admin',
} = {}) {
  return {
    data: {
      status: 'success',
      data: {
        token,
        role,
        app_name: appName,
        dashboard_url: dashboardUrl,
        expires_in: 3600,
      },
    },
  }
}

function meResponse(user = adminUser) {
  return {
    data: {
      status: 'success',
      data: user,
    },
  }
}

describe('Sprint 4 authentication frontend', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('keeps the landing page public and exposes login navigation', async () => {
    renderApp()

    expect(
      await screen.findByRole('heading', { name: 'API Integrator Gateway' }),
    ).toBeInTheDocument()
    expect(screen.getAllByRole('link', { name: 'Login' }).length).toBeGreaterThan(0)
    expect(screen.queryByText('Login segera hadir')).not.toBeInTheDocument()
  })

  it('preserves accessible landing navigation and mobile controls', async () => {
    renderApp()
    await screen.findByRole('heading', { name: 'API Integrator Gateway' })

    const navigation = screen.getByRole('navigation', {
      name: 'Navigasi utama',
    })
    expect(within(navigation).getByRole('link', { name: 'Manfaat' }))
      .toHaveAttribute('href', '#manfaat')

    const menuButton = screen.getByRole('button', {
      name: 'Buka menu navigasi',
    })
    fireEvent.click(menuButton)
    expect(
      screen.getByRole('navigation', { name: 'Navigasi mobile' }),
    ).toBeInTheDocument()
    expect(
      within(screen.getByRole('navigation', { name: 'Navigasi mobile' }))
        .getByRole('link', { name: 'Login' }),
    ).toHaveAttribute('href', '/login')
  })

  it('renders a keyboard-accessible login form with all official applications', async () => {
    renderApp({ path: '/login' })

    expect(
      await screen.findByRole('heading', { name: 'Masuk ke API Integrator' }),
    ).toBeInTheDocument()
    expect(screen.getByLabelText('Username')).toBeRequired()
    expect(screen.getByLabelText('Password')).toBeRequired()

    const applicationSelect = screen.getByLabelText('Aplikasi')
    const options = within(applicationSelect).getAllByRole('option')
    expect(options.map((option) => option.textContent)).toEqual([
      'Pilih aplikasi',
      'SmartBank',
      'Marketplace',
      'POS',
      'SupplierHub',
      'LogistiKita',
      'UMKM Insight',
      'API Gateway',
    ])
  })

  it('logs in, stores the token, validates the session, and routes by role', async () => {
    const apiClient = createApiClient()
    apiClient.post.mockResolvedValue(loginResponse())
    apiClient.get.mockResolvedValue(meResponse())
    renderApp({ path: '/login', apiClient })

    fireEvent.change(await screen.findByLabelText('Username'), {
      target: { value: 'admin' },
    })
    fireEvent.change(screen.getByLabelText('Password'), {
      target: { value: 'admin-development-password' },
    })
    fireEvent.change(screen.getByLabelText('Aplikasi'), {
      target: { value: 'API Gateway' },
    })
    fireEvent.click(screen.getByRole('button', { name: 'Masuk' }))

    expect(
      await screen.findByRole('heading', { name: 'Dashboard Admin Gateway' }),
    ).toBeInTheDocument()
    expect(localStorage.getItem(ACCESS_TOKEN_KEY)).toBe('signed-token')
    expect(apiClient.post).toHaveBeenCalledWith('/auth/login', {
      username: 'admin',
      password: 'admin-development-password',
      app_name: 'API Gateway',
    })
    expect(apiClient.get).toHaveBeenCalledWith('/auth/me')
  })

  it('shows a generic message for invalid credentials', async () => {
    const apiClient = createApiClient()
    apiClient.post.mockRejectedValue({
      response: {
        status: 401,
        data: {
          status: 'error',
          error: { code: 'invalid_credentials' },
        },
      },
    })
    renderApp({ path: '/login', apiClient })

    fireEvent.change(await screen.findByLabelText('Username'), {
      target: { value: 'admin' },
    })
    fireEvent.change(screen.getByLabelText('Password'), {
      target: { value: 'wrong-password' },
    })
    fireEvent.change(screen.getByLabelText('Aplikasi'), {
      target: { value: 'API Gateway' },
    })
    fireEvent.click(screen.getByRole('button', { name: 'Masuk' }))

    expect(
      await screen.findByRole('alert'),
    ).toHaveTextContent('Username, password, atau aplikasi tidak valid.')
    expect(localStorage.getItem(ACCESS_TOKEN_KEY)).toBeNull()
  })

  it('rejects a login response whose dashboard does not match its role', async () => {
    const apiClient = createApiClient()
    apiClient.post.mockResolvedValue(loginResponse({
      dashboardUrl: '/dashboard/user',
    }))
    renderApp({ path: '/login', apiClient })

    fireEvent.change(await screen.findByLabelText('Username'), {
      target: { value: 'admin' },
    })
    fireEvent.change(screen.getByLabelText('Password'), {
      target: { value: 'admin-development-password' },
    })
    fireEvent.change(screen.getByLabelText('Aplikasi'), {
      target: { value: 'API Gateway' },
    })
    fireEvent.click(screen.getByRole('button', { name: 'Masuk' }))

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Respons autentikasi tidak valid.',
    )
    expect(localStorage.getItem(ACCESS_TOKEN_KEY)).toBeNull()
    expect(apiClient.get).not.toHaveBeenCalled()
  })

  it('restores a valid session after refresh', async () => {
    localStorage.setItem(ACCESS_TOKEN_KEY, 'stored-token')
    const apiClient = createApiClient()
    apiClient.get.mockResolvedValue(meResponse())
    renderApp({ path: '/dashboard/admin', apiClient })

    expect(
      await screen.findByRole('heading', { name: 'Dashboard Admin Gateway' }),
    ).toBeInTheDocument()
    // AdminDashboardPage menampilkan "admin · API Gateway" di satu span
    expect(screen.getByText(/API Gateway/)).toBeInTheDocument()
    expect(apiClient.get).toHaveBeenCalledWith('/auth/me')
  })

  it('removes an invalid session and redirects protected routes to login', async () => {
    localStorage.setItem(ACCESS_TOKEN_KEY, 'expired-token')
    const apiClient = createApiClient()
    apiClient.get.mockRejectedValue({ response: { status: 401 } })
    renderApp({ path: '/dashboard/admin', apiClient })

    expect(
      await screen.findByRole('heading', { name: 'Masuk ke API Integrator' }),
    ).toBeInTheDocument()
    expect(localStorage.getItem(ACCESS_TOKEN_KEY)).toBeNull()
  })

  it('redirects an authenticated user away from a dashboard for another role', async () => {
    localStorage.setItem(ACCESS_TOKEN_KEY, 'app-user-token')
    const apiClient = createApiClient()
    apiClient.get.mockResolvedValue(meResponse({
      user_id: '2',
      username: 'marketplace',
      role: 'app_user',
      app_name: 'Marketplace',
    }))
    renderApp({ path: '/dashboard/admin', apiClient })

    // UserDashboardPage: heading = "Dashboard {app_name}" (loading state — app_name from user object)
    expect(
      await screen.findByRole('heading', { name: /Dashboard/ }),
    ).toBeInTheDocument()
  })

  it('renders the protected monitoring placeholder for a monitoring user', async () => {
    localStorage.setItem(ACCESS_TOKEN_KEY, 'monitoring-token')
    const apiClient = createApiClient()
    apiClient.get.mockResolvedValue(meResponse({
      user_id: '3',
      username: 'insight',
      role: 'monitoring_user',
      app_name: 'UMKM Insight',
    }))
    renderApp({ path: '/dashboard/monitoring', apiClient })

    // MonitoringDashboardPage: heading = "Monitoring Gateway"
    expect(
      await screen.findByRole('heading', { name: 'Monitoring Gateway' }),
    ).toBeInTheDocument()
  })

  it('shows dashboard and logout actions on the landing page for an active session', async () => {
    localStorage.setItem(ACCESS_TOKEN_KEY, 'stored-token')
    const apiClient = createApiClient()
    apiClient.get.mockResolvedValue(meResponse())
    renderApp({ path: '/', apiClient })

    expect(
      await screen.findByRole('link', { name: 'Dashboard' }),
    ).toHaveAttribute('href', '/dashboard/admin')
    expect(screen.getByRole('button', { name: 'Logout' })).toBeInTheDocument()
    expect(screen.queryByRole('link', { name: 'Login' })).not.toBeInTheDocument()
  })

  it('logs out locally and returns to login', async () => {
    localStorage.setItem(ACCESS_TOKEN_KEY, 'stored-token')
    const apiClient = createApiClient()
    apiClient.get.mockResolvedValue(meResponse())
    renderApp({ path: '/dashboard/admin', apiClient })

    fireEvent.click(await screen.findByRole('button', { name: 'Logout' }))

    expect(
      await screen.findByRole('heading', { name: 'Masuk ke API Integrator' }),
    ).toBeInTheDocument()
    expect(localStorage.getItem(ACCESS_TOKEN_KEY)).toBeNull()
  })

  it('redirects an authenticated visitor from login to their dashboard', async () => {
    localStorage.setItem(ACCESS_TOKEN_KEY, 'stored-token')
    const apiClient = createApiClient()
    apiClient.get.mockResolvedValue(meResponse())
    renderApp({ path: '/login', apiClient })

    await waitFor(() => {
      expect(
        screen.getByRole('heading', { name: 'Dashboard Admin Gateway' }),
      ).toBeInTheDocument()
    })
  })
})
