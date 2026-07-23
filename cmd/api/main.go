package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	httpapi "github.com/Talen400/sp_b2b/internal/handler/http"
	"github.com/Talen400/sp_b2b/internal/repository/sqlite"
	"github.com/Talen400/sp_b2b/internal/seed"
)

func main() {
	migrateOnly := flag.Bool("migrate", false, "Apenas rodar as migrations e sair")
	seedOnly := flag.Bool("seed", false, "Popular banco com cenário de demonstração")
	port := flag.String("port", "8080", "Porta do servidor HTTP")
	dbPath := flag.String("db", "data/split.db", "Caminho do arquivo SQLite")
	flag.Parse()

	db, err := sqlite.Open(*dbPath)
	if err != nil {
		log.Fatalf("Erro ao abrir banco: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		log.Fatalf("Erro ao aplicar migrations: %v", err)
	}

	if *migrateOnly {
		fmt.Println("Migrations aplicadas com sucesso.")
		return
	}

	if *seedOnly {
		if err := seed.Run(
			sqlite.NewCompanyRepository(db),
			sqlite.NewTransactionRepository(db),
		); err != nil {
			log.Fatalf("Erro ao popular banco: %v", err)
		}
		fmt.Println("Banco populado com cenário de demonstração.")
		return
	}

	if _, err := os.Stat(*dbPath); err != nil {
		log.Printf("Banco não encontrado em %s, criando...", *dbPath)
	}

	handler := httpapi.NewHandler(
		sqlite.NewCompanyRepository(db),
		sqlite.NewTransactionRepository(db),
	)
	mux := handler.Router()

	addr := fmt.Sprintf(":%s", *port)
	log.Printf("Servidor rodando em http://localhost%s", addr)
	log.Printf("Endpoints:")
	log.Printf("  GET  /api/v1/healthz")
	log.Printf("  POST /api/v1/companies")
	log.Printf("  GET  /api/v1/companies/{cnpj}")
	log.Printf("  GET  /api/v1/companies")
	log.Printf("  POST /api/v1/transactions")
	log.Printf("  GET  /api/v1/transactions/{id}")
	log.Printf("  GET  /api/v1/transactions?cnpj={cnpj}")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
