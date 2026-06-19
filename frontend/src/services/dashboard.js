/**
 * Dashboard API service.
 * Mengambil data analytics dari endpoint /dashboard/*.
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

/**
 * Fetch data dashboard user (difilter berdasarkan app_name dari token).
 * @param {Object} apiClient - Axios instance dengan auth interceptor
 * @param {{ page?: number, limit?: number }} [params]
 * @returns {Promise<Object>} data: my_app, traffic_summary, service_status, recent_logs, total_logs
 */
export async function fetchUserDashboard(apiClient, params = {}) {
  const { page = 1, limit = 20 } = params
  const response = await apiClient.get('/dashboard/user', {
    params: { page, limit },
  })
  return response.data.data
}

/**
 * Fetch data dashboard monitoring (read-only, semua aplikasi).
 * @param {Object} apiClient - Axios instance dengan auth interceptor
 * @returns {Promise<Object>} data: traffic_summary, service_indicators, app_breakdown
 */
export async function fetchMonitoringDashboard(apiClient) {
  const response = await apiClient.get('/dashboard/monitoring')
  return response.data.data
}
