package server

import "github.com/gofiber/fiber/v3"

type landingBenefit struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type landingServiceOverview struct {
	Name        string           `json:"name"`
	Tagline     string           `json:"tagline"`
	Description string           `json:"description"`
	Benefits    []landingBenefit `json:"benefits"`
}

type landingApplicationRole struct {
	Application string `json:"application"`
	Role        string `json:"role"`
	Interaction string `json:"interaction"`
}

type landingIntegrationStep struct {
	Step        int    `json:"step"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type landingContactInfo struct {
	RepositoryURL string `json:"repository_url"`
	LoginStatus   string `json:"login_status"`
}

type landingData struct {
	ServiceOverview  landingServiceOverview   `json:"service_overview"`
	ApplicationRoles []landingApplicationRole `json:"application_roles"`
	IntegrationFlow  []landingIntegrationStep `json:"integration_flow"`
	ContactInfo      landingContactInfo       `json:"contact_info"`
}

func landingHandler(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "success",
		"data": landingData{
			ServiceOverview: landingServiceOverview{
				Name:    "API Integrator Gateway",
				Tagline: "Satu pintu aman untuk setiap komunikasi antar aplikasi.",
				Description: "Gateway untuk routing, validasi, logging, dan standardisasi " +
					"request dalam ekosistem UMKM.",
				Benefits: []landingBenefit{
					{
						Title:       "Keamanan terpusat",
						Description: "Validasi identitas dan akses dilakukan sebelum request diteruskan.",
					},
					{
						Title:       "Routing konsisten",
						Description: "Setiap aplikasi berkomunikasi melalui kontrak yang seragam.",
					},
					{
						Title:       "Operasional terpantau",
						Description: "Request dan response dicatat untuk audit dan diagnosis.",
					},
				},
			},
			ApplicationRoles: []landingApplicationRole{
				{
					Application: "SmartBank",
					Role:        "Core keuangan",
					Interaction: "Memproses pembayaran, saldo, fee, dan ledger.",
				},
				{
					Application: "Marketplace",
					Role:        "Perdagangan digital",
					Interaction: "Mengirim checkout dan payment request melalui gateway.",
				},
				{
					Application: "POS",
					Role:        "Transaksi offline",
					Interaction: "Mengirim invoice kasir dan menerima status pembayaran.",
				},
				{
					Application: "SupplierHub",
					Role:        "Rantai pasok",
					Interaction: "Mengirim order bahan dan pembayaran supplier.",
				},
				{
					Application: "LogistiKita",
					Role:        "Distribusi",
					Interaction: "Mengelola request pengiriman, ongkir, dan status.",
				},
				{
					Application: "UMKM Insight",
					Role:        "Analitik read-only",
					Interaction: "Membaca data untuk insight tanpa mengubah transaksi.",
				},
				{
					Application: "API Gateway",
					Role:        "Orkestrasi",
					Interaction: "Menangani routing, validasi, logging, dan standardisasi.",
				},
			},
			IntegrationFlow: []landingIntegrationStep{
				{
					Step:        1,
					Title:       "Aplikasi mengirim request",
					Description: "Layanan mengirim payload melalui satu endpoint gateway.",
				},
				{
					Step:        2,
					Title:       "Gateway memvalidasi",
					Description: "Identitas, hak akses, dan payload diperiksa.",
				},
				{
					Step:        3,
					Title:       "Request diarahkan",
					Description: "Gateway menentukan dan memanggil layanan tujuan.",
				},
				{
					Step:        4,
					Title:       "Response distandarkan",
					Description: "Hasil dicatat dan dikembalikan dalam format JSON konsisten.",
				},
			},
			ContactInfo: landingContactInfo{
				RepositoryURL: "https://github.com/airdanapi/API_Integrator_gateway",
				LoginStatus:   "coming_soon",
			},
		},
	})
}
