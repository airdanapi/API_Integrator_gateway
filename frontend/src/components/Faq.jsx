import SectionHeading from './SectionHeading'
import { faqs } from '../data/landingContent'

function Faq() {
  return (
    <section
      id="faq"
      className="scroll-mt-24 border-y border-slate-200 bg-slate-50 py-20 sm:py-24"
    >
      <div className="mx-auto max-w-4xl px-5 sm:px-8">
        <SectionHeading
          eyebrow="FAQ"
          title="Pertanyaan yang sering diajukan"
          description="Ringkasan mengenai posisi gateway, cakupan layanan, dan cara akses portal."
          align="center"
        />
        <div className="mt-12 space-y-3">
          {faqs.map((faq, index) => (
            <details
              key={faq.question}
              className="group rounded-2xl border border-slate-200 bg-white p-5 shadow-sm open:border-blue-200"
              open={index === 0}
            >
              <summary className="flex cursor-pointer list-none items-center justify-between gap-4 font-bold text-slate-900 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-blue-700">
                {faq.question}
                <span className="grid h-8 w-8 shrink-0 place-items-center rounded-full bg-slate-100 text-blue-700 transition group-open:rotate-45">
                  +
                </span>
              </summary>
              <p className="mt-4 max-w-3xl pr-10 leading-7 text-slate-600">
                {faq.answer}
              </p>
            </details>
          ))}
        </div>
      </div>
    </section>
  )
}

export default Faq
