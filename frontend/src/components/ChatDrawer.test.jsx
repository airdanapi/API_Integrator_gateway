import { act, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import ChatDrawer from './ChatDrawer'

const adminUser = { username: 'admin', app_name: 'API Gateway', role: 'admin_gateway' }
const appUser = { username: 'marketplace', app_name: 'Marketplace', role: 'app_user' }

function conversations(totalUnread = 2) {
  return {
    conversations: [
      {
        conversation_id: 'admin__marketplace__Marketplace',
        target_username: 'marketplace',
        target_app_name: 'Marketplace',
        unread_count: totalUnread,
        latest_message: {
          id: 2,
          conversation_id: 'admin__marketplace__Marketplace',
          from_user: 'marketplace',
          to_user: 'admin',
          message: 'Butuh bantuan integrasi',
          timestamp: '2026-06-21T10:01:00Z',
          is_read: false,
        },
      },
    ],
    total_unread: totalUnread,
  }
}

function history() {
  return {
    messages: [
      {
        id: 1,
        conversation_id: 'admin__marketplace__Marketplace',
        from_user: 'admin',
        to_user: 'marketplace',
        message: 'Halo Marketplace',
        timestamp: '2026-06-21T10:00:00Z',
        is_read: true,
      },
      {
        id: 2,
        conversation_id: 'admin__marketplace__Marketplace',
        from_user: 'marketplace',
        to_user: 'admin',
        message: 'Butuh bantuan integrasi',
        timestamp: '2026-06-21T10:01:00Z',
        is_read: false,
      },
    ],
    total_unread: 2,
    page: 1,
    limit: 20,
  }
}

function createClient({ empty = false, fail = false } = {}) {
  if (fail) {
    return {
      get: vi.fn().mockRejectedValue(new Error('network')),
      post: vi.fn(),
    }
  }
  return {
    get: vi.fn((url) => {
      if (url === '/chat/conversations') {
        return Promise.resolve({ data: { data: empty ? { conversations: [], total_unread: 0 } : conversations() } })
      }
      if (url === '/chat/history') {
        return Promise.resolve({ data: { data: history() } })
      }
      return Promise.reject(new Error(`unexpected ${url}`))
    }),
    post: vi.fn((url) => {
      if (url === '/chat/read') {
        return Promise.resolve({ data: { data: { total_unread: 0 } } })
      }
      if (url === '/chat/message') {
        return Promise.resolve({ data: { data: { message: { id: 3, message: 'Balasan admin' }, total_unread: 0 } } })
      }
      return Promise.reject(new Error(`unexpected ${url}`))
    }),
  }
}

describe('ChatDrawer', () => {
  beforeEach(() => {
    vi.useRealTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('renders unread badge and opens conversation history', async () => {
    const apiClient = createClient()
    render(<ChatDrawer apiClient={apiClient} currentUser={adminUser} />)

    expect(await screen.findByText('2')).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: /chat/i }))
    fireEvent.click(await screen.findByRole('button', { name: /marketplace/i }))

    expect(await screen.findByText('Halo Marketplace')).toBeInTheDocument()
    expect(screen.getByText('Butuh bantuan integrasi')).toBeInTheDocument()
    expect(apiClient.post).toHaveBeenCalledWith('/chat/read', {
      conversation_id: 'admin__marketplace__Marketplace',
    })
  })

  it('sends a message in the active conversation', async () => {
    const apiClient = createClient()
    render(<ChatDrawer apiClient={apiClient} currentUser={adminUser} />)
    fireEvent.click(await screen.findByRole('button', { name: /chat/i }))
    fireEvent.click(await screen.findByRole('button', { name: /marketplace/i }))

    fireEvent.change(await screen.findByLabelText('Pesan'), { target: { value: 'Balasan admin' } })
    fireEvent.click(screen.getByRole('button', { name: /kirim/i }))

    await waitFor(() => expect(apiClient.post).toHaveBeenCalledWith('/chat/message', {
      to_username: 'marketplace',
      to_app_name: 'Marketplace',
      message: 'Balasan admin',
    }))
  })

  it('renders empty and error states', async () => {
    const emptyClient = createClient({ empty: true })
    const { rerender } = render(<ChatDrawer apiClient={emptyClient} currentUser={adminUser} />)
    fireEvent.click(screen.getByRole('button', { name: /chat/i }))
    expect(await screen.findByText('Belum ada percakapan.')).toBeInTheDocument()

    const failingClient = createClient({ fail: true })
    rerender(<ChatDrawer apiClient={failingClient} currentUser={adminUser} />)
    fireEvent.click(screen.getByRole('button', { name: /chat/i }))
    expect(await screen.findByRole('alert')).toHaveTextContent('Gagal memuat chat.')
  })

  it('polls conversations every 2 seconds', async () => {
    vi.useFakeTimers()
    const apiClient = createClient()
    render(<ChatDrawer apiClient={apiClient} currentUser={appUser} />)

    await act(async () => {
      await vi.advanceTimersByTimeAsync(0)
    })
    expect(apiClient.get).toHaveBeenCalledWith('/chat/conversations')

    await act(async () => {
      await vi.advanceTimersByTimeAsync(2_000)
    })
    const conversationCalls = apiClient.get.mock.calls.filter(([url]) => url === '/chat/conversations')
    expect(conversationCalls.length).toBeGreaterThanOrEqual(2)
  })
})
