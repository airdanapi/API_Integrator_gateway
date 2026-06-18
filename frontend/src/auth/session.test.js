import { beforeEach, describe, expect, it } from 'vitest'
import {
  ACCESS_TOKEN_KEY,
  clearAccessToken,
  getAccessToken,
  setAccessToken,
} from './session'

describe('auth session storage', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('stores only the access token under the agreed key', () => {
    setAccessToken('signed-token')

    expect(ACCESS_TOKEN_KEY).toBe('access_token')
    expect(getAccessToken()).toBe('signed-token')
    expect(localStorage).toHaveLength(1)
  })

  it('clears the stored access token', () => {
    setAccessToken('signed-token')

    clearAccessToken()

    expect(getAccessToken()).toBeNull()
  })
})
