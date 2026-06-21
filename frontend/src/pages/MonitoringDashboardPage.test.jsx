import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import MonitoringDashboardPage from './MonitoringDashboardPage'

vi.mock('../auth/auth-context', () => ({
  useAuth: () => ({
    user: { username: 'umkm-monitor', app_name: 'UMKM Insight', role: 'monitoring_user' },
    logout: vi.fn(),
  }),
}))


vi.mock('../components/NotificationBell', () => ({
  default: () => <button type="button">Notifikasi</button>,
}))
function renderPage(fetchData) {
  return render(
    <MemoryRouter>
      <MonitoringDashboardPage fetchData={fetchData} />
    </MemoryRouter>
  )
}

function waitForDataLoaded() {
  return waitFor(() => expect(screen.queryByRole('status')).not.toBeInTheDocument(), {
    timeout: 3000,
  })
}

const mockData = {
  traffic_summary: {
    total_requests: 65,
    success_count: 50,
    // error_count diset 8 agar tidak tabrakan dengan POS total_requests=15
    error_count: 8,
    success_rate_pct: 76.9,
  },
  service_indicators: [
    { app_name: 'Marketplace', status: 'active', last_request: new Date().toISOString() },
    { app_name: 'POS', status: 'active', last_request: new Date().toISOString() },
    { app_name: 'SupplierHub', status: 'inactive', last_request: new Date().toISOString() },
    { app_name: 'LogistiKita', status: 'active', last_request: new Date().toISOString() },
    { app_name: 'SmartBank', status: 'active', last_request: new Date().toISOString() },
  ],
  app_breakdown: [
    { app_name: 'Marketplace', total_requests: 20, success_rate_pct: 0 },
    { app_name: 'POS', total_requests: 15, success_rate_pct: 0 },
    { app_name: 'SupplierHub', total_requests: 0, success_rate_pct: 0 },
    { app_name: 'LogistiKita', total_requests: 18, success_rate_pct: 0 },
    { app_name: 'SmartBank', total_requests: 12, success_rate_pct: 0 },
  ],
}

describe('MonitoringDashboardPage', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('shows loading spinner initially', () => {
    const fetchData = vi.fn(() => new Promise(() => {}))
    renderPage(fetchData)
    expect(screen.getByRole('status')).toBeInTheDocument()
  })

  it('renders overall traffic summary cards', async () => {
    const fetchData = vi.fn().mockResolvedValue(mockData)
    renderPage(fetchData)
    await waitForDataLoaded()

    // nilai unik di traffic summary
    expect(screen.getByText('65')).toBeInTheDocument()  // total_requests
    expect(screen.getByText('50')).toBeInTheDocument()  // success_count
    expect(screen.getByText('8')).toBeInTheDocument()   // error_count (unik, tidak tabrakan)
    expect(screen.getByText('76.9%')).toBeInTheDocument()
  })

  it('renders service indicators for all 5 apps', async () => {
    const fetchData = vi.fn().mockResolvedValue(mockData)
    renderPage(fetchData)
    await waitForDataLoaded()

    // Setiap nama app muncul minimal 1x (bisa ada di kedua section: indicators + breakdown)
    expect(screen.getAllByText('Marketplace').length).toBeGreaterThanOrEqual(1)
    expect(screen.getAllByText('POS').length).toBeGreaterThanOrEqual(1)
    expect(screen.getAllByText('SupplierHub').length).toBeGreaterThanOrEqual(1)
    expect(screen.getAllByText('LogistiKita').length).toBeGreaterThanOrEqual(1)
    expect(screen.getAllByText('SmartBank').length).toBeGreaterThanOrEqual(1)
  })

  it('renders app breakdown with request counts', async () => {
    const fetchData = vi.fn().mockResolvedValue(mockData)
    renderPage(fetchData)
    await waitForDataLoaded()

    // Nilai unik per app (tidak bertabrakan dengan traffic summary)
    expect(screen.getAllByText('20').length).toBeGreaterThanOrEqual(1)  // Marketplace
    expect(screen.getAllByText('15').length).toBeGreaterThanOrEqual(1)  // POS
    expect(screen.getAllByText('12').length).toBeGreaterThanOrEqual(1)  // SmartBank
  })

  it('shows Tidak Aktif badge for inactive service', async () => {
    const fetchData = vi.fn().mockResolvedValue(mockData)
    renderPage(fetchData)
    await waitForDataLoaded()

    expect(screen.getByText('Tidak Aktif')).toBeInTheDocument()
  })

  it('shows error alert when API fails', async () => {
    const fetchData = vi.fn().mockRejectedValue(new Error('Server Error'))
    renderPage(fetchData)
    await waitFor(() => expect(screen.getByRole('alert')).toBeInTheDocument())
  })

  it('calls fetchData once on mount', async () => {
    const fetchData = vi.fn().mockResolvedValue(mockData)
    renderPage(fetchData)
    await waitForDataLoaded()
    expect(fetchData).toHaveBeenCalledTimes(1)
  })
})
