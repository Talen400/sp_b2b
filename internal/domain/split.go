package domain

import (
	"fmt"
	"math"
)

type SplitResult struct {
	Liquido  int64
	ValorIBS int64
	ValorCBS int64
}

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
