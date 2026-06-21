import { useState, useEffect, useCallback } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../auth/auth-context'
import { api } from '../services/api'
import { fetchMonitoringDashboard } from '../services/dashboard'
import NotificationBell from '../components/NotificationBell'

const POLL_INTERVAL_MS = 30_000

// ─── Sub-components ─────────────────────────────────────────────────────────

function SummaryCard({ label, value, sub }) {
  return (
    <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
      <p className="text-xs font-semibold uppercase tracking-wider text-slate-400">{label}</p>
      <p className="mt-2 text-3xl font-black text-slate-900">{value}</p>
      {sub && <p className="mt-1 text-xs text-slate-500">{sub}</p>}
    </div>
  )
}

function StatusDot({ status }) {
  return status === 'active' ? (
    <span className="flex h-2.5 w-2.5 rounded-full bg-emerald-500 shadow-[0_0_6px_2px] shadow-emerald-300" />
  ) : (
    <span className="flex h-2.5 w-2.5 rounded-full bg-slate-300" />
  )
}

function AppBreakdownTable({ breakdown }) {
  const maxRequests = Math.max(1, ...breakdown.map((a) => a.total_requests))
  return (
    <div className="rounded-2xl border border-slate-200 bg-white shadow-sm">
      <div className="border-b border-slate-100 px-6 py-4">
        <h2 className="text-xs font-bold uppercase tracking-widest text-slate-400">
          Traffic Per Aplikasi (7 Hari)
        </h2>
      </div>
      <ul className="divide-y divide-slate-50">
        {breakdown.map((app) => (
          <li key={app.app_name} className="flex items-center gap-4 px-6 py-4">
            <StatusDot status={app.total_requests > 0 ? 'active' : 'inactive'} />
            <span className="w-28 shrink-0 text-sm font-semibold text-slate-700">{app.app_name}</span>
            <div className="flex-1">
              <div className="h-2 overflow-hidden rounded-full bg-slate-100">
                <div
                  className="h-2 rounded-full bg-blue-500 transition-all duration-500"
                  style={{ width: `${(app.total_requests / maxRequests) * 100}%` }}
                />
              </div>
            </div>
            <span className="w-10 text-right text-sm font-bold tabular-nums text-slate-700">
              {app.total_requests}
            </span>
          </li>
        ))}
      </ul>
    </div>
  )
}

function ServiceIndicatorGrid({ indicators }) {
  return (
    <div className="rounded-2xl border border-slate-200 bg-white shadow-sm">
      <div className="border-b border-slate-100 px-6 py-4">
        <h2 className="text-xs font-bold uppercase tracking-widest text-slate-400">
          Status Layanan
        </h2>
      </div>
      <ul className="divide-y divide-slate-50">
        {indicators.map((ind) => (
          <li key={ind.app_name} className="flex items-center justify-between px-6 py-3.5">
            <div className="flex items-center gap-3">
              <StatusDot status={ind.status} />
              <span className="text-sm font-medium text-slate-700">{ind.app_name}</span>
            </div>
            {ind.status === 'active' ? (
              <span className="rounded-full bg-emerald-100 px-2.5 py-0.5 text-xs font-bold text-emerald-700">
                Aktif
              </span>
            ) : (
              <span className="rounded-full bg-slate-100 px-2.5 py-0.5 text-xs font-bold text-slate-500">
                Tidak Aktif
              </span>
            )}
          </li>
        ))}
      </ul>
    </div>
  )
}

function BrandMark({ compact }) {
  if (compact)
    return (
      <span className="text-lg font-black tracking-tight text-slate-800">
        API<span className="text-blue-600">·</span>GW
      </span>
    )
  return (
    <span className="text-2xl font-black tracking-tight text-slate-800">
      API Integrator<span className="text-blue-600"> Gateway</span>
    </span>
  )
}

// ─── MonitoringDashboardPage ─────────────────────────────────────────────────

