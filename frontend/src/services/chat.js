/**
 * Chat service – fetch conversations, history, send message, and mark read.
 */

export async function fetchChatConversations(apiClient) {
  const { data } = await apiClient.get('/chat/conversations')
  return data.data
}

export async function fetchChatHistory(apiClient, conversationId, { page = 1, limit = 20 } = {}) {
  const { data } = await apiClient.get('/chat/history', {
    params: { conversation_id: conversationId, page, limit },
  })
  return data.data
}

export async function sendChatMessage(apiClient, { to_username, to_app_name, message }) {
  const { data } = await apiClient.post('/chat/message', { to_username, to_app_name, message })
  return data.data
}

export async function markChatRead(apiClient, conversationId) {
  const { data } = await apiClient.post('/chat/read', { conversation_id: conversationId })
  return data.data
}
