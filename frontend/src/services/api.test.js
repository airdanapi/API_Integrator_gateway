import { AxiosError } from 'axios'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import {
  AUTH_UNAUTHORIZED_EVENT,
  createApiClient,
} from './api'
import { getAccessToken, setAccessToken } from '../auth/session'

describe('API client', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('attaches the stored Bearer token to requests', async () => {
    const client = createApiClient('https://gateway.example.test')
    let requestConfig
    setAccessToken('signed-token')
    client.defaults.adapter = async (config) => {
      requestConfig = config
      return {
        config,
        data: { status: 'success' },
        headers: {},
        status: 200,
        statusText: 'OK',
      }
    }

    await client.get('/auth/me')

    expect(requestConfig.headers.Authorization).toBe('Bearer signed-token')
  })

  it('clears the token and announces unauthorized responses', async () => {
    const client = createApiClient('https://gateway.example.test')
    const unauthorizedListener = vi.fn()
    window.addEventListener(AUTH_UNAUTHORIZED_EVENT, unauthorizedListener)
    setAccessToken('expired-token')
    client.defaults.adapter = async (config) => {
      throw new AxiosError(
        'Unauthorized',
        'ERR_BAD_REQUEST',
        config,
        null,
        {
          config,
          data: {
            status: 'error',
            error: { code: 'unauthorized' },
          },
          headers: {},
          status: 401,
          statusText: 'Unauthorized',
        },
      )
    }

    await expect(client.get('/auth/me')).rejects.toBeInstanceOf(AxiosError)

    expect(getAccessToken()).toBeNull()
    expect(unauthorizedListener).toHaveBeenCalledOnce()
    window.removeEventListener(AUTH_UNAUTHORIZED_EVENT, unauthorizedListener)
  })
})
