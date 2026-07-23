# Fontes e Links Oficiais

## TL;DR
Este vault se baseia nas fontes primárias da Reforma Tributária (leis e manuais
oficiais). Aqui estão os links e referências para consulta direta, organizados
por prioridade conforme `REFERENCES.md` do `agent/`.

## 1. Base legal

### EC 132/2023
- **Emenda Constitucional nº 132, de 20/12/2023**
- Cria a base constitucional do IBS, CBS e Imposto Seletivo.
- [Publicação oficial no Planalto](https://www.planalto.gov.br/ccivil_03/constituicao/emendas/emc/emc132.htm)
- ✅ `Fonte primária`

### LC 214/2025
- **Lei Complementar nº 214, de 16/01/2025**
- Regulamenta o modelo: alíquotas reduzidas, regimes especiais, split payment.
- [Publicação oficial no Planalto](https://www.planalto.gov.br/ccivil_03/leis/lcp/lcp214.htm)
- ✅ `Fonte primária`

### LC 227/2026
- **Lei Complementar nº 227, de 2026**
- Institui o CGIBS, processo administrativo do IBS.
- [Publicação oficial no Planalto](https://www.planalto.gov.br/ccivil_03/leis/lcp/lcp227.htm)
- ✅ `Fonte primária`

## 2. Documentação técnica oficial

### Manuais do Split Payment (Ato Conjunto RFB/CGIBS nº 02/2026)
- **Manual de Operações — Split Payment** (v. preliminar, jun/2026)
  — fluxos de negócio, arranjos de pagamento, MOC.
- **Manual de Integração — Plataforma Pública de Split Payment v1.0**
  — dicionário de campos, endpoints, política de erros.
- **OpenAPI v0.0.10** (`openapi-v0_0_10.json`)
  — especificação machine-readable dos endpoints.
- 📚 `Disponíveis em docs-oficiais/ no repositório`

## 3. Portal oficial do governo
- [gov.br/fazenda — Reforma Tributária](https://www.gov.br/fazenda/pt-br/acesso-a-informacao/acoes-e-programas/reforma-tributaria)
- 📚 `Portal oficial — cronograma, consultas públicas, comunicados`

## 4. Imprensa especializada (contexto)
- [Migalhas — Reforma Tributária](https://www.migalhas.com.br)
- [Contábeis](https://www.contabeis.com.br)
- [Jettax](https://www.jettax.com.br)
- [IOB](https://www.iob.com.br)
- 📚 `Usar apenas para contexto — lei/manual sempre prevalecem`

## 5. Nosso simulador Go (projeto irmão)
- O código do simulador está em `cmd/` e `internal/` na raiz do repositório.
- ⚠️ `Não é fonte de verdade tributária — é simplificação didática.
  Ver AGENT.md e TASKS.md do simulador.`

---

> [!note]- 🔍 Aprofundando: Como usar este vault
> Cada nota deste vault cita sua fonte no final usando as tags definidas em
> `AGENT.md` (regra 3). Se você encontrar uma afirmação sem fonte ou marcada
> como `⚠️ não verificado`, questione antes de usar em produção. Sempre
> consulte a fonte primária quando a precisão for crítica.

## Perguntas de autoavaliação
1. Qual a diferença entre o Manual de Operações e o Manual de Integração?
2. Onde encontrar a especificação OpenAPI da Plataforma Pública?
3. O simulador Go pode ser usado como referência para uma implementação real?

## Fontes
- ✅ Planalto — EC 132/2023, LC 214/2025, LC 227/2026
- ✅ Manual de Operações e Manual de Integração (docs-oficiais/)
- 📚 gov.br/fazenda — Reforma Tributária
- ⚠️ nosso simulador (simplificação didática)
