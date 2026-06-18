import SectionHeading from './SectionHeading'
import { applications, integrationSteps } from '../data/landingContent'

const toneClasses = {
  blue: 'bg-blue-50 text-blue-700 ring-blue-100',
  violet: 'bg-violet-50 text-violet-700 ring-violet-100',
  amber: 'bg-amber-50 text-amber-700 ring-amber-100',
  emerald: 'bg-emerald-50 text-emerald-700 ring-emerald-100',
  rose: 'bg-rose-50 text-rose-700 ring-rose-100',
  cyan: 'bg-cyan-50 text-cyan-700 ring-cyan-100',
  indigo: 'bg-indigo-50 text-indigo-700 ring-indigo-100',
}

function IntegrationFlow() {
  return (
    <section
      id="alur-integrasi"
      className="scroll-mt-24 border-y border-slate-200 bg-slate-50 py-20 sm:py-24"
    >
      <div className="mx-auto max-w-7xl px-5 sm:px-8">
        <SectionHeading
          eyebrow="Cara kerja"
          title="Satu alur untuk seluruh ekosistem"
          description="Setiap request melewati tahapan yang dapat dipahami, diuji, dan diaudit tanpa mengubah tanggung jawab bisnis masing-masing aplikasi."
        />
        <ol className="mt-12 grid gap-4 lg:grid-cols-4">
          {integrationSteps.map((item, index) => (
            <li
              key={item.step}
              className="relative rounded-2xl border border-slate-200 bg-white p-6 shadow-sm"
            >
              {index < integrationSteps.length - 1 ? (
                <span
                  className="absolute left-[calc(100%+1px)] top-10 hidden h-px w-4 bg-blue-300 lg:block"
                  aria-hidden="true"
                />
              ) : null}
              <span className="text-sm font-black text-blue-700">{item.step}</span>
              <h3 className="mt-5 text-lg font-extrabold text-slate-950">
                {item.title}
              </h3>
              <p className="mt-3 text-sm leading-6 text-slate-600">
                {item.description}
              </p>
            </li>
          ))}
        </ol>

        <div className="mt-16 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm sm:p-8">
          <div className="flex flex-col justify-between gap-3 sm:flex-row sm:items-end">
            <div>
              <p className="text-xs font-bold uppercase tracking-[0.2em] text-blue-700">
                Node ekosistem
              </p>
              <h3 className="mt-2 text-2xl font-extrabold text-slate-950">
                Setiap layanan tetap fokus pada perannya
              </h3>
            </div>
            <p className="max-w-md text-sm leading-6 text-slate-500">
              Gateway menghubungkan layanan tanpa mengambil alih proses bisnis
              internal mereka.
            </p>
          </div>
          <div className="mt-8 grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
            {applications.map((application) => (
              <article
                key={application.name}
                className="rounded-2xl border border-slate-200 p-4"
              >
                <span
                  className={`inline-flex rounded-full px-2.5 py-1 text-[0.68rem] font-bold ring-1 ring-inset ${toneClasses[application.tone]}`}
                >
                  {application.category}
                </span>
                <h4 className="mt-4 font-extrabold text-slate-900">
                  {application.name}
                </h4>
                <p className="mt-2 text-sm leading-6 text-slate-600">
                  {application.description}
                </p>
              </article>
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}

export default IntegrationFlow
