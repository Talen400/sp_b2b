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

// Client is an HTTP client for the Plataforma Pública do Split Payment.
// It injects required headers and handles RFC 7807 errors.
type Client struct {
	baseURL    string
	httpClient *http.Client
	tenantID   string
}

// New creates a new PP client. baseURL is the PP endpoint (e.g. "http://localhost:4010"
// for the local Prism mock). tenantID is the PSP identifier sent as Tenant-Id header.
func New(baseURL, tenantID string) *Client {
	return &Client{
		baseURL:    baseURL,
		tenantID:   tenantID,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// SendInformeTransacaoIniciada sends an Informe de Transação Iniciada to the PP.
// arrangement must be one of: "boleto", "pix-dinamico", "pix-automatico".
// Returns HTTP status code, error response (if any), and error.
func (c *Client) SendInformeTransacaoIniciada(ctx context.Context, arrangement string, req *InformeTransacaoIniciadaRequest) (int, *PPErrorResponse, error) {
	path := fmt.Sprintf("/api/v1/%s", arrangement)
	return c.doRequest(ctx, http.MethodPost, path, req)
}

// SendInformeSegregacao sends an Informe de Segregação (batch) to the PP.
func (c *Client) SendInformeSegregacao(ctx context.Context, req *InformeSegregacaoRequest) (int, *PPErrorResponse, error) {
	return c.doRequest(ctx, http.MethodPost, "/api/v1/segregacao", req)
}

// StartLongPolling initiates a long polling stream for Retorno Super Inteligente.
func (c *Client) StartLongPolling(ctx context.Context, arrangement, pspID string) (*LongPollingResponse, *PPErrorResponse, error) {
	path := fmt.Sprintf("/api/v1/%s/%s/tributos/stream/start", arrangement, pspID)
	return c.doLongPolling(ctx, path)
}

// ContinueLongPolling continues a long polling stream using the token from the previous response.
func (c *Client) ContinueLongPolling(ctx context.Context, token string) (*LongPollingResponse, *PPErrorResponse, error) {
	return c.doLongPolling(ctx, token)
}

// EndLongPolling closes a long polling stream.
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

func decodeError(r io.Reader) (*PPErrorResponse, error) {
	var ppErr PPErrorResponse
	if err := json.NewDecoder(r).Decode(&ppErr); err != nil {
		return nil, fmt.Errorf("decodificar erro PP: %w", err)
	}
	return &ppErr, nil
}
