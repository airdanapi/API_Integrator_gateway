import { useState, useEffect, useCallback } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../auth/auth-context'
import { api } from '../services/api'
import { fetchUserDashboard } from '../services/dashboard'

const POLL_INTERVAL_MS = 30_000

// ─── Sub-components ─────────────────────────────────────────────────────────

function StatCard({ label, value, color }) {
  const colorMap = {
    blue: 'border-blue-200 bg-blue-50 text-blue-700',
    green: 'border-emerald-200 bg-emerald-50 text-emerald-700',
    red: 'border-red-200 bg-red-50 text-red-700',
    purple: 'border-violet-200 bg-violet-50 text-violet-700',
  }
  const textMap = {
    blue: 'text-blue-900',
    green: 'text-emerald-900',
    red: 'text-red-900',
    purple: 'text-violet-900',
  }
  return (
    <div className={`rounded-2xl border p-6 ${colorMap[color]}`}>
      <p className="text-xs font-semibold uppercase tracking-wider opacity-70">{label}</p>
      <p className={`mt-2 text-4xl font-black ${textMap[color]}`}>{value}</p>
    </div>
  )
}

function ServiceStatusBadge({ status }) {
  return status === 'active' ? (
    <span className="inline-flex items-center gap-1.5 rounded-full bg-emerald-100 px-3 py-1 text-xs font-bold text-emerald-700">
      <span className="h-2 w-2 rounded-full bg-emerald-500" />
      Aktif
    </span>
  ) : (
    <span className="inline-flex items-center gap-1.5 rounded-full bg-red-100 px-3 py-1 text-xs font-bold text-red-700">
      <span className="h-2 w-2 rounded-full bg-red-400" />
      Tidak Aktif
    </span>
  )
}

function HttpStatusBadge({ status }) {
  const isSuccess = status >= 200 && status < 300
  return (
    <span
      className={`inline-block rounded-md px-2 py-0.5 text-xs font-bold tabular-nums ${
        isSuccess ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-700'
      }`}
    >
      {status}
    </span>
  )
}

