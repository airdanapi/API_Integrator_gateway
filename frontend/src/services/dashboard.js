/**
 * Dashboard API service.
 * Mengambil data analytics dari endpoint /dashboard/admin.
 */

/**
 * Fetch data dashboard admin.
 * @param {Object} apiClient - Axios instance dengan auth interceptor
 * @param {{ page?: number, limit?: number }} [params]
 * @returns {Promise<Object>} data dari backend (traffic_summary, service_indicators, audit_logs)
 */
export async function fetchAdminDashboard(apiClient, params = {}) {
  const { page = 1, limit = 20 } = params
  const response = await apiClient.get('/dashboard/admin', {
    params: { page, limit },
  })
  return response.data.data
}
