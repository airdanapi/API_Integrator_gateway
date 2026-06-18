import BrandMark from './BrandMark'
import { navigationItems, repositoryUrl } from '../data/landingContent'

function Footer() {
  return (
    <footer className="border-t border-slate-200 bg-slate-50">
      <div className="mx-auto grid max-w-7xl gap-10 px-5 py-12 sm:px-8 md:grid-cols-[1fr_auto] md:items-start">
        <div>
          <BrandMark compact />
          <p className="mt-5 max-w-md text-sm leading-6 text-slate-500">
            Jalur tunggal yang menjaga keamanan dan konsistensi komunikasi antar
            aplikasi dalam ekosistem ekonomi UMKM.
          </p>
        </div>
        <div className="grid grid-cols-2 gap-8 sm:grid-cols-[auto_auto]">
          <div>
            <p className="text-xs font-bold uppercase tracking-[0.18em] text-slate-400">
              Navigasi
            </p>
            <ul className="mt-4 space-y-3">
              {navigationItems.map((item) => (
                <li key={item.href}>
                  <a
                    href={item.href}
                    className="text-sm font-semibold text-slate-600 hover:text-blue-700"
                  >
                    {item.label}
                  </a>
                </li>
              ))}
            </ul>
          </div>
          <div>
            <p className="text-xs font-bold uppercase tracking-[0.18em] text-slate-400">
              Proyek
            </p>
            <a
              href={repositoryUrl}
              target="_blank"
              rel="noreferrer"
              className="mt-4 inline-block text-sm font-semibold text-blue-700 hover:text-blue-900"
            >
              Lihat repositori
            </a>
            <p className="mt-3 text-xs text-slate-500">Login: Segera hadir</p>
          </div>
        </div>
      </div>
      <div className="border-t border-slate-200">
        <div className="mx-auto flex max-w-7xl flex-col gap-2 px-5 py-5 text-xs text-slate-500 sm:flex-row sm:items-center sm:justify-between sm:px-8">
          <p>© 2026 API Integrator Gateway.</p>
          <p>Dibangun untuk ekosistem UMKM.</p>
        </div>
      </div>
    </footer>
  )
}

export default Footer
