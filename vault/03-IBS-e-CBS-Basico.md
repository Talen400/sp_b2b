# IBS e CBS — Básico

## TL;DR
O **IBS** (Imposto sobre Bens e Serviços) substitui ICMS (estadual) e ISS
(municipal) — é administrado pelo **CGIBS** (Comitê Gestor do IBS). A **CBS**
(Contribuição sobre Bens e Serviços) substitui PIS e Cofins — é federal,
administrada pela **RFB** (Receita Federal). Ambos são IVAs não-cumulativos,
cobrados no destino.

## Por que isso importa pra quem programa
Seu sistema precisa saber que são **dois tributos distintos**, com administradores
diferentes, prazos de transição diferentes e alíquotas que variam por localização
do comprador (não do vendedor). Um erro comum é tratar IBS e CBS como um tributo
único ou ignorar a diferença de alíquota por UF/município.

## Funcionamento

### O que cada um substitui

| Tributo Novo | Substitui | Administrador | Esfera |
|---|---|---|---|
| **IBS** | ICMS (estadual) + ISS (municipal) | CGIBS | Estadual/Municipal |
| **CBS** | PIS + Cofins (federal) | RFB | Federal |

✅ `Fonte: EC 132/2023, arts. 156-A (IBS) e 195 (CBS)`

### Por que dois tributos e não um
A reforma unifica cinco tributos em dois, não em um, porque a Constituição exige
que estados/municípios tenham autonomia tributária. O IBS dá a estados e
municípios o poder de definir alíquotas (dentro de faixas determinadas pelo
Senado). A CBS é uniforme em todo o território nacional.
✅ `Fonte: LC 214/2025, arts. 7–12 (IBS) e 13–18 (CBS)`

### Regras comuns
- **Não-cumulatividade plena:** crédito integral na aquisição de insumos.
  ✅ `Fonte: LC 214/2025, art. 28`
- **Cobrança no destino:** o tributo fica com o estado/município do comprador.
  ✅ `Fonte: EC 132/2023, art. 156-A, §1º`
- **Split payment:** modalidade de recolhimento obrigatória (2027+).
  ✅ `Fonte: LC 214/2025, arts. 54–65`

> [!note]- 🔍 Aprofundando: Alíquotas de referência
> O IBS tem **alíquota de referência** definida anualmente pelo Senado Federal
> para cada estado (RF) e cada município (RFM). A alíquota efetiva aplicada a
> uma operação depende da localização do **comprador** (destino). Para a CBS,
> a alíquota é única federal (8,8% em 2027, conforme Anexo II da LC 214/2025).
> ⚠️ As alíquotas definitivas do IBS só se consolidam em 2033 e ainda podem
> mudar com normas complementares.

## Perguntas de autoavaliação
1. IBS substitui quais tributos? E CBS?
2. Por que a reforma criou dois tributos em vez de um só?
3. O que significa "cobrança no destino" na prática?

## Fontes
- ✅ EC 132/2023, arts. 156-A, 195
- ✅ LC 214/2025, arts. 7–12, 13–18, 28, 54–65, Anexo II
