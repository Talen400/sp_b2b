// Pacote repository define as interfaces de persistência que o domínio
// e os handlers HTTP usam. A implementação concreta fica em sqlite/.
// Isso permite trocar de banco sem alterar as camadas acima.
package repository

import "github.com/Talen400/sp_b2b/internal/domain"

// CompanyRepo é a interface de persistência para empresas.
type CompanyRepo interface {
	// Create cadastra uma nova empresa com saldo de crédito inicial zero.
	// Retorna erro se o CNPJ já existir.
	Create(cnpj, nome string) error

	// Get consulta uma empresa pelo CNPJ.
	// Retorna erro "não encontrada" se não existir.
	Get(cnpj string) (domain.Company, error)

	// List retorna todas as empresas cadastradas, ordenadas por nome.
	List() ([]domain.Company, error)

	// UpdateCredit atualiza o saldo de crédito de uma empresa.
	// Usado após uma transação para refletir o crédito usado/gerado.
	UpdateCredit(cnpj string, novoSaldo int64) error
}

// TransactionRepo é a interface de persistência para transações.
type TransactionRepo interface {
	// Create persiste uma nova transação (split + crédito já calculados).
	Create(t domain.Transaction) error

	// List retorna o histórico de transações.
	// Se cnpjFilter não for vazio, filtra por CNPJ (vendedor ou comprador).
	List(cnpjFilter string) ([]domain.Transaction, error)

	// Get consulta uma transação pelo ID.
	// Retorna erro "não encontrada" se não existir.
	Get(id string) (domain.Transaction, error)
}
