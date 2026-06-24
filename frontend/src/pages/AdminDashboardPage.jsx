import { useCallback, useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import BrandMark from '../components/BrandMark'
import NotificationBell from '../components/NotificationBell'
import ChatDrawer from '../components/ChatDrawer'
import { useAuth } from '../auth/auth-context'
import { api } from '../services/api'
import { fetchAdminDashboard } from '../services/dashboard'
import { PieChart, Pie, Cell, Tooltip as RechartsTooltip, ResponsiveContainer, Legend } from 'recharts'

const POLL_INTERVAL_MS = 30_000

// ─── Status Badge ────────────────────────────────────────────────────────────

function StatusBadge({ status }) {
  const isActive = status === 'active'
  return (
    <span
      className={[
        'inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-bold',
        isActive
          ? 'bg-emerald-100 text-emerald-700 ring-1 ring-emerald-200'
          : 'bg-red-100 text-red-700 ring-1 ring-red-200',
      ].join(' ')}
    >
      <span
        className={[
          'h-1.5 w-1.5 rounded-full',
          isActive ? 'bg-emerald-500' : 'bg-red-500',
        ].join(' ')}
      />
      {isActive ? 'Aktif' : 'Tidak Aktif'}
    </span>
  )
}

// ─── Charts ──────────────────────────────────────────────────────────────────

function TrafficCompositionChart({ summary }) {
  const data = [
    { name: 'Sukses', value: summary?.success_count || 0 },
    { name: 'Error', value: summary?.error_count || 0 },
  ]
  const COLORS = ['#10b981', '#ef4444']

  return (
    <div className="rounded-2xl border border-slate-200 bg-white shadow-sm p-5 h-[320px] flex flex-col">
      <h2 className="text-sm font-bold uppercase tracking-wider text-slate-500 mb-4">Komposisi Traffic</h2>
      <div className="flex-1 min-h-0">
        <ResponsiveContainer width="100%" height="100%">
          <PieChart>
            <Pie
              data={data}
              innerRadius={60}
              outerRadius={80}
              paddingAngle={5}
              dataKey="value"
            >
              {data.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))}
            </Pie>
            <RechartsTooltip formatter={(value) => [`${value} Request`, '']} />
            <Legend verticalAlign="bottom" height={36} iconType="circle" />
          </PieChart>
        </ResponsiveContainer>
      </div>
    </div>
  )
}

function ServiceHealthChart({ indicators }) {
  const activeCount = indicators?.filter(i => i.status === 'active').length || 0
  const inactiveCount = indicators?.filter(i => i.status !== 'active').length || 0
  
  const data = [
    { name: 'Aktif', value: activeCount },
    { name: 'Tidak Aktif', value: inactiveCount },
  ]
  const COLORS = ['#10b981', '#94a3b8']

  return (
    <div className="rounded-2xl border border-slate-200 bg-white shadow-sm p-5 h-[320px] flex flex-col">
      <h2 className="text-sm font-bold uppercase tracking-wider text-slate-500 mb-4">Kesehatan Layanan</h2>
      <div className="flex-1 min-h-0">
        <ResponsiveContainer width="100%" height="100%">
          <PieChart>
            <Pie
              data={data}
              innerRadius={0}
              outerRadius={80}
              dataKey="value"
            >
              {data.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))}
            </Pie>
            <RechartsTooltip formatter={(value) => [`${value} Layanan`, '']} />
            <Legend verticalAlign="bottom" height={36} iconType="circle" />
          </PieChart>
        </ResponsiveContainer>
      </div>
    </div>
  )
}

// ─── Traffic Summary Cards ────────────────────────────────────────────────────

function TrafficSummaryCards({ summary }) {
  const cards = [
    { label: 'Total Request', value: summary.total_requests, color: 'blue' },
    { label: 'Sukses', value: summary.success_count, color: 'emerald' },
    { label: 'Error', value: summary.error_count, color: 'red' },
    { label: 'Success Rate', value: `${summary.success_rate_pct?.toFixed(1)}%`, color: 'violet' },
  ]

  return (
    <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
      {cards.map(({ label, value, color }) => (
        <div
          key={label}
          className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
        >
          <p className="text-xs font-bold uppercase tracking-wider text-slate-500">{label}</p>
          <p className={`mt-2 text-3xl font-black text-${color}-600`}>{value ?? '—'}</p>
        </div>
      ))}
    </div>
  )
}

// ─── Service Indicators ──────────────────────────────────────────────────────

function ServiceIndicatorList({ indicators }) {
  return (
    <div className="rounded-2xl border border-slate-200 bg-white shadow-sm">
      <div className="border-b border-slate-100 px-6 py-4">
        <h2 className="text-sm font-bold uppercase tracking-wider text-slate-500">
          Status Layanan
        </h2>
      </div>
      <ul className="divide-y divide-slate-100">
        {indicators.map((ind) => (
          <li
            key={ind.app_name}
            className="flex items-center justify-between px-6 py-3"
          >
            <span className="text-sm font-semibold text-slate-800">{ind.app_name}</span>
            <StatusBadge status={ind.status} />
          </li>
        ))}
      </ul>
    </div>
  )
}

