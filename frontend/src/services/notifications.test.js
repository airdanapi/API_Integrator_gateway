import { describe, expect, it, vi } from 'vitest'
import {
  fetchNotifications,
  markAllNotificationsRead,
  markNotificationRead,
} from './notifications'

describe('notifications service', () => {
  it('fetches notifications with default pagination', async () => {
    const mockClient = {
      get: vi.fn().mockResolvedValue({
        data: { status: 'success', data: { notifications: [], unread_count: 0, page: 1, limit: 10 } },
      }),
    }

    const result = await fetchNotifications(mockClient)

    expect(mockClient.get).toHaveBeenCalledWith('/notifications', {
      params: { page: 1, limit: 10 },
    })
    expect(result.unread_count).toBe(0)
  })

  it('marks one notification as read', async () => {
    const mockClient = {
      post: vi.fn().mockResolvedValue({ data: { status: 'success', data: { unread_count: 2 } } }),
    }

    const result = await markNotificationRead(mockClient, 7)

    expect(mockClient.post).toHaveBeenCalledWith('/notifications/read', { notification_id: 7 })
    expect(result.unread_count).toBe(2)
  })

  it('marks all visible notifications as read', async () => {
    const mockClient = {
      post: vi.fn().mockResolvedValue({ data: { status: 'success', data: { unread_count: 0 } } }),
    }

    const result = await markAllNotificationsRead(mockClient)

    expect(mockClient.post).toHaveBeenCalledWith('/notifications/read', { all: true })
    expect(result.unread_count).toBe(0)
  })
})