// Package pp implements an HTTP client for the Plataforma Pública do Split Payment
// (RFB/CGIBS/Serpro). It communicates with either a local mock (Prism) or a real
// PP instance, sending informes and receiving Retorno Super Inteligente via long polling.
//
// Usage (dev/example):
//
//	client := pp.New("http://localhost:4010", "PSP-RFB-123456")
//	status, errResp, err := client.SendInformeTransacaoIniciada(ctx, "boleto", &req)
//
// The client injects required headers (Message-Id, Correlation-Id, Tenant-Id, Timestamp)
// automatically on every request.
package pp
