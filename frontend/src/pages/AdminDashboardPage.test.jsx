import { render, screen, waitFor, within } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { describe, expect, it, vi } from 'vitest'
import AdminDashboardPage from './AdminDashboardPage'

// ─── mock auth-context ────────────────────────────────────────────────────────

vi.mock('../auth/auth-context', () => ({
  useAuth: () => ({
    user: { username: 'admin', app_name: 'API Gateway', role: 'admin_gateway' },
    logout: vi.fn(),
  }),
}))

// ─── sample data ─────────────────────────────────────────────────────────────

const sampleData = {
  traffic_summary: {
    total_requests: 100,
    success_count: 84,
    error_count: 16,
    success_rate_pct: 84.0,
    avg_duration_ms: 120,
  },
  service_indicators: [
    { app_name: 'Marketplace', status: 'inactive', last_request: '2026-06-10T00:00:00Z' },
    { app_name: 'POS', status: 'active', last_request: '2026-06-19T07:00:00Z' },
    { app_name: 'SupplierHub', status: 'active', last_request: '2026-06-19T06:00:00Z' },
    { app_name: 'LogistiKita', status: 'inactive', last_request: '2026-06-10T00:00:00Z' },
    { app_name: 'SmartBank', status: 'active', last_request: '2026-06-19T08:00:00Z' },
  ],
  audit_logs: {
    items: [
      { id: 1, source_app: 'POS', endpoint: '/gateway/payment', method: 'POST', status: 200, timestamp: '2026-06-19T07:00:00Z' },
      { id: 2, source_app: 'Marketplace', endpoint: '/gateway/marketplace', method: 'POST', status: 400, timestamp: '2026-06-18T12:00:00Z' },
    ],
    total: 50,
    page: 1,
    limit: 20,
  },
}

// ─── helpers ─────────────────────────────────────────────────────────────────


vi.mock('../components/NotificationBell', () => ({
  default: () => <button type="button">Notifikasi</button>,
}))
function renderPage(fetchData) {
  return render(
    <MemoryRouter>
      <AdminDashboardPage fetchData={fetchData} />
    </MemoryRouter>,
  )
}

// Helper: tunggu traffic card total_requests muncul = data sudah selesai dimuat
async function waitForData() {
  return screen.findByText('100', {}, { timeout: 3000 })
}

// ─── tests ───────────────────────────────────────────────────────────────────

describe('AdminDashboardPage', () => {
  it('shows loading spinner initially', () => {
    const neverResolves = vi.fn(() => new Promise(() => {}))
    renderPage(neverResolves)
    expect(screen.getByRole('status')).toBeInTheDocument()
  })

  it('renders traffic summary cards with correct data', async () => {
    const mockFetch = vi.fn().mockResolvedValue(sampleData)
    renderPage(mockFetch)

    expect(await waitForData()).toBeInTheDocument()
    expect(screen.getByText('84')).toBeInTheDocument()
    expect(screen.getByText('84.0%')).toBeInTheDocument()
    expect(screen.getByText('16')).toBeInTheDocument()
  })

  it('renders service indicators with correct status badges', async () => {
    const mockFetch = vi.fn().mockResolvedValue(sampleData)
    renderPage(mockFetch)

    await waitForData()

    // 'Marketplace' muncul di service indicator LIST dan di audit log TABLE → pakai getAllByText
    const marketplaceCells = screen.getAllByText('Marketplace')
    expect(marketplaceCells.length).toBeGreaterThanOrEqual(1)

    const posCells = screen.getAllByText('POS')
    expect(posCells.length).toBeGreaterThanOrEqual(1)

    expect(screen.getByText('SupplierHub')).toBeInTheDocument()

    const inactiveBadges = screen.getAllByText('Tidak Aktif')
    expect(inactiveBadges.length).toBeGreaterThanOrEqual(2)

    const activeBadges = screen.getAllByText('Aktif')
    expect(activeBadges.length).toBeGreaterThanOrEqual(3)
  })

  it('renders audit log table with items', async () => {
    const mockFetch = vi.fn().mockResolvedValue(sampleData)
    renderPage(mockFetch)

    const table = await screen.findByRole('table')
    expect(within(table).getByText('POS')).toBeInTheDocument()
    expect(within(table).getByText('/gateway/payment')).toBeInTheDocument()
    expect(within(table).getByText('200')).toBeInTheDocument()
    expect(within(table).getByText('400')).toBeInTheDocument()
  })

  it('shows total entries count in audit logs', async () => {
    const mockFetch = vi.fn().mockResolvedValue(sampleData)
    renderPage(mockFetch)

    await waitForData()
    expect(screen.getByText(/1–20 dari 50/)).toBeInTheDocument()
  })

  it('shows error alert when API fails', async () => {
    const mockFetch = vi.fn().mockRejectedValue(new Error('network error'))
    renderPage(mockFetch)

    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument()
    })
  })

  it('calls fetchData with correct default pagination params', async () => {
    const mockFetch = vi.fn().mockResolvedValue(sampleData)
    renderPage(mockFetch)

    await waitForData()
    expect(mockFetch).toHaveBeenCalledWith(
      expect.anything(),
      { page: 1, limit: 20 },
    )
  })
})