// ─── HTTP Status Badge ────────────────────────────────────────────────────────

function HttpStatusBadge({ code }) {
  const isSuccess = code >= 200 && code < 300
  return (
    <span
      className={[
        'inline-block rounded px-2 py-0.5 font-mono text-xs font-bold',
        isSuccess ? 'bg-emerald-50 text-emerald-700' : 'bg-red-50 text-red-700',
      ].join(' ')}
    >
      {code}
    </span>
  )
}

// ─── Audit Log Table ─────────────────────────────────────────────────────────

function AuditLogTable({ auditLogs, onPageChange }) {
  const { items = [], total = 0, page = 1, limit = 20 } = auditLogs ?? {}
  const from = total === 0 ? 0 : (page - 1) * limit + 1
  const to = Math.min(page * limit, total)
  const totalPages = Math.ceil(total / limit)

  return (
    <div className="rounded-2xl border border-slate-200 bg-white shadow-sm">
      <div className="flex items-center justify-between border-b border-slate-100 px-6 py-4">
        <h2 className="text-sm font-bold uppercase tracking-wider text-slate-500">
          Audit Log
        </h2>
        <span className="text-xs text-slate-400">
          {from}–{to} dari {total}
        </span>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-slate-100 bg-slate-50 text-xs font-bold uppercase tracking-wider text-slate-500">
              <th className="px-6 py-3 text-left">Aplikasi</th>
              <th className="px-6 py-3 text-left">Endpoint</th>
              <th className="px-6 py-3 text-left">Method</th>
              <th className="px-6 py-3 text-left">Status</th>
              <th className="px-6 py-3 text-left">Waktu</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {items.map((log) => (
              <tr key={log.id} className="hover:bg-slate-50">
                <td className="px-6 py-3 font-semibold text-slate-800">{log.source_app}</td>
                <td className="px-6 py-3 font-mono text-slate-600">{log.endpoint}</td>
                <td className="px-6 py-3 text-slate-500">{log.method}</td>
                <td className="px-6 py-3">
                  <HttpStatusBadge code={log.status} />
                </td>
                <td className="px-6 py-3 text-slate-400">
                  {new Date(log.timestamp).toLocaleString('id-ID', { dateStyle: 'short', timeStyle: 'short' })}
                </td>
              </tr>
            ))}
            {items.length === 0 && (
              <tr>
                <td colSpan={5} className="px-6 py-10 text-center text-slate-400">
                  Tidak ada data log.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-between border-t border-slate-100 px-6 py-3">
          <button
            type="button"
            disabled={page <= 1}
            onClick={() => onPageChange(page - 1)}
            className="rounded-lg border border-slate-200 px-3 py-1.5 text-xs font-semibold text-slate-600 disabled:opacity-40"
          >
            Sebelumnya
          </button>
          <span className="text-xs text-slate-500">
            Halaman {page} dari {totalPages}
          </span>
          <button
            type="button"
            disabled={page >= totalPages}
            onClick={() => onPageChange(page + 1)}
            className="rounded-lg border border-slate-200 px-3 py-1.5 text-xs font-semibold text-slate-600 disabled:opacity-40"
          >
            Berikutnya
          </button>
        </div>
      )}
    </div>
  )
}

// ─── AdminDashboardPage ──────────────────────────────────────────────────────

function AdminDashboardPage({ apiClient: propApiClient, fetchData = fetchAdminDashboard }) {
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

  // Initial load + polling setiap 30 detik
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

  function handleRefresh() {
    load(true)
  }

  function handleLogout() {
    logout()
    navigate('/login', { replace: true })
  }

  return (
    <div className="min-h-screen bg-slate-100 text-slate-900">
      {/* Header */}
      <header className="border-b border-slate-200 bg-white">
        <div className="mx-auto flex max-w-6xl items-center justify-between gap-4 px-5 py-5 sm:px-8">
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
            <NotificationBell apiClient={client} />
            <ChatDrawer apiClient={client} currentUser={user} />
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
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-black tracking-tight text-slate-900">
              Dashboard Admin Gateway
            </h1>
            <p className="mt-1 text-sm text-slate-500">
              Data 7 hari terakhir · refresh otomatis setiap 30 detik
            </p>
          </div>
          {!loading && (
            <button
              type="button"
              onClick={handleRefresh}
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
            <TrafficSummaryCards summary={data.traffic_summary} />
            <div className="grid gap-6 lg:grid-cols-2">
              <TrafficCompositionChart summary={data.traffic_summary} />
              <ServiceHealthChart indicators={data.service_indicators ?? []} />
            </div>
            <div className="grid gap-6 lg:grid-cols-3">
              <div className="lg:col-span-1">
                <ServiceIndicatorList indicators={data.service_indicators ?? []} />
              </div>
              <div className="lg:col-span-2">
                <AuditLogTable
                  auditLogs={data.audit_logs}
                  onPageChange={(p) => setPage(p)}
                />
              </div>
            </div>
          </>
        )}
      </main>
    </div>
  )
}

export default AdminDashboardPage
