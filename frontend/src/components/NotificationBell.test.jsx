import { fireEvent, render, screen, waitFor, act } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import NotificationBell from './NotificationBell'

function response(unreadCount = 1, notifications = [
  {
    id: 1,
    created_at: '2026-06-21T10:00:00Z',
    app_name: 'Marketplace',
    type: 'api_inactive',
    message: 'API Marketplace tidak aktif selama lebih dari 1 minggu.',
    is_read: false,
  },
]) {
  return {
    status: 'success',
    notifications,
    unread_count: unreadCount,
    page: 1,
    limit: 10,
  }
}

describe('NotificationBell', () => {
  beforeEach(() => {
    vi.useRealTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('renders unread badge and dropdown notifications', async () => {
    const apiClient = {
      get: vi.fn().mockResolvedValue({ data: { data: response(2) } }),
      post: vi.fn(),
    }
    render(<NotificationBell apiClient={apiClient} />)

    expect(await screen.findByText('2')).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: /notifikasi/i }))

    expect(screen.getByText('API Marketplace tidak aktif selama lebih dari 1 minggu.')).toBeInTheDocument()
    expect(screen.getByText('Marketplace')).toBeInTheDocument()
  })

  it('marks a single notification as read', async () => {
    const apiClient = {
      get: vi.fn().mockResolvedValue({ data: { data: response(1) } }),
      post: vi.fn().mockResolvedValue({ data: { data: { unread_count: 0 } } }),
    }
    render(<NotificationBell apiClient={apiClient} />)
    await screen.findByText('1')
    fireEvent.click(screen.getByRole('button', { name: /notifikasi/i }))

    fireEvent.click(screen.getByRole('button', { name: /tandai dibaca/i }))

    await waitFor(() => expect(apiClient.post).toHaveBeenCalledWith('/notifications/read', { notification_id: 1 }))
    expect(screen.queryByText('1')).not.toBeInTheDocument()
  })

  it('marks all notifications as read', async () => {
    const apiClient = {
      get: vi.fn().mockResolvedValue({ data: { data: response(3) } }),
      post: vi.fn().mockResolvedValue({ data: { data: { unread_count: 0 } } }),
    }
    render(<NotificationBell apiClient={apiClient} />)
    await screen.findByText('3')
    fireEvent.click(screen.getByRole('button', { name: /notifikasi/i }))

    fireEvent.click(screen.getByRole('button', { name: /tandai semua/i }))

    await waitFor(() => expect(apiClient.post).toHaveBeenCalledWith('/notifications/read', { all: true }))
  })

  it('renders empty state', async () => {
    const apiClient = {
      get: vi.fn().mockResolvedValue({ data: { data: response(0, []) } }),
      post: vi.fn(),
    }
    render(<NotificationBell apiClient={apiClient} />)
    fireEvent.click(screen.getByRole('button', { name: /notifikasi/i }))

    expect(await screen.findByText('Belum ada notifikasi.')).toBeInTheDocument()
  })

  it('renders error state', async () => {
    const failingClient = {
      get: vi.fn().mockRejectedValue(new Error('network')),
      post: vi.fn(),
    }
    render(<NotificationBell apiClient={failingClient} />)
    await waitFor(() => expect(failingClient.get).toHaveBeenCalled())
    fireEvent.click(screen.getByRole('button', { name: /notifikasi/i }))

    expect(await screen.findByRole('alert')).toHaveTextContent('Gagal memuat notifikasi.')
  })

  it('polls notifications every 60 seconds', async () => {
    vi.useFakeTimers()
    const apiClient = {
      get: vi.fn().mockResolvedValue({ data: { data: response(0, []) } }),
      post: vi.fn(),
    }
    render(<NotificationBell apiClient={apiClient} />)

    await act(async () => {
      await vi.advanceTimersByTimeAsync(0)
    })
    expect(apiClient.get).toHaveBeenCalledTimes(1)

    await act(async () => {
      await vi.advanceTimersByTimeAsync(60_000)
    })
    expect(apiClient.get).toHaveBeenCalledTimes(2)
  })
})