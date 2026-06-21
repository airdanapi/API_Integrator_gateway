import { useCallback, useEffect, useState } from 'react'
import { api } from '../services/api'
import {
  fetchNotifications,
  markAllNotificationsRead,
  markNotificationRead,
} from '../services/notifications'

const POLL_INTERVAL_MS = 60_000

const typeLabels = {
  api_inactive: 'API tidak aktif',
  error_rate: 'Error rate',
  response_time: 'Response time',
  system: 'Sistem',
}

function formatDate(value) {
  return new Date(value).toLocaleString('id-ID', {
    day: '2-digit',
    month: '2-digit',
    year: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function BellIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="none" className="h-5 w-5" aria-hidden="true">
      <path
        d="M15 17H9m9-2.5-1.2-1.7V9a4.8 4.8 0 0 0-9.6 0v3.8L6 14.5V16h12v-1.5ZM10 19a2 2 0 0 0 4 0"
        stroke="currentColor"
        strokeWidth="1.8"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}

function NotificationBell({ apiClient = api }) {
  const [open, setOpen] = useState(false)
  const [notifications, setNotifications] = useState([])
  const [unreadCount, setUnreadCount] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [pendingId, setPendingId] = useState(null)
  const [markingAll, setMarkingAll] = useState(false)

  const load = useCallback(async () => {
    try {
      setLoading(true)
      const result = await fetchNotifications(apiClient)
      setNotifications(result.notifications ?? [])
      setUnreadCount(result.unread_count ?? 0)
      setError(null)
    } catch (err) {
      setError(err)
    } finally {
      setLoading(false)
    }
  }, [apiClient])

  useEffect(() => {
    const initialLoad = setTimeout(() => {
      load()
    }, 0)
    const timer = setInterval(() => load(), POLL_INTERVAL_MS)
    return () => {
      clearTimeout(initialLoad)
      clearInterval(timer)
    }
  }, [load])

  async function handleMarkOne(notificationId) {
    setPendingId(notificationId)
    try {
      const result = await markNotificationRead(apiClient, notificationId)
      setNotifications((items) => items.map((item) => (
        item.id === notificationId ? { ...item, is_read: true } : item
      )))
      setUnreadCount(result.unread_count ?? 0)
      setError(null)
    } catch (err) {
      setError(err)
    } finally {
      setPendingId(null)
    }
  }

  async function handleMarkAll() {
    setMarkingAll(true)
    try {
      const result = await markAllNotificationsRead(apiClient)
      setNotifications((items) => items.map((item) => ({ ...item, is_read: true })))
      setUnreadCount(result.unread_count ?? 0)
      setError(null)
    } catch (err) {
      setError(err)
    } finally {
      setMarkingAll(false)
    }
  }

  const badgeText = unreadCount > 99 ? '99+' : String(unreadCount)

  return (
    <div className="relative">
      <button
        type="button"
        aria-label={unreadCount > 0 ? `Notifikasi, ${unreadCount} belum dibaca` : 'Notifikasi'}
        aria-expanded={open}
        onClick={() => setOpen((value) => !value)}
        className="relative grid h-11 w-11 place-items-center rounded-xl border border-slate-200 bg-white text-slate-600 transition hover:bg-slate-50 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-700"
      >
        <BellIcon />
        {unreadCount > 0 && (
          <span className="absolute -right-1 -top-1 min-w-5 rounded-full bg-red-600 px-1.5 py-0.5 text-center text-[10px] font-black leading-none text-white ring-2 ring-white">
            {badgeText}
          </span>
        )}
      </button>

      {open && (
        <div className="absolute right-0 z-50 mt-3 w-[min(22rem,calc(100vw-2rem))] overflow-hidden rounded-xl border border-slate-200 bg-white shadow-xl shadow-slate-900/10">
          <div className="flex items-center justify-between gap-3 border-b border-slate-100 px-4 py-3">
            <div>
              <h2 className="text-sm font-black text-slate-900">Notifikasi</h2>
              <p className="text-xs text-slate-500">{unreadCount} belum dibaca</p>
            </div>
            {unreadCount > 0 && (
              <button
                type="button"
                onClick={handleMarkAll}
                disabled={markingAll}
                className="rounded-lg px-2.5 py-1.5 text-xs font-bold text-blue-700 transition hover:bg-blue-50 disabled:opacity-50"
              >
                {markingAll ? 'Memproses...' : 'Tandai semua'}
              </button>
            )}
          </div>

          {error && (
            <div role="alert" className="border-b border-red-100 bg-red-50 px-4 py-3 text-sm font-semibold text-red-700">
              Gagal memuat notifikasi.
            </div>
          )}

          <div className="max-h-96 overflow-y-auto">
            {loading && (
              <div className="flex items-center gap-3 px-4 py-5 text-sm text-slate-500" role="status">
                <span className="h-4 w-4 animate-spin rounded-full border-2 border-slate-200 border-t-blue-600" />
                Memuat notifikasi...
              </div>
            )}

            {!loading && notifications.length === 0 && !error && (
              <p className="px-4 py-8 text-center text-sm text-slate-500">
                Belum ada notifikasi.
              </p>
            )}

            {!loading && notifications.length > 0 && (
              <ul className="divide-y divide-slate-100">
                {notifications.map((item) => (
                  <li key={item.id} className="px-4 py-3">
                    <div className="flex items-start gap-3">
                      <span className={[
                        'mt-1 h-2.5 w-2.5 shrink-0 rounded-full',
                        item.is_read ? 'bg-slate-300' : 'bg-blue-600',
                      ].join(' ')} />
                      <div className="min-w-0 flex-1">
                        <div className="flex flex-wrap items-center gap-2">
                          <span className="rounded-full bg-slate-100 px-2 py-0.5 text-[11px] font-bold text-slate-600">
                            {typeLabels[item.type] ?? item.type}
                          </span>
                          <span className="text-xs font-bold text-slate-500">{item.app_name}</span>
                        </div>
                        <p className="mt-2 text-sm font-semibold leading-5 text-slate-800">
                          {item.message}
                        </p>
                        <div className="mt-2 flex items-center justify-between gap-3">
                          <span className="text-xs text-slate-400">{formatDate(item.created_at)}</span>
                          {!item.is_read && (
                            <button
                              type="button"
                              onClick={() => handleMarkOne(item.id)}
                              disabled={pendingId === item.id}
                              className="rounded-md px-2 py-1 text-xs font-bold text-blue-700 transition hover:bg-blue-50 disabled:opacity-50"
                            >
                              {pendingId === item.id ? 'Memproses...' : 'Tandai dibaca'}
                            </button>
                          )}
                        </div>
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

export default NotificationBell