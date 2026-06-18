import { applications } from '../data/landingContent'

function HeroDiagram() {
  const satelliteApps = applications.filter(
    (application) => application.name !== 'API Gateway',
  )

  return (
    <div className="relative mx-auto w-full max-w-xl" aria-label="Diagram ekosistem aplikasi">
      <div className="absolute inset-12 rounded-full border border-dashed border-blue-200 sm:inset-16" />
      <div className="relative grid grid-cols-2 gap-3 rounded-[2rem] border border-slate-200 bg-white p-4 shadow-2xl shadow-blue-900/10 sm:grid-cols-3 sm:p-6">
        {satelliteApps.map((application) => (
          <div
            key={application.name}
            className="flex min-h-24 flex-col justify-between rounded-2xl border border-slate-200 bg-slate-50 p-3"
          >
            <span className="h-2 w-2 rounded-full bg-blue-500" />
            <span className="mt-5 text-xs font-bold text-slate-700 sm:text-sm">
              {application.name}
            </span>
          </div>
        ))}
        <div className="col-span-2 flex min-h-24 items-center gap-3 rounded-2xl bg-blue-700 p-4 text-white sm:col-span-3">
          <span className="grid h-11 w-11 shrink-0 place-items-center rounded-xl bg-white/15">
            <svg
              viewBox="0 0 24 24"
              fill="none"
              className="h-6 w-6"
              aria-hidden="true"
            >
              <path
                d="M5 8h14M5 16h14M8 5v6M16 13v6"
                stroke="currentColor"
                strokeWidth="1.8"
                strokeLinecap="round"
              />
            </svg>
          </span>
          <span>
            <span className="block text-xs font-semibold uppercase tracking-[0.18em] text-blue-200">
              Single entry point
            </span>
            <span className="mt-1 block font-extrabold">API Gateway</span>
          </span>
          <span className="ml-auto rounded-full bg-emerald-400/20 px-3 py-1 text-xs font-bold text-emerald-100">
            Terhubung
          </span>
        </div>
      </div>
    </div>
  )
}

function Hero() {
  return (
    <section
      id="beranda"
      className="relative overflow-hidden border-b border-slate-200 bg-slate-50"
    >
      <div
        className="pointer-events-none absolute inset-0 opacity-60"
        aria-hidden="true"
        style={{
          backgroundImage:
            'radial-gradient(circle at 20% 10%, rgba(37,99,235,.12), transparent 28%), radial-gradient(circle at 90% 70%, rgba(14,165,233,.09), transparent 28%)',
        }}
      />
      <div className="relative mx-auto grid max-w-7xl items-center gap-14 px-5 py-20 sm:px-8 sm:py-24 lg:grid-cols-[1.05fr_.95fr] lg:py-28">
        <div>
          <div className="mb-7 inline-flex items-center gap-2 rounded-full border border-blue-200 bg-white px-4 py-2 text-xs font-bold text-blue-800 shadow-sm">
            <span className="h-2 w-2 rounded-full bg-emerald-500" />
            Infrastruktur integrasi untuk ekosistem UMKM
          </div>
          <h1 className="text-4xl font-black tracking-[-0.04em] text-slate-950 sm:text-5xl lg:text-6xl">
            API Integrator Gateway
          </h1>
          <p className="mt-5 max-w-2xl text-xl font-semibold leading-8 text-blue-800 sm:text-2xl">
            Satu pintu aman untuk setiap komunikasi antar aplikasi.
          </p>
          <p className="mt-5 max-w-2xl text-base leading-7 text-slate-600 sm:text-lg">
            Menyatukan routing, validasi, logging, dan standardisasi request agar
            SmartBank dan seluruh layanan UMKM dapat terhubung secara konsisten.
          </p>
          <div className="mt-9 flex flex-col gap-3 sm:flex-row">
            <a
              href="#alur-integrasi"
              className="inline-flex min-h-12 items-center justify-center gap-2 rounded-xl bg-blue-700 px-6 py-3 text-sm font-bold text-white shadow-lg shadow-blue-700/20 transition hover:-translate-y-0.5 hover:bg-blue-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-700"
            >
              Pelajari alur integrasi
              <span aria-hidden="true">→</span>
            </a>
            <a
              href="#manfaat"
              className="inline-flex min-h-12 items-center justify-center rounded-xl border border-slate-300 bg-white px-6 py-3 text-sm font-bold text-slate-700 transition hover:border-blue-300 hover:text-blue-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-700"
            >
              Lihat manfaat
            </a>
          </div>
          <dl className="mt-10 grid max-w-xl grid-cols-3 divide-x divide-slate-200">
            <div className="pr-4">
              <dt className="text-2xl font-black text-slate-950">7</dt>
              <dd className="mt-1 text-xs leading-5 text-slate-500">Layanan ekosistem</dd>
            </div>
            <div className="px-4">
              <dt className="text-2xl font-black text-slate-950">1</dt>
              <dd className="mt-1 text-xs leading-5 text-slate-500">Pintu integrasi</dd>
            </div>
            <div className="pl-4">
              <dt className="text-2xl font-black text-slate-950">JSON</dt>
              <dd className="mt-1 text-xs leading-5 text-slate-500">Kontrak konsisten</dd>
            </div>
          </dl>
        </div>
        <HeroDiagram />
      </div>
    </section>
  )
}

export default Hero
