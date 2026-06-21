import { describe, expect, it, vi } from 'vitest'
import {
  fetchChatConversations,
  fetchChatHistory,
  markChatRead,
  sendChatMessage,
} from './chat'

describe('chat service', () => {
  it('fetches conversations', async () => {
    const apiClient = {
      get: vi.fn().mockResolvedValue({ data: { data: { conversations: [], total_unread: 0 } } }),
    }

    const result = await fetchChatConversations(apiClient)

    expect(apiClient.get).toHaveBeenCalledWith('/chat/conversations')
    expect(result.total_unread).toBe(0)
  })

  it('fetches conversation history with pagination', async () => {
    const apiClient = {
      get: vi.fn().mockResolvedValue({ data: { data: { messages: [], page: 2, limit: 20 } } }),
    }

    const result = await fetchChatHistory(apiClient, 'admin__marketplace__Marketplace', { page: 2, limit: 20 })

    expect(apiClient.get).toHaveBeenCalledWith('/chat/history', {
      params: { conversation_id: 'admin__marketplace__Marketplace', page: 2, limit: 20 },
    })
    expect(result.page).toBe(2)
  })

  it('sends and marks chat messages', async () => {
    const apiClient = {
      post: vi.fn()
        .mockResolvedValueOnce({ data: { data: { message: { id: 7 } } } })
        .mockResolvedValueOnce({ data: { data: { total_unread: 0 } } }),
    }

    await sendChatMessage(apiClient, { to_username: 'marketplace', to_app_name: 'Marketplace', message: 'Halo' })
    await markChatRead(apiClient, 'admin__marketplace__Marketplace')

    expect(apiClient.post).toHaveBeenNthCalledWith(1, '/chat/message', {
      to_username: 'marketplace',
      to_app_name: 'Marketplace',
      message: 'Halo',
    })
    expect(apiClient.post).toHaveBeenNthCalledWith(2, '/chat/read', {
      conversation_id: 'admin__marketplace__Marketplace',
    })
  })
})
