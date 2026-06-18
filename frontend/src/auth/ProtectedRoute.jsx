import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { useAuth } from './auth-context'

function LoadingScreen() {
  return (
    <main
      className="grid min-h-screen place-items-center bg-slate-950 px-6 text-white"
      aria-live="polite"
    >
      <div className="text-center">
        <div
          className="mx-auto h-10 w-10 animate-spin rounded-full border-4 border-blue-300 border-t-transparent"
          aria-hidden="true"
        />
        <p className="mt-4 text-sm font-semibold text-slate-200">
          Memeriksa sesi...
        </p>
      </div>
    </main>
  )
}

function ProtectedRoute({ role }) {
  const { status, user, dashboardPath } = useAuth()
  const location = useLocation()

  if (status === 'loading') {
    return <LoadingScreen />
  }
  if (status !== 'authenticated') {
    return (
      <Navigate
        to="/login"
        replace
        state={{ from: location.pathname }}
      />
    )
  }
  if (user.role !== role) {
    return <Navigate to={dashboardPath} replace />
  }
  return <Outlet />
}

export default ProtectedRoute
