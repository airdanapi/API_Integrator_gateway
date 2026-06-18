// Command seed adalah tool CLI untuk menjalankan database seeder Sprint 5 secara mandiri,
// tanpa harus menjalankan server utama.
//
// Penggunaan:
//
//	go run ./cmd/seed [flags]
//
// Flags:
//
//	-verbose         Tampilkan detail log setiap record (default: false)
//	-logs  <N>       Jumlah request_logs yang di-seed (default: 50)
//	-chats <N>       Jumlah chat_messages yang di-seed (default: 20)
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/config"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/database"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/database/seed"
)

func main() {
	verbose := flag.Bool("verbose", false, "tampilkan detail log setiap record")
	logCount := flag.Int("logs", 50, "jumlah request_logs yang di-seed")
	chatCount := flag.Int("chats", 20, "jumlah chat_messages yang di-seed")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	log.Println("[seed] menghubungkan ke database...")
	db, err := database.Open(ctx, cfg)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	log.Println("[seed] menjalankan migrasi...")
	if err := database.Migrate(ctx, db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	opts := seed.Options{
		Verbose:          *verbose,
		RequestLogCount:  *logCount,
		ChatMessageCount: *chatCount,
	}

	log.Printf("[seed] mulai seeding Sprint 5 (logs=%d, chats=%d, verbose=%v)...",
		opts.RequestLogCount, opts.ChatMessageCount, opts.Verbose)

	if err := seed.Run(ctx, db, opts); err != nil {
		log.Fatalf("seeding gagal: %v", err)
	}

	log.Println("[seed] seeding Sprint 5 selesai dengan sukses.")
}
