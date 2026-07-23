package seed

import (
	"fmt"
	"time"

	"github.com/Talen400/sp_b2b/internal/domain"
	"github.com/Talen400/sp_b2b/internal/repository"
)

func Run(cr repository.CompanyRepo, tr repository.TransactionRepo) error {
	companies := []struct {
		cnpj string
		nome string
	}{
		{"11.111.111/0001-11", "Fazenda Boa Vista"},
		{"22.222.222/0001-22", "Fábrica de Sucos SA"},
		{"33.333.333/0001-33", "Mercado Central"},
	}

	for _, c := range companies {
		if err := cr.Create(c.cnpj, c.nome); err != nil {
			return fmt.Errorf("cadastrar %s: %w", c.nome, err)
		}
		fmt.Printf("✅ Empresa cadastrada: %s (%s)\n", c.nome, c.cnpj)
	}

	vendas := []struct {
		desc       string
		vendedor   string
		comprador  string
		valorBruto int64
		ibs        float64
		cbs        float64
	}{
		{
			"Fazenda → Fábrica (R$ 1.000,00)",
			"11.111.111/0001-11", "22.222.222/0001-22",
			100000, 0.12, 0.03,
		},
		{
			"Fábrica → Mercado (R$ 3.000,00)",
			"22.222.222/0001-22", "33.333.333/0001-33",
			300000, 0.12, 0.03,
		},
	}

	for _, v := range vendas {
		split, err := domain.CalculateSplit(v.valorBruto, v.ibs, v.cbs)
		if err != nil {
			return fmt.Errorf("calcular split %s: %w", v.desc, err)
		}

		vendedor, _ := cr.Get(v.vendedor)
		comprador, _ := cr.Get(v.comprador)

		creditoUsado := domain.UseCredit(&vendedor, split)
		creditoGerado := domain.ApplyCredit(&comprador, split)

		id := fmt.Sprintf("TXN-DEMO-%d", time.Now().UnixNano())
		txn := domain.Transaction{
			ID:            id,
			VendedorCNPJ:  v.vendedor,
			CompradorCNPJ: v.comprador,
			ValorBruto:    v.valorBruto,
			AliquotaIBS:   v.ibs,
			AliquotaCBS:   v.cbs,
			ValorLiquido:  split.Liquido,
			ValorIBS:      split.ValorIBS,
			ValorCBS:      split.ValorCBS,
			CreditoUsado:  creditoUsado,
			Timestamp:     time.Now(),
		}

		if err := tr.Create(txn); err != nil {
			return fmt.Errorf("salvar transação %s: %w", v.desc, err)
		}
		if err := cr.UpdateCredit(v.vendedor, vendedor.SaldoCredito); err != nil {
			return fmt.Errorf("atualizar crédito vendedor %s: %w", v.desc, err)
		}
		if err := cr.UpdateCredit(v.comprador, comprador.SaldoCredito); err != nil {
			return fmt.Errorf("atualizar crédito comprador %s: %w", v.desc, err)
		}

		fmt.Printf("✅ %s\n", v.desc)
		fmt.Printf("   Líquido: R$ %.2f | IBS: R$ %.2f | CBS: R$ %.2f\n",
			float64(split.Liquido)/100, float64(split.ValorIBS)/100, float64(split.ValorCBS)/100)
		fmt.Printf("   Crédito usado: R$ %.2f | Crédito gerado: R$ %.2f\n",
			float64(creditoUsado)/100, float64(creditoGerado)/100)
	}

	fmt.Println("\n✅ Seed concluído!")
	for _, c := range companies {
		emp, _ := cr.Get(c.cnpj)
		fmt.Printf("   %s: saldo de crédito R$ %.2f\n", emp.Nome, float64(emp.SaldoCredito)/100)
	}

	return nil
}
