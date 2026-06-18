import { Navigate, Route, Routes } from 'react-router-dom'
import ProtectedRoute from './auth/ProtectedRoute'
import DashboardPage from './pages/DashboardPage'
import LandingPage from './pages/LandingPage'
import LoginPage from './pages/LoginPage'

function App() {
  return (
    <Routes>
      <Route path="/" element={<LandingPage />} />
      <Route path="/login" element={<LoginPage />} />

      <Route element={<ProtectedRoute role="admin_gateway" />}>
        <Route path="/dashboard/admin" element={<DashboardPage />} />
      </Route>
      <Route element={<ProtectedRoute role="app_user" />}>
        <Route path="/dashboard/user" element={<DashboardPage />} />
      </Route>
      <Route element={<ProtectedRoute role="monitoring_user" />}>
        <Route path="/dashboard/monitoring" element={<DashboardPage />} />
      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default App
