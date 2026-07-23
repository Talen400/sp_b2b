package split

import "testing"

func TestCalculateSplit_Normal(t *testing.T) {
	result, err := CalculateSplit(100000, 0.12, 0.03)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if result.ValorIBS != 12000 {
		t.Errorf("ValorIBS esperado 12000, got %d", result.ValorIBS)
	}
	if result.ValorCBS != 3000 {
		t.Errorf("ValorCBS esperado 3000, got %d", result.ValorCBS)
	}
	if result.Liquido != 85000 {
		t.Errorf("Liquido esperado 85000, got %d", result.Liquido)
	}
}

func TestCalculateSplit_ValorZero(t *testing.T) {
	result, err := CalculateSplit(0, 0.12, 0.03)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if result.ValorIBS != 0 || result.ValorCBS != 0 || result.Liquido != 0 {
		t.Errorf("esperado tudo zero, got %+v", result)
	}
}

func TestCalculateSplit_AliquotaZero(t *testing.T) {
	result, err := CalculateSplit(100000, 0, 0)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if result.ValorIBS != 0 || result.ValorCBS != 0 {
		t.Errorf("esperado imposto zero, got IBS=%d CBS=%d", result.ValorIBS, result.ValorCBS)
	}
	if result.Liquido != 100000 {
		t.Errorf("Liquido esperado 100000, got %d", result.Liquido)
	}
}

func TestCalculateSplit_SomaAliquotaExcede100(t *testing.T) {
	_, err := CalculateSplit(100000, 0.6, 0.5)
	if err == nil {
		t.Fatal("esperava erro para soma > 100%, mas não houve erro")
	}
}

func TestCalculateSplit_SomaAliquotaExata100(t *testing.T) {
	_, err := CalculateSplit(100000, 0.7, 0.3)
	if err == nil {
		t.Fatal("esperava erro para soma == 100%, mas não houve erro")
	}
}

func TestCalculateSplit_AliquotaNegativa(t *testing.T) {
	_, err := CalculateSplit(100000, -0.1, 0.03)
	if err == nil {
		t.Fatal("esperava erro para alíquota negativa")
	}
}

func TestCalculateSplit_Arredondamento(t *testing.T) {
	result, err := CalculateSplit(3, 0.3333, 0.3333)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	// 3 * 0.3333 = 0.9999 -> arredonda para 1
	// 3 * 0.3333 = 0.9999 -> arredonda para 1
	// líquido = 3 - 1 - 1 = 1
	if result.ValorIBS != 1 || result.ValorCBS != 1 || result.Liquido != 1 {
		t.Errorf("resultado inesperado: %+v", result)
	}
}
