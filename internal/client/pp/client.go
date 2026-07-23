// Pacote pp implementa um cliente HTTP para a Plataforma Pública do Split
// Payment. Injeta automaticamente os headers obrigatórios (Message-Id,
// Correlation-Id, Tenant-Id, Timestamp) e trata erros no formato RFC 7807.
package pp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Client é um cliente HTTP para a Plataforma Pública do Split Payment.
// Injeta automaticamente os headers obrigatórios e trata erros RFC 7807.
//
// Em ambiente dev, aponte baseURL para o mock Prism (ex: http://localhost:4010).
// Em produção, aponte para o endpoint real da PP (requer mTLS + OAuth 2.0,
// não implementados neste cliente simplificado).
type Client struct {
	baseURL    string
	httpClient *http.Client
	tenantID   string
}

// New cria um novo Client para a Plataforma Pública.
//
// baseURL é o endpoint da PP (ex: "http://localhost:4010" para o mock Prism).
// tenantID é o identificador do PSP enviado como header Tenant-Id.
func New(baseURL, tenantID string) *Client {
	return &Client{
		baseURL:    baseURL,
		tenantID:   tenantID,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// SendInformeTransacaoIniciada envia um Informe de Transação Iniciada para a PP.
//
// arrangement deve ser um dos arranjos Super Inteligente:
// "boleto", "pix-dinamico", "pix-automatico".
//
// Retorna o HTTP status code, um PPErrorResponse (se a PP rejeitou), e erro
// (se houve falha de comunicação/parse).
//
// Fonte: Manual de Integração, seção 6 — endpoints por arranjo.
func (c *Client) SendInformeTransacaoIniciada(ctx context.Context, arrangement string, req *InformeTransacaoIniciadaRequest) (int, *PPErrorResponse, error) {
	path := fmt.Sprintf("/api/v1/%s", arrangement)
	return c.doRequest(ctx, http.MethodPost, path, req)
}

// SendInformeSegregacao envia um Informe de Segregação (lote) para a PP.
// Gera obrigação de Repasse Financeiro.
//
// Fonte: Manual de Operações, seção 4.3 — Informe de Segregação como
// comunicação definitiva e vinculante.
func (c *Client) SendInformeSegregacao(ctx context.Context, req *InformeSegregacaoRequest) (int, *PPErrorResponse, error) {
	return c.doRequest(ctx, http.MethodPost, "/api/v1/segregacao", req)
}

// StartLongPolling inicia um stream de long polling para receber mensagens
// do Retorno Super Inteligente. A PP mantém a conexão aberta até surgir
// uma mensagem ou atingir o timeout interno (long polling HTTP).
//
// arrangement e pspID compõem o path: /api/v1/{arrangement}/{pspID}/tributos/stream/start.
//
// Retorna LongPollingResponse com as mensagens (pode ser NoContent=true se
// não houver mensagens no momento) e o header proximoToken para continuar.
//
// Fonte: Manual de Integração, seção 3.6 — fluxo de long polling com token.
func (c *Client) StartLongPolling(ctx context.Context, arrangement, pspID string) (*LongPollingResponse, *PPErrorResponse, error) {
	path := fmt.Sprintf("/api/v1/%s/%s/tributos/stream/start", arrangement, pspID)
	return c.doLongPolling(ctx, path)
}

// ContinueLongPolling continua um stream de long polling usando o token
// recebido no header proximoToken da resposta anterior.
//
// O token é usado como path: GET /{token}.
// O PSP deve repetir ContinueLongPolling até receber NoContent=true.
//
// Fonte: Manual de Integração, seção 3.6 — continuação com token de posição.
func (c *Client) ContinueLongPolling(ctx context.Context, token string) (*LongPollingResponse, *PPErrorResponse, error) {
	return c.doLongPolling(ctx, token)
}

// EndLongPolling encerra um stream de long polling.
// Deve ser chamado quando o PSP não quiser mais receber mensagens.
// Envia DELETE /{token} para a PP.
//
// Fonte: Manual de Integração, seção 3.6 — encerramento do ciclo.
func (c *Client) EndLongPolling(ctx context.Context, token string) error {
	req, err := c.newRequest(ctx, http.MethodDelete, token, nil)
	if err != nil {
		return fmt.Errorf("criar requisição: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executar requisição: %w", err)
	}
	resp.Body.Close()
	return nil
}

// doRequest executa uma requisição HTTP genérica contra a PP.
// Se a resposta for 2xx, retorna apenas o status code.
// Se a resposta for 4xx/5xx, tenta decodificar como PPErrorResponse (RFC 7807).
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (int, *PPErrorResponse, error) {
	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return 0, nil, fmt.Errorf("criar requisição: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("executar requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return resp.StatusCode, nil, nil
	}

	ppErr, err := decodeError(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("decodificar erro: %w", err)
	}
	return resp.StatusCode, ppErr, nil
}

// doLongPolling executa uma requisição de long polling.
// Trata 204 No Content como resposta válida (sem mensagens no momento).
// Extrai o header proximoToken para continuar o ciclo.
func (c *Client) doLongPolling(ctx context.Context, path string) (*LongPollingResponse, *PPErrorResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("criar requisição: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("executar requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return &LongPollingResponse{NoContent: true, ProximoToken: resp.Header.Get("proximoToken")}, nil, nil
	}

	if resp.StatusCode >= 400 {
		ppErr, err := decodeError(resp.Body)
		if err != nil {
			return nil, nil, fmt.Errorf("decodificar erro: %w", err)
		}
		return nil, ppErr, nil
	}

	var messages RetornoSuperInteligente
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, nil, fmt.Errorf("decodificar resposta: %w", err)
	}

	return &LongPollingResponse{
		Messages:     &messages,
		ProximoToken: resp.Header.Get("proximoToken"),
	}, nil, nil
}

// newRequest cria um http.Request com todos os headers obrigatórios da PP.
//
// Headers injetados:
//   - Content-Type: application/json
//   - Message-Id: UUID4 (novo a cada chamada, idempotência)
//   - Correlation-Id: UUID4 (rastreio ponta a ponta)
//   - Tenant-Id: identificador do PSP (configurado no New)
//   - Timestamp: RFC 3339 (momento da geração da requisição)
//
// Fonte: Manual de Integração, seção 4.1 — headers obrigatórios.
func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := c.baseURL + path

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("codificar body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, &buf)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	messageID := uuid.New().String()
	correlationID := uuid.New().String()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Message-Id", messageID)
	req.Header.Set("Correlation-Id", correlationID)
	req.Header.Set("Tenant-Id", c.tenantID)
	req.Header.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	return req, nil
}

// decodeError decodifica o body de uma resposta de erro no formato RFC 7807.
// Fonte: Manual de Integração, seção 5 — application/problem+json.
func decodeError(r io.Reader) (*PPErrorResponse, error) {
	var ppErr PPErrorResponse
	if err := json.NewDecoder(r).Decode(&ppErr); err != nil {
		return nil, fmt.Errorf("decodificar erro PP: %w", err)
	}
	return &ppErr, nil
}
