// Pacote sqlite implementa as interfaces repository.CompanyRepo e
// repository.TransactionRepo usando SQLite puro-Go (modernc.org/sqlite,
// sem cgo). As migrations são embutidas no binário via //go:embed.
package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/Talen400/sp_b2b/internal/domain"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// DB é um wrapper sobre *sql.DB com limite de 1 conexão simultânea
// (SQLite não lida bem com concorrência de escrita).
type DB struct {
	db *sql.DB
}

// Open abre (ou cria) o arquivo SQLite no path informado.
// O driver usado é modernc.org/sqlite, implementação pura em Go.
func Open(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("abrir banco: %w", err)
	}
	db.SetMaxOpenConns(1)
	return &DB{db: db}, nil
}

// Close fecha a conexão com o banco.
func (d *DB) Close() error {
	return d.db.Close()
}

// Migrate aplica todos os arquivos .sql embutidos em migrations/.
// Executa na ordem alfabética dos nomes de arquivo.
func (d *DB) Migrate() error {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("ler migrations: %w", err)
	}
	for _, e := range entries {
		data, err := migrationsFS.ReadFile("migrations/" + e.Name())
		if err != nil {
			return fmt.Errorf("ler %s: %w", e.Name(), err)
		}
		if _, err := d.db.Exec(string(data)); err != nil {
			return fmt.Errorf("executar %s: %w", e.Name(), err)
		}
	}
	return nil
}

// CompanyRepository implementa repository.CompanyRepo com SQLite.
type CompanyRepository struct {
	db *DB
}

// TransactionRepository implementa repository.TransactionRepo com SQLite.
type TransactionRepository struct {
	db *DB
}

