package pp

import "fmt"

// Decimal18_2 represents a monetary value in centavos (int64) that serializes to
// a Decimal(18,2) string (e.g., 100050 → "1000.50") in JSON, as required by the
// Plataforma Pública specification.
type Decimal18_2 int64

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

func (d Decimal18_2) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

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

// InformeTransacaoIniciadaRequest is the payload for POST /api/v1/{arrj}
// (Informe de Transação Iniciada).
type InformeTransacaoIniciadaRequest struct {
	NSUId       string      `json:"nsuId"`
	CNPJRec     string      `json:"cnpjRec"`
	VLInf       Decimal18_2 `json:"vlInf"`
	VLIbs       Decimal18_2 `json:"vlIbs"`
	VLCbs       Decimal18_2 `json:"vlCbs"`
	DtHrCriacao string      `json:"dtHrCriacao"`
}

// InformeSegregacaoRequest is the payload for POST /api/v1/segregacao.
type InformeSegregacaoRequest struct {
	CNPJRaizPspRecDir string               `json:"cnpjRaizPspRecDir"`
	DtHrMsg           string               `json:"dtHrMsg"`
	IDInfSegr         string               `json:"idInfSegr"`
	Transacoes        []TransacaoSegregada `json:"transacoes"`
	TotalTrans        int                  `json:"totalTrans"`
	ValorTotalCbs     Decimal18_2          `json:"valorTotalCbs"`
	ValorTotalIbs     Decimal18_2          `json:"valorTotalIbs"`
}

// TransacaoSegregada represents a single transaction within a segregação batch.
type TransacaoSegregada struct {
	NSUId    string      `json:"nsuId"`
	VLIbsSeg Decimal18_2 `json:"vlIbsSeg"`
	VLCbsSeg Decimal18_2 `json:"vlCbsSeg"`
}

// RetornoSuperInteligente is the response payload from long polling
// (Retorno Super Inteligente).
type RetornoSuperInteligente struct {
	Tributos []TributoRetorno `json:"tributos"`
}

// TributoRetorno represents a single corrected tax value from the Super Inteligente.
type TributoRetorno struct {
	NSUId       string       `json:"nsuId"`
	CodMsg      string       `json:"codMsg"`
	VLInf       Decimal18_2  `json:"vlInf"`
	VLCbsCorr   *Decimal18_2 `json:"vlCbsCorr,omitempty"`
	VLIbsCorr   *Decimal18_2 `json:"vlIbsCorr,omitempty"`
	VLCbsAberto *Decimal18_2 `json:"vlCbsAberto,omitempty"`
	VLIbsAberto *Decimal18_2 `json:"vlIbsAberto,omitempty"`
}

// LongPollingResponse wraps the response from a long polling request.
type LongPollingResponse struct {
	Messages     *RetornoSuperInteligente
	ProximoToken string
	NoContent    bool
}

// PPErrorResponse represents an RFC 7807 error from the PP API.
type PPErrorResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func (e *PPErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s (HTTP %d)", e.Title, e.Detail, e.Status)
}
