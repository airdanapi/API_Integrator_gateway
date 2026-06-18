import Benefits from './components/Benefits'
import ContactCta from './components/ContactCta'
import Faq from './components/Faq'
import Footer from './components/Footer'
import Header from './components/Header'
import Hero from './components/Hero'
import IntegrationFlow from './components/IntegrationFlow'
import UseCases from './components/UseCases'

function App() {
  return (
    <div className="min-h-screen bg-white text-slate-900">
      <a
        href="#konten-utama"
        className="fixed left-4 top-4 z-[100] -translate-y-24 rounded-lg bg-slate-950 px-4 py-3 text-sm font-bold text-white transition focus:translate-y-0"
      >
        Lewati ke konten utama
      </a>
      <Header />
      <main id="konten-utama">
        <Hero />
        <Benefits />
        <IntegrationFlow />
        <UseCases />
        <Faq />
        <ContactCta />
      </main>
      <Footer />
    </div>
  )
}

export default App
