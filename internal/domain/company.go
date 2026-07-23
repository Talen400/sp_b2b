// Pacote domain contém as entidades e regras de negócio puras do split payment.
package domain

// Company representa uma empresa participante da simulação.
//
// Cada empresa tem um CNPJ (string livre — não validamos dígitos nem formato
// alfanumérico), um nome fantasia e um saldo de crédito tributário acumulado
// em centavos (int64).
//
// O saldo de crédito é a simplificação didática que fizemos para representar
// a não-cumulatividade do IBS/CBS: quando a empresa compra, acumula crédito;
// quando vende, abate o crédito do imposto devido (ver ApplyCredit e UseCredit).
//
// No mundo real:
//   - CNPJ pode conter letras a partir de jul/2026 (IN RFB 2.229/2024).
//   - O crédito tributário não é um saldo simples como neste modelo — é
//     apurado contabilmente por período (geralmente mensal) com regras
//     complexas de creditamento.
type Company struct {
	CNPJ         string `json:"cnpj"`
	Nome         string `json:"nome"`
	SaldoCredito int64  `json:"saldo_credito"`
}
