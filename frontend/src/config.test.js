import { describe, expect, it } from 'vitest'
import { getApiBaseUrl } from './config'

describe('getApiBaseUrl', () => {
  it('uses the configured API URL', () => {
    expect(
      getApiBaseUrl({ VITE_API_BASE_URL: 'https://gateway.example.test/' }),
    ).toBe('https://gateway.example.test')
  })

  it('defaults to the local backend', () => {
    expect(getApiBaseUrl({})).toBe('http://localhost:8080')
  })
})
