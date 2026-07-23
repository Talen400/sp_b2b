package repository

import "github.com/Talen400/sp_b2b/internal/domain"

type CompanyRepo interface {
	Create(cnpj, nome string) error
	Get(cnpj string) (domain.Company, error)
	List() ([]domain.Company, error)
	UpdateCredit(cnpj string, novoSaldo int64) error
}

type TransactionRepo interface {
	Create(t domain.Transaction) error
	List(cnpjFilter string) ([]domain.Transaction, error)
	Get(id string) (domain.Transaction, error)
}
