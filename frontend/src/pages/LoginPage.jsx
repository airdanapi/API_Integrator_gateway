import { useState } from 'react'
import { Link, Navigate, useNavigate } from 'react-router-dom'
import BrandMark from '../components/BrandMark'
import { useAuth } from '../auth/auth-context'
import { applications } from '../data/landingContent'

function LoginPage() {
  const { status, dashboardPath, login } = useAuth()
  const navigate = useNavigate()
  const [form, setForm] = useState({
    username: '',
    password: '',
    app_name: '',
  })
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  if (status === 'authenticated') {
    return <Navigate to={dashboardPath} replace />
  }

  function updateField(event) {
    const { name, value } = event.target
    setForm((current) => ({ ...current, [name]: value }))
  }

  async function handleSubmit(event) {
    event.preventDefault()
    setError('')
    setSubmitting(true)
    try {
      const destination = await login({
        username: form.username.trim(),
        password: form.password,
        app_name: form.app_name,
      })
      navigate(destination, { replace: true })
    } catch (loginError) {
      setError(loginError.message)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <main className="min-h-screen bg-slate-950 text-slate-900 lg:grid lg:grid-cols-[1.05fr_0.95fr]">
      <section className="relative hidden overflow-hidden bg-blue-700 px-12 py-14 text-white lg:flex lg:flex-col lg:justify-between">
        <div
          className="absolute -right-32 -top-32 h-96 w-96 rounded-full bg-cyan-300/20 blur-3xl"
          aria-hidden="true"
        />
        <div
          className="absolute -bottom-40 -left-28 h-[28rem] w-[28rem] rounded-full bg-indigo-950/40 blur-3xl"
          aria-hidden="true"
        />
        <div className="relative">
          <BrandMark inverse />
        </div>
        <div className="relative max-w-xl">
          <p className="text-sm font-bold uppercase tracking-[0.25em] text-blue-100">
            Portal terintegrasi
          </p>
          <h1 className="mt-5 text-5xl font-black leading-tight">
            Satu identitas untuk mengakses ekosistem UMKM.
          </h1>
          <p className="mt-6 max-w-lg text-lg leading-8 text-blue-100">
            Autentikasi terpusat memastikan setiap pengguna diarahkan ke ruang
            kerja yang sesuai dengan peran dan aplikasinya.
          </p>
        </div>
        <p className="relative text-sm text-blue-100">
          JWT dilindungi dan divalidasi oleh API Integrator Gateway.
        </p>
      </section>

      <section className="flex min-h-screen items-center justify-center bg-slate-50 px-5 py-10 sm:px-8">
        <div className="w-full max-w-md">
          <div className="mb-8 lg:hidden">
            <BrandMark />
          </div>
          <Link
            to="/"
            className="mb-8 inline-flex items-center gap-2 rounded-lg text-sm font-bold text-blue-700 transition hover:text-blue-900 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-blue-700"
          >
            <span aria-hidden="true">←</span>
            Kembali ke beranda
          </Link>

          <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-xl shadow-slate-900/5 sm:p-9">
            <p className="text-sm font-bold uppercase tracking-[0.2em] text-blue-700">
              Authentication
            </p>
            <h2 className="mt-3 text-3xl font-black tracking-tight text-slate-950">
              Masuk ke API Integrator
            </h2>
            <p className="mt-3 text-sm leading-6 text-slate-600">
              Gunakan kredensial aplikasi yang telah terdaftar pada gateway.
            </p>

            <form className="mt-8 space-y-5" onSubmit={handleSubmit}>
              <div>
                <label
                  htmlFor="username"
                  className="mb-2 block text-sm font-bold text-slate-800"
                >
                  Username
                </label>
                <input
                  id="username"
                  name="username"
                  type="text"
                  autoComplete="username"
                  required
                  value={form.username}
                  onChange={updateField}
                  className="w-full rounded-xl border border-slate-300 bg-white px-4 py-3 text-base outline-none transition placeholder:text-slate-400 focus:border-blue-600 focus:ring-4 focus:ring-blue-100"
                  placeholder="Masukkan username"
                />
              </div>

              <div>
                <label
                  htmlFor="password"
                  className="mb-2 block text-sm font-bold text-slate-800"
                >
                  Password
                </label>
                <input
                  id="password"
                  name="password"
                  type="password"
                  autoComplete="current-password"
                  required
                  value={form.password}
                  onChange={updateField}
                  className="w-full rounded-xl border border-slate-300 bg-white px-4 py-3 text-base outline-none transition placeholder:text-slate-400 focus:border-blue-600 focus:ring-4 focus:ring-blue-100"
                  placeholder="Masukkan password"
                />
              </div>

              <div>
                <label
                  htmlFor="app_name"
                  className="mb-2 block text-sm font-bold text-slate-800"
                >
                  Aplikasi
                </label>
                <select
                  id="app_name"
                  name="app_name"
                  required
                  value={form.app_name}
                  onChange={updateField}
                  className="w-full rounded-xl border border-slate-300 bg-white px-4 py-3 text-base outline-none transition focus:border-blue-600 focus:ring-4 focus:ring-blue-100"
                >
                  <option value="">Pilih aplikasi</option>
                  {applications.map((application) => (
                    <option key={application.name} value={application.name}>
                      {application.name}
                    </option>
                  ))}
                </select>
              </div>

              {error ? (
                <p
                  role="alert"
                  className="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm font-semibold text-red-800"
                >
                  {error}
                </p>
              ) : null}

              <button
                type="submit"
                disabled={submitting || status === 'loading'}
                className="flex w-full items-center justify-center rounded-xl bg-blue-700 px-4 py-3.5 text-sm font-black text-white shadow-lg shadow-blue-700/20 transition hover:bg-blue-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-blue-700 disabled:cursor-not-allowed disabled:bg-slate-400 disabled:shadow-none"
              >
                {submitting ? 'Memproses...' : 'Masuk'}
              </button>
            </form>
          </div>

          <p className="mt-6 text-center text-xs leading-5 text-slate-500">
            Jangan gunakan kredensial development pada environment bersama.
          </p>
        </div>
      </section>
    </main>
  )
}

export default LoginPage
