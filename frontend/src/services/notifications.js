export async function fetchNotifications(apiClient, params = {}) {
  const { page = 1, limit = 10 } = params
  const response = await apiClient.get('/notifications', {
    params: { page, limit },
  })
  return response.data.data
}

export async function markNotificationRead(apiClient, notificationId) {
  const response = await apiClient.post('/notifications/read', {
    notification_id: notificationId,
  })
  return response.data.data
}

export async function markAllNotificationsRead(apiClient) {
  const response = await apiClient.post('/notifications/read', { all: true })
  return response.data.data
}