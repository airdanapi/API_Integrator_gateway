import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../auth/auth-context'
import BrandMark from './BrandMark'
import { navigationItems } from '../data/landingContent'

function NavigationLinks({ onNavigate }) {
  return navigationItems.map((item) => (
    <a
      key={item.href}
      href={item.href}
      className="rounded-lg px-3 py-2 text-sm font-semibold text-slate-600 transition hover:bg-blue-50 hover:text-blue-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-700"
      onClick={onNavigate}
    >
      {item.label}
    </a>
  ))
}

function Header() {
  const [menuOpen, setMenuOpen] = useState(false)
  const { status, dashboardPath, logout } = useAuth()
  const navigate = useNavigate()
  const authenticated = status === 'authenticated'

  function handleLogout() {
    logout()
    setMenuOpen(false)
    navigate('/login')
  }

  const authActions = authenticated ? (
    <>
      <Link
        to={dashboardPath}
        className="rounded-xl border border-blue-200 px-4 py-2.5 text-sm font-bold text-blue-800 transition hover:bg-blue-50 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-700"
      >
        Dashboard
      </Link>
      <button
        type="button"
        onClick={handleLogout}
        className="rounded-xl bg-slate-950 px-4 py-2.5 text-sm font-bold text-white transition hover:bg-slate-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-950"
      >
        Logout
      </button>
    </>
  ) : (
    <Link
      to="/login"
      className="rounded-xl border border-blue-200 px-4 py-2.5 text-sm font-bold text-blue-800 transition hover:bg-blue-50 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-700"
    >
      Login
    </Link>
  )

  return (
    <header className="sticky top-0 z-50 border-b border-slate-200/80 bg-white/95 backdrop-blur">
      <div className="mx-auto flex h-20 max-w-7xl items-center justify-between px-5 sm:px-8">
        <a
          href="#beranda"
          className="rounded-xl focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-blue-700"
          aria-label="API Integrator Gateway, kembali ke beranda"
        >
          <BrandMark compact />
        </a>

        <nav
          className="hidden items-center gap-1 lg:flex"
          aria-label="Navigasi utama"
        >
          <NavigationLinks />
        </nav>

        <div className="hidden items-center gap-3 lg:flex">
          {authActions}
          <a
            href="#kontak"
            className="rounded-xl bg-blue-700 px-4 py-2.5 text-sm font-bold text-white shadow-sm transition hover:bg-blue-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-700"
          >
            Hubungi kami
          </a>
        </div>

        <button
          type="button"
          className="grid h-11 w-11 place-items-center rounded-xl border border-slate-200 text-slate-700 transition hover:bg-slate-50 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-700 lg:hidden"
          aria-expanded={menuOpen}
          aria-controls="mobile-navigation"
          aria-label={menuOpen ? 'Tutup menu navigasi' : 'Buka menu navigasi'}
          onClick={() => setMenuOpen((isOpen) => !isOpen)}
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            className="h-6 w-6"
            aria-hidden="true"
          >
            {menuOpen ? (
              <path
                d="m6 6 12 12M18 6 6 18"
                stroke="currentColor"
                strokeWidth="1.8"
                strokeLinecap="round"
              />
            ) : (
              <path
                d="M5 7h14M5 12h14M5 17h14"
                stroke="currentColor"
                strokeWidth="1.8"
                strokeLinecap="round"
              />
            )}
          </svg>
        </button>
      </div>

      {menuOpen ? (
        <nav
          id="mobile-navigation"
          className="border-t border-slate-200 bg-white px-5 py-4 lg:hidden"
          aria-label="Navigasi mobile"
        >
          <div className="mx-auto flex max-w-7xl flex-col gap-1">
            <NavigationLinks onNavigate={() => setMenuOpen(false)} />
            {authenticated ? (
              <>
                <Link
                  to={dashboardPath}
                  className="mt-3 rounded-xl border border-blue-200 px-4 py-3 text-center text-sm font-bold text-blue-800"
                  onClick={() => setMenuOpen(false)}
                >
                  Dashboard
                </Link>
                <button
                  type="button"
                  className="rounded-xl border border-slate-300 px-4 py-3 text-center text-sm font-bold text-slate-700"
                  onClick={handleLogout}
                >
                  Logout
                </button>
              </>
            ) : (
              <Link
                to="/login"
                className="mt-3 rounded-xl border border-blue-200 px-4 py-3 text-center text-sm font-bold text-blue-800"
                onClick={() => setMenuOpen(false)}
              >
                Login
              </Link>
            )}
            <a
              href="#kontak"
              className="rounded-xl bg-blue-700 px-4 py-3 text-center text-sm font-bold text-white"
              onClick={() => setMenuOpen(false)}
            >
              Hubungi kami
            </a>
          </div>
        </nav>
      ) : null}
    </header>
  )
}

export default Header
