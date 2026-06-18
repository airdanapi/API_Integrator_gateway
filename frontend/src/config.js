const DEFAULT_API_BASE_URL = 'http://localhost:8080'

export function getApiBaseUrl(env = {}) {
  const configuredUrl = env.VITE_API_BASE_URL?.trim()
  return (configuredUrl || DEFAULT_API_BASE_URL).replace(/\/+$/, '')
}

export const API_BASE_URL = getApiBaseUrl(import.meta.env)
