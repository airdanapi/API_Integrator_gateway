import axios from 'axios'
import { API_BASE_URL } from '../config'
import {
  clearAccessToken,
  getAccessToken,
} from '../auth/session'

export const AUTH_UNAUTHORIZED_EVENT = 'auth:unauthorized'

function announceUnauthorized() {
  if (typeof window !== 'undefined') {
    window.dispatchEvent(new Event(AUTH_UNAUTHORIZED_EVENT))
  }
}

export function createApiClient(baseURL) {
  const client = axios.create({ baseURL })

  client.interceptors.request.use((config) => {
    const token = getAccessToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  })

  client.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response?.status === 401) {
        clearAccessToken()
        announceUnauthorized()
      }
      return Promise.reject(error)
    },
  )

  return client
}

export const api = createApiClient(API_BASE_URL)
