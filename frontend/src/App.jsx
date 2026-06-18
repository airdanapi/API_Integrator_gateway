import { API_BASE_URL } from './config'

function App() {
  return (
    <main className="min-h-screen bg-slate-950 px-6 py-16 text-slate-100">
      <section className="mx-auto flex max-w-3xl flex-col gap-8 rounded-3xl border border-slate-800 bg-slate-900 p-8 shadow-2xl shadow-cyan-950/30 sm:p-12">
        <div className="flex items-center gap-3">
          <span className="h-3 w-3 rounded-full bg-emerald-400 shadow-lg shadow-emerald-400/50" />
          <span className="text-sm font-semibold uppercase tracking-[0.2em] text-emerald-300">
            Development ready
          </span>
        </div>

        <div className="space-y-4">
          <p className="text-sm font-medium uppercase tracking-[0.25em] text-cyan-300">
            Ekosistem UMKM
          </p>
          <h1 className="text-4xl font-bold tracking-tight sm:text-5xl">
            API Integrator Gateway
          </h1>
          <p className="max-w-2xl text-lg leading-8 text-slate-300">
            Sprint 1 infrastructure is ready.
          </p>
        </div>

        <dl className="grid gap-4 rounded-2xl bg-slate-950/70 p-5 text-sm sm:grid-cols-2">
          <div>
            <dt className="text-slate-500">Frontend</dt>
            <dd className="mt-1 font-mono text-slate-200">localhost:5173</dd>
          </div>
          <div>
            <dt className="text-slate-500">Backend API</dt>
            <dd className="mt-1 break-all font-mono text-slate-200">
              {API_BASE_URL}
            </dd>
          </div>
        </dl>
      </section>
    </main>
  )
}

export default App
