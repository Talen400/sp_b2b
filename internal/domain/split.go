// Pacote domain contém as regras de negócio puras do split payment.
// Nenhuma função aqui importa net/http, database/sql ou qualquer
// infraestrutura. Todas as regras são testadas exaustivamente nos
// arquivos _test.go correspondentes.
package domain

import (
	"fmt"
	"math"
)

// SplitResult é o retorno de CalculateSplit. Todos os valores em centavos (int64).
type SplitResult struct {
	Liquido  int64 // valor que o vendedor recebe de fato (bruto - IBS - CBS)
	ValorIBS int64 // valor do IBS segregado na transação
	ValorCBS int64 // valor da CBS segregada na transação
}

// CalculateSplit calcula a segregação de IBS e CBS sobre um valor bruto.
//
// Parâmetros:
//   - valorBruto: valor total da transação em centavos (int64).
//   - aliquotaIBS: alíquota do IBS em decimal (ex: 0.12 = 12%).
//   - aliquotaCBS: alíquota da CBS em decimal (ex: 0.03 = 3%).
//
// Retorna SplitResult com os valores segregados e o líquido.
// Erro se valorBruto < 0, se alíquota < 0, ou se soma das alíquotas >= 100%.
//
// Algoritmo:
//   - valorIBS = round(valorBruto * aliquotaIBS)
//   - valorCBS = round(valorBruto * aliquotaCBS)
//   - líquido  = valorBruto - valorIBS - valorCBS
//
// Fonte: a lógica de cálculo é do simulador original (pré-pivot).
// O arredondamento segue math.Round (bancário não se aplica neste
// contexto simplificado). No mundo real, os valores de tributo são
// definidos pelo Documento Fiscal ou pelo Super Inteligente, não
// calculados pelo PSP.
func CalculateSplit(valorBruto int64, aliquotaIBS, aliquotaCBS float64) (SplitResult, error) {
	if valorBruto < 0 {
		return SplitResult{}, fmt.Errorf("valor bruto não pode ser negativo: %d", valorBruto)
	}
	if aliquotaIBS < 0 || aliquotaCBS < 0 {
		return SplitResult{}, fmt.Errorf("alíquotas não podem ser negativas: IBS=%.4f, CBS=%.4f", aliquotaIBS, aliquotaCBS)
	}
	if aliquotaIBS+aliquotaCBS >= 1.0 {
		return SplitResult{}, fmt.Errorf("soma das alíquotas (%.2f%%) não pode ser >= 100%%", (aliquotaIBS+aliquotaCBS)*100)
	}

	valorIBS := int64(math.Round(float64(valorBruto) * aliquotaIBS))
	valorCBS := int64(math.Round(float64(valorBruto) * aliquotaCBS))
	liquido := valorBruto - valorIBS - valorCBS

	return SplitResult{
		Liquido:  liquido,
		ValorIBS: valorIBS,
		ValorCBS: valorCBS,
	}, nil
}
