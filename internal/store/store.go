package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/Talen400/sp_b2b/internal/company"
	"github.com/Talen400/sp_b2b/internal/split"
)

type Store struct {
	mu         sync.RWMutex
	empresas   map[string]*company.Company
	transacoes []split.Transaction
	nextID     int64
}

func NewStore() *Store {
	return &Store{
		empresas:   make(map[string]*company.Company),
		transacoes: make([]split.Transaction, 0),
		nextID:     1,
	}
}

func (s *Store) AddCompany(cnpj, nome string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.empresas[cnpj]; exists {
		return fmt.Errorf("empresa com CNPJ %s já está cadastrada", cnpj)
	}

	s.empresas[cnpj] = company.New(cnpj, nome)
	return nil
}

func (s *Store) GetCompany(cnpj string) (*company.Company, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	emp, exists := s.empresas[cnpj]
	if !exists {
		return nil, fmt.Errorf("empresa com CNPJ %s não encontrada", cnpj)
	}
	return emp, nil
}

func (s *Store) GetBalance(cnpj string) (int64, error) {
	emp, err := s.GetCompany(cnpj)
	if err != nil {
		return 0, err
	}
	return emp.SaldoCredito, nil
}

func (s *Store) ListTransactions() []split.Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]split.Transaction, len(s.transacoes))
	copy(result, s.transacoes)
	return result
}

type ProcessTransactionResult struct {
	Transaction   split.Transaction
	CreditoGerado int64
}

func (s *Store) ProcessTransaction(vendedorCNPJ, compradorCNPJ string, valorBruto int64, aliquotaIBS, aliquotaCBS float64) (ProcessTransactionResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	vendedor, exists := s.empresas[vendedorCNPJ]
	if !exists {
		return ProcessTransactionResult{}, fmt.Errorf("vendedor com CNPJ %s não encontrado", vendedorCNPJ)
	}
	comprador, exists := s.empresas[compradorCNPJ]
	if !exists {
		return ProcessTransactionResult{}, fmt.Errorf("comprador com CNPJ %s não encontrado", compradorCNPJ)
	}

	splitResult, err := split.CalculateSplit(valorBruto, aliquotaIBS, aliquotaCBS)
	if err != nil {
		return ProcessTransactionResult{}, fmt.Errorf("erro no cálculo do split: %w", err)
	}

	impostoTotal := splitResult.ValorIBS + splitResult.ValorCBS
	creditoUsado := vendedor.SaldoCredito
	if creditoUsado > impostoTotal {
		creditoUsado = impostoTotal
	}
	vendedor.SaldoCredito -= creditoUsado

	creditoGerado := impostoTotal
	comprador.SaldoCredito += creditoGerado

	id := fmt.Sprintf("TXN-%05d", s.nextID)
	s.nextID++

	txn := split.Transaction{
		ID:            id,
		VendedorCNPJ:  vendedorCNPJ,
		CompradorCNPJ: compradorCNPJ,
		ValorBruto:    valorBruto,
		AliquotaIBS:   aliquotaIBS,
		AliquotaCBS:   aliquotaCBS,
		ValorIBS:      splitResult.ValorIBS,
		ValorCBS:      splitResult.ValorCBS,
		Liquido:       splitResult.Liquido,
		CreditoUsado:  creditoUsado,
		Timestamp:     time.Now(),
	}

	s.transacoes = append(s.transacoes, txn)

	return ProcessTransactionResult{
		Transaction:   txn,
		CreditoGerado: creditoGerado,
	}, nil
}
