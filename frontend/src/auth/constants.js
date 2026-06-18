export const ROLE_DASHBOARD_PATHS = Object.freeze({
  admin_gateway: '/dashboard/admin',
  app_user: '/dashboard/user',
  monitoring_user: '/dashboard/monitoring',
})

export const ROLE_LABELS = Object.freeze({
  admin_gateway: 'Admin Gateway',
  app_user: 'Pengguna Aplikasi',
  monitoring_user: 'Monitoring Read-only',
})

export function getDashboardPath(role) {
  return ROLE_DASHBOARD_PATHS[role] ?? null
}
