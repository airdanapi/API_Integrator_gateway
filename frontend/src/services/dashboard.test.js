import { describe, expect, it, vi } from 'vitest'
import { fetchAdminDashboard, fetchUserDashboard, fetchMonitoringDashboard } from './dashboard'

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

describe('fetchUserDashboard', () => {
  it('calls GET /dashboard/user with correct params and returns data', async () => {
    const mockClient = {
      get: vi.fn().mockResolvedValue({
        data: {
          status: 'success',
          data: {
            my_app: 'Marketplace',
            service_status: 'active',
            traffic_summary: { total_requests: 15, success_count: 12, error_count: 3, success_rate_pct: 80 },
            recent_logs: [],
            total_logs: 15,
            page: 1,
            limit: 20,
          },
        },
      }),
    }

    const result = await fetchUserDashboard(mockClient, { page: 1, limit: 20 })

    expect(mockClient.get).toHaveBeenCalledWith('/dashboard/user', {
      params: { page: 1, limit: 20 },
    })
    expect(result.my_app).toBe('Marketplace')
    expect(result.service_status).toBe('active')
    expect(result.traffic_summary.total_requests).toBe(15)
  })

  it('uses default page=1 limit=20 when no params given', async () => {
    const mockClient = {
      get: vi.fn().mockResolvedValue({
        data: { status: 'success', data: { my_app: 'POS', service_status: 'active', traffic_summary: {}, recent_logs: [], total_logs: 0, page: 1, limit: 20 } },
      }),
    }
    await fetchUserDashboard(mockClient)
    expect(mockClient.get).toHaveBeenCalledWith('/dashboard/user', {
      params: { page: 1, limit: 20 },
    })
  })

  it('propagates errors from the API client', async () => {
    const mockClient = {
      get: vi.fn().mockRejectedValue(new Error('unauthorized')),
    }
    await expect(fetchUserDashboard(mockClient)).rejects.toThrow('unauthorized')
  })
})

describe('fetchMonitoringDashboard', () => {
  it('calls GET /dashboard/monitoring and returns data', async () => {
    const mockClient = {
      get: vi.fn().mockResolvedValue({
        data: {
          status: 'success',
          data: {
            traffic_summary: { total_requests: 65, success_count: 50, error_count: 15, success_rate_pct: 76.9 },
            service_indicators: [
              { app_name: 'Marketplace', status: 'active' },
            ],
            app_breakdown: [
              { app_name: 'Marketplace', total_requests: 20 },
            ],
          },
        },
      }),
    }

    const result = await fetchMonitoringDashboard(mockClient)

    expect(mockClient.get).toHaveBeenCalledWith('/dashboard/monitoring')
    expect(result.traffic_summary.total_requests).toBe(65)
    expect(result.service_indicators).toHaveLength(1)
    expect(result.app_breakdown[0].app_name).toBe('Marketplace')
  })

  it('propagates errors from the API client', async () => {
    const mockClient = {
      get: vi.fn().mockRejectedValue(new Error('forbidden')),
    }
    await expect(fetchMonitoringDashboard(mockClient)).rejects.toThrow('forbidden')
  })
})