function RecentLogsTable({ logs, totalLogs, page, limit, onPageChange }) {
  const from = logs.length === 0 ? 0 : (page - 1) * limit + 1
  const to = (page - 1) * limit + logs.length

  function fmt(ts) {
    return new Date(ts).toLocaleString('id-ID', {
      day: '2-digit', month: '2-digit', year: '2-digit',
      hour: '2-digit', minute: '2-digit',
    })
  }

  return (
    <div className="rounded-2xl border border-slate-200 bg-white shadow-sm">
      <div className="flex items-center justify-between border-b border-slate-100 px-6 py-4">
        <h2 className="text-xs font-bold uppercase tracking-widest text-slate-400">
          Riwayat Request Saya
        </h2>
        <span className="text-xs text-slate-400">
          {logs.length === 0 ? 'Tidak ada data' : `${from}–${to} dari ${totalLogs}`}
        </span>
      </div>
      <div className="overflow-x-auto">
        <table className="min-w-full text-sm">
          <thead>
            <tr className="border-b border-slate-100 text-xs uppercase tracking-wider text-slate-400">
              <th className="py-3 pl-6 pr-3 text-left font-semibold">Endpoint</th>
              <th className="px-3 py-3 text-left font-semibold">Method</th>
              <th className="px-3 py-3 text-left font-semibold">Status</th>
              <th className="py-3 pl-3 pr-6 text-left font-semibold">Waktu</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-50">
            {logs.map((log) => (
              <tr key={log.id} className="hover:bg-slate-50/60">
                <td className="py-3 pl-6 pr-3 font-mono text-xs text-slate-600">{log.endpoint}</td>
                <td className="px-3 py-3">
                  <span className="rounded bg-slate-100 px-1.5 py-0.5 text-xs font-bold text-slate-500">
                    {log.method}
                  </span>
                </td>
                <td className="px-3 py-3">
                  <HttpStatusBadge status={log.status} />
                </td>
                <td className="py-3 pl-3 pr-6 text-xs text-slate-500">{fmt(log.timestamp)}</td>
              </tr>
            ))}
            {logs.length === 0 && (
              <tr>
                <td colSpan={4} className="py-12 text-center text-sm text-slate-400">
                  Belum ada riwayat request dari aplikasi ini.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
      {totalLogs > limit && (
        <div className="flex items-center justify-end gap-2 border-t border-slate-100 px-6 py-3">
          <button
            type="button"
            disabled={page <= 1}
            onClick={() => onPageChange(page - 1)}
            className="rounded-lg border border-slate-200 px-3 py-1.5 text-xs font-semibold text-slate-600 transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-40"
          >
            Sebelumnya
          </button>
          <button
            type="button"
            disabled={(page - 1) * limit + logs.length >= totalLogs}
            onClick={() => onPageChange(page + 1)}
            className="rounded-lg border border-slate-200 px-3 py-1.5 text-xs font-semibold text-slate-600 transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-40"
          >
            Berikutnya
          </button>
        </div>
      )}
    </div>
  )
}

// ─── BrandMark (dipinjam dari komponen lain) ─────────────────────────────────
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

// ─── UserDashboardPage ────────────────────────────────────────────────────────

function UserDashboardPage({ apiClient: propApiClient, fetchData = fetchUserDashboard }) {
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  const client = propApiClient ?? api

  const [data, setData] = useState(null)
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [error, setError] = useState(null)
  const [page, setPage] = useState(1)

  const load = useCallback(async (isManual = false) => {
    if (isManual) setRefreshing(true)
    try {
      const result = await fetchData(client, { page, limit: 20 })
      setData(result)
      setError(null)
    } catch (err) {
      setError(err)
    } finally {
      setLoading(false)
      if (isManual) setRefreshing(false)
    }
  }, [fetchData, client, page])

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
        <div className="mx-auto flex max-w-5xl items-center justify-between gap-4 px-5 py-5 sm:px-8">
          <Link
            to="/"
            className="rounded-xl focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-blue-700"
            aria-label="API Integrator Gateway, kembali ke beranda"
          >
            <BrandMark compact />
          </Link>
          <div className="flex items-center gap-3">
            <span className="hidden text-sm text-slate-500 sm:block">
              {user?.username} · {user?.app_name}
            </span>
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
      <main className="mx-auto max-w-5xl space-y-6 px-5 py-10 sm:px-8">
        {/* Title + Refresh */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-black tracking-tight text-slate-900">
              Dashboard {data?.my_app ?? user?.app_name}
            </h1>
            <p className="mt-1 flex items-center gap-2 text-sm text-slate-500">
              Data 7 hari terakhir · refresh otomatis setiap 30 detik
              {!loading && (
                <ServiceStatusBadge status={data?.service_status ?? 'inactive'} />
              )}
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
          <div
            role="alert"
            className="rounded-2xl border border-red-200 bg-red-50 px-6 py-5 text-sm text-red-800"
          >
            Gagal memuat data dashboard. Pastikan koneksi ke server aktif, lalu coba refresh.
          </div>
        )}

        {/* Content */}
        {!loading && !error && data && (
          <>
            {/* Traffic Stats */}
            <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
              <StatCard label="Total Request" value={summary.total_requests ?? 0} color="blue" />
              <StatCard label="Sukses" value={summary.success_count ?? 0} color="green" />
              <StatCard label="Error" value={summary.error_count ?? 0} color="red" />
              <StatCard
                label="Success Rate"
                value={`${(summary.success_rate_pct ?? 0).toFixed(1)}%`}
                color="purple"
              />
            </div>

            {/* Recent Logs */}
            <RecentLogsTable
              logs={data.recent_logs ?? []}
              totalLogs={data.total_logs ?? 0}
              page={page}
              limit={20}
              onPageChange={(p) => setPage(p)}
            />
          </>
        )}
      </main>
    </div>
  )
}

export default UserDashboardPage
