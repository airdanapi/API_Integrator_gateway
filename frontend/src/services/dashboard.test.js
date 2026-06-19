import { describe, expect, it, vi } from 'vitest'
import { fetchAdminDashboard } from './dashboard'

describe('fetchAdminDashboard', () => {
  it('calls GET /dashboard/admin with correct params and returns data', async () => {
    const mockClient = {
      get: vi.fn().mockResolvedValue({
        data: {
          status: 'success',
          data: {
            traffic_summary: {
              total_requests: 50,
              success_count: 42,
              error_count: 8,
              success_rate_pct: 84.0,
              avg_duration_ms: 120,
            },
            service_indicators: [
              { app_name: 'Marketplace', status: 'inactive' },
              { app_name: 'POS', status: 'active' },
            ],
            audit_logs: {
              items: [{ id: 1, source_app: 'POS', endpoint: '/gateway/payment', status: 200 }],
              total: 50,
              page: 1,
              limit: 20,
            },
          },
        },
      }),
    }

    const result = await fetchAdminDashboard(mockClient, { page: 1, limit: 20 })

    expect(mockClient.get).toHaveBeenCalledWith('/dashboard/admin', {
      params: { page: 1, limit: 20 },
    })
    expect(result.traffic_summary.total_requests).toBe(50)
    expect(result.service_indicators).toHaveLength(2)
    expect(result.audit_logs.total).toBe(50)
  })

  it('uses default page=1 limit=20 when no params given', async () => {
    const mockClient = {
      get: vi.fn().mockResolvedValue({
        data: { status: 'success', data: { traffic_summary: {}, service_indicators: [], audit_logs: { items: [], total: 0, page: 1, limit: 20 } } },
      }),
    }
    await fetchAdminDashboard(mockClient)
    expect(mockClient.get).toHaveBeenCalledWith('/dashboard/admin', {
      params: { page: 1, limit: 20 },
    })
  })

  it('propagates errors from the API client', async () => {
    const mockClient = {
      get: vi.fn().mockRejectedValue(new Error('network error')),
    }
    await expect(fetchAdminDashboard(mockClient)).rejects.toThrow('network error')
  })
})
