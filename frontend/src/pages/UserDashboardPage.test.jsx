import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import UserDashboardPage from './UserDashboardPage'

vi.mock('../auth/auth-context', () => ({
  useAuth: () => ({
    user: { username: 'marketplace-user', app_name: 'Marketplace', role: 'app_user' },
    logout: vi.fn(),
  }),
}))

function renderPage(fetchData) {
  return render(
    <MemoryRouter>
      <UserDashboardPage fetchData={fetchData} />
    </MemoryRouter>
  )
}

function waitForDataLoaded() {
  return waitFor(() => expect(screen.queryByRole('status')).not.toBeInTheDocument(), {
    timeout: 3000,
  })
}

const mockData = {
  my_app: 'Marketplace',
  service_status: 'active',
  traffic_summary: {
    total_requests: 25,
    success_count: 20,
    error_count: 5,
    success_rate_pct: 80.0,
  },
  recent_logs: [
    {
      id: 1,
      source_app: 'Marketplace',
      endpoint: '/gateway/payment',
      method: 'POST',
      status: 200,
      timestamp: new Date().toISOString(),
    },
    {
      id: 2,
      source_app: 'Marketplace',
      endpoint: '/gateway/payment',
      method: 'POST',
      status: 400,
      timestamp: new Date().toISOString(),
    },
  ],
  total_logs: 25,
  page: 1,
  limit: 20,
}

describe('UserDashboardPage', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('shows loading spinner initially', () => {
    const fetchData = vi.fn(() => new Promise(() => {}))
    renderPage(fetchData)
    expect(screen.getByRole('status')).toBeInTheDocument()
  })

  it('renders traffic summary cards with app data', async () => {
    const fetchData = vi.fn().mockResolvedValue(mockData)
    renderPage(fetchData)
    await waitForDataLoaded()

    expect(screen.getByText('25')).toBeInTheDocument()
    expect(screen.getByText('20')).toBeInTheDocument()
    expect(screen.getByText('5')).toBeInTheDocument()
    expect(screen.getByText('80.0%')).toBeInTheDocument()
  })

  it('renders service status badge active', async () => {
    const fetchData = vi.fn().mockResolvedValue(mockData)
    renderPage(fetchData)
    await waitForDataLoaded()

    expect(screen.getByText('Aktif')).toBeInTheDocument()
  })

  it('renders recent logs table with entries', async () => {
    const fetchData = vi.fn().mockResolvedValue(mockData)
    renderPage(fetchData)
    await waitForDataLoaded()

    // '/gateway/payment' muncul di 2 baris — gunakan getAllByText
    const endpoints = screen.getAllByText('/gateway/payment')
    expect(endpoints.length).toBeGreaterThanOrEqual(1)

    // Status 200 dan 400 muncul sebagai badge
    expect(screen.getAllByText('200').length).toBeGreaterThanOrEqual(1)
    expect(screen.getAllByText('400').length).toBeGreaterThanOrEqual(1)
  })

  it('shows inactive status when no requests', async () => {
    const inactiveData = {
      ...mockData,
      service_status: 'inactive',
      traffic_summary: { total_requests: 0, success_count: 0, error_count: 0, success_rate_pct: 0 },
      recent_logs: [],
      total_logs: 0,
    }
    const fetchData = vi.fn().mockResolvedValue(inactiveData)
    renderPage(fetchData)
    await waitForDataLoaded()

    expect(screen.getByText('Tidak Aktif')).toBeInTheDocument()
    expect(screen.getByText('Belum ada riwayat request dari aplikasi ini.')).toBeInTheDocument()
  })

  it('shows error alert when API fails', async () => {
    const fetchData = vi.fn().mockRejectedValue(new Error('Network Error'))
    renderPage(fetchData)
    await waitFor(() => expect(screen.getByRole('alert')).toBeInTheDocument())
  })

  it('calls fetchData with default pagination params', async () => {
    const fetchData = vi.fn().mockResolvedValue(mockData)
    renderPage(fetchData)
    await waitForDataLoaded()

    expect(fetchData).toHaveBeenCalledWith(
      expect.anything(),
      expect.objectContaining({ page: 1, limit: 20 })
    )
  })
})
