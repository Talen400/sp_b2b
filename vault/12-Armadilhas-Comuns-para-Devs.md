# Armadilhas Comuns para Pessoas Desenvolvedoras

## TL;DR
Modelar split payment em código tem armadilhas específicas validadas pelos manuais oficiais da RFB/CGIBS. Esta nota lista as mais críticas observadas no Manual de Integração e no Manual de Operações.

## Armadilhas

### 1. Usar float/double para valores monetários
A PP especifica Decimal(18,2) representado como **string** no JSON. Em Go, `float64` causa erros de arredondamento.
✅ Solução: `int64` (centavos). Conversão para string com 2 casas decimais na serialização.

### 2. Ignorar os headers obrigatórios
`Message-Id`, `Correlation-Id`, `Tenant-Id` e `Timestamp` são **obrigatórios** em toda requisição à PP. Faltar um deles causa erro 400.
✅ Solução: cliente HTTP que injeta os headers automaticamente. Message-Id e Correlation-Id como UUID4.

### 3. Confundir Informe de Segregação com Informe Preliminar de Pagamento
- **Preliminar:** informativo, não vinculante, uma transação por vez, enviado logo após o pagamento. ⚠️ É descartado quando o Informe de Segregação chega.
- **Segregação:** definitivo, vinculante, em lote (2x/dia útil), gera obrigação de Repasse Financeiro.
✅ Solução: implementar ambos os fluxos corretamente — o Preliminar é quase um "ping" pro governo.

### 4. Assumir endpoint único para todos os arranjos
Cada arranjo tem seu próprio endpoint:
- `POST /api/v1/boleto` (não `/api/v1/transacao-unica`)
- `POST /api/v1/pix-dinamico`
- `POST /api/v1/pix-automatico`
- `POST /api/v1/pix-estatico`
- `POST /api/v1/ted`
- `POST /api/v1/tef`
- `POST /api/v1/segregacao` (este sim é único para lote)
✅ Solução: mapear o arranjo correto no cliente — não existe "endpoint genérico de transação".

### 5. Não tratar retornabilidade
A PP sinaliza se um erro é retentável ou não. Erro 400/422/409 **não deve** ser retentado. Erro 500/503 **deve** ser retentado com backoff exponencial.
✅ Solução: implementar retry logic condicional baseado no `status` e no campo `retornabilidade`.

### 6. Misturar centavos com reais no campo `vlInf`
O campo `vlInf` no payload da PP espera Decimal(18,2) como string. Enviar `1000` (significando R$10,00) quando deveria enviar `"1000.00"` causa rejeição.
✅ Solução: serializador que converte `int64` centavos para string Decimal(18,2). Ex: `100000` → `"1000.00"`.

### 7. Ignorar o Retorno Super Inteligente
Nos arranjos Super Inteligente (Boleto, Pix Dinâmico, Pix Automático), o governo pode corrigir os valores depois do Informe de Transação Iniciada. Ignorar o long polling significa usar valores incorretos no Informe de Segregação.
✅ Solução: implementar o ciclo de long polling (start → continue → delete) antes de enviar o Informe de Segregação.

### 8. Tratar token de posição como opcional
O header `proximoToken` no response do long polling é **obrigatório** para continuar a consulta. Perdê-lo significa recomeçar do início.
✅ Solução: extrair e armazenar o token a cada response. Usar na URL da próxima requisição.

### 9. Não tratar CNPJ alfanumérico
A partir de julho/2026 (IN RFB 2.229/2024), CNPJ pode conter letras. Validar como "apenas 14 dígitos" quebra.
✅ Solução: tratar CNPJ como string livre, sem validação de formato rígido de dígitos.

### 10. Ignorar a diferença de modelo Inteligente vs Super Inteligente
- **Inteligente** (Pix Estático, TED, TEF): sem retorno do governo. Valor informado = valor segregado.
- **Super Inteligente** (Boleto, Pix Dinâmico, Pix Automático): governo pode corrigir. Requer long polling.
Misturar os fluxos causa inconsistência.
✅ Solução: no cliente, bifurcar o comportamento por arranjo (Inteligente: informa e pronto; Super Inteligente: informa + consulta retorno antes de segregar).

### 11. Subestimar o MOC (Mecanismo de Ocorrências)
Falhas de split (ex: split não realizado, valor divergente) precisam ser registradas no MOC. Ignorá-las pode gerar sanções regulatórias.
✅ Solução: após cada informe, verificar se houve erro e registrar ocorrência no endpoint apropriado da PP. `Manual de Operações, seção 7`.

### 12. Esquecer que o PSP Recebedor Direto é o responsável
Mesmo que exista um PSP Recebedor Indireto na relação comercial, o **Direto** é quem responde perante a PP. Implementar a comunicação no PSP errado causa rejeição.
✅ Solução: no `Tenant-Id`, usar sempre o identificador do PSP Recebedor Direto.

## Fontes
- ✅ Manual de Integração v1.0, seções 3.6 (long polling), 4 (dicionário), 4.2 (Decimal), 5 (erros)
- ✅ Manual de Operações, seções 2.1–2.2 (modelos), 3.1 (categorias de valor), 4 (informes), 7 (MOC)
- ✅ IN RFB 2.229/2024 (CNPJ alfanumérico)
