package domain

import "time"

// Transaction representa uma transação de venda entre duas empresas na simulação.
//
// Cada transação registra:
//   - Quem vendeu (VendedorCNPJ) e quem comprou (CompradorCNPJ).
//   - O valor bruto e as alíquotas de IBS/CBS aplicadas.
//   - O resultado do split: valor líquido, IBS segregado, CBS segregada.
//   - O crédito tributário usado pelo vendedor (abatido do imposto devido)
//     e gerado para o comprador (acumulado para abater em vendas futuras).
//
// O ID é gerado no momento da criação (formato "TXN-{nanotimestamp}").
// O timestamp é registrado no momento da criação no servidor.
//
// No mundo real, uma transação teria dezenas de campos adicionais:
// Documento Fiscal, NSU, arranjo de pagamento, PSPs envolvidos,
// valores segregados por categoria (Informado, Corrigido, Em Aberto,
// Segregado, Aplicado), etc. Esta é uma simplificação didática.
type Transaction struct {
	ID            string    `json:"id"`
	VendedorCNPJ  string    `json:"vendedor_cnpj"`
	CompradorCNPJ string    `json:"comprador_cnpj"`
	ValorBruto    int64     `json:"valor_bruto"`
	AliquotaIBS   float64   `json:"aliquota_ibs"`
	AliquotaCBS   float64   `json:"aliquota_cbs"`
	ValorLiquido  int64     `json:"valor_liquido"`
	ValorIBS      int64     `json:"valor_ibs"`
	ValorCBS      int64     `json:"valor_cbs"`
	CreditoUsado  int64     `json:"credito_usado"`
	Timestamp     time.Time `json:"timestamp"`
}
