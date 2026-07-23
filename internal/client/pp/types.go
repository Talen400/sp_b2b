// Pacote pp implementa um cliente HTTP para a Plataforma Pública do Split
// Payment (RFB/CGIBS/Serpro). Os tipos e funções seguem o dicionário de
// campos do Manual de Integração v1.0 e o OpenAPI v0.0.10.
//
// Fonte principal: Manual de Integração, seção 4 (dicionário de campos) e
// seção 6 (endpoints). Decimal(18,2) conforme formato especificado na seção 4.2.
package pp

import "fmt"

// Decimal18_2 representa um valor monetário em centavos (int64) que serializa
// para string Decimal(18,2) no JSON (ex: 100050 → "1000.50"), conforme exigido
// pela especificação da Plataforma Pública.
//
// O formato Decimal(18,2) tem 18 dígitos no total, sendo 2 após a vírgula.
// Em centavos int64, o valor máximo é 9.223.372.036.854.775.807 (int64 max),
// o que equivale a R$ 92.233.720.368.547.758,07 — mais que suficiente para
// qualquer transação individual.
//
// Fonte: Manual de Integração, seção 4.2 — campos monetários como string.
type Decimal18_2 int64

// String converte centavos para o formato Decimal(18,2).
// Ex: 100050 → "1000.50", -500 → "-5.00".
func (d Decimal18_2) String() string {
	abs := int64(d)
	sign := ""
	if abs < 0 {
		sign = "-"
		abs = -abs
	}
	reais := abs / 100
	centavos := abs % 100
	return fmt.Sprintf("%s%d.%02d", sign, reais, centavos)
}

// MarshalJSON implementa json.Marshaler para Decimal18_2.
// Serializa como string entre aspas: 100000 → "\"1000.00\"".
func (d Decimal18_2) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

// UnmarshalJSON implementa json.Unmarshaler para Decimal18_2.
// Aceita string no formato "\"1000.50\"" e converte para centavos int64.
// Suporta valores negativos.
func (d *Decimal18_2) UnmarshalJSON(data []byte) error {
	s := string(data)
	if len(s) < 2 {
		return fmt.Errorf("Decimal18_2: string muito curta: %s", s)
	}
	s = s[1 : len(s)-1]

	var reais, centavos int64
	n, err := fmt.Sscanf(s, "%d.%02d", &reais, &centavos)
	if err != nil || n != 2 {
		return fmt.Errorf("Decimal18_2: formato inválido: %s", s)
	}
	if reais < 0 {
		*d = Decimal18_2(reais*100 - centavos)
	} else {
		*d = Decimal18_2(reais*100 + centavos)
	}
	return nil
}

// InformeTransacaoIniciadaRequest é o payload do Informe de Transação Iniciada
// (POST /api/v1/{arrj} para Boleto, Pix Dinâmico, Pix Automático).
//
// Fonte: Manual de Integração, seção 4.2 — dicionário de campos do Informe
// de Transação Iniciada.
type InformeTransacaoIniciadaRequest struct {
	NSUId       string      `json:"nsuId"`       // Número Sequencial Único da transação
	CNPJRec     string      `json:"cnpjRec"`     // CNPJ do Recebedor (Fornecedor)
	VLInf       Decimal18_2 `json:"vlInf"`       // Valor total Informado da transação
	VLIbs       Decimal18_2 `json:"vlIbs"`       // Valor de IBS informado
	VLCbs       Decimal18_2 `json:"vlCbs"`       // Valor de CBS informado
	DtHrCriacao string      `json:"dtHrCriacao"` // Data e hora da criação (RFC 3339)
}