function MonitoringDashboardPage({ apiClient: propApiClient, fetchData = fetchMonitoringDashboard }) {
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  const client = propApiClient ?? api

  const [data, setData] = useState(null)
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [error, setError] = useState(null)

  const load = useCallback(async (isManual = false) => {
    if (isManual) setRefreshing(true)
    try {
      const result = await fetchData(client)
      setData(result)
      setError(null)
    } catch (err) {
      setError(err)
    } finally {
      setLoading(false)
      if (isManual) setRefreshing(false)
    }
  }, [fetchData, client])

  useEffect(() => {
    const initialLoad = setTimeout(() => {
      load()
    }, 0)
    const timer = setInterval(() => load(false), POLL_INTERVAL_MS)
    return () => {
      clearTimeout(initialLoad)
      clearInterval(timer)
    }
  }, [load])

  function handleLogout() {
    logout()
    navigate('/login', { replace: true })
  }

  const summary = data?.traffic_summary ?? {}

  return (
    <div className="min-h-screen bg-slate-100 text-slate-900">
      {/* Header */}
      <header className="border-b border-slate-200 bg-white">
        <div className="mx-auto flex max-w-6xl items-center justify-between gap-4 px-5 py-5 sm:px-8">
          <Link
            to="/"
            aria-label="API Integrator Gateway, kembali ke beranda"
            className="rounded-xl focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-blue-700"
          >
            <BrandMark compact />
          </Link>
          <div className="flex items-center gap-3">
            <span className="hidden items-center gap-2 text-sm text-slate-500 sm:flex">
              <span className="rounded-full bg-indigo-100 px-2.5 py-0.5 text-xs font-bold text-indigo-700">
                MONITORING
              </span>
              {user?.username} · {user?.app_name}
            </span>
            <NotificationBell apiClient={client} />
            <button
              type="button"
              onClick={handleLogout}
              className="rounded-xl border border-slate-300 px-4 py-2.5 text-sm font-bold text-slate-700 transition hover:border-red-200 hover:bg-red-50 hover:text-red-700"
            >
              Logout
            </button>
          </div>
        </div>
      </header>

      {/* Main */}
      <main className="mx-auto max-w-6xl space-y-6 px-5 py-10 sm:px-8">
        {/* Title + Refresh */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-black tracking-tight text-slate-900">
              Monitoring Gateway
            </h1>
            <p className="mt-1 text-sm text-slate-500">
              Read-only · Data 7 hari terakhir · refresh otomatis setiap 30 detik
            </p>
          </div>
          {!loading && (
            <button
              type="button"
              onClick={() => load(true)}
              disabled={refreshing}
              className="inline-flex items-center gap-2 rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm font-semibold text-slate-600 shadow-sm transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {refreshing ? (
                <>
                  <span className="h-3.5 w-3.5 animate-spin rounded-full border-2 border-slate-300 border-t-slate-600" />
                  Memperbarui...
                </>
              ) : (
                'Refresh'
              )}
            </button>
          )}
        </div>

        {/* Loading */}
        {loading && (
          <div className="flex items-center justify-center py-24" role="status" aria-label="Memuat data">
            <div className="h-8 w-8 animate-spin rounded-full border-4 border-blue-200 border-t-blue-600" />
          </div>
        )}

        {/* Error */}
        {!loading && error && (
          <div role="alert" className="rounded-2xl border border-red-200 bg-red-50 px-6 py-5 text-sm text-red-800">
            Gagal memuat data monitoring. Pastikan koneksi ke server aktif, lalu coba refresh.
          </div>
        )}

        {/* Content */}
        {!loading && !error && data && (
          <>
            {/* Summary Cards */}
            <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
              <SummaryCard
                label="Total Request"
                value={summary.total_requests ?? 0}
                sub="Semua aplikasi, 7 hari"
              />
              <SummaryCard
                label="Sukses"
                value={summary.success_count ?? 0}
              />
              <SummaryCard
                label="Error"
                value={summary.error_count ?? 0}
              />
              <SummaryCard
                label="Success Rate"
                value={`${(summary.success_rate_pct ?? 0).toFixed(1)}%`}
              />
            </div>

            {/* 2-col: Service Indicators + App Breakdown */}
            <div className="grid gap-6 lg:grid-cols-2">
              <ServiceIndicatorGrid indicators={data.service_indicators ?? []} />
              <AppBreakdownTable breakdown={data.app_breakdown ?? []} />
            </div>
          </>
        )}
      </main>
    </div>
  )
}

export default MonitoringDashboardPage
