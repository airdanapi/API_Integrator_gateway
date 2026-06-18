import { Link, useNavigate } from 'react-router-dom'
import BrandMark from '../components/BrandMark'
import { ROLE_LABELS } from '../auth/constants'
import { useAuth } from '../auth/auth-context'

const dashboardContent = {
  admin_gateway: {
    title: 'Dashboard Admin Gateway',
    description:
      'Monitoring traffic, audit log, dan indikator layanan akan tersedia pada Sprint 7.',
  },
  app_user: {
    title: 'Dashboard Pengguna Aplikasi',
    description:
      'Status layanan, riwayat request, dan grafik performa akan tersedia pada Sprint 8.',
  },
  monitoring_user: {
    title: 'Dashboard Monitoring',
    description:
      'Ringkasan analitik read-only UMKM Insight akan tersedia pada Sprint 8.',
  },
}

function DashboardPage() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  const content = dashboardContent[user.role]

  function handleLogout() {
    logout()
    navigate('/login', { replace: true })
  }

  return (
    <div className="min-h-screen bg-slate-100 text-slate-900">
      <header className="border-b border-slate-200 bg-white">
        <div className="mx-auto flex max-w-6xl items-center justify-between gap-4 px-5 py-5 sm:px-8">
          <Link
            to="/"
            className="rounded-xl focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-blue-700"
            aria-label="API Integrator Gateway, kembali ke beranda"
          >
            <BrandMark compact />
          </Link>
          <button
            type="button"
            onClick={handleLogout}
            className="rounded-xl border border-slate-300 px-4 py-2.5 text-sm font-bold text-slate-700 transition hover:border-red-200 hover:bg-red-50 hover:text-red-700 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-700"
          >
            Logout
          </button>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-5 py-12 sm:px-8 sm:py-16">
        <div className="overflow-hidden rounded-3xl bg-slate-950 text-white shadow-2xl shadow-slate-900/10">
          <div className="border-b border-white/10 px-6 py-8 sm:px-10">
            <span className="inline-flex rounded-full bg-blue-500/15 px-3 py-1 text-xs font-bold uppercase tracking-[0.18em] text-blue-200 ring-1 ring-inset ring-blue-400/30">
              Sesi aktif
            </span>
            <h1 className="mt-5 text-3xl font-black tracking-tight sm:text-4xl">
              {content.title}
            </h1>
            <p className="mt-4 max-w-2xl leading-7 text-slate-300">
              {content.description}
            </p>
          </div>

          <dl className="grid gap-px bg-white/10 sm:grid-cols-3">
            <div className="bg-slate-900 px-6 py-7 sm:px-10">
              <dt className="text-xs font-bold uppercase tracking-wider text-slate-400">
                Username
              </dt>
              <dd className="mt-2 text-lg font-bold">{user.username}</dd>
            </div>
            <div className="bg-slate-900 px-6 py-7 sm:px-10">
              <dt className="text-xs font-bold uppercase tracking-wider text-slate-400">
                Aplikasi
              </dt>
              <dd className="mt-2 text-lg font-bold">{user.app_name}</dd>
            </div>
            <div className="bg-slate-900 px-6 py-7 sm:px-10">
              <dt className="text-xs font-bold uppercase tracking-wider text-slate-400">
                Role
              </dt>
              <dd className="mt-2 text-lg font-bold">
                {ROLE_LABELS[user.role]}
              </dd>
            </div>
          </dl>
        </div>

        <div className="mt-8 rounded-2xl border border-blue-200 bg-blue-50 px-6 py-5 text-sm leading-6 text-blue-950">
          Halaman ini adalah placeholder autentikasi Sprint 4. Data operasional
          belum dimuat sampai endpoint dan UI dashboard selesai pada sprint
          berikutnya.
        </div>
      </main>
    </div>
  )
}

export default DashboardPage
