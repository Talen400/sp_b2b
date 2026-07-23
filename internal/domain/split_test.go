package domain

import "testing"

func TestCalculateSplit_Normal(t *testing.T) {
	r, err := CalculateSplit(100000, 0.12, 0.03)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if r.ValorIBS != 12000 || r.ValorCBS != 3000 || r.Liquido != 85000 {
		t.Errorf("resultado inesperado: %+v", r)
	}
}

func TestCalculateSplit_ValorZero(t *testing.T) {
	r, err := CalculateSplit(0, 0.12, 0.03)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if r.ValorIBS != 0 || r.ValorCBS != 0 || r.Liquido != 0 {
		t.Errorf("esperado tudo zero: %+v", r)
	}
}

func TestCalculateSplit_AliquotaZero(t *testing.T) {
	r, err := CalculateSplit(100000, 0, 0)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if r.Liquido != 100000 {
		t.Errorf("Liquido esperado 100000, got %d", r.Liquido)
	}
}

func TestCalculateSplit_SomaAliquotaExcede100(t *testing.T) {
	_, err := CalculateSplit(100000, 0.6, 0.5)
	if err == nil {
		t.Fatal("esperava erro para soma > 100%")
	}
}

func TestCalculateSplit_SomaAliquotaExata100(t *testing.T) {
	_, err := CalculateSplit(100000, 0.7, 0.3)
	if err == nil {
		t.Fatal("esperava erro para soma == 100%")
	}
}

func TestCalculateSplit_AliquotaNegativa(t *testing.T) {
	_, err := CalculateSplit(100000, -0.1, 0.03)
	if err == nil {
		t.Fatal("esperava erro para alíquota negativa")
	}
}

func TestCalculateSplit_ValorNegativo(t *testing.T) {
	_, err := CalculateSplit(-1, 0.12, 0.03)
	if err == nil {
		t.Fatal("esperava erro para valor negativo")
	}
}

func TestCalculateSplit_Arredondamento(t *testing.T) {
	r, err := CalculateSplit(3, 0.3333, 0.3333)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if r.ValorIBS != 1 || r.ValorCBS != 1 || r.Liquido != 1 {
		t.Errorf("resultado inesperado: %+v", r)
	}
}