// NewCompanyRepository cria um CompanyRepository a partir de uma conexão DB.
func NewCompanyRepository(db *DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// NewTransactionRepository cria um TransactionRepository a partir de uma conexão DB.
func NewTransactionRepository(db *DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create insere uma nova empresa com saldo de crédito inicial zero.
// Retorna erro se o CNPJ já existir (violação da constraint UNIQUE).
func (r *CompanyRepository) Create(cnpj, nome string) error {
	_, err := r.db.db.Exec(
		"INSERT INTO companies (cnpj, nome, saldo_credito) VALUES (?, ?, 0)",
		cnpj, nome,
	)
	if err != nil {
		if isUniqueConstraintErr(err) {
			return fmt.Errorf("empresa com CNPJ %s já existe", cnpj)
		}
		return fmt.Errorf("inserir empresa: %w", err)
	}
	return nil
}

// Get consulta uma empresa pelo CNPJ. Retorna erro se não existir.
func (r *CompanyRepository) Get(cnpj string) (domain.Company, error) {
	row := r.db.db.QueryRow(
		"SELECT cnpj, nome, saldo_credito FROM companies WHERE cnpj = ?", cnpj,
	)
	var c domain.Company
	if err := row.Scan(&c.CNPJ, &c.Nome, &c.SaldoCredito); err != nil {
		if err == sql.ErrNoRows {
			return c, fmt.Errorf("empresa com CNPJ %s não encontrada", cnpj)
		}
		return c, fmt.Errorf("consultar empresa: %w", err)
	}
	return c, nil
}

// List retorna todas as empresas ordenadas por nome.
func (r *CompanyRepository) List() ([]domain.Company, error) {
	rows, err := r.db.db.Query("SELECT cnpj, nome, saldo_credito FROM companies ORDER BY nome")
	if err != nil {
		return nil, fmt.Errorf("listar empresas: %w", err)
	}
	defer rows.Close()

	var companies []domain.Company
	for rows.Next() {
		var c domain.Company
		if err := rows.Scan(&c.CNPJ, &c.Nome, &c.SaldoCredito); err != nil {
			return nil, fmt.Errorf("ler empresa: %w", err)
		}
		companies = append(companies, c)
	}
	return companies, rows.Err()
}

// UpdateCredit atualiza o saldo de crédito de uma empresa.
// Retorna erro se o CNPJ não for encontrado.
func (r *CompanyRepository) UpdateCredit(cnpj string, novoSaldo int64) error {
	res, err := r.db.db.Exec("UPDATE companies SET saldo_credito = ? WHERE cnpj = ?", novoSaldo, cnpj)
	if err != nil {
		return fmt.Errorf("atualizar crédito: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("empresa com CNPJ %s não encontrada", cnpj)
	}
	return nil
}

// Create persiste uma transação com todos os campos do split já calculados.
// O timestamp é armazenado como string RFC 3339.
func (r *TransactionRepository) Create(t domain.Transaction) error {
	_, err := r.db.db.Exec(
		`INSERT INTO transactions
		 (id, vendedor_cnpj, comprador_cnpj, valor_bruto, aliquota_ibs, aliquota_cbs,
		  valor_liquido, valor_ibs, valor_cbs, credito_usado, timestamp)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.VendedorCNPJ, t.CompradorCNPJ, t.ValorBruto,
		t.AliquotaIBS, t.AliquotaCBS, t.ValorLiquido,
		t.ValorIBS, t.ValorCBS, t.CreditoUsado,
		t.Timestamp.Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("inserir transação: %w", err)
	}
	return nil
}

// List retorna as transações ordenadas por timestamp.
// Se cnpjFilter não for vazio, filtra por CNPJ (vendedor ou comprador).
func (r *TransactionRepository) List(cnpjFilter string) ([]domain.Transaction, error) {
	var query string
	var args []interface{}

	if cnpjFilter != "" {
		query = `SELECT id, vendedor_cnpj, comprador_cnpj, valor_bruto,
		 aliquota_ibs, aliquota_cbs, valor_liquido, valor_ibs, valor_cbs,
		 credito_usado, timestamp FROM transactions
		 WHERE vendedor_cnpj = ? OR comprador_cnpj = ? ORDER BY timestamp`
		args = append(args, cnpjFilter, cnpjFilter)
	} else {
		query = `SELECT id, vendedor_cnpj, comprador_cnpj, valor_bruto,
		 aliquota_ibs, aliquota_cbs, valor_liquido, valor_ibs, valor_cbs,
		 credito_usado, timestamp FROM transactions ORDER BY timestamp`
	}

	rows, err := r.db.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("listar transações: %w", err)
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		var ts string
		if err := rows.Scan(&t.ID, &t.VendedorCNPJ, &t.CompradorCNPJ,
			&t.ValorBruto, &t.AliquotaIBS, &t.AliquotaCBS,
			&t.ValorLiquido, &t.ValorIBS, &t.ValorCBS,
			&t.CreditoUsado, &ts); err != nil {
			return nil, fmt.Errorf("ler transação: %w", err)
		}
		t.Timestamp, _ = time.Parse(time.RFC3339, ts)
		transactions = append(transactions, t)
	}
	return transactions, rows.Err()
}

// Get consulta uma transação pelo ID. Retorna erro se não existir.
func (r *TransactionRepository) Get(id string) (domain.Transaction, error) {
	row := r.db.db.QueryRow(
		`SELECT id, vendedor_cnpj, comprador_cnpj, valor_bruto,
		 aliquota_ibs, aliquota_cbs, valor_liquido, valor_ibs, valor_cbs,
		 credito_usado, timestamp FROM transactions WHERE id = ?`,
		id,
	)
	var t domain.Transaction
	var ts string
	if err := row.Scan(&t.ID, &t.VendedorCNPJ, &t.CompradorCNPJ,
		&t.ValorBruto, &t.AliquotaIBS, &t.AliquotaCBS,
		&t.ValorLiquido, &t.ValorIBS, &t.ValorCBS,
		&t.CreditoUsado, &ts); err != nil {
		if err == sql.ErrNoRows {
			return t, fmt.Errorf("transação %s não encontrada", id)
		}
		return t, fmt.Errorf("consultar transação: %w", err)
	}
	t.Timestamp, _ = time.Parse(time.RFC3339, ts)
	return t, nil
}

// isUniqueConstraintErr verifica se o erro do SQLite é violação de UNIQUE.
func isUniqueConstraintErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE")
}
