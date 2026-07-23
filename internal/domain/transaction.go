package domain

import "time"

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
