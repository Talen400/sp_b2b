package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Talen400/sp_b2b/internal/store"
)

func main() {
	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)
	flag.CommandLine.Usage = printUsage

	if len(os.Args) > 1 && os.Args[1] != "repl" {
		runCommand(strings.Join(os.Args[1:], " "))
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	if !isInteractive() {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			runCommand(line)
		}
		return
	}

	fmt.Println("Split Payment Simulator — digite um comando ou 'help'")
	fmt.Print("> ")
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			fmt.Print("> ")
			continue
		}
		if line == "exit" || line == "quit" {
			break
		}
		runCommand(line)
		fmt.Print("> ")
	}
}

var st = store.NewStore()

func runCommand(input string) {
	args := parseArgs(input)
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "add-company":
		cmdAddCompany(args[1:])
	case "sell":
		cmdSell(args[1:])
	case "balance":
		cmdBalance(args[1:])
	case "history":
		cmdHistory()
	case "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "comando desconhecido: %s\n", args[0])
	}
}

func cmdAddCompany(args []string) {
	fs := flag.NewFlagSet("add-company", flag.ContinueOnError)
	cnpj := fs.String("cnpj", "", "CNPJ da empresa")
	nome := fs.String("nome", "", "Nome da empresa")
	if err := fs.Parse(args); err != nil {
		return
	}

	if *cnpj == "" || *nome == "" {
		fmt.Fprintln(os.Stderr, "uso: add-company --cnpj=CNPJ --nome=NOME")
		return
	}

	if err := st.AddCompany(*cnpj, *nome); err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		return
	}

	fmt.Printf("Empresa cadastrada: %s (%s)\n", *nome, *cnpj)
}

func cmdSell(args []string) {
	fs := flag.NewFlagSet("sell", flag.ContinueOnError)
	vendedor := fs.String("vendedor", "", "CNPJ do vendedor")
	comprador := fs.String("comprador", "", "CNPJ do comprador")
	valor := fs.Int64("valor", 0, "Valor bruto em centavos")
	aliquotaIBS := fs.Float64("aliquota-ibs", 0.12, "Alíquota IBS")
	aliquotaCBS := fs.Float64("aliquota-cbs", 0.03, "Alíquota CBS")
	if err := fs.Parse(args); err != nil {
		return
	}

	if *vendedor == "" || *comprador == "" || *valor <= 0 {
		fmt.Fprintln(os.Stderr, "uso: sell --vendedor=CNPJ --comprador=CNPJ --valor=CENTAVOS")
		return
	}

	result, err := st.ProcessTransaction(*vendedor, *comprador, *valor, *aliquotaIBS, *aliquotaCBS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		return
	}

	t := result.Transaction
	impostoTotal := t.ValorIBS + t.ValorCBS
	repassado := impostoTotal - t.CreditoUsado

	fmt.Println("--- Transação concluída ---")
	fmt.Printf("  ID:               %s\n", t.ID)
	fmt.Printf("  Vendedor:         %s\n", t.VendedorCNPJ)
	fmt.Printf("  Comprador:        %s\n", t.CompradorCNPJ)
	fmt.Printf("  Valor bruto:      R$ %.2f\n", centsToReal(t.ValorBruto))
	fmt.Printf("  Líquido vend.:    R$ %.2f\n", centsToReal(t.Liquido))
	fmt.Printf("  IBS (%.0f%%):      R$ %.2f\n", t.AliquotaIBS*100, centsToReal(t.ValorIBS))
	fmt.Printf("  CBS (%.0f%%):      R$ %.2f\n", t.AliquotaCBS*100, centsToReal(t.ValorCBS))
	fmt.Printf("  Crédito usado:    R$ %.2f\n", centsToReal(t.CreditoUsado))
	fmt.Printf("  Repassado fisco:  R$ %.2f\n", centsToReal(repassado))
	fmt.Printf("  Crédito gerado:   R$ %.2f\n", centsToReal(result.CreditoGerado))
}

func cmdBalance(args []string) {
	fs := flag.NewFlagSet("balance", flag.ContinueOnError)
	cnpj := fs.String("cnpj", "", "CNPJ da empresa")
	if err := fs.Parse(args); err != nil {
		return
	}

	if *cnpj == "" {
		fmt.Fprintln(os.Stderr, "uso: balance --cnpj=CNPJ")
		return
	}

	emp, err := st.GetCompany(*cnpj)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		return
	}

	fmt.Printf("Empresa: %s (%s)\n", emp.Nome, emp.CNPJ)
	fmt.Printf("Saldo de crédito: R$ %.2f\n", centsToReal(emp.SaldoCredito))
}

func cmdHistory() {
	txns := st.ListTransactions()
	if len(txns) == 0 {
		fmt.Println("Nenhuma transação registrada.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tVendedor\tComprador\tBruto\tLíquido\tIBS\tCBS\tCrédito Usado\tCrédito Gerado")
	fmt.Fprintln(w, "---\t--------\t---------\t-----\t-------\t---\t---\t-------------\t-------------")

	for _, t := range txns {
		impostoTotal := t.ValorIBS + t.ValorCBS
		fmt.Fprintf(w, "%s\t%s\t%s\tR$ %.2f\tR$ %.2f\tR$ %.2f\tR$ %.2f\tR$ %.2f\tR$ %.2f\n",
			t.ID,
			t.VendedorCNPJ,
			t.CompradorCNPJ,
			centsToReal(t.ValorBruto),
			centsToReal(t.Liquido),
			centsToReal(t.ValorIBS),
			centsToReal(t.ValorCBS),
			centsToReal(t.CreditoUsado),
			centsToReal(impostoTotal),
		)
	}
	w.Flush()
}

func printUsage() {
	fmt.Println(`Comandos:
  add-company --cnpj=CNPJ --nome=NOME
  sell --vendedor=CNPJ --comprador=CNPJ --valor=CENTAVOS [--aliquota-ibs=0.12] [--aliquota-cbs=0.03]
  balance --cnpj=CNPJ
  history
  help
  exit`)
}

func parseArgs(input string) []string {
	var args []string
	current := strings.Builder{}
	inQuote := false

	for _, r := range input {
		switch {
		case r == '"':
			inQuote = !inQuote
		case r == ' ' && !inQuote:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}

func centsToReal(c int64) float64 {
	return float64(c) / 100.0
}

func isInteractive() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) != 0
}
