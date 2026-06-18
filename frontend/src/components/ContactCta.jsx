import { repositoryUrl } from '../data/landingContent'

function ContactCta() {
  return (
    <section id="kontak" className="scroll-mt-24 bg-white py-20 sm:py-24">
      <div className="mx-auto max-w-7xl px-5 sm:px-8">
        <div className="relative overflow-hidden rounded-3xl bg-blue-700 px-6 py-12 text-white shadow-2xl shadow-blue-900/20 sm:px-12 sm:py-14">
          <div
            className="absolute -right-32 -top-32 h-80 w-80 rounded-full border-[64px] border-white/5"
            aria-hidden="true"
          />
          <div className="relative grid gap-10 lg:grid-cols-[1fr_auto] lg:items-end">
            <div className="max-w-2xl">
              <p className="text-xs font-bold uppercase tracking-[0.22em] text-blue-200">
                Mulai dari fondasi yang sama
              </p>
              <h2 className="mt-4 text-3xl font-black tracking-tight sm:text-4xl">
                Siap menghubungkan layanan Anda?
              </h2>
              <p className="mt-5 text-base leading-7 text-blue-100 sm:text-lg">
                Tinjau implementasi, kontrak API, dan roadmap pengembangan pada
                repositori resmi proyek.
              </p>
            </div>
            <div className="flex flex-col gap-3 sm:flex-row lg:flex-col">
              <a
                href={repositoryUrl}
                target="_blank"
                rel="noreferrer"
                className="inline-flex min-h-12 items-center justify-center rounded-xl bg-white px-6 py-3 text-sm font-extrabold text-blue-800 transition hover:bg-blue-50 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
              >
                Lihat repositori
              </a>
              <button
                type="button"
                disabled
                aria-disabled="true"
                aria-label="Login segera hadir"
                className="inline-flex min-h-12 cursor-not-allowed items-center justify-center gap-2 rounded-xl border border-white/25 bg-white/10 px-6 py-3 text-sm font-bold text-blue-100"
              >
                Login
                <span className="rounded-full bg-amber-300 px-2 py-0.5 text-[0.65rem] font-black uppercase tracking-wide text-amber-950">
                  Segera hadir
                </span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}

export default ContactCta
