import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import App from './App'

describe('App', () => {
  it('renders the Sprint 1 application shell', () => {
    render(<App />)

    expect(
      screen.getByRole('heading', { name: 'API Integrator Gateway' }),
    ).toBeInTheDocument()
    expect(
      screen.getByText('Sprint 1 infrastructure is ready.'),
    ).toBeInTheDocument()
  })
})
