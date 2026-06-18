export const repositoryUrl =
  'https://github.com/airdanapi/API_Integrator_gateway'

export const navigationItems = [
  { label: 'Manfaat', href: '#manfaat' },
  { label: 'Alur integrasi', href: '#alur-integrasi' },
  { label: 'Use case', href: '#use-case' },
  { label: 'FAQ', href: '#faq' },
]

export const benefits = [
  {
    icon: 'shield',
    title: 'Keamanan terpusat',
    description:
      'Validasi token dan kebijakan akses diterapkan di satu pintu sebelum request mencapai layanan tujuan.',
  },
  {
    icon: 'route',
    title: 'Routing konsisten',
    description:
      'Setiap aplikasi berkomunikasi melalui kontrak yang seragam, mudah ditelusuri, dan mudah dikembangkan.',
  },
  {
    icon: 'pulse',
    title: 'Operasional terpantau',
    description:
      'Request dan response dicatat untuk mendukung audit, diagnosis kendala, dan pemantauan layanan.',
  },
]

export const integrationSteps = [
  {
    step: '01',
    title: 'Aplikasi mengirim request',
    description:
      'Marketplace, POS, SupplierHub, atau LogistiKita mengirim payload melalui satu endpoint gateway.',
  },
  {
    step: '02',
    title: 'Gateway memvalidasi',
    description:
      'Identitas, hak akses, dan bentuk payload diperiksa sebelum request dapat diteruskan.',
  },
  {
    step: '03',
    title: 'Request diarahkan',
    description:
      'Gateway menentukan layanan tujuan tanpa mengambil alih logika bisnis milik aplikasi tersebut.',
  },
  {
    step: '04',
    title: 'Response distandarkan',
    description:
      'Hasil dari layanan tujuan dicatat dan dikembalikan dalam format JSON yang konsisten.',
  },
]

export const applications = [
  {
    name: 'SmartBank',
    category: 'Core keuangan',
    description:
      'Memproses pembayaran, saldo, fee, dan ledger sebagai sumber kebenaran transaksi.',
    tone: 'blue',
  },
  {
    name: 'Marketplace',
    category: 'Perdagangan digital',
    description:
      'Mengelola katalog dan checkout, lalu meneruskan payment request melalui gateway.',
    tone: 'violet',
  },
  {
    name: 'POS',
    category: 'Transaksi offline',
    description:
      'Membuat tagihan toko fisik dan mengirim permintaan pembayaran secara terintegrasi.',
    tone: 'amber',
  },
  {
    name: 'SupplierHub',
    category: 'Rantai pasok',
    description:
      'Mengelola order bahan baku dan pembayaran supplier tanpa mengubah saldo langsung.',
    tone: 'emerald',
  },
  {
    name: 'LogistiKita',
    category: 'Distribusi',
    description:
      'Mengelola pengiriman, ongkir, dan status distribusi untuk aplikasi dalam ekosistem.',
    tone: 'rose',
  },
  {
    name: 'UMKM Insight',
    category: 'Analitik read-only',
    description:
      'Menyajikan insight dari data ledger tanpa mengubah data transaksi operasional.',
    tone: 'cyan',
  },
  {
    name: 'API Gateway',
    category: 'Orkestrasi',
    description:
      'Menjadi jalur tunggal untuk routing, validasi, logging, dan standardisasi komunikasi.',
    tone: 'indigo',
  },
]

export const useCases = [
  {
    number: '01',
    title: 'Checkout marketplace',
    description:
      'Payment request dari Marketplace divalidasi gateway sebelum diteruskan ke SmartBank.',
  },
  {
    number: '02',
    title: 'Pembayaran di kasir',
    description:
      'POS mengirim invoice melalui jalur yang sama dan menerima status transaksi yang konsisten.',
  },
  {
    number: '03',
    title: 'Order bahan baku',
    description:
      'SupplierHub mengoordinasikan order, pembayaran, dan pengiriman tanpa akses saldo langsung.',
  },
]

export const faqs = [
  {
    question: 'Apakah API Integrator memproses transaksi keuangan?',
    answer:
      'Tidak. Gateway mengatur validasi, routing, dan pencatatan. Seluruh transaksi keuangan tetap diproses oleh SmartBank.',
  },
  {
    question: 'Aplikasi apa saja yang terhubung?',
    answer:
      'Ekosistem mencakup SmartBank, Marketplace, POS, SupplierHub, LogistiKita, UMKM Insight, dan API Gateway sebagai penghubungnya.',
  },
  {
    question: 'Apakah landing page membutuhkan login?',
    answer:
      'Tidak. Landing page bersifat publik. Autentikasi diperlukan pada fitur operasional dan dashboard yang akan dibangun pada sprint berikutnya.',
  },
  {
    question: 'Bagaimana gateway menjaga konsistensi integrasi?',
    answer:
      'Gateway menggunakan kontrak endpoint, validasi payload, format response JSON, dan logging yang seragam untuk setiap komunikasi.',
  },
]
