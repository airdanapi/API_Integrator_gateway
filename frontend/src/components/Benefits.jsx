import SectionHeading from './SectionHeading'
import { benefits } from '../data/landingContent'

function BenefitIcon({ name }) {
  const paths = {
    shield:
      'M12 3 5.5 5.5v5.25c0 4.2 2.7 7.85 6.5 9.25 3.8-1.4 6.5-5.05 6.5-9.25V5.5L12 3Zm-3 9 2 2 4-4',
    route:
      'M6.5 6.5h4m3 0h4M6.5 17.5h4m3 0h4M8.5 6.5a2 2 0 1 1-4 0 2 2 0 0 1 4 0Zm11 0a2 2 0 1 1-4 0 2 2 0 0 1 4 0Zm-11 11a2 2 0 1 1-4 0 2 2 0 0 1 4 0Zm11 0a2 2 0 1 1-4 0 2 2 0 0 1 4 0ZM12 8.5v7',
    pulse:
      'M4 12h3l2-5 3.5 10 2.5-6 1.5 3H20M12 3.5a8.5 8.5 0 1 1-8.5 8.5',
  }

  return (
    <svg viewBox="0 0 24 24" fill="none" className="h-6 w-6" aria-hidden="true">
      <path
        d={paths[name]}
        stroke="currentColor"
        strokeWidth="1.7"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}

function Benefits() {
  return (
    <section id="manfaat" className="scroll-mt-24 bg-white py-20 sm:py-24">
      <div className="mx-auto max-w-7xl px-5 sm:px-8">
        <SectionHeading
          eyebrow="Fondasi yang dapat dipercaya"
          title="Integrasi yang aman, terukur, dan konsisten"
          description="Gateway memisahkan tanggung jawab integrasi dari logika bisnis, sehingga setiap layanan dapat berkembang tanpa kehilangan kontrol."
          align="center"
        />
        <div className="mt-12 grid gap-5 md:grid-cols-3">
          {benefits.map((benefit) => (
            <article
              key={benefit.title}
              className="group rounded-2xl border border-slate-200 bg-white p-6 shadow-sm transition hover:-translate-y-1 hover:border-blue-200 hover:shadow-xl hover:shadow-blue-900/5 sm:p-7"
            >
              <span className="grid h-12 w-12 place-items-center rounded-xl bg-blue-50 text-blue-700 transition group-hover:bg-blue-700 group-hover:text-white">
                <BenefitIcon name={benefit.icon} />
              </span>
              <h3 className="mt-6 text-xl font-extrabold text-slate-950">
                {benefit.title}
              </h3>
              <p className="mt-3 leading-7 text-slate-600">
                {benefit.description}
              </p>
            </article>
          ))}
        </div>
      </div>
    </section>
  )
}

export default Benefits
