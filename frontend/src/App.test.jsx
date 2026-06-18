import { fireEvent, render, screen, within } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import App from './App'

describe('App', () => {
  it('renders every public landing page section', () => {
    render(<App />)

    expect(
      screen.getByRole('heading', { name: 'API Integrator Gateway' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('heading', { name: 'Integrasi yang aman, terukur, dan konsisten' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('heading', { name: 'Satu alur untuk seluruh ekosistem' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('heading', { name: 'Dibangun untuk kebutuhan nyata UMKM' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('heading', { name: 'Pertanyaan yang sering diajukan' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('heading', { name: 'Siap menghubungkan layanan Anda?' }),
    ).toBeInTheDocument()
  })

  it('provides anchor navigation to the main sections', () => {
    render(<App />)

    const desktopNavigation = screen.getByRole('navigation', {
      name: 'Navigasi utama',
    })

    expect(
      within(desktopNavigation).getByRole('link', { name: 'Manfaat' }),
    ).toHaveAttribute('href', '#manfaat')
    expect(
      within(desktopNavigation).getByRole('link', { name: 'Alur integrasi' }),
    ).toHaveAttribute('href', '#alur-integrasi')
    expect(
      within(desktopNavigation).getByRole('link', { name: 'Use case' }),
    ).toHaveAttribute('href', '#use-case')
    expect(
      within(desktopNavigation).getByRole('link', { name: 'FAQ' }),
    ).toHaveAttribute('href', '#faq')
  })

  it('opens and closes the accessible mobile navigation', () => {
    render(<App />)

    const menuButton = screen.getByRole('button', {
      name: 'Buka menu navigasi',
    })
    expect(menuButton).toHaveAttribute('aria-expanded', 'false')
    expect(screen.queryByRole('navigation', { name: 'Navigasi mobile' })).not.toBeInTheDocument()

    fireEvent.click(menuButton)

    expect(
      screen.getByRole('button', { name: 'Tutup menu navigasi' }),
    ).toHaveAttribute('aria-expanded', 'true')
    expect(
      screen.getByRole('navigation', { name: 'Navigasi mobile' }),
    ).toBeInTheDocument()

    fireEvent.click(
      screen.getByRole('button', { name: 'Tutup menu navigasi' }),
    )

    expect(screen.queryByRole('navigation', { name: 'Navigasi mobile' })).not.toBeInTheDocument()
  })

  it('marks login as coming soon without exposing a dead route', () => {
    render(<App />)

    const loginCta = screen.getByRole('button', { name: 'Login segera hadir' })
    expect(loginCta).toBeDisabled()
    expect(loginCta).toHaveAttribute('aria-disabled', 'true')
    expect(screen.queryByRole('link', { name: /login/i })).not.toBeInTheDocument()
    expect(screen.getAllByText('Segera hadir').length).toBeGreaterThan(0)
  })

  it('renders FAQ content with native disclosure controls', () => {
    render(<App />)

    expect(
      screen.getByText('Apakah API Integrator memproses transaksi keuangan?'),
    ).toBeInTheDocument()
    expect(
      screen.getByText(
        'Tidak. Gateway mengatur validasi, routing, dan pencatatan. Seluruh transaksi keuangan tetap diproses oleh SmartBank.',
      ),
    ).toBeInTheDocument()
  })

  it('lists every application in the integration ecosystem', () => {
    render(<App />)

    for (const application of [
      'SmartBank',
      'Marketplace',
      'POS',
      'SupplierHub',
      'LogistiKita',
      'UMKM Insight',
      'API Gateway',
    ]) {
      expect(screen.getAllByText(application).length).toBeGreaterThan(0)
    }
  })

  it('links contact actions to the official repository', () => {
    render(<App />)

    const repositoryLinks = screen.getAllByRole('link', {
      name: 'Lihat repositori',
    })
    expect(repositoryLinks.length).toBeGreaterThan(0)
    for (const repositoryLink of repositoryLinks) {
      expect(repositoryLink).toHaveAttribute(
        'href',
        'https://github.com/airdanapi/API_Integrator_gateway',
      )
    }
  })
})
