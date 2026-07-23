package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Talen400/sp_b2b/internal/split"
	"github.com/Talen400/sp_b2b/internal/store"
)

func main() {
	s := store.NewStore()

	// 1. Cadastrar empresas
	fmt.Println("╔══════════════════════════════════════════════╗")
	fmt.Println("║  Split Payment Simulator — Demonstração B2B  ║")
	fmt.Println("╚══════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println(">>> Cadastrando empresas...")
	must(s.AddCompany("11.111.111/0001-11", "Fazenda Boa Vista"))
	must(s.AddCompany("22.222.222/0001-22", "Fábrica de Sucos SA"))
	must(s.AddCompany("33.333.333/0001-33", "Mercado Central"))
	fmt.Println("  ✅ Fazenda Boa Vista  (11.111.111/0001-11)")
	fmt.Println("  ✅ Fábrica de Sucos SA (22.222.222/0001-22)")
	fmt.Println("  ✅ Mercado Central     (33.333.333/0001-33)")
	fmt.Println()

	// 2. Venda 1: Fazenda → Fábrica
	fmt.Println("════════════════════════════════════════════════")
	fmt.Println(">>> Venda 1: Fazenda Boa Vista → Fábrica de Sucos SA")
	fmt.Println("    Valor bruto: R$ 1.000,00 | IBS: 12% | CBS: 3%")
	r1, err := s.ProcessTransaction("11.111.111/0001-11", "22.222.222/0001-22", 100000, 0.12, 0.03)
	must(err)
	printTransaction(r1)
	fmt.Println()

	// 3. Venda 2: Fábrica → Mercado
	fmt.Println("════════════════════════════════════════════════")
	fmt.Println(">>> Venda 2: Fábrica de Sucos SA → Mercado Central")
	fmt.Println("    Valor bruto: R$ 3.000,00 | IBS: 12% | CBS: 3%")
	r2, err := s.ProcessTransaction("22.222.222/0001-22", "33.333.333/0001-33", 300000, 0.12, 0.03)
	must(err)
	printTransaction(r2)
	fmt.Println()

	// 4. Resultado final
	fmt.Println("════════════════════════════════════════════════")
	fmt.Println(">>> Saldo de crédito final:")
	for _, cnpj := range []string{"11.111.111/0001-11", "22.222.222/0001-22", "33.333.333/0001-33"} {
		emp, _ := s.GetCompany(cnpj)
		fmt.Printf("  %-25s  R$ %8.2f\n", emp.Nome, centToReal(emp.SaldoCredito))
	}
	fmt.Println()

	fmt.Println("════════════════════════════════════════════════")
	fmt.Println(">>> Histórico completo de transações:")
	printHistory(s.ListTransactions())
	fmt.Println()

	fmt.Println("╔══════════════════════════════════════════════╗")
	fmt.Println("║  Demonstração concluída com sucesso!         ║")
	fmt.Println("╚══════════════════════════════════════════════╝")
}

func printTransaction(r store.ProcessTransactionResult) {
	t := r.Transaction
	imposto := t.ValorIBS + t.ValorCBS
	repassado := imposto - t.CreditoUsado

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "  ID:\t%s\n", t.ID)
	fmt.Fprintf(w, "  Vendedor:\t%s\n", t.VendedorCNPJ)
	fmt.Fprintf(w, "  Comprador:\t%s\n", t.CompradorCNPJ)
	fmt.Fprintf(w, "  Valor bruto:\tR$ %.2f\n", centToReal(t.ValorBruto))
	fmt.Fprintf(w, "  Líquido do vendedor:\tR$ %.2f\n", centToReal(t.Liquido))
	fmt.Fprintf(w, "  IBS:\tR$ %.2f\n", centToReal(t.ValorIBS))
	fmt.Fprintf(w, "  CBS:\tR$ %.2f\n", centToReal(t.ValorCBS))
	fmt.Fprintf(w, "  ─────────────────────────────\n")
	fmt.Fprintf(w, "  Crédito usado (abate):\tR$ %.2f\n", centToReal(t.CreditoUsado))
	fmt.Fprintf(w, "  Repassado ao fisco:\tR$ %.2f\n", centToReal(repassado))
	fmt.Fprintf(w, "  Crédito gerado p/ comprador:\tR$ %.2f\n", centToReal(r.CreditoGerado))
	w.Flush()
}

func printHistory(txns []split.Transaction) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tVendedor\tComprador\tBruto\tLíquido\tIBS\tCBS\tCrédito Usado\tCrédito Gerado")
	fmt.Fprintln(w, "---\t--------\t---------\t-----\t-------\t---\t---\t-------------\t-------------")

	for _, t := range txns {
		imposto := t.ValorIBS + t.ValorCBS
		fmt.Fprintf(w, "%s\t%s\t%s\tR$ %.2f\tR$ %.2f\tR$ %.2f\tR$ %.2f\tR$ %.2f\tR$ %.2f\n",
			t.ID, t.VendedorCNPJ, t.CompradorCNPJ,
			centToReal(t.ValorBruto), centToReal(t.Liquido),
			centToReal(t.ValorIBS), centToReal(t.ValorCBS),
			centToReal(t.CreditoUsado), centToReal(imposto),
		)
	}
	w.Flush()
}

func centToReal(c int64) float64 {
	return float64(c) / 100.0
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro fatal: %v\n", err)
		os.Exit(1)
	}
}
