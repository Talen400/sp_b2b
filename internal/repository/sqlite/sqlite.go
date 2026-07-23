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

type DB struct {
	db *sql.DB
}

func Open(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("abrir banco: %w", err)
	}
	db.SetMaxOpenConns(1)
	return &DB{db: db}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

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

type CompanyRepository struct {
	db *DB
}

type TransactionRepository struct {
	db *DB
}

func NewCompanyRepository(db *DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func NewTransactionRepository(db *DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

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

func isUniqueConstraintErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE")
}