// InformeSegregacaoRequest é o payload do Informe de Segregação
// (POST /api/v1/segregacao). Enviado em lote 2x por dia útil.
// Gera obrigação de Repasse Financeiro.
//
// Fonte: Manual de Integração, seção 4.2 — dicionário de campos do Informe
// de Segregação. Manual de Operações, seção 4.3 — diferença entre Preliminar
// e Segregação.
type InformeSegregacaoRequest struct {
	CNPJRaizPspRecDir string               `json:"cnpjRaizPspRecDir"` // Raiz do CNPJ do PSP Recebedor Direto
	DtHrMsg           string               `json:"dtHrMsg"`           // Data e hora da mensagem (RFC 3339)
	IDInfSegr         string               `json:"idInfSegr"`         // Identificador único do informe
	Transacoes        []TransacaoSegregada `json:"transacoes"`        // Lista de transações segregadas
	TotalTrans        int                  `json:"totalTrans"`        // Total de transações no lote
	ValorTotalCbs     Decimal18_2          `json:"valorTotalCbs"`     // Valor total de CBS segregado
	ValorTotalIbs     Decimal18_2          `json:"valorTotalIbs"`     // Valor total de IBS segregado
}

// TransacaoSegregada representa uma transação individual dentro de um lote
// de segregação. Contém apenas os campos necessários para o repasse.
type TransacaoSegregada struct {
	NSUId    string      `json:"nsuId"`    // NSU da transação
	VLIbsSeg Decimal18_2 `json:"vlIbsSeg"` // Valor de IBS efetivamente segregado
	VLCbsSeg Decimal18_2 `json:"vlCbsSeg"` // Valor de CBS efetivamente segregado
}

// RetornoSuperInteligente é a resposta do long polling do Super Inteligente.
// O governo envia correções de valores de tributo para os arranjos Super
// Inteligente (Boleto, Pix Dinâmico, Pix Automático).
//
// Fonte: Manual de Integração, seção 3.6 — Retorno Super Inteligente.
// Manual de Operações, seção 2.2 — Modelo Super Inteligente.
type RetornoSuperInteligente struct {
	Tributos []TributoRetorno `json:"tributos"`
}

// TributoRetorno representa um valor de tributo corrigido pelo Super Inteligente.
// vlCbsCorr/vlIbsCorr: valores corrigidos (diferem do informado).
// vlCbsAberto/vlIbsAberto: valores em aberto (parcela não liquidada).
//
// Fonte: Manual de Operações, seção 3.1 — as 5 categorias de valor de tributo.
type TributoRetorno struct {
	NSUId       string       `json:"nsuId"`                 // NSU da transação
	CodMsg      string       `json:"codMsg"`                // Código da mensagem (ex: CORRECAO, ABERTO)
	VLInf       Decimal18_2  `json:"vlInf"`                 // Valor informado original
	VLCbsCorr   *Decimal18_2 `json:"vlCbsCorr,omitempty"`   // CBS corrigido (opcional)
	VLIbsCorr   *Decimal18_2 `json:"vlIbsCorr,omitempty"`   // IBS corrigido (opcional)
	VLCbsAberto *Decimal18_2 `json:"vlCbsAberto,omitempty"` // CBS em aberto (opcional)
	VLIbsAberto *Decimal18_2 `json:"vlIbsAberto,omitempty"` // IBS em aberto (opcional)
}

// LongPollingResponse encapsula a resposta de uma requisição de long polling.
// Se NoContent = true, não há mensagens no momento (204 No Content da PP).
// ProximoToken deve ser usado na próxima chamada (ContinueLongPolling).
//
// Fonte: Manual de Integração, seção 3.6 — fluxo de long polling com
// token de posição.
type LongPollingResponse struct {
	Messages     *RetornoSuperInteligente // mensagens recebidas (nil se NoContent)
	ProximoToken string                   // token para a próxima consulta (header proximoToken)
	NoContent    bool                     // true se PP respondeu 204
}

// PPErrorResponse representa um erro RFC 7807 retornado pela PP.
// O campo retornabilidade (não incluído aqui por simplicidade) indica
// se o PSP pode ou não retentar a requisição.
//
// Fonte: Manual de Integração, seção 5 — Política de tratamento de erros.
type PPErrorResponse struct {
	Type     string `json:"type"`     // URI do tipo de erro
	Title    string `json:"title"`    // Título legível do erro
	Status   int    `json:"status"`   // HTTP status code
	Detail   string `json:"detail"`   // Descrição detalhada
	Instance string `json:"instance"` // Path do endpoint que gerou o erro
}

func (e *PPErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s (HTTP %d)", e.Title, e.Detail, e.Status)
}
