CREATE TABLE IF NOT EXISTS companies (
    cnpj         TEXT PRIMARY KEY,
    nome         TEXT NOT NULL,
    saldo_credito INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS transactions (
    id            TEXT PRIMARY KEY,
    vendedor_cnpj TEXT NOT NULL REFERENCES companies(cnpj),
    comprador_cnpj TEXT NOT NULL REFERENCES companies(cnpj),
    valor_bruto   INTEGER NOT NULL,
    aliquota_ibs  REAL NOT NULL,
    aliquota_cbs  REAL NOT NULL,
    valor_liquido INTEGER NOT NULL,
    valor_ibs     INTEGER NOT NULL,
    valor_cbs     INTEGER NOT NULL,
    credito_usado INTEGER NOT NULL DEFAULT 0,
    timestamp     TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_transactions_vendedor ON transactions(vendedor_cnpj);
CREATE INDEX IF NOT EXISTS idx_transactions_comprador ON transactions(comprador_cnpj);
