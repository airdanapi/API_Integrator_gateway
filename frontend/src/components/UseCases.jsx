import SectionHeading from './SectionHeading'
import { useCases } from '../data/landingContent'

function UseCases() {
  return (
    <section id="use-case" className="scroll-mt-24 bg-white py-20 sm:py-24">
      <div className="mx-auto max-w-7xl px-5 sm:px-8">
        <SectionHeading
          eyebrow="Use case"
          title="Dibangun untuk kebutuhan nyata UMKM"
          description="Pola integrasi yang sama dapat melayani kanal digital, transaksi toko fisik, dan rantai pasok."
        />
        <div className="mt-12 grid gap-5 lg:grid-cols-3">
          {useCases.map((useCase) => (
            <article
              key={useCase.number}
              className="relative overflow-hidden rounded-2xl border border-slate-200 bg-slate-950 p-7 text-white"
            >
              <span
                className="absolute -right-4 -top-10 text-[9rem] font-black leading-none text-white/[0.04]"
                aria-hidden="true"
              >
                {useCase.number}
              </span>
              <span className="relative text-sm font-black text-blue-300">
                {useCase.number}
              </span>
              <h3 className="relative mt-14 text-2xl font-extrabold">
                {useCase.title}
              </h3>
              <p className="relative mt-4 leading-7 text-slate-300">
                {useCase.description}
              </p>
            </article>
          ))}
        </div>
      </div>
    </section>
  )
}

export default UseCases
