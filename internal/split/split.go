package split

import (
	"fmt"
	"math"
	"time"
)

type SplitResult struct {
	Liquido  int64
	ValorIBS int64
	ValorCBS int64
}

type Transaction struct {
	ID            string
	VendedorCNPJ  string
	CompradorCNPJ string
	ValorBruto    int64
	AliquotaIBS   float64
	AliquotaCBS   float64
	ValorIBS      int64
	ValorCBS      int64
	Liquido       int64
	CreditoUsado  int64
	Timestamp     time.Time
}

func CalculateSplit(valorBruto int64, aliquotaIBS, aliquotaCBS float64) (SplitResult, error) {
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
