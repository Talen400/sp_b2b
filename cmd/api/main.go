// main é o entrypoint da API.
//
// Flags:
//
//	-port: porta do servidor (default 8080)
//	-db: caminho do arquivo SQLite (default data/split.db)
//	-migrate: só roda as migrations e sai
//	-seed: popula banco com cenário de demonstração e sai
//	-pp-url: URL da Plataforma Pública (mock Prism). Ex: http://localhost:4010
//	-pp-tenant: Tenant-Id para enviar à PP (default PSP-SIMULADOR-001)
//
// Fluxo normal: migrate → servidor HTTP.
// Se -pp-url for fornecido, notifica a PP a cada POST /transactions.
//
// Uso:
//
//	go run ./cmd/api
//	go run ./cmd/api -migrate
//	go run ./cmd/api -seed
//	go run ./cmd/api -pp-url http://localhost:4010
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Talen400/sp_b2b/internal/client/pp"
	"github.com/Talen400/sp_b2b/internal/domain"
	httpapi "github.com/Talen400/sp_b2b/internal/handler/http"
	"github.com/Talen400/sp_b2b/internal/repository/sqlite"
	"github.com/Talen400/sp_b2b/internal/seed"
)

func main() {
	migrateOnly := flag.Bool("migrate", false, "Apenas rodar as migrations e sair")
	seedOnly := flag.Bool("seed", false, "Popular banco com cenário de demonstração")
	port := flag.String("port", "8080", "Porta do servidor HTTP")
	dbPath := flag.String("db", "data/split.db", "Caminho do arquivo SQLite")
	ppURL := flag.String("pp-url", "", "URL da Plataforma Pública (mock Prism). Ex: http://localhost:4010")
	ppTenant := flag.String("pp-tenant", "PSP-SIMULADOR-001", "Tenant-Id para a PP")
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

	handler := newHandlerWithPP(
		sqlite.NewCompanyRepository(db),
		sqlite.NewTransactionRepository(db),
		*ppURL, *ppTenant,
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
	if *ppURL != "" {
		log.Printf("PP notificação ativada -> %s", *ppURL)
	}

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}

// newHandlerWithPP cria um Handler HTTP, opcionalmente com notificação PP.
//
// Se ppURL não for vazio:
//   - Cria um cliente pp.Client apontando para ppURL.
//   - Cria um PPNotifyFunc que, a cada POST /transactions, monta um
//     InformeTransacaoIniciadaRequest e envia para o mock.
//   - O arranjo é fixo em "boleto" (simplificação didática).
//   - Timeout de 5s por notificação; se falhar, o response da API ainda
//     é 201 com pp_notification.status = "failed".
//
// Se ppURL for vazio: retorna um Handler sem notificação.
func newHandlerWithPP(cr *sqlite.CompanyRepository, tr *sqlite.TransactionRepository, ppURL, ppTenant string) *httpapi.Handler {
	if ppURL == "" {
		return httpapi.NewHandler(cr, tr)
	}

	ppClient := pp.New(ppURL, ppTenant)

	notifyFunc := func(ctx context.Context, txn domain.Transaction) *httpapi.PPNotification {
		req := &pp.InformeTransacaoIniciadaRequest{
			NSUId:       txn.ID,
			CNPJRec:     txn.VendedorCNPJ,
			VLInf:       pp.Decimal18_2(txn.ValorBruto),
			VLIbs:       pp.Decimal18_2(txn.ValorIBS),
			VLCbs:       pp.Decimal18_2(txn.ValorCBS),
			DtHrCriacao: txn.Timestamp.Format(time.RFC3339),
		}

		ppCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		status, ppErr, err := ppClient.SendInformeTransacaoIniciada(ppCtx, "boleto", req)
		if err != nil {
			return &httpapi.PPNotification{
				Status:      "failed",
				Arrangement: "boleto",
				Error:       err.Error(),
			}
		}
		if ppErr != nil {
			return &httpapi.PPNotification{
				Status:      "failed",
				Arrangement: "boleto",
				Error:       ppErr.Error(),
			}
		}
		_ = status
		return &httpapi.PPNotification{
			Status:      "sent",
			Arrangement: "boleto",
		}
	}

	return httpapi.NewHandlerWithPP(cr, tr, notifyFunc)
}
