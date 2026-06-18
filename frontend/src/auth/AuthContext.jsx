import {
  useCallback,
  useEffect,
  useMemo,
  useState,
} from 'react'
import {
  api,
  AUTH_UNAUTHORIZED_EVENT,
} from '../services/api'
import { getDashboardPath } from './constants'
import {
  clearAccessToken,
  getAccessToken,
  setAccessToken,
} from './session'
import { AuthContext } from './auth-context'

const anonymousState = {
  status: 'anonymous',
  user: null,
  dashboardPath: null,
}

class AuthError extends Error {
  constructor(message, code) {
    super(message)
    this.name = 'AuthError'
    this.code = code
  }
}

function validateUser(user) {
  const dashboardPath = getDashboardPath(user?.role)
  if (!dashboardPath || !user?.username || !user?.app_name) {
    throw new AuthError(
      'Respons autentikasi tidak valid.',
      'invalid_contract',
    )
  }
  return dashboardPath
}

function validateLoginResult(result) {
  const dashboardPath = getDashboardPath(result?.role)
  if (
    !result?.token ||
    !dashboardPath ||
    result.dashboard_url !== dashboardPath ||
    !result.app_name
  ) {
    throw new AuthError(
      'Respons autentikasi tidak valid.',
      'invalid_contract',
    )
  }
  return dashboardPath
}

export function AuthProvider({ children, apiClient = api }) {
  const [authState, setAuthState] = useState({
    status: 'loading',
    user: null,
    dashboardPath: null,
  })

  const becomeAnonymous = useCallback(() => {
    clearAccessToken()
    setAuthState(anonymousState)
  }, [])

  useEffect(() => {
    const handleUnauthorized = () => {
      becomeAnonymous()
    }
    window.addEventListener(AUTH_UNAUTHORIZED_EVENT, handleUnauthorized)
    return () => {
      window.removeEventListener(AUTH_UNAUTHORIZED_EVENT, handleUnauthorized)
    }
  }, [becomeAnonymous])

  useEffect(() => {
    let active = true

    async function restoreSession() {
      if (!getAccessToken()) {
        if (active) {
          setAuthState(anonymousState)
        }
        return
      }

      try {
        const response = await apiClient.get('/auth/me')
        const user = response.data?.data
        const dashboardPath = validateUser(user)
        if (active) {
          setAuthState({
            status: 'authenticated',
            user,
            dashboardPath,
          })
        }
      } catch {
        clearAccessToken()
        if (active) {
          setAuthState(anonymousState)
        }
      }
    }

    restoreSession()
    return () => {
      active = false
    }
  }, [apiClient])

  const login = useCallback(async (credentials) => {
    try {
      const response = await apiClient.post('/auth/login', credentials)
      const loginResult = response.data?.data
      const dashboardPath = validateLoginResult(loginResult)
      setAccessToken(loginResult.token)

      const meResponse = await apiClient.get('/auth/me')
      const user = meResponse.data?.data
      const verifiedDashboardPath = validateUser(user)
      if (
        verifiedDashboardPath !== dashboardPath ||
        user.role !== loginResult.role ||
        user.app_name !== loginResult.app_name
      ) {
        throw new AuthError(
          'Respons autentikasi tidak valid.',
          'invalid_contract',
        )
      }

      setAuthState({
        status: 'authenticated',
        user,
        dashboardPath,
      })
      return dashboardPath
    } catch (error) {
      clearAccessToken()
      setAuthState(anonymousState)
      if (error instanceof AuthError) {
        throw error
      }
      if (
        error.response?.status === 401 &&
        error.response?.data?.error?.code === 'invalid_credentials'
      ) {
        throw new AuthError(
          'Username, password, atau aplikasi tidak valid.',
          'invalid_credentials',
        )
      }
      throw new AuthError(
        'Login gagal. Periksa koneksi dan coba lagi.',
        'login_failed',
      )
    }
  }, [apiClient])

  const logout = useCallback(() => {
    becomeAnonymous()
  }, [becomeAnonymous])

  const value = useMemo(() => ({
    ...authState,
    login,
    logout,
  }), [authState, login, logout])

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}
